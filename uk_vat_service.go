package vat

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ukVATService is service that calls a UK VAT API to validate UK VAT numbers.
type ukVATService struct{}

// Validate checks if the given VAT number exists and is active. If no error is returned, then it is.
func (s *ukVATService) Validate(vatNumber string) error {
	vatNumber = strings.ToUpper(vatNumber)

	// Only VAT numbers starting with "GB" are supported by this service. All others should go through the VIES service.
	if !strings.HasPrefix(vatNumber, "GB") {
		return ErrInvalidCountryCode
	}

	client := http.Client{
		Timeout: serviceTimeout,
	}
	response, err := client.Get(fmt.Sprintf(ukVATServiceURL, vatNumber[2:]))
	if err != nil {
		return ErrServiceUnavailable{Err: err}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode == http.StatusBadRequest {
		return ErrInvalidVATNumberFormat
	}
	if response.StatusCode == http.StatusNotFound {
		return ErrVATNumberNotFound
	}
	if response.StatusCode != http.StatusOK {
		return ErrServiceUnavailable{
			Err: fmt.Errorf("unexpected status code from UK VAT API: %d", response.StatusCode),
		}
	}

	// If we receive a valid 200 response from this API, it means the VAT number exists and is valid
	return nil
}

// API Documentation:
// https://developer.service.hmrc.gov.uk/api-documentation/docs/api/service/vat-registered-companies-api/1.0
const ukVATServiceURL = "https://api.service.hmrc.gov.uk/organisations/vat/check-vat-number/lookup/%s"
