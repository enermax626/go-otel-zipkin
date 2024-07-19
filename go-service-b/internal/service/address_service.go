package service

import (
	"github.com/enermax626/go-postalcode-temperature/internal/dao"
	"github.com/enermax626/go-postalcode-temperature/internal/model"
	"strconv"
)

const postalCodeSize = 8

type AddressServiceInterface interface {
	FindByPostalCode(postalCode string) (*model.Address, error)
}

type AddressService struct {
	addressDao dao.AddressDaoInterface
}

func NewAddressService(dao dao.AddressDaoInterface) *AddressService {
	return &AddressService{
		addressDao: dao,
	}
}

func (s *AddressService) FindByPostalCode(postalCode string) (*model.Address, error) {
	err := s.isValidPostalCode(postalCode)
	if err != nil {
		return nil, err
	}

	address, err := s.addressDao.FindByPostalCode(postalCode)
	if err != nil {
		return nil, err
	}

	return address, nil
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
