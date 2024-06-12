package vat

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"
)

// LookupServiceInterface is an interface for the service that calls external services to validate VATs.
type LookupServiceInterface interface {
	Validate(vatNumber string) error
}

// viesService validates EU VAT numbers with the VIES service
type viesService struct{}

// Validate returns whether the given VAT number is valid or not
func (s *viesService) Validate(vatNumber string) error {
	if len(vatNumber) < 3 {
		return ErrInvalidVATNumberFormat
	}

	res, err := s.lookup(s.getEnvelope(vatNumber))
	if err != nil {
		return ErrServiceUnavailable{Err: err}
	}
	defer func() {
		_ = res.Body.Close()
	}()

	xmlRes, err := io.ReadAll(res.Body)
	if err != nil {
		return ErrUnexpected{Err: err}
	}

	// check if response contains "INVALID_INPUT" string
	if bytes.Contains(xmlRes, []byte("INVALID_INPUT")) {
		return ErrInvalidVATNumberFormat
	}

	// check if response contains "MS_UNAVAILABLE" string
	if bytes.Contains(xmlRes, []byte("MS_UNAVAILABLE")) ||
		bytes.Contains(xmlRes, []byte("MS_MAX_CONCURRENT_REQ")) {
		return ErrServiceUnavailable{Err: nil}
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
		return ErrServiceUnavailable{Err: err} // assume if response data doesn't match the struct, the service is down
	}

	r := &viesResponse{
		CountryCode: rd.Soap.Soap.CountryCode,
		VATNumber:   rd.Soap.Soap.VATNumber,
		RequestDate: rd.Soap.Soap.RequestDate,
		Valid:       rd.Soap.Soap.Valid,
		Name:        rd.Soap.Soap.Name,
		Address:     rd.Soap.Soap.Address,
	}

	if !r.Valid {
		return ErrVATNumberNotFound
	}
	return nil
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
		Timeout: time.Duration(ViesServiceTimeout) * time.Second,
	}
	return client.Post(viesServiceURL, "text/xml;charset=UTF-8", envelopeBuffer)
}

const viesServiceURL = "https://ec.europa.eu/taxation_customs/vies/services/checkVatService"

// ViesServiceTimeout is the timeout for the VIES service
const ViesServiceTimeout = 10

// viesResponse holds the response data from the Vies call
type viesResponse struct {
	CountryCode string
	VATNumber   string
	RequestDate string
	Valid       bool
	Name        string
	Address     string
}
