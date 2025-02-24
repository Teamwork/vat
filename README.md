Package vat
===

![Build](https://github.com/Teamwork/vat/actions/workflows/build.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/teamwork/vat/v3)](https://goreportcard.com/report/github.com/teamwork/vat/v3)
[![GoDoc](https://godoc.org/github.com/teamwork/vat/v3?status.svg)](https://godoc.org/github.com/teamwork/vat/v3)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/teamwork/vat/master/LICENSE)

Package for validating VAT numbers & retrieving VAT rates (
from [ibericode/vat-rates](https://github.com/ibericode/vat-rates)) in Go.

Based on https://github.com/dannyvankooten/vat

## Installation

Use go get.

``` 
go get github.com/teamwork/vat
```

Then import the package into your own code.

```
import "github.com/teamwork/vat"
```

## Usage

### Validating VAT numbers

VAT numbers can be validated by format, existence or both.

EU VAT numbers are looked up using the [VIES VAT validation API](http://ec.europa.eu/taxation_customs/vies/).

UK VAT numbers are looked up
using [UK GOV VAT validation API](https://developer.service.hmrc.gov.uk/api-documentation/docs/api/service/vat-registered-companies-api/1.0)
(requires [signing up for the UK API](#accessing-the-uk-vat-api)).

```go
package main

import "github.com/teamwork/vat"

func main() {
	// These validation functions return an error if the VAT number is invalid. If no error, then it is valid.

	// Validate number by format + existence
	err := vat.Validate("NL123456789B01")

	// Validate number format
	err := vat.ValidateFormat("NL123456789B01")

	// Validate number existence
	err := vat.ValidateExists("NL123456789B01")
}
```

### Retrieving VAT rates

> This package relies on a [community maintained repository of vat rates](https://github.com/ibericode/vat-rates). We
> invite you to toggle notifications for that repository and contribute changes to VAT rates in your country once they
> are announced.

To get VAT rate periods for a country, first get a CountryRates struct using the country's ISO-3166-1-alpha2 code.

You can get the rate that is currently in effect using the `GetRate` function.

```go
package main

import (
	"fmt"
	"github.com/teamwork/vat"
)

func main() {
	c, err := vat.GetCountryRates("IE")
	r, err := c.GetRate("standard")

	fmt.Printf("Standard VAT rate for IE is %.2f", r)
	// Output: Standard VAT rate for IE is 23.00
}
```

# Accessing the UK VAT API

For validating VAT numbers that begin with "GB" you will need
to [sign up to gain access to the UK government's VAT API](https://developer.service.hmrc.gov.uk/api-documentation/docs/using-the-hub).

Once you have signed up and acquire a client ID and client secret, here's how to generate and use an access token:

```go
package main

import "github.com/teamwork/vat"

func main() {
	var ukAccessToken *vat.UKAccessToken
	ukAccessToken, err := vat.GenerateUKAccessToken(vat.ValidatorOpts{
		UKClientID:     "yourClientID",
		UKClientSecret: "yourClientSecret",
		IsUKTest:       true, // set this if you are testing this in their sandbox test API
	})
	if err != nil {
		panic(err)
	}
	// Recommended to cache the access token until it expires

	err := vat.Validate("GB123456789", vat.ValidatorOpts{
		UKAccessToken: ukAccessToken.Token,
		IsUKTest:      true, // if token created in test mode, run validation in test mode
	})
}

```

## License

MIT licensed. See the LICENSE file for details.
