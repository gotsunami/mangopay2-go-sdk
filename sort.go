// Copyright 2016 Go Tsunami. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"errors"
	"fmt"
	"net/url"
)

type SortKey string
type SortDirection string

const (
	CreationDate  SortKey = "CreationDate"
	ExecutionDate         = "ExecutionDate"
	Date                  = "Date"
)

const (
	Ascending  SortDirection = ":asc"
	Descending SortDirection = ":desc"
)

// Basic method to add a sort
func addSort(params *url.Values, key SortKey, direction SortDirection) error {
	if direction != "" || direction != Ascending || direction != Descending {
		return errors.New(fmt.Sprintf("Invalid direction for sort '%s'", direction))
	}
	params.Add("Sort", fmt.Sprintf("%s%s", string(key), string(direction)))
	return nil
}

// Helpers
// Only available sorts have a method.
// See https://docs.mangopay.com/api-references/sort-lists/
func SortWalletTransactionsByCreationDate(params *url.Values, direction SortDirection) error {
	return addSort(params, CreationDate, direction)
}

func SortUserTransactionsByCreationDate(params *url.Values, direction SortDirection) error {
	return addSort(params, CreationDate, direction)
}

func SortWalletTransactionsByExecutionDate(params *url.Values, direction SortDirection) error {
	return addSort(params, ExecutionDate, direction)
}

func SortUserTransactionsByExecutionDate(params *url.Values, direction SortDirection) error {
	return addSort(params, ExecutionDate, direction)
}

func SortKYCByCreationDate(params *url.Values, direction SortDirection) error {
	return addSort(params, CreationDate, direction)
}

func SortCardsByCreationDate(params *url.Values, direction SortDirection) error {
	return addSort(params, CreationDate, direction)
}

func SortUsersByCreationDate(params *url.Values, direction SortDirection) error {
	return addSort(params, CreationDate, direction)
}

func SortBankAccountsByCreationDate(params *url.Values, direction SortDirection) error {
	return addSort(params, CreationDate, direction)
}

func SortEventsByDate(params *url.Values, direction SortDirection) error {
	return addSort(params, Date, direction)
}
