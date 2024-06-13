package vat

import (
	"fmt"
	"regexp"
	"strings"
)

// Validate validates a VAT number by both format and existence. If no error then it is valid.
func Validate(vatNumber string) error {
	err := ValidateFormat(vatNumber)
	if err != nil {
		return err
	}
	return ValidateExists(vatNumber)
}

// ValidateFormat validates a VAT number by its format. If no error is returned then it is valid.
func ValidateFormat(vatNumber string) error {
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
		// Supposedly the regex for GB numbers is `[0-9]{9}|[0-9]{12}|(GD|HA)[0-9]{3}`,
		// but our validator service only accepts numbers with 9 or 12 digits following the country code.
		// Seems like the official site only accepts 9 digits... https://www.gov.uk/check-uk-vat-number
		"GB": `([0-9]{9}|[0-9]{12})`,
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
		"XI": `([0-9]{9}|[0-9]{12})`, // Northern Ireland, same format as GB
	}

	if len(vatNumber) < 3 {
		return ErrInvalidVATNumberFormat
	}

	vatNumber = strings.ToUpper(vatNumber)

	pattern, ok := patterns[vatNumber[0:2]]
	if !ok {
		return ErrInvalidCountryCode
	}

	matched, err := regexp.MatchString(fmt.Sprintf("^%s$", pattern), vatNumber[2:])
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidVATNumberFormat
	}
	return nil
}

// ValidateExists validates that the given VAT number exists in the external lookup service.
func ValidateExists(vatNumber string) error {
	if len(vatNumber) < 3 {
		return ErrInvalidVATNumberFormat
	}

	vatNumber = strings.ToUpper(vatNumber)

	lookupService := ViesLookupService
	if strings.HasPrefix(vatNumber, "GB") {
		lookupService = UKVATLookupService
	}

	return lookupService.Validate(vatNumber)
}
