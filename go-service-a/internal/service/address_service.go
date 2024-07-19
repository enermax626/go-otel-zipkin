package service

import (
	"service_a/internal/model"
	"strconv"
)

const postalCodeSize = 8

type AddressService struct {
}

func NewAddressService() *AddressService {
	return &AddressService{}
}

func (s *AddressService) isValidPostalCode(postalCode string) error {
	if len(postalCode) != postalCodeSize {
		return model.ErrInvalidPostalCode
	}
	_, err := strconv.Atoi(postalCode)
	if err != nil {
		return model.ErrInvalidPostalCode
	}
	return nil
}
