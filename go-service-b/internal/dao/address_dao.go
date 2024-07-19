package dao

import (
	"encoding/json"
	"fmt"
	"github.com/enermax626/go-postalcode-temperature/internal/model"
	"github.com/valyala/fastjson"
	"io"
	"log"
	"net/http"
	"time"
)

type AddressDaoInterface interface {
	FindByPostalCode(postalCode string) (*model.Address, error)
}

type AddressDao struct {
	client http.Client
}

func NewAddressDao() *AddressDao {
	return &AddressDao{
		client: http.Client{
			Timeout: time.Second * 10,
		},
	}
}

const AddressBaseUrl = "https://viacep.com.br"

func (s *AddressDao) FindByPostalCode(postalCode string) (*model.Address, error) {
	response, err := s.client.Get(fmt.Sprintf(AddressBaseUrl+"/ws/%s/json/", postalCode))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("Error closing http body")
		}
	}(response.Body)

	byteResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsedResponse, err := fastjson.ParseBytes(byteResponse)
	if err != nil {
		return nil, err
	}

	if parsedResponse.Get("erro") != nil {
		return nil, model.ErrPostalCodeNotFound
	}

	var address *model.Address
	err = json.Unmarshal(byteResponse, &address)
	if err != nil {
		return nil, err
	}

	return address, nil
}
