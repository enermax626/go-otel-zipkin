package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/enermax626/go-postalcode-temperature/internal/dao"
	"github.com/enermax626/go-postalcode-temperature/internal/model"
	"github.com/enermax626/go-postalcode-temperature/internal/service"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	trace2 "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	log2 "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func main() {
	if err := run(); err != nil {
		log2.Fatalln(err)
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	tracer := otel.Tracer("service_a_tracer")

	srv := &http.Server{
		Addr:         ":8081",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(tracer),
	}
	srvErr := make(chan error, 1)
	go func() {
		fmt.Println("Server listening...")
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	err = srv.Shutdown(context.Background())
	return
}

func newHTTPHandler(tracer trace2.Tracer) http.Handler {
	mux := http.NewServeMux()

	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	handleFunc("/weather/{postalCode}", findWeatherByPostalCode(tracer))

	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

func findWeatherByPostalCode(tracer trace2.Tracer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		carrier := propagation.HeaderCarrier(r.Header)
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		_, span := tracer.Start(ctx, "service-b_bypostalcode-handler")

		defer span.End()

		var addressService = service.NewAddressService(dao.NewAddressDao())
		weatherService := service.NewWeatherService(addressService, dao.NewWeatherDao())

		postalCode := r.PathValue("postalCode")
		weatherResponse, err := weatherService.FindWeatherByPostalCode(postalCode)
		if err != nil {
			switch err {
			case model.ErrPostalCodeNotFound:
				w.WriteHeader(http.StatusNotFound)
				w.Write(MarshalResponse(ErrorResponse{
					Message: err.Error(),
				}))
			case model.ErrInvalidPostalCode:
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write(MarshalResponse(ErrorResponse{
					Message: err.Error(),
				}))
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(MarshalResponse(ErrorResponse{
					Message: "Internal server error",
				}))
			}
			return
		}
		_, _ = w.Write(MarshalResponse(weatherResponse))
	}
}

func MarshalResponse(res interface{}) []byte {
	marshalledResponse, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return []byte{}
	}
	return marshalledResponse
}

func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	return
}

func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("service_b")))

	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient("otel-collector:4317", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))

	if err != nil {
		return nil, err
	}

	bsp := trace.NewBatchSpanProcessor(traceExporter)

	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
	return traceProvider, nil
}
