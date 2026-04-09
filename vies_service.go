package vat

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"
)

// LookupServiceInterface is an interface for the service that calls external services to validate VATs.
type LookupServiceInterface interface {
	Validate(vatNumber string, opts ValidatorOpts) error
}

// lookupServiceWithResponse is implemented by services that can return the full response.
type lookupServiceWithResponse interface {
	validateWithResponse(vatNumber string, opts ValidatorOpts) (*LookupResponse, error)
}

// viesService validates EU VAT numbers with the VIES service
type viesService struct{}

// Validate returns whether the given VAT number is valid or not
func (s *viesService) Validate(vatNumber string, opts ValidatorOpts) error {
	_, err := s.validateWithResponse(vatNumber, opts)
	return err
}

// validateWithResponse performs validation and returns the full VIES response.
func (s *viesService) validateWithResponse(vatNumber string, _ ValidatorOpts) (*LookupResponse, error) {
	if len(vatNumber) < 3 {
		return nil, ErrInvalidVATNumberFormat
	}

	res, err := s.lookup(s.getEnvelope(vatNumber))
	if err != nil {
		return nil, ErrServiceUnavailable{Err: err}
	}
	defer func() {
		_ = res.Body.Close()
	}()

	xmlRes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ErrServiceUnavailable{Err: err} // assume if we can't read the body then VIES gave us a bad response
	}

	// check if response contains "INVALID_INPUT" string
	if bytes.Contains(xmlRes, []byte("INVALID_INPUT")) {
		return nil, ErrInvalidVATNumberFormat
	}

	// check if response contains "MS_UNAVAILABLE" string
	if bytes.Contains(xmlRes, []byte("MS_UNAVAILABLE")) {
		return nil, ErrServiceUnavailable{Err: errors.New("vies reports service is unavailable")}
	} else if bytes.Contains(xmlRes, []byte("MS_MAX_CONCURRENT_REQ")) {
		return nil, ErrServiceUnavailable{Err: errors.New("max concurrent requests limit hit")}
	}

	var rd struct {
		XMLName xml.Name `xml:"Envelope"`
		Soap    struct {
			XMLName xml.Name `xml:"Body"`
			Soap    struct {
				XMLName     xml.Name `xml:"checkVatResponse"`
				CountryCode string   `xml:"countryCode"`
				VATNumber   string   `xml:"vatNumber"`
				RequestDate string   `xml:"requestDate"` // 2015-03-06+01:00
				Valid       bool     `xml:"valid"`
				Name        string   `xml:"name"`
				Address     string   `xml:"address"`
			}
		}
	}
	if err = xml.Unmarshal(xmlRes, &rd); err != nil {
		return nil, ErrServiceUnavailable{Err: err} // assume if response data doesn't match the struct, the service is down
	}

	r := &VIESResponse{
		CountryCode: rd.Soap.Soap.CountryCode,
		VATNumber:   rd.Soap.Soap.VATNumber,
		RequestDate: rd.Soap.Soap.RequestDate,
		Valid:       rd.Soap.Soap.Valid,
		Name:        rd.Soap.Soap.Name,
		Address:     rd.Soap.Soap.Address,
	}

	if !r.Valid {
		return nil, ErrVATNumberNotFound
	}

	return &LookupResponse{VIESResponse: r}, nil
}

// getEnvelope parses VIES lookup envelope template
func (s *viesService) getEnvelope(n string) string {
	n = strings.ToUpper(n)
	countryCode := n[0:2]
	vatNumber := n[2:]
	const envelopeTemplate = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
<soapenv:Header/>
<soapenv:Body>
  <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
	<countryCode>{{.countryCode}}</countryCode>
	<vatNumber>{{.vatNumber}}</vatNumber>
  </checkVat>
</soapenv:Body>
</soapenv:Envelope>`

	e := envelopeTemplate
	e = strings.Replace(e, "{{.countryCode}}", countryCode, 1)
	e = strings.Replace(e, "{{.vatNumber}}", vatNumber, 1)
	return e
}

// lookup calls the VIES service to get info about the VAT number
func (s *viesService) lookup(envelope string) (*http.Response, error) {
	envelopeBuffer := bytes.NewBufferString(envelope)
	client := http.Client{
		Timeout: serviceTimeout,
	}
	return client.Post(viesServiceURL, "text/xml;charset=UTF-8", envelopeBuffer)
}

const viesServiceURL = "https://ec.europa.eu/taxation_customs/vies/services/checkVatService"

// VIESResponse holds the response data from the VIES service.
type VIESResponse struct {
	CountryCode string `json:"countryCode"`
	VATNumber   string `json:"vatNumber"`
	RequestDate string `json:"requestDate"`
	Valid       bool   `json:"valid"`
	Name        string `json:"name"`
	Address     string `json:"address"`
}
