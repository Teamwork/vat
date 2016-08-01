package vat

import "testing"

func Test_Sanitize(t *testing.T) {
	type out struct {
		result  string
		err     error
		message string
	}

	var tests = []struct {
		in  string
		out out
	}{
		{
			in: "CY99999999L",
			out: out{
				result:  "CY99999999L",
				err:     nil,
				message: "cyprus no reformat",
			},
		},
		{
			in: "CY-99999999-L",
			out: out{
				result:  "CY99999999L",
				err:     nil,
				message: "cyprus dashes",
			},
		},
		{
			in: "CY_99999999_L",
			out: out{
				result:  "CY99999999L",
				err:     nil,
				message: "cyprus underscore",
			},
		},
		{
			in: " CY 99999999 L ",
			out: out{
				result:  "CY99999999L",
				err:     nil,
				message: "cyprus spaces",
			},
		},
		{
			in: "CY+99999999+L",
			out: out{
				result:  "CY+99999999+L",
				err:     ErrInvalidVAT{"CY+99999999+L", "CY+99999999+L"},
				message: "cyprus invalid chars",
			},
		},
		{
			in: "",
			out: out{
				result:  "",
				err:     ErrInvalidVAT{"", ""},
				message: "blank",
			},
		},
	}

	for _, test := range tests {
		result, err := Sanitize(test.in)

		if err != test.out.err {
			t.Errorf("expected error to be %q: got %q for %s", test.out.err, err, test.out.message)
		}

		if result != test.out.result {
			t.Errorf("expected result to be %s: got %s for %s", test.out.result, result, test.out.message)
		}
	}
}
