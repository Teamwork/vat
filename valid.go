package vat

import (
	"fmt"
	"regexp"
	"strings"
)

// ErrInvalidVATRegion implements the error interface and is returned when
// ValidInRegion fails to find the requested region.  Blank "" regions will not
// fail the region lookup
type ErrInvalidVATRegion struct {
	region string
}

func (e ErrInvalidVATRegion) Error() string {
	return fmt.Sprintf("invalid vat region: %s.  valid regions are (case-insensitive): %q", e.region, validRegions)
}

var regionPatterns = map[string][]*regexp.Regexp{
	"AT":  []*regexp.Regexp{regexp.MustCompile(`^(AT)?U(\d{8})$`)},
	"BE":  []*regexp.Regexp{regexp.MustCompile(`^(BE)?(0?\d{9})$`)},
	"BG":  []*regexp.Regexp{regexp.MustCompile(`^(BG)?(\d{9,10})$`)},
	"CHE": []*regexp.Regexp{regexp.MustCompile(`^(CHE)?(\d{9})(MWST)?$`)},
	"CY":  []*regexp.Regexp{regexp.MustCompile(`^(CY)?([0-5|9]\d{7}[A-Z])$`)},
	"CZ":  []*regexp.Regexp{regexp.MustCompile(`^(CZ)?(\d{8,10})(\d{3})?$`)},
	"DE":  []*regexp.Regexp{regexp.MustCompile(`^(DE)?([1-9]\d{8})$`)},
	"DK":  []*regexp.Regexp{regexp.MustCompile(`^(DK)?(\d{8})$`)},
	"EE":  []*regexp.Regexp{regexp.MustCompile(`^(EE)?(10\d{7})$`)},
	"EL":  []*regexp.Regexp{regexp.MustCompile(`^(EL)?(\d{9})$`)},
	"GR":  []*regexp.Regexp{regexp.MustCompile(`^(EL)?(\d{9})$`)},
	"ES": []*regexp.Regexp{
		regexp.MustCompile(`^(ES)?([A-Z]\d{8})$`),
		regexp.MustCompile(`^(ES)?([A-H|N-S|W]\d{7}[A-J])$`),
		regexp.MustCompile(`^(ES)?([0-9|Y|Z]\d{7}[A-Z])$`),
		regexp.MustCompile(`^(ES)?([K|L|M|X]\d{7}[A-Z])$`),
	},
	"FI": []*regexp.Regexp{regexp.MustCompile(`^(FI)?(\d{8})$`)},
	"FR": []*regexp.Regexp{
		regexp.MustCompile(`^(FR)?(\d{11})$`),
		regexp.MustCompile(`^(FR)?([(A-H)|(J-N)|(P-Z)]\d{10})$`),
		regexp.MustCompile(`^(FR)?(\d[(A-H)|(J-N)|(P-Z)]\d{9})$`),
		regexp.MustCompile(`^(FR)?([(A-H)|(J-N)|(P-Z)]{2}\d{9})$`),
	},
	"GB": []*regexp.Regexp{
		regexp.MustCompile(`^(GB)?(\d{9})$`),
		regexp.MustCompile(`^(GB)?(\d{12})$`),
		regexp.MustCompile(`^(GB)?(GD\d{3})$`),
		regexp.MustCompile(`^(GB)?(HA\d{3})$`),
	},
	"HR": []*regexp.Regexp{
		regexp.MustCompile(`^(HR)?(\d{11})$`),
	},
	"HU": []*regexp.Regexp{
		regexp.MustCompile(`^(HU)?(\d{8})$`),
	},
	"IE": []*regexp.Regexp{
		regexp.MustCompile(`^(IE)?(\d{7}[A-W])(W)?$`),
		regexp.MustCompile(`^(IE)?([7-9][A-Z\*\+)]\d{5}[A-W])$`),
		regexp.MustCompile(`^(IE)?(\d{7}[A-W][AH])$`),
	},
	"IT": []*regexp.Regexp{
		regexp.MustCompile(`^(IT)?(\d{11})$`),
	},
	"LV": []*regexp.Regexp{
		regexp.MustCompile(`^(LV)?(\d{11})$`),
	},
	"LT": []*regexp.Regexp{
		regexp.MustCompile(`^(LT)?(\d{9}|\d{12})$`),
	},
	"LU": []*regexp.Regexp{
		regexp.MustCompile(`^(LU)?(\d{8})$`),
	},
	"MT": []*regexp.Regexp{
		regexp.MustCompile(`^(MT)?([1-9]\d{7})$`),
	},
	"NL": []*regexp.Regexp{
		regexp.MustCompile(`^(NL)?(\d{9})B\d{2}$`),
	},
	"NO": []*regexp.Regexp{
		regexp.MustCompile(`^(NO)?(\d{9})$`),
	},
	"PL": []*regexp.Regexp{
		regexp.MustCompile(`^(PL)?(\d{10})$`),
	},
	"PT": []*regexp.Regexp{
		regexp.MustCompile(`^(PT)?(\d{9})$`),
	},
	"RO": []*regexp.Regexp{
		regexp.MustCompile(`^(RO)?([1-9]\d{1,9})$`),
	},
	"RU": []*regexp.Regexp{
		regexp.MustCompile(`^(RU)?(\d{10}|\d{12})$`),
	},
	"RS": []*regexp.Regexp{
		regexp.MustCompile(`^(RS)?(\d{9})$`),
	},
	"SI": []*regexp.Regexp{
		regexp.MustCompile(`^(SI)?([1-9]\d{7})$`),
	},
	"SK": []*regexp.Regexp{
		regexp.MustCompile(`^(SK)?([1-9]\d[(2-4)|(6-9)]\d{7})$`),
	},
	"SE": []*regexp.Regexp{
		regexp.MustCompile(`^(SE)?(\d{10}01)$`),
	},
	"EU": []*regexp.Regexp{
		regexp.MustCompile(`^(EU)?(\d{9})$`),
	},
}

var validRegions []string

func init() {
	for k := range regionPatterns {
		validRegions = append(validRegions, k)
	}
}

// ValidInRegion is like Valid except it checks only validates against a (case
// insensitive) specific EU region.  Providing a blank EU region will run a
// check against the default EU VAT format.  If the requested region is invalid
// then the validation will fail and ErrInvalidVATRegion will be returned. See
// https://en.wikipedia.org/wiki/VAT_identification_number for more information
// on formats and region prefixes.
func ValidInRegion(vat, region string) (bool, error) {
	region = strings.TrimSpace(strings.ToUpper(region))
	vat = strings.TrimSpace(vat)

	// Blank region checks against the general EU format
	if region == "" {
		region = "EU"
	}

	m, ok := regionPatterns[region]

	if !ok {
		return false, ErrInvalidVATRegion{region}
	}

	return matches(vat, m), nil
}

// Valid checks a VAT number against all known EU regions formats.  Only exact
// matches (start to end of string) will be matched, though the VAT will have
// space trimmed from each end before attempting a match.
func Valid(vat string) bool {
	vat = strings.TrimSpace(vat)
	if vat == "" {
		return false
	}

	for _, m := range regionPatterns {
		if matches(vat, m) {
			return true
		}
	}

	return false
}

func matches(vat string, m []*regexp.Regexp) bool {
	for _, t := range m {
		if t.MatchString(vat) {
			return true
		}
	}

	return false
}
