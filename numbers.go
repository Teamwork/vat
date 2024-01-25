package vat

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"regexp"
	"strings"
)

// ViesResponse holds the response data from the Vies call
type ViesResponse struct {
	CountryCode string
	VATNumber   string
	RequestDate string
	Valid       bool
	Name        string
	Address     string
}

// ErrInvalidVATNumber will be returned when an invalid VAT number is passed to a function that validates existence.
var ErrInvalidVATNumber = errors.New("vat: vat number is invalid")

// ErrCountryNotFound indicates that this package could not find a country matching the VAT number prefix.
var ErrCountryNotFound = errors.New("vat: country not found")

// ErrServiceUnavailable will be returned when VIES VAT validation API or jsonvat.com is unreachable.
var ErrServiceUnavailable = errors.New("vat: service is unreachable")

// ValidateNumber validates a VAT number by both format and existence.
// The existence check uses the VIES VAT validation SOAP API and will only run when format validation passes.
func ValidateNumber(n string) (bool, error) {
	isValidFormat, err := ValidateNumberFormat(n)
	existence := false

	if isValidFormat {
		existence, err = ValidateNumberExistence(n)
	}

	return isValidFormat && existence, err
}

// ValidateNumberFormat validates a VAT number by its format.
func ValidateNumberFormat(n string) (bool, error) {
	patterns := map[string]string{
		"AT": `U[A-Z0-9]{8}`,
		"BE": `(0[0-9]{9}|[0-9]{10})`,
		"BG": `[0-9]{9,10}`,
		"CH": `(?:E(?:-| )[0-9]{3}(?:\.| )[0-9]{3}(?:\.| )[0-9]{3}( MWST)?|E[0-9]{9}(?:MWST)?)`,
		"CY": `[0-9]{8}[A-Z]`,
		"CZ": `[0-9]{8,10}`,
		"DE": `[0-9]{9}`,
		"DK": `[0-9]{8}`,
		"EE": `[0-9]{9}`,
		"EL": `[0-9]{9}`,
		"ES": `[A-Z][0-9]{7}[A-Z]|[0-9]{8}[A-Z]|[A-Z][0-9]{8}`,
		"FI": `[0-9]{8}`,
		"FR": `([A-Z]{2}|[0-9]{2})[0-9]{9}`,
		"GB": `[0-9]{9}|[0-9]{12}|(GD|HA)[0-9]{3}`,
		"HR": `[0-9]{11}`,
		"HU": `[0-9]{8}`,
		"IE": `[A-Z0-9]{7}[A-Z]|[A-Z0-9]{7}[A-W][A-I]`,
		"IT": `[0-9]{11}`,
		"LT": `([0-9]{9}|[0-9]{12})`,
		"LU": `[0-9]{8}`,
		"LV": `[0-9]{11}`,
		"MT": `[0-9]{8}`,
		"NL": `[0-9]{9}B[0-9]{2}`,
		"PL": `[0-9]{10}`,
		"PT": `[0-9]{9}`,
		"RO": `[0-9]{2,10}`,
		"SE": `[0-9]{12}`,
		"SI": `[0-9]{8}`,
		"SK": `[0-9]{10}`,
	}

	if len(n) < 3 {
		return false, nil
	}

	n = strings.ToUpper(n)
	pattern, ok := patterns[n[0:2]]
	if !ok {
		return false, ErrCountryNotFound
	}

	matched, err := regexp.MatchString(pattern, n[2:])
	return matched, err
}

// ValidateNumberExistence validates a VAT number by its existence using the VIES VAT API (using SOAP)
func ValidateNumberExistence(vatNumber string) (bool, error) {
	r, err := Lookup(vatNumber, &ViesService{})
	if err != nil {
		return false, err
	}
	return r.Valid, nil
}

// Lookup returns *ViesResponse for a VAT number
func Lookup(vatNumber string, service ViesServiceInterface) (*ViesResponse, error) {
	if len(vatNumber) < 3 {
		return nil, ErrInvalidVATNumber
	}

	res, err := service.Lookup(getEnvelope(vatNumber))
	if err != nil {
		return nil, ErrServiceUnavailable
	}
	defer func() {
		_ = res.Body.Close()
	}()

	xmlRes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// check if response contains "INVALID_INPUT" string
	if bytes.Contains(xmlRes, []byte("INVALID_INPUT")) {
		return nil, ErrInvalidVATNumber
	}

	// check if response contains "MS_UNAVAILABLE" string
	if bytes.Contains(xmlRes, []byte("MS_UNAVAILABLE")) ||
		bytes.Contains(xmlRes, []byte("MS_MAX_CONCURRENT_REQ")) {
		return nil, ErrServiceUnavailable
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
		return nil, err
	}

	r := &ViesResponse{
		CountryCode: rd.Soap.Soap.CountryCode,
		VATNumber:   rd.Soap.Soap.VATNumber,
		RequestDate: rd.Soap.Soap.RequestDate,
		Valid:       rd.Soap.Soap.Valid,
		Name:        rd.Soap.Soap.Name,
		Address:     rd.Soap.Soap.Address,
	}

	return r, nil
}

// getEnvelope parses envelope template
func getEnvelope(n string) string {
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
