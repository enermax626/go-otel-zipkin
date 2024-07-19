package service

import (
	"github.com/enermax626/go-postalcode-temperature/internal/dao"
	"github.com/enermax626/go-postalcode-temperature/internal/dto"
	"net/http"
	"time"
)

type WeatherService struct {
	client         http.Client
	weatherDao     dao.WeatherDaoInterface
	addressService AddressServiceInterface
}

func NewWeatherService(addressService AddressServiceInterface, weatherDao dao.WeatherDaoInterface) *WeatherService {
	return &WeatherService{
		weatherDao:     weatherDao,
		addressService: addressService,
		client: http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *WeatherService) FindWeatherByPostalCode(postalCode string) (*dto.WeatherTemperatureResponse, error) {
	address, err := s.addressService.FindByPostalCode(postalCode)
	if err != nil {
		return nil, err
	}

	normalizedLocalidade, err := NormalizeString(address.Localidade)
	if err != nil {
		return nil, err
	}
	weather, err := s.weatherDao.FindByLocalidade(normalizedLocalidade)
	if err != nil {
		return nil, err
	}
	return dto.NewWeatherTemperatureResponse(weather.Current.TempC, weather.Current.TempF, address.Localidade), nil
}
