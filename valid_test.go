package vat

import "testing"

func Test_ValidInRegion(t *testing.T) {
	type in struct {
		vat, region string
	}

	type out struct {
		valid   bool
		err     error
		message string
	}

	var tests = []struct {
		in  in
		out out
	}{
		{
			in: in{
				vat:    "ABC",
				region: "USA",
			},
			out: out{
				valid:   false,
				err:     ErrInvalidVATRegion{"USA"},
				message: "Invalid region",
			},
		},
		{
			in: in{
				vat:    "EU123456789",
				region: "",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "Blank region",
			},
		},
		{
			in: in{
				vat:    "ATU12345678",
				region: "AT",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "Austrian VAT",
			},
		},
		{
			in: in{
				vat:    "ATU12345678",
				region: " AT ",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "Region with spaces",
			},
		},
		{
			in: in{
				vat:    " ATU12345678 ",
				region: "AT",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "VAT with spaces",
			},
		},
		{
			in: in{
				vat:    "ATU1234567",
				region: "AT",
			},
			out: out{
				valid:   false,
				err:     nil,
				message: "Invalid Austrian VAT",
			},
		},
		{
			in: in{
				vat:    "AT12345678",
				region: "AT",
			},
			out: out{
				valid:   false,
				err:     nil,
				message: "Invalid Austrian VAT (No U)",
			},
		},
		{
			in: in{
				vat:    "IE1234567T",
				region: "IE",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "ireland",
			},
		},
		{
			in: in{
				vat:    "IE1234567TW",
				region: "IE",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "ireland woman",
			},
		},
		{
			in: in{
				vat:    "IE1234567FA",
				region: "IE",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "ireland January 2013",
			},
		},
		{
			in: in{
				vat:    "1234567WA",
				region: "IE",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "ireland company",
			},
		},
		{
			in: in{
				vat:    "NL123456789B01",
				region: "NL",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "netherlands",
			},
		},
		{
			in: in{
				vat:    "NL123456789C01",
				region: "NL",
			},
			out: out{
				valid:   false,
				err:     nil,
				message: "netherlands invalid",
			},
		},
		{
			in: in{
				vat:    "ESX12345678",
				region: "ES",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "spain (1) alpha at start",
			},
		},
		{
			in: in{
				vat:    "ES12345678X",
				region: "ES",
			},
			out: out{
				valid:   true,
				err:     nil,
				message: "spain (1) alpha at end",
			},
		},
	}

	for _, test := range tests {
		valid, err := ValidInRegion(test.in.vat, test.in.region)
		if err != test.out.err {
			t.Errorf("expected error to be %q: got %q for %s", test.out.err, err, test.out.message)
		}

		if valid != test.out.valid {
			t.Errorf("expected valid to be %t: got %t for %s", test.out.valid, valid, test.out.message)
		}
	}
}

func Test_Valid(t *testing.T) {
	type out struct {
		valid   bool
		message string
	}

	var tests = []struct {
		in  string
		out out
	}{
		{
			in: "",
			out: out{
				valid:   false,
				message: "blank vat",
			},
		},
		{
			in: "CY99999999L",
			out: out{
				valid:   true,
				message: "cyprus",
			},
		},
		{
			in: "EU123456789",
			out: out{
				valid:   true,
				message: "general eu",
			},
		},
		{
			in: " EU123456789 ",
			out: out{
				valid:   true,
				message: "general eu spaces",
			},
		},
	}

	for _, test := range tests {
		valid := Valid(test.in)
		if valid != test.out.valid {
			t.Errorf("expected valid to be %t: got %t for %s", test.out.valid, valid, test.out.message)
		}
	}
}
