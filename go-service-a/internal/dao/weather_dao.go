package dao

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"service_a/internal/dto"
	"service_a/internal/model"
	"time"
)

var serviceBUrl = "http://service_b:8081"

//var serviceBUrl = "http://localhost:8081"

type WeatherDaoInterface interface {
	FindByPostalCode(ctx context.Context, postalCode string) (*dto.WeatherTemperatureResponse, error)
}

type WeatherDao struct {
	client http.Client
}

func NewWeatherDao() WeatherDaoInterface {
	return &WeatherDao{
		client: http.Client{
			Timeout: time.Second * 3,
		},
	}
}

func (d *WeatherDao) FindByPostalCode(ctx context.Context, postalCode string) (*dto.WeatherTemperatureResponse, error) {
	req, err := http.NewRequest("GET", serviceBUrl+"/weather/"+postalCode, nil)
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, model.ErrPostalCodeNotFound
	case http.StatusUnprocessableEntity:
		return nil, model.ErrInvalidPostalCode
	}

	var weatherResponse dto.WeatherTemperatureResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return nil, err
	}

	return &weatherResponse, nil
}
