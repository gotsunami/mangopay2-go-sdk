// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"errors"
	"fmt"
)

// List of wallets.
type WalletList []*Wallet

// List of wallet's owners.
type ConsumerList []Consumer

// Money specifies which currency and amount (in cents!) to use in
// a payment transaction.
type Money struct {
	Currency string
	Amount   int // in cents, i.e 120 for 1.20 EUR
}

func (b *Money) String() string {
	return fmt.Sprintf("%.2f %s", float64(b.Amount/100), b.Currency)
}

// Wallet stores all payins and tranfers from users in order to
// collect money.
type Wallet struct {
	ProcessIdent
	Owners      []string
	Description string
	Currency    string
	Balance     Money
	service     *MangoPay
}

func (u *Wallet) String() string {
	return fmt.Sprintf(`
Id            : %s
Tag           : %s
Owners        : %s
Description   : %s
Currency      : %s
Balance       : %s
Creation date : %s
`, u.Id, u.Tag, u.Owners, u.Description, u.Currency, u.Balance.String(), unixTimeToString(u.CreationDate))
}

// NewWallet creates a new wallet. Owners must have a well-defined Id. Empty Ids will
// return an error.
func (m *MangoPay) NewWallet(owners ConsumerList, desc string, currency string) (*Wallet, error) {
	all := []string{}
	for k, o := range owners {
		id := ""
		switch o.(type) {
		case *LegalUser:
			id = o.(*LegalUser).Id
		case *NaturalUser:
			id = o.(*NaturalUser).Id
		}
		if id == "" {
			return nil, errors.New(fmt.Sprintf("Empty Id for owner %d. Unable to create wallet.", k))
		}
		all = append(all, id)
	}
	w := &Wallet{
		Owners:      all,
		Description: desc,
		Currency:    currency,
	}
	w.service = m
	return w, nil
}

// Save creates or updates a legal user. The Create API is used
// if the user's Id is an empty string. The Edit API is used when
// the Id is a non-empty string.
func (w *Wallet) Save() error {
	var action mangoAction
	if w.Id == "" {
		action = actionCreateWallet
	} else {
		action = actionEditWallet
	}

	data := JsonObject{}
	j, err := json.Marshal(w)
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

	// Fields not allowed when creating a wallet.
	if action == actionCreateWallet {
		delete(data, "Id")
	}
	delete(data, "CreationDate")
	delete(data, "Balance")

	if action == actionEditWallet {
		// Delete empty values so that existing ones don't get
		// overwritten with empty values.
		for k, v := range data {
			switch v.(type) {
			case string:
				if v.(string) == "" {
					delete(data, k)
				}
			case int:
				if v.(int) == 0 {
					delete(data, k)
				}
			}
		}
	}

	wallet, err := w.service.walletRequest(action, data)
	if err != nil {
		return err
	}
	*w = *wallet
	return nil
}

func (w *Wallet) Transactions() TransferList {
	return nil
}

// Wallet finds a legal user using the user_id attribute.
func (m *MangoPay) Wallet(id string) (*Wallet, error) {
	w, err := m.walletRequest(actionFetchWallet, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (m *MangoPay) walletRequest(action mangoAction, data JsonObject) (*Wallet, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}
	u := new(Wallet)
	if err := m.unMarshalJSONResponse(resp, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (m *MangoPay) walletListRequest(action mangoAction, data JsonObject) (WalletList, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}
	ws := WalletList{}
	if err := m.unMarshalJSONResponse(resp, &ws); err != nil {
		return nil, err
	}
	return ws, nil
}

func (m *MangoPay) wallets(u Consumer) (WalletList, error) {
	id := ""
	switch u.(type) {
	case *LegalUser:
		id = u.(*LegalUser).Id
	case *NaturalUser:
		id = u.(*NaturalUser).Id
	}
	if id == "" {
		return nil, errors.New("user has empty Id")
	}
	trs, err := m.walletListRequest(actionFetchUserWallets, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return trs, nil
}

// Wallet finds all user's wallets. Provided for convenience.
func (m *MangoPay) Wallets(user Consumer) (WalletList, error) {
	ws, err := m.wallets(user)
	return ws, err
}
