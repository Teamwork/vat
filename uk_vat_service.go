package vat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ukVATService is service that calls a UK VAT API to validate UK VAT numbers.
type ukVATService struct{}

// Validate checks if the given VAT number exists and is active. If no error is returned, then it is.
// When opts.IncludeResponse is true, the full UK VAT API response is returned in LookupResponse.UKVATResponse.
func (s *ukVATService) Validate(vatNumber string, opts ValidatorOpts) (*LookupResponse, error) {
	if opts.UKAccessToken == nil || opts.UKAccessToken.IsExpired() {
		// if no access token is provided or if it's expired, try to generate one
		// (it is recommended to generate one separately and cache it and pass it in as an option here)
		accessToken, err := GenerateUKAccessToken(opts)
		if err != nil {
			return nil, ErrMissingUKAccessToken
		}
		opts.UKAccessToken = accessToken
	}

	vatNumber = strings.ToUpper(vatNumber)

	// Only VAT numbers starting with "GB" are supported by this service. All others should go through the VIES service.
	if !strings.HasPrefix(vatNumber, "GB") {
		return nil, ErrInvalidCountryCode
	}

	apiURL := fmt.Sprintf(
		"%s/organisations/vat/check-vat-number/lookup/%s",
		ukVatServiceURL(opts.IsUKTest),
		vatNumber[2:],
	)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, ErrServiceUnavailable{Err: err}
	}

	req.Header.Set("Accept", "application/vnd.hmrc.2.0+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", opts.UKAccessToken.Token))

	client := &http.Client{
		Timeout: serviceTimeout,
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, ErrServiceUnavailable{Err: err}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode == http.StatusBadRequest {
		return nil, ErrInvalidVATNumberFormat
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, ErrVATNumberNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, ErrServiceUnavailable{
			Err: fmt.Errorf("unexpected status code from UK VAT API: %d", response.StatusCode),
		}
	}

	// If we receive a valid 200 response from this API, it means the VAT number exists and is valid
	if opts.IncludeResponse {
		var ukResp UKVATResponse
		if err := json.NewDecoder(response.Body).Decode(&ukResp); err != nil {
			return nil, ErrServiceUnavailable{Err: err}
		}
		return &LookupResponse{UKVATResponse: &ukResp}, nil
	}
	return nil, nil
}

// UKAccessToken contains access token information used to authenticate with the UK VAT API.
type UKAccessToken struct {
	Token               string    `json:"access_token"`
	SecondsUntilExpires int64     `json:"expires_in"`
	ExpiresAt           time.Time `json:"expires_at"`
	IsTest              bool      `json:"-"`
}

// IsExpired checks if the access token is expired.
func (t UKAccessToken) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

// GenerateUKAccessToken generates an access token from given client credentials for use with the UK VAT API.
func GenerateUKAccessToken(opts ValidatorOpts) (*UKAccessToken, error) {
	if opts.UKClientID == "" || opts.UKClientSecret == "" {
		return nil, ErrUnableToGenerateUKAccessToken{Err: errors.New("missing client ID or secret")}
	}

	data := url.Values{}
	data.Set("client_secret", opts.UKClientSecret)
	data.Set("client_id", opts.UKClientID)
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "read:vat")

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oauth/token", ukVatServiceURL(opts.IsUKTest)),
		bytes.NewBufferString(data.Encode()),
	)
	if err != nil {
		return nil, ErrUnableToGenerateUKAccessToken{Err: err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrUnableToGenerateUKAccessToken{Err: err}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ErrUnableToGenerateUKAccessToken{Err: fmt.Errorf("unexpected status code: %d", resp.StatusCode)}
		}
		return nil, ErrUnableToGenerateUKAccessToken{
			Err: fmt.Errorf("unexpected status code: %d: %s", resp.StatusCode, respBody),
		}
	}

	var token UKAccessToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, ErrUnableToGenerateUKAccessToken{Err: err}
	}
	// set the expiration time to 1 minute before the actual expiration for safety
	token.ExpiresAt = time.Now().Add(time.Duration(token.SecondsUntilExpires-60) * time.Second)
	if opts.IsUKTest {
		token.IsTest = true
	}

	return &token, nil
}

func ukVatServiceURL(isTest bool) string {
	if isTest {
		return fmt.Sprintf("https://test-%s", ukVATServiceDomain)
	}
	return fmt.Sprintf("https://%s", ukVATServiceDomain)
}

// API Documentation:
// https://developer.service.hmrc.gov.uk/api-documentation/docs/api/service/vat-registered-companies-api/2.0/oas/page
const ukVATServiceDomain = "api.service.hmrc.gov.uk"

// UKVATResponse holds the response data from the UK VAT API.
type UKVATResponse struct {
	Target struct {
		Name      string       `json:"name"`
		VATNumber string       `json:"vatNumber"`
		Address   UKVATAddress `json:"address"`
	} `json:"target"`
	Requester          string `json:"requester,omitempty"`
	ConsultationNumber string `json:"consultationNumber,omitempty"`
	ProcessingDate     string `json:"processingDate"`
}

// UKVATAddress holds the address data from the UK VAT API response.
type UKVATAddress struct {
	Line1       string `json:"line1"`
	Line2       string `json:"line2,omitempty"`
	Line3       string `json:"line3,omitempty"`
	Line4       string `json:"line4,omitempty"`
	Postcode    string `json:"postcode"`
	CountryCode string `json:"countryCode"`
}
