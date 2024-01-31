package vat

import (
	"bytes"
	"net/http"
	"time"
)

// ViesServiceInterface is an interface for the service that calls VIES.
type ViesServiceInterface interface {
	Lookup(envelope string) (*http.Response, error)
}

// ViesService implements ViesServiceInterface using the HTTP client.
type ViesService struct{}

// Lookup calls the VIES service to get info about the VAT number
func (s *ViesService) Lookup(envelope string) (*http.Response, error) {
	envelopeBuffer := bytes.NewBufferString(envelope)
	client := http.Client{
		Timeout: time.Duration(ServiceTimeout) * time.Second,
	}
	return client.Post(serviceURL, "text/xml;charset=UTF-8", envelopeBuffer)
}

const serviceURL = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"

// ServiceTimeout indicates the number of seconds before a service request times out.
var ServiceTimeout = 10
