package mango

import (
	"encoding/json"
	"errors"
)

// List of banking aliases.
type BankingAliasList []*BankingAlias

type BankingAlias struct {
	ProcessIdent

	CreditedUserId string
	WalletId       string
	Type           string // IBAN
	Country        string
	OwnerName      string
	Active         bool
	IBAN           string
	BIC            string

	service *MangoPay
}

// NewBankingAlias creates a new banking alias.
//
// See https://docs.mangopay.com/endpoints/v2.01/banking-aliases
func (m *MangoPay) NewBankingAlias(wallet Wallet, ownerName string, country string) (*BankingAlias, error) {
	if wallet.Id == "" {
		return nil, errors.New("wallet has empty Id")
	}
	b := &BankingAlias{
		ProcessIdent: ProcessIdent{},
		OwnerName:    ownerName,
		Country:      country,
		WalletId:     wallet.Id,
		service:      m,
	}
	return b, nil
}

// Save sends the HTTP query to create the bank alias.
func (b *BankingAlias) Save() error {
	data := JsonObject{}
	j, err := json.Marshal(b)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	// Data fields to remove before sending the HTTP request.
	for _, field := range []string{"Id", "CreationDate"} {
		delete(data, field)
	}

	ba, err := b.service.anyRequest(new(BankingAlias), actionCreateBankingAlias, data)
	if err != nil {
		return err
	}
	serv := b.service
	*b = *(ba.(*BankingAlias))
	b.service = serv

	return nil
}

// BankingAlias returns a user's banking alias.
func (m *MangoPay) BankingAlias(id string) (*BankingAlias, error) {
	w, err := m.anyRequest(new(BankingAlias), actionFetchBankingAlias,
		JsonObject{"BankingAliasId": id})
	if err != nil {
		return nil, err
	}
	return w.(*BankingAlias), nil
}

// BankingAliases finds all user's bank aliases.
func (m *MangoPay) BankingAliases(wallet Wallet) (BankAccountList, error) {
	if wallet.Id == "" {
		return nil, errors.New("wallet has empty Id")
	}
	accs, err := m.anyRequest(new(BankingAliasList), actionFetchBankingAliases,
		JsonObject{"WalletId": wallet.Id})
	if err != nil {
		return nil, err
	}
	return *(accs.(*BankAccountList)), nil
}
