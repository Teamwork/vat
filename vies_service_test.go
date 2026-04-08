//go:build integration
// +build integration

package vat

import (
	"errors"
	"testing"
	"time"
)

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
		_, err := ViesLookupService.Validate(test.vatNumber, ValidatorOpts{})
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Expected <%v> for %v, got <%v>", test.expectedError, test.vatNumber, err)
		}
		time.Sleep(time.Second * 2) // delay to prevent rate limiting
	}
}

func TestViesServiceIncludeResponse(t *testing.T) {
	resp, err := ViesLookupService.Validate("BE0472429986", ValidatorOpts{IncludeResponse: true})
	if err != nil {
		t.Fatalf("expected no error for valid VAT, got %v", err)
	}
	if resp == nil || resp.VIESResponse == nil {
		t.Fatal("expected VIESResponse to be populated when IncludeResponse is true")
	}
	if resp.VIESResponse.CountryCode != "BE" {
		t.Errorf("expected CountryCode 'BE', got %q", resp.VIESResponse.CountryCode)
	}
	if resp.VIESResponse.VATNumber != "0472429986" {
		t.Errorf("expected VATNumber '0472429986', got %q", resp.VIESResponse.VATNumber)
	}
	if !resp.VIESResponse.Valid {
		t.Error("expected Valid to be true")
	}
	if resp.VIESResponse.Name == "" {
		t.Error("expected Name to be populated")
	}
}
