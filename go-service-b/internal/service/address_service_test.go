package service

import (
	"github.com/enermax626/go-postalcode-temperature/internal/dao"
	"github.com/enermax626/go-postalcode-temperature/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindByPostalCodeWhenValidPostalCodeThenGetAddress(t *testing.T) {

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

	addressDaoMock := &dao.AddressDaoMock{}
	addressDaoMock.On("FindByPostalCode", postalCode).Return(expectedAddress, nil)

	addressService := NewAddressService(addressDaoMock)

	address, err := addressService.FindByPostalCode(postalCode)

	assert.Nil(t, err)
	assert.Equal(t, expectedAddress.Cep, address.Cep)
}

func TestFindByPostalCodeWhenInvalidPostalCodeThenErrInvalidPostalCode(t *testing.T) {

	postalCode := "123453"

	addressDaoMock := &dao.AddressDaoMock{}
	//addressDaoMock.On("FindByPostalCode", postalCode).Return(expectedAddress, nil)

	addressService := NewAddressService(addressDaoMock)

	address, err := addressService.FindByPostalCode(postalCode)

	assert.Nil(t, address)
	assert.Equal(t, model.ErrInvalidPostalCode, err)
}

func TestFindByPostalCodeWhenPostalCodeNotFoundThenErrPostalCodeNotFound(t *testing.T) {

	postalCode := "12345678"

	addressDaoMock := &dao.AddressDaoMock{}
	addressDaoMock.On("FindByPostalCode", postalCode).Return(nil, model.ErrPostalCodeNotFound)

	addressService := NewAddressService(addressDaoMock)

	address, err := addressService.FindByPostalCode(postalCode)

	assert.Nil(t, address)
	assert.Equal(t, model.ErrPostalCodeNotFound, err)
}
