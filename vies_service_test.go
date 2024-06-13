//go:build integration
// +build integration

package vat

import "testing"

var viesTests = []struct {
	vatNumber     string
	expectedError error
}{
	{"BE0472429986", nil}, // valid number
	{"NL123456789B01", ErrVATNumberNotFound},
	{"Hi", ErrInvalidVATNumberFormat},
	{"INVALID INPUT", ErrInvalidVATNumberFormat},
}

// TestViesService tests the VIES service. Just meant to be a quick way to check that this service is working.
// The external VIES calls are not always reliable so sometimes these tests may fail. Do not include them in CI/CD.
func TestViesService(t *testing.T) {
	for _, test := range viesTests {
		err := ViesLookupService.Validate(test.vatNumber)
		if err != test.expectedError {
			t.Errorf("Expected <%v> for %v, got <%v>", test.expectedError, test.vatNumber, err)
		}
	}
}
