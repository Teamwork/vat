//go:build integration
// +build integration

package vat

import (
	"errors"
	"testing"
)

// test numbers to use with the UK VAT service API in their sandbox environment:
// https://github.com/hmrc/vat-registered-companies-api/blob/main/public/api/conf/2.0/test-data/vrn.csv

var ukTests = []struct {
	vatNumber     string
	expectedError error
}{
	{"GB553557881", nil}, // valid VAT number in the sandbox environment
	{"GB0472429986", ErrInvalidVATNumberFormat},
	{"Hi", ErrInvalidCountryCode},
	{"GB333289453", ErrVATNumberNotFound},
}

// TestUKVATService tests the UK VAT service. Just meant to be a quick way to check that this service is working.
// Makes external calls that sometimes might fail. Do not include them in CI/CD.
func TestUKVATService(t *testing.T) {
	opts := ValidatorOpts{
		UKClientID:     "yourClientID", // insert your own client ID and secret here. Do not commit real values.
		UKClientSecret: "yourClientSecret",
		IsUKTest:       true,
	}
	token, err := GenerateUKAccessToken(opts)
	if err != nil {
		t.Fatal(err)
	}
	opts.UKAccessToken = token.Token

	for _, test := range ukTests {
		err := UKVATLookupService.Validate(test.vatNumber, opts)
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Expected <%v> for %v, got <%v>", test.expectedError, test.vatNumber, err)
		}
	}
}
