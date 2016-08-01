package vat

import (
	"fmt"
	"regexp"
)

// ErrInvalidVAT implements the error interface and is returned when Sanitize
// fails to clean and validate a VAT number
type ErrInvalidVAT struct {
	vat, sanitized string
}

var sanitizeVAT = regexp.MustCompile(`[-_ ]`)

func (e ErrInvalidVAT) Error() string {
	return fmt.Sprintf("invalid vat: %s.  sanitized: %s", e.vat, e.sanitized)
}

// Sanitize removes whitespace and invalid characters from a VAT number (often
// provided by end users to "enhance readabilty").  After replacement the VAT
// number is validated.  If the VAT validation fails ErrInvalidVat will be
// returned.  All (and only) spaces, dashes and underscores will be removed.
func Sanitize(vat string) (string, error) {
	sanitized := sanitizeVAT.ReplaceAllString(vat, "")

	if Valid(sanitized) {
		return sanitized, nil
	}

	return sanitized, ErrInvalidVAT{vat, sanitized}
}
