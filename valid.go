// Package vat validates and sanitizes VAT numbers belonging to EU (and non-EU)
// regions according to the current formats listed here:
// https://en.wikipedia.org/wiki/VAT_identification_number
package vat

import (
	"regexp"
	"strings"

	vatapi "github.com/mattes/vat"
)

// ValidInRegion is like Valid except it checks only validates against a (case
// insensitive) specific EU region.  Providing a blank EU region will run a
// check against the default EU VAT format.  If the requested region is invalid
// then the validation will fail and ErrInvalidVATRegion will be returned. See
// https://en.wikipedia.org/wiki/VAT_identification_number for more information
// on formats and region prefixes.
func ValidInRegion(vat, region string) (bool, error) {
	region = strings.TrimSpace(strings.ToUpper(region))
	vat = strings.TrimSpace(vat)

	// Blank region checks against the general EU format
	if region == "" {
		region = "EU"
	}

	m, ok := regionPatterns[region]

	if !ok {
		return false, ErrInvalidVATRegion{region}
	}

	return matches(vat, m), nil
}

// Valid checks a VAT number against all known EU regions formats.  Only exact
// matches (start to end of string) will be matched, though the VAT will have
// space trimmed from each end before attempting a match.  Note: Valid provides
// a loose match as exact algoritms for each country are technically not
// published and even if they were you cannot guarantee that the number has
// actually been issued. If a guarantee is needed, use MustValidate
func Valid(vat string) bool {
	vat = strings.TrimSpace(vat)
	if vat == "" {
		return false
	}

	for _, m := range regionPatterns {
		if matches(vat, m) {
			return true
		}
	}

	return false
}

// MustValidate provides a real-time lookup of the europa API to ensure that the
// VAT number provided is currently active.  Unfortunately the API is not
// entirely reliable so a retry should be issued on error.  The default timeout
// for the API call is 10 seconds.  Unlike stdlib functions prefixed with
// `Must`, this function does not panic on error, but more signifies the
// importance of a true real-time lookup.
func MustValidate(vat string) (*vatapi.VATresponse, error) {
	vat, err := Sanitize(vat)
	if err != nil {
		return nil, err
	}

	return vatapi.CheckVAT(vat)
}

func matches(vat string, m []*regexp.Regexp) bool {
	for _, t := range m {
		if t.MatchString(vat) {
			return true
		}
	}

	return false
}
