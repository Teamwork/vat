//go:build integration
// +build integration

package vat

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
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
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	opts := ValidatorOpts{
		UKClientID:     os.Getenv("CLIENT_ID"),
		UKClientSecret: os.Getenv("SECRET"),
		IsUKTest:       true,
	}

	if opts.UKClientID == "" || opts.UKClientSecret == "" {
		t.Fatal("CLIENT_ID and SECRET must be set in .env file")
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
		_, err := UKVATLookupService.Validate(test.vatNumber, opts)
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Expected <%v> for %v, got <%v>", test.expectedError, test.vatNumber, err)
		}
	}
}

func TestUKVATServiceIncludeResponse(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	opts := ValidatorOpts{
		UKClientID:      os.Getenv("CLIENT_ID"),
		UKClientSecret:  os.Getenv("SECRET"),
		IsUKTest:        true,
		IncludeResponse: true,
	}

	if opts.UKClientID == "" || opts.UKClientSecret == "" {
		t.Fatal("CLIENT_ID and SECRET must be set in .env file")
	}

	token, err := GenerateUKAccessToken(opts)
	if err != nil {
		t.Fatal(err)
	}
	opts.UKAccessToken = token

	resp, err := UKVATLookupService.Validate("GB553557881", opts)
	if err != nil {
		t.Fatalf("expected no error for valid VAT, got %v", err)
	}
	if resp == nil || resp.UKVATResponse == nil {
		t.Fatal("expected UKVATResponse to be populated when IncludeResponse is true")
	}
	if resp.UKVATResponse.Target.VATNumber != "553557881" {
		t.Errorf("expected Target.VATNumber '553557881', got %q", resp.UKVATResponse.Target.VATNumber)
	}
	if resp.UKVATResponse.Target.Name == "" {
		t.Error("expected Target.Name to be populated")
	}
	if resp.UKVATResponse.ProcessingDate == "" {
		t.Error("expected ProcessingDate to be populated")
	}
}
