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
	IBAN          string
	BIC           string // For IBAN, OTHER
	AccountNumber string // For GB, US, CA, OTHER
	// Required for GB type
	SortCode string
	// Required for US type
	ABA string
	// Required for CA type
	BankName          string
	InstitutionNumber string
	BranchCode        string
	// Required for OTHER type
	Country string

	service *MangoPay
	atype   AccountType
}

type IBANBankAccount struct {
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
		atype:        t,
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

	// Data fields to remove before sending the HTTP request.
	ignore := []string{"Id", "CreationDate"}
	switch b.atype {
	case IBAN:
		if b.IBAN == "" || b.BIC == "" {
			return errors.New("missing full IBAN information")
		}
		ignore = append(ignore, "AccountNumber", "SortCode", "ABA", "BankName",
			"InstitutionNumber", "BranchCode", "Country")
	case GB:
		if b.AccountNumber == "" || b.SortCode == "" {
			return errors.New("missing full GB information")
		}
		ignore = append(ignore, "IBAN", "BIC", "ABA", "BankName",
			"InstitutionNumber", "BranchCode", "Country")
	case US:
		if b.AccountNumber == "" || b.ABA == "" {
			return errors.New("missing full US information")
		}
		ignore = append(ignore, "IBAN", "BIC", "SortCode", "BankName",
			"InstitutionNumber", "BranchCode", "Country")
	case CA:
		if b.BankName == "" || b.InstitutionNumber == "" || b.BranchCode == "" ||
			b.AccountNumber == "" {
			return errors.New("missing full CA information")
		}
		ignore = append(ignore, "IBAN", "BIC", "SortCode", "ABA", "Country")
	case OTHER:
		if b.AccountNumber == "" || b.BIC == "" {
			return errors.New("missing full OTHER information")
		}
		ignore = append(ignore, "IBAN", "ABA", "SortCode", "BankName",
			"InstitutionNumber", "BranchCode")
	}

	for _, field := range ignore {
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

// BankAccount returns a user's bank account.
func (m *MangoPay) BankAccount(user Consumer, id string) (*BankAccount, error) {
	userId := ""
	switch user.(type) {
	case *LegalUser:
		userId = user.(*LegalUser).Id
	case *NaturalUser:
		userId = user.(*NaturalUser).Id
	}
	if userId == "" {
		return nil, errors.New("user has empty Id")
	}
	w, err := m.anyRequest(new(BankAccount), actionFetchBankAccount,
		JsonObject{"Id": id, "UserId": userId})
	if err != nil {
		return nil, err
	}
	return w.(*BankAccount), nil
}

// BankAccounts finds all user's bank accounts.
func (m *MangoPay) BankAccounts(user Consumer) (BankAccountList, error) {
	userId := ""
	switch user.(type) {
	case *LegalUser:
		userId = user.(*LegalUser).Id
	case *NaturalUser:
		userId = user.(*NaturalUser).Id
	}
	if userId == "" {
		return nil, errors.New("user has empty Id")
	}
	accs, err := m.anyRequest(new(BankAccountList), actionFetchUserBankAccounts, JsonObject{"Id": userId})
	if err != nil {
		return nil, err
	}
	return *(accs.(*BankAccountList)), nil
}
