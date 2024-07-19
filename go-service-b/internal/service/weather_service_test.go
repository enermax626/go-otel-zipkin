package service

import (
	"github.com/enermax626/go-postalcode-temperature/internal/dao"
	"github.com/enermax626/go-postalcode-temperature/internal/dto"
	"github.com/enermax626/go-postalcode-temperature/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindByPostalCodeWhenValidPostalCodeThenGetWeather(t *testing.T) {

	postalCode := "12345123"
	expectedAddress := &model.Address{
		Cep:         "12345123",
		Logradouro:  "Rua Exemplo",
		Complemento: "Apto 101",
		Unidade:     "1",
		Bairro:      "Centro",
		Localidade:  "Cidade Exemplo",
		Uf:          "EX",
		Ibge:        "1234567",
		Gia:         "1234",
		Ddd:         "12",
		Siafi:       "1234",
	}

	expectedWeatherLocationResponse := &dto.WeatherLocationResponse{
		Current: dto.CurrentTemperature{
			TempC: 28.6,
			TempF: 83.5,
		},
	}

	expectedTempK := 301.6

	addressServiceMock := &AddressServiceMock{}
	weatherDaoMock := &dao.WeatherDaoMock{}
	weatherService := NewWeatherService(addressServiceMock, weatherDaoMock)

	addressServiceMock.On("FindByPostalCode", postalCode).Return(expectedAddress, nil)
	weatherDaoMock.On("FindByLocalidade", expectedAddress.Localidade).Return(expectedWeatherLocationResponse, nil)

	weatherTemperatureResponse, err := weatherService.FindWeatherByPostalCode(postalCode)

	assert.Nil(t, err)
	assert.Equal(t, expectedWeatherLocationResponse.Current.TempC, weatherTemperatureResponse.TempC)
	assert.Equal(t, expectedWeatherLocationResponse.Current.TempF, weatherTemperatureResponse.TempF)
	assert.Equal(t, expectedTempK, weatherTemperatureResponse.TempK)
}
