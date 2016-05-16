// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// How to use filters.
// 1. Create an empty *url.Values
// 2. Add filters and sorts using Add<Filter_or_Sort>
// 3. Pass the url.Values as the final parameter to a request() call
//
// Exemple:
//
// m, err := NewManGoPay(config, mode)
// params := new(url.Values)
// AddNatureFilter(params, NatureRegular)
// m.Transfers(user, params)
//
// Note: not all methods accept query parameters. However, we might have forgotten some of them,
// do not hesitate to open tickets or make pull requests.

package mango

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Filter int
type FilterKey string
type FilterValue string

type mangoFilter struct {
	key    FilterKey
	values []FilterValue
}

const (
	TransactionNature Filter = iota
	TransactionStatus
	TransactionType
	KYCStatus
	KYCType
)

const (
	Nature     FilterKey = "Nature"
	Status               = "Status"
	Type                 = "Type"
	BeforeDate           = "BeforeDate"
	AfterDate            = "AfterDate"
)

const (
	NatureRegular     FilterValue = "REGULAR"
	NatureRefund                  = "REFUND"
	NatureRepudiation             = "REPUDIATION"

	StatusCreated         = "CREATED"
	StatusSucceeded       = "SUCCEEDED"
	StatusFailed          = "FAILED"
	StatusValidationAsked = "VALIDATION_ASKED"
	StatusValidated       = "VALIDATED"
	StatusRefused         = "REFUSED"

	TypePayin                  = "PAYIN"
	TypePayout                 = "PAYOUT"
	TypeTransfer               = "TRANSFER"
	TypeIdentityProof          = "IDENTITY_PROOF"
	TypeRegistrationProof      = "REGISTRATION_PROOF"
	TypeArticlesOfAssociation  = "ARTICLES_OF_ASSOCIATION"
	TypeShareholderDeclaration = "SHAREHOLDER_DECLARATION"
	TypeAddressProof           = "ADDRESS_PROOF"
)

// List of valid values for each filter
var mangoFilters = map[Filter]mangoFilter{
	TransactionNature: {
		key:    Nature,
		values: []FilterValue{NatureRegular, NatureRefund, NatureRepudiation},
	},
	TransactionStatus: {
		key:    Status,
		values: []FilterValue{StatusCreated, StatusSucceeded, StatusFailed},
	},
	TransactionType: {
		key:    Type,
		values: []FilterValue{TypePayin, TypePayout, TypeTransfer},
	},
	KYCStatus: {
		key:    Status,
		values: []FilterValue{StatusValidationAsked, StatusValidated, StatusRefused},
	},
	KYCType: {
		key:    Type,
		values: []FilterValue{TypeIdentityProof, TypeRegistrationProof, TypeArticlesOfAssociation, TypeShareholderDeclaration, TypeAddressProof},
	},
}

// Basic method to check and add a filter based on list of keywords
func addFilter(params *url.Values, filter Filter, val FilterValue) error {
	f, validKey := mangoFilters[filter]
	if !validKey {
		return errors.New(fmt.Sprintf("invalid filter"))
	}
	for _, valid := range f.values {
		if val == valid {
			params.Add(string(f.key), string(val))
			return nil
		}
	}
	return errors.New(fmt.Sprintf("invalid value %s for key %s", val, f.key))
}

// Basic method to add a date filter
func addDateRangeFilter(params *url.Values, key FilterKey, val time.Time) error {
	params.Add(string(key), fmt.Sprintf("%d", val.Unix()))
	return nil
}

// Helpers
func AddTransactionNatureFilter(params *url.Values, val FilterValue) error {
	return addFilter(params, TransactionNature, val)
}

func AddTransactionStatusFilter(params *url.Values, val FilterValue) error {
	return addFilter(params, TransactionStatus, val)
}

func AddTransactionTypeFilter(params *url.Values, val FilterValue) error {
	return addFilter(params, TransactionType, val)
}

func AddKYCStatusFilter(params *url.Values, val FilterValue) error {
	return addFilter(params, KYCStatus, val)
}

func AddKYCTypeFilter(params *url.Values, val FilterValue) error {
	return addFilter(params, KYCType, val)
}

func AddBeforeDateFilter(params *url.Values, val time.Time) error {
	return addDateRangeFilter(params, BeforeDate, val)
}

func AddAfterDateFilter(params *url.Values, val time.Time) error {
	return addDateRangeFilter(params, AfterDate, val)

}
