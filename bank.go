// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"errors"
)

// Bank account type.
type AccountType int

const (
	IBAN AccountType = iota
	GB
	US
	CA
	OTHER
)

var accountTypes = map[AccountType]string{
	IBAN:  "IBAN",
	GB:    "GB",
	US:    "US",
	CA:    "CA",
	OTHER: "OTHER",
}

// List of bank accounts.
type BankAccountList []*BankAccount

// BankAccount is an item mainly used for pay-out bank wire request. It is
// used as a generic bank account container for all supported account
// types: IBAN, GB, US, CA or OTHER.
//
// This way, only one structure is used to unmarshal any JSON response related
// to bank accounts.
type BankAccount struct {
	ProcessIdent
	Type         string // IBAN, GB, US, CA or OTHER
	OwnerName    string
	OwnerAddress string
	UserId       string
	// Required for IBAN type
	Iban          string
	Bic           string // For IBAN, OTHER
	AccountNumber string // For GB, US, CA, OTHER
	// Required for GB type
	SortCode string
	// Required for US type
	Aba string
	// Required for CA type
	BankName          string
	InstitutionNumber string
	BranchCode        string
	// Required for OTHER type
	Country string

	service *MangoPay
}

type IbanBankAccount struct {
}

func (b *BankAccount) String() string {
	return struct2string(b)
}

// NewBankAccount creates a new bank account. Note that depending on the account's
// type, some fields of the newly BankAccount instance must be filled (they are
// required) before a call to Save().
//
// See http://docs.mangopay.com/api-references/bank-accounts/
func (m *MangoPay) NewBankAccount(user Consumer, ownerName, ownerAddress string, t AccountType) (*BankAccount, error) {
	id := ""
	switch user.(type) {
	case *LegalUser:
		id = user.(*LegalUser).Id
	case *NaturalUser:
		id = user.(*NaturalUser).Id
	}
	if id == "" {
		return nil, errors.New("user has empty Id")
	}
	b := &BankAccount{
		ProcessIdent: ProcessIdent{},
		Type:         accountTypes[t],
		OwnerName:    ownerName,
		OwnerAddress: ownerAddress,
		UserId:       id,
		service:      m,
	}
	return b, nil
}

func (b *BankAccount) Save() error {
	data := JsonObject{}
	j, err := json.Marshal(b)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	// Force float64 to int conversion after unmarshalling.
	for _, field := range []string{"CreationDate"} {
		data[field] = int(data[field].(float64))
	}

	// Fields not allowed when creating a tranfer.
	for _, field := range []string{"Id", "CreationDate", "UserId"} {
		delete(data, field)
	}

	ba, err := b.service.anyRequest(new(BankAccount), actionCreateBankAccount, data)
	if err != nil {
		return err
	}
	serv := b.service
	*b = *(ba.(*BankAccount))
	b.service = serv

	return nil
}
