//go:build integration
// +build integration

package vat

import "testing"

var ukTests = []struct {
	vatNumber     string
	expectedError error
}{
	{"GB333289454", nil}, // valid number (as of 2024-04-26)
	{"GB0472429986", ErrInvalidVATNumberFormat},
	{"Hi", ErrInvalidCountryCode},
	{"GB333289453", ErrVATNumberNotFound},
}

// TestUKVATService tests the UK VAT service. Just meant to be a quick way to check that this service is working.
// Makes external calls that sometimes might fail. Do not include them in CI/CD.
func TestUKVATService(t *testing.T) {
	for _, test := range ukTests {
		err := UKVATLookupService.Validate(test.vatNumber)
		if err != test.expectedError {
			t.Errorf("Expected <%v> for %v, got <%v>", test.expectedError, test.vatNumber, err)
		}
	}
}
