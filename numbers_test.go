package vat

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/teamwork/vat/mocks"
)

var tests = []struct {
	number string
	valid  bool
}{
	{"", false},
	{"A", false},
	{"AB123A01", false},
	{"ATU12345678", true},
	{"ATU15673009", true},
	{"ATU1234567", false},
	{"BE0123456789", true},
	{"BE1234567891", true},
	{"BE0999999999", true},
	{"BE9999999999", true},
	{"BE012345678", false},
	{"BE123456789", false},
	{"BG123456789", true},
	{"BG1234567890", true},
	{"BG1234567", false},
	{"CHE-156.730.098 MWST", true},
	{"CHE-156.730.098", true},
	{"CHE156730098MWST", true},
	{"CHE156730098", true},
	{"CY12345678X", true},
	{"CY15673009L", true},
	{"CY1234567X", false},
	{"CZ12345678", true},
	{"CZ1234567", false},
	{"DE123456789", true},
	{"DE12345678", false},
	{"DK12345678", true},
	{"DK1234567", false},
	{"EE123456789", true},
	{"EE12345678", false},
	{"EL123456789", true},
	{"EL12345678", false},
	{"ESX12345678", true},
	{"ESX1234567", false},
	{"FI1234567", false},
	{"FI12345678", true},
	{"FR12345678901", true},
	{"FR1234567890", false},
	{"GB999999973", true},
	{"GB156730098481", true},
	{"GBGD549", true},
	{"GBHA549", true},
	{"GB99999997", false},
	{"HU12345678", true},
	{"HU1234567", false},
	{"HR12345678901", true},
	{"HR1234567890", false},
	{"IE1234567X", true},
	{"IE123456X", false},
	{"IT12345678901", true},
	{"IT1234567890", false},
	{"LT123456789", true},
	{"LT12345678", false},
	{"LU26375245", true},
	{"LU12345678", true},
	{"LU1234567", false},
	{"LV12345678901", true},
	{"LV1234567890", false},
	{"MT12345678", true},
	{"MT1234567", false},
	{"NL123456789B01", true},
	{"NL123456789B12", true},
	{"NL12345678B12", false},
	{"PL1234567890", true},
	{"PL123456789", false},
	{"PT123456789", true},
	{"PT12345678", false},
	{"RO123456789", true},
	{"RO1", false}, // Romania has a really weird VAT format...
	{"SE123456789012", true},
	{"SE12345678901", false},
	{"SI12345678", true},
	{"SI1234567", false},
	{"SK1234567890", true},
	{"SK123456789", false},
	{"KL123456789B12", false},
	{"NL12343678B12", false},
	{"PL1224567890", true},
	{"PL123456989", false},
	{"PT125456789", true},
	{"PT16345678", false},
	{"RO123456789", true},
	{"KT123456789", false},
	{"LT12335678", false},
	{"LU26275245", true},
	{"LU14345678", true},
	{"LU1234567", false},
	{"LQ12345678901", false},
	{"QE123456789", false},
	{"DE12375670", false},
	{"DK123365678", true},
	{"DO1231567", false},
	{"EE123456789", true},
	{"EL123434678", true},
	{"EL122456789", true},
	{"EL12245678", false},

	{"CY1134567X", false},
	{"CZ17345678", true},
	{"CZ1239567", false},
	{"DE123456789", true},
	{"DE12325678", false},
	{"DK12385678", true},
	{"DK1234767", false},
	{"EE123456689", true},
	{"EE12343678", false},
	{"EL123426789", true},
	{"EL12342678", false},
	{"ESX12315678", true},
	{"ESX1634567", false},
	{"FI1274567", false},
	{"FI12385678", true},
	{"FR12349678901", true},
	{"FR1234597890", false},
	{"GB999979973", true},
	{"GB15673098481", true},
	{"GBGD529", true},
	{"GBHA519", true},
	{"GB99997997", false},
	{"HU12342678", true},
	{"HU1234167", false},
	{"HR12341678901", true},
	{"HR1234267890", false},
	{"IE1234167X", true},
	{"IE123356X", false},
	{"IT12245678901", true},
	{"IT1231567890", false},
	{"LT123156789", true},
	{"LT12343678", false},
	{"LU26371245", true},
	{"LU12325678", true},
	{"LU1214567", false},
	{"LV12345578901", true},
	{"LV1234367890", false},
	{"MT12325678", true},
	{"MT1234167", false},
	{"NL123156789B01", true},
	{"NL113456789B12", true},
	{"NL12315678B12", false},
	{"PO1234567890", false},
}

func BenchmarkValidateFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ValidateNumberFormat("NL" + strconv.Itoa(i))
	}
}

func TestValidateNumberFormat(t *testing.T) {
	for _, test := range tests {
		valid, err := ValidateNumberFormat(test.number)
		if err != nil {
			if err.Error() != ErrCountryNotFound.Error() {
				panic(err)
			}
		}

		if test.valid != valid {
			t.Errorf("Expected %v for %v, got %v", test.valid, test.number, valid)
		}

	}
}

var lookupTests = []struct {
	fullVatNumber    string
	vatNumber        string
	countryCode      string
	isValid          bool
	expectedResponse string
	expectedError    *error
}{
	{"BE0472429986", "BE", "0472429986", true, "", nil},
	{"NL123456789B01", "NL", "123456789B01", false, "", nil},
	{"Hi", "", "Hi", false, "", &ErrInvalidVATNumber},
	{"INVALID INPUT", "IN", "VALID INPUT", false, "INVALID_INPUT", &ErrInvalidVATNumber},
	{"BE0472429986", "BE", "0472429986", true, "MS_UNAVAILABLE", &ErrServiceUnavailable},
	{"BE0472429986", "BE", "0472429986", true, "MS_MAX_CONCURRENT_REQ", &ErrServiceUnavailable},
}

func TestLookupValidateNumberExistence(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockViesService := mocks.NewMockViesServiceInterface(ctrl)

	for _, test := range lookupTests {
		if len(test.fullVatNumber) >= 3 {
			expectedViesResponseBody := test.expectedResponse
			if expectedViesResponseBody == "" {
				expectedViesResponseBody = expectedViesResponse(test.vatNumber, test.countryCode, test.isValid)
			}
			expectedViesResponse := &http.Response{
				Body:       io.NopCloser(bytes.NewBufferString(expectedViesResponseBody)),
				StatusCode: http.StatusOK,
			}

			mockViesService.EXPECT().
				Lookup(expectedViesRequestEnvelope(test.vatNumber, test.countryCode)).
				Return(expectedViesResponse, nil)
		}

		viesResponse, err := Lookup(test.fullVatNumber, mockViesService)
		if err != nil {
			if test.expectedError == nil {
				t.Error(fmt.Errorf("unexpected error: %w", err))
			} else if err != *test.expectedError {
				t.Error(fmt.Errorf("Expected error: %w\nGot: %s", *test.expectedError, err.Error()))
			}
		} else {
			if viesResponse == nil {
				t.Error("expected a ViesResponse, none received")
			}
			if test.isValid && !viesResponse.Valid {
				t.Error(fmt.Sprintf("expected %s to be a valid VAT number.", test.fullVatNumber))
			}
			if !test.isValid && viesResponse.Valid {
				t.Error(fmt.Sprintf("expected %s to not be a valid VAT number.", test.fullVatNumber))
			}
		}
	}
}

func expectedViesRequestEnvelope(vatNumber string, countryCode string) string {
	return fmt.Sprintf(
		`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
<soapenv:Header/>
<soapenv:Body>
  <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
	<countryCode>%s</countryCode>
	<vatNumber>%s</vatNumber>
  </checkVat>
</soapenv:Body>
</soapenv:Envelope>`,
		vatNumber,
		countryCode,
	)
}

func expectedViesResponse(vatNumber string, countryCode string, isValid bool) string {
	isValidString := "true"
	if !isValid {
		isValidString = "false"
	}

	return fmt.Sprintf(
		`<env:Envelope xmlns:env="http://schemas.xmlsoap.org/soap/envelope/">
			<env:Header/>
			<env:Body>
				<ns2:checkVatResponse xmlns:ns2="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
					<ns2:countryCode>%s</ns2:countryCode>
					<ns2:vatNumber>%s</ns2:vatNumber>
					<ns2:requestDate>2024-01-24+01:00</ns2:requestDate>
					<ns2:valid>%s</ns2:valid>
					<ns2:name>VAT HOLDER'S NAME'</ns2:name>
					<ns2:address>VAT HOLDER'S ADDRESS</ns2:address>
				</ns2:checkVatResponse>
			</env:Body>
		</env:Envelope>`,
		vatNumber,
		countryCode,
		isValidString,
	)
}
