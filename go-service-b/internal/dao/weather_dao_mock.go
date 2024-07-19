package dao

import (
	"github.com/enermax626/go-postalcode-temperature/internal/dto"
	"github.com/stretchr/testify/mock"
)

type WeatherDaoMock struct {
	mock.Mock
}

func (m *WeatherDaoMock) FindByLocalidade(localidade string) (*dto.WeatherLocationResponse, error) {
	args := m.Mock.Called(localidade)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.WeatherLocationResponse), args.Error(1)
	}
	return nil, args.Error(1)
}
