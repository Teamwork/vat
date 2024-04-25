//go:generate mockgen -destination=mocks/mock_lookup_service.go -package=mocks github.com/teamwork/vat/v3 LookupServiceInterface

package vat

import _ "github.com/golang/mock/mockgen/model" // Blank import to ensure the model is included in the binary.
