package dao

import (
	"encoding/json"
	"fmt"
	"github.com/enermax626/go-postalcode-temperature/internal/dto"
	"github.com/enermax626/go-postalcode-temperature/internal/model"
	"io"
	"log"
	"net/http"
	"time"
)

const WeatherBaseUrl = "https://api.weatherapi.com"

// API Key used for demonstration purposes, not valid anymore.
const WeatherAPIKey = "c32de2806aeb4644b8b145420241107"

type WeatherDaoInterface interface {
	FindByLocalidade(localidade string) (*dto.WeatherLocationResponse, error)
}

type WeatherDao struct {
	client http.Client
}

func NewWeatherDao() *WeatherDao {
	return &WeatherDao{
		client: http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *WeatherDao) FindByLocalidade(localidade string) (*dto.WeatherLocationResponse, error) {
	response, err := s.client.Get(fmt.Sprintf(WeatherBaseUrl+"/v1/current.json?q=%s&key=%s", localidade, WeatherAPIKey))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, model.ErrPostalCodeNotFound
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("Error closing http body")
		}
	}(response.Body)

	var weatherResponse dto.WeatherLocationResponse
	err = json.NewDecoder(response.Body).Decode(&weatherResponse)
	if err != nil {
		return nil, err
	}
	return &weatherResponse, nil
}
