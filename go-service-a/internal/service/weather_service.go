package service

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"service_a/internal/dao"
	"service_a/internal/dto"
	"time"
)

type WeatherServiceInterface interface {
	FindByPostalCode(ctx context.Context, postalCode string) (*dto.WeatherTemperatureResponse, error)
}

type WeatherService struct {
	tracer         trace.Tracer
	addressService AddressService
	weatherDao     dao.WeatherDaoInterface
}

func NewWeatherService(weatherDao dao.WeatherDaoInterface, addressService *AddressService, tracer trace.Tracer) WeatherServiceInterface {
	return &WeatherService{
		tracer:         tracer,
		addressService: *addressService,
		weatherDao:     weatherDao,
	}
}

func (s *WeatherService) FindByPostalCode(ctx context.Context, postalCode string) (*dto.WeatherTemperatureResponse, error) {
	_, span := s.tracer.Start(ctx, "weatherbypostalcode-weather-service")
	defer span.End()
	err := s.addressService.isValidPostalCode(postalCode)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	weather, err := s.weatherDao.FindByPostalCode(ctx, postalCode)
	if err != nil {
		return nil, err
	}
	return weather, nil
}
