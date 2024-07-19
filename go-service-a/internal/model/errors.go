package model

import "errors"

var ErrPostalCodeNotFound = errors.New("can not find zipcode")
var ErrInvalidPostalCode = errors.New("invalid zipcode")
