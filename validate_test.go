package vat

import (
	"errors"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
)

var tests = []struct {
	number        string
	expectedError error
}{
	{"", ErrInvalidVATNumberFormat},
	{"A", ErrInvalidVATNumberFormat},
	{"AB123A01", ErrInvalidCountryCode},
	{"ATU12345678", nil},
	{"ATU15673009", nil},
	{"ATU1234567", ErrInvalidVATNumberFormat},
	{"BE0123456789", nil},
	{"BE1234567891", nil},
	{"BE0999999999", nil},
	{"BE9999999999", nil},
	{"BE012345678", ErrInvalidVATNumberFormat},
	{"BE123456789", ErrInvalidVATNumberFormat},
	{"BG123456789", nil},
	{"BG1234567890", nil},
	{"BG1234567", ErrInvalidVATNumberFormat},
	{"CHE-156.730.098 MWST", nil},
	{"CHE-156.730.098", nil},
	{"CHE156730098MWST", nil},
	{"CHE156730098", nil},
	{"CY12345678X", nil},
	{"CY15673009L", nil},
	{"CY1234567X", ErrInvalidVATNumberFormat},
	{"CZ12345678", nil},
	{"CZ1234567", ErrInvalidVATNumberFormat},
	{"DE123456789", nil},
	{"DE12345678", ErrInvalidVATNumberFormat},
	{"DK12345678", nil},
	{"DK1234567", ErrInvalidVATNumberFormat},
	{"EE123456789", nil},
	{"EE12345678", ErrInvalidVATNumberFormat},
	{"EL123456789", nil},
	{"EL12345678", ErrInvalidVATNumberFormat},
	{"ESX12345678", nil},
	{"ESX1234567", ErrInvalidVATNumberFormat},
	{"FI1234567", ErrInvalidVATNumberFormat},
	{"FI12345678", nil},
	{"FR12345678901", nil},
	{"FR1B345678901", nil},
	{"FRAB345678901", nil},
	{"FRABC45678901", ErrInvalidVATNumberFormat},
	{"FR1234567890", ErrInvalidVATNumberFormat},
	{"GB999999973", nil},
	{"GB156730098481", nil},
	{"GBGD549", ErrInvalidVATNumberFormat},
	{"GBHA549", ErrInvalidVATNumberFormat},
	{"GB99999997", ErrInvalidVATNumberFormat},
	{"HU12345678", nil},
	{"HU1234567", ErrInvalidVATNumberFormat},
	{"HR12345678901", nil},
	{"HR1234567890", ErrInvalidVATNumberFormat},
	{"IE1234567X", nil},
	{"IE123456X", ErrInvalidVATNumberFormat},
	{"IT12345678901", nil},
	{"IT1234567890", ErrInvalidVATNumberFormat},
	{"LT123456789", nil},
	{"LT12345678", ErrInvalidVATNumberFormat},
	{"LU26375245", nil},
	{"LU12345678", nil},
	{"LU1234567", ErrInvalidVATNumberFormat},
	{"LV12345678901", nil},
	{"LV1234567890", ErrInvalidVATNumberFormat},
	{"MT12345678", nil},
	{"MT1234567", ErrInvalidVATNumberFormat},
	{"NL123456789B01", nil},
	{"NL123456789B12", nil},
	{"NL12345678B12", ErrInvalidVATNumberFormat},
	{"PL1234567890", nil},
	{"PL123456789", ErrInvalidVATNumberFormat},
	{"PT123456789", nil},
	{"PT12345678", ErrInvalidVATNumberFormat},
	{"RO123456789", nil},
	{"RO1", ErrInvalidVATNumberFormat}, // Romania has a really weird VAT format...
	{"SE123456789012", nil},
	{"SE12345678901", ErrInvalidVATNumberFormat},
	{"SI12345678", nil},
	{"SI1234567", ErrInvalidVATNumberFormat},
	{"SK1234567890", nil},
	{"SK123456789", ErrInvalidVATNumberFormat},
	{"KL123456789B12", ErrInvalidCountryCode},
	{"NL12343678B12", ErrInvalidVATNumberFormat},
	{"PL1224567890", nil},
	{"PL123456989", ErrInvalidVATNumberFormat},
	{"PT125456789", nil},
	{"PT16345678", ErrInvalidVATNumberFormat},
	{"RO123456789", nil},
	{"KT123456789", ErrInvalidCountryCode},
	{"LT12335678", ErrInvalidVATNumberFormat},
	{"LU26275245", nil},
	{"LU14345678", nil},
	{"LU1234567", ErrInvalidVATNumberFormat},
	{"LQ12345678901", ErrInvalidCountryCode},
	{"QE123456789", ErrInvalidCountryCode},
	{"DE12375670", ErrInvalidVATNumberFormat},
	{"DK123365678", ErrInvalidVATNumberFormat},
	{"DO1231567", ErrInvalidCountryCode},
	{"EE123456789", nil},
	{"EL123434678", nil},
	{"EL122456789", nil},
	{"EL12245678", ErrInvalidVATNumberFormat},
	{"CY1134567X", ErrInvalidVATNumberFormat},
	{"CZ17345678", nil},
	{"CZ1239567", ErrInvalidVATNumberFormat},
	{"DE123456789", nil},
	{"DE12325678", ErrInvalidVATNumberFormat},
	{"DK12385678", nil},
	{"DK1234767", ErrInvalidVATNumberFormat},
	{"EE123456689", nil},
	{"EE12343678", ErrInvalidVATNumberFormat},
	{"EL123426789", nil},
	{"EL12342678", ErrInvalidVATNumberFormat},
	{"ESX12315678", nil},
	{"ESX1634567", ErrInvalidVATNumberFormat},
	{"FI1274567", ErrInvalidVATNumberFormat},
	{"FI12385678", nil},
	{"FR12349678901", nil},
	{"FR1234597890", ErrInvalidVATNumberFormat},
	{"GB999979973", nil},
	{"GB156730984812", nil},
	{"GBGD529", ErrInvalidVATNumberFormat},
	{"GBHA519", ErrInvalidVATNumberFormat},
	{"GB99997997", ErrInvalidVATNumberFormat},
	{"GB15673098481", ErrInvalidVATNumberFormat},
	{"HU12342678", nil},
	{"HU1234167", ErrInvalidVATNumberFormat},
	{"HR12341678901", nil},
	{"HR1234267890", ErrInvalidVATNumberFormat},
	{"IE1234167X", nil},
	{"IE123356X", ErrInvalidVATNumberFormat},
	{"IT12245678901", nil},
	{"IT122456789013333", ErrInvalidVATNumberFormat},
	{"IT1231567890", ErrInvalidVATNumberFormat},
	{"LT123156789", nil},
	{"LT12343678", ErrInvalidVATNumberFormat},
	{"LU26371245", nil},
	{"LU12325678", nil},
	{"LU1214567", ErrInvalidVATNumberFormat},
	{"LV12345578901", nil},
	{"LV1234367890", ErrInvalidVATNumberFormat},
	{"MT12325678", nil},
	{"MT1234167", ErrInvalidVATNumberFormat},
	{"NL123156789B01", nil},
	{"NL113456789B12", nil},
	{"NL12315678B12", ErrInvalidVATNumberFormat},
	{"PO1234567890", ErrInvalidCountryCode},
}

func BenchmarkValidateFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ValidateFormat("NL" + strconv.Itoa(i))
	}
}

func TestValidateFormat(t *testing.T) {
	for _, test := range tests {
		err := ValidateFormat(test.number)
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Expected <%v> for %v, got <%v>", test.expectedError, test.number, err)
		}
	}
}

func TestValidateExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockViesService := NewMockLookupServiceInterface(ctrl)
	mockUKVATService := NewMockLookupServiceInterface(ctrl)
	ViesLookupService = mockViesService
	UKVATLookupService = mockUKVATService

	defer restoreLookupServices()

	var lookupTests = []struct {
		vatNumber     string
		service       *MockLookupServiceInterface
		expectedError error
	}{
		{"BE0472429986", mockViesService, ErrVATNumberNotFound},
		{"NL123456789B01", mockViesService, nil},
		{"Hi", mockViesService, ErrInvalidVATNumberFormat},
		{"INVALID INPUT", mockViesService, ErrInvalidVATNumberFormat},
		{"BE0472429986", mockViesService, ErrServiceUnavailable{Err: nil}},
		{"GB333289454", mockUKVATService, nil},
		// XI is Northern Ireland, which while part of the UK is actually still validated by VIES
		{"XI0472429986", mockViesService, nil},
	}

	for _, test := range lookupTests {
		if len(test.vatNumber) >= 3 {
			test.service.EXPECT().Validate(test.vatNumber, ValidatorOpts{}).Return(test.expectedError)
		}

		err := ValidateExists(test.vatNumber)
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Expected <%v> for %v, got <%v>", test.expectedError, test.vatNumber, err)
		}
	}
}

func restoreLookupServices() {
	ViesLookupService = &viesService{}
	UKVATLookupService = &ukVATService{}
}
