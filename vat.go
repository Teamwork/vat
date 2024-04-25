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

var ViesLookupService LookupServiceInterface = &viesService{}
var UKVATLookupService LookupServiceInterface = &ukVATService{}
