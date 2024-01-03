Package vat
===

![Build](https://github.com/Teamwork/vat/actions/workflows/build.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Teamwork/vat)](https://goreportcard.com/report/github.com/Teamwork/vat)
[![GoDoc](https://godoc.org/github.com/Teamwork/vat?status.svg)](https://godoc.org/github.com/Teamwork/vat)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/teamwork/vat/master/LICENSE)

Package for validating VAT numbers & retrieving VAT rates (from [ibericode/vat-rates](https://github.com/ibericode/vat-rates)) in Go.

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

VAT numbers can be validated by format, existence or both. VAT numbers are looked up using the [VIES VAT validation API](http://ec.europa.eu/taxation_customs/vies/).

```go
package main

import "github.com/teamwork/vat"

func main() {
  // Validate number by format + existence
  validity, err := vat.ValidateNumber("NL123456789B01")

  // Validate number format
  validity, err := vat.ValidateNumberFormat("NL123456789B01")

  // Validate number existence
  validity, err := vat.ValidateNumberExistence("NL123456789B01")
}
```

### Retrieving VAT rates

> This package relies on a [community maintained repository of vat rates](https://github.com/ibericode/vat-rates). We invite you to toggle notifications for that repository and contribute changes to VAT rates in your country once they are announced.

To get VAT rate periods for a country, first get a CountryRates struct using the country's ISO-3166-1-alpha2 code.

You can get the rate that is currently in effect using the `GetRate` function.

```go
package main

import (
  "fmt"
  "github.com/teamwork/vat"
)

func main() {
  c, err := vat.GetCountryRates("NL")
  r, err := c.GetRate("standard")

  fmt.Printf("Standard VAT rate for NL is %.2f", r)
  // Output: Standard VAT rate for NL is 21.00
}
```

## License

MIT licensed. See the LICENSE file for details.
