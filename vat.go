/*
Package vat helps you deal with European VAT in Go.

It offers VAT number validation using the VIES VAT validation API & VAT rates retrieval using jsonvat.com
It also offers UK VAT number validation using a UK VAT API.

Validate a VAT number

	err := vat.Validate("NL123456789B01")

Get VAT rate that is currently in effect for a given country

	c, _ := vat.GetCountryRates("NL")
	r, _ := c.GetRate("standard")
*/
package vat

import "time"

// ViesLookupService is the interface for the VIES VAT number validation service
var ViesLookupService LookupServiceInterface = &viesService{}

// UKVATLookupService is the interface for the UK VAT number validation service
var UKVATLookupService LookupServiceInterface = &ukVATService{}

var serviceTimeout = time.Duration(60) * time.Second

// SetServiceTimeout sets the timeout in seconds for the external VAT lookup services.
func SetServiceTimeout(seconds int) {
	serviceTimeout = time.Duration(seconds) * time.Second
}
