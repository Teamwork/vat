package vat

import (
	"errors"
	"fmt"
)

// ErrInvalidVATNumberFormat will be returned if a VAT with an invalid format is given
var ErrInvalidVATNumberFormat = errors.New("vat: VAT number format is invalid")

// ErrVATNumberNotFound will be returned if the given VAT number is not found in the external lookup service
var ErrVATNumberNotFound = errors.New("vat: number not found as an existing active VAT number")

// ErrInvalidCountryCode indicates that this package could not find a country matching the VAT number prefix
var ErrInvalidCountryCode = errors.New("vat: unknown country code")

// ErrInvalidRateLevel will be returned when getting wrong rate level
var ErrInvalidRateLevel = errors.New("vat: unknown rate level")

// ErrServiceUnavailable will be returned when the service is unreachable
type ErrServiceUnavailable struct {
	Err error
}

// Error returns the error message
func (e ErrServiceUnavailable) Error() string {
	return fmt.Sprintf("vat: service is unreachable: %v", e.Err)
}

// ErrUnableToGenerateUKAccessToken will be returned if the UK API Access token could not be generated
type ErrUnableToGenerateUKAccessToken struct {
	Err error
}

// Error returns the error message
func (e ErrUnableToGenerateUKAccessToken) Error() string {
	return fmt.Sprintf("vat: Error generating UK API Access token: %v", e.Err)
}

// ErrMissingUKAccessToken will be returned if the UK API Access token is missing
var ErrMissingUKAccessToken = errors.New(
	"vat: missing UK API Access token. Run `vat.GenerateUKAccessToken` to generate one",
)
