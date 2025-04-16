//go:build integration
// +build integration

package vat

import (
	"errors"
	"testing"
	"time"
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
	// expect the token to be valid for 14400 seconds (4 hours),
	// but with ExpireAt 60 seconds earlier to help ensure we generate a new one before old one expires
	expectedExpiresAt := time.Now().Add(time.Duration(14400-60) * time.Second)
	isExpiresAtCorrect := expectedExpiresAt.Sub(token.ExpiresAt).Seconds() <= 1 // 1 second leeway for test to run

	if token.SecondsUntilExpires != 14400 || !isExpiresAtCorrect {
		t.Errorf("Expected token to be valid for 14400 seconds, got %d", token.SecondsUntilExpires)
	}
	opts.UKAccessToken = token

	for _, test := range ukTests {
		err := UKVATLookupService.Validate(test.vatNumber, opts)
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Expected <%v> for %v, got <%v>", test.expectedError, test.vatNumber, err)
		}
	}
}
