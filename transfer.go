// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ErrTransactionFailed struct {
	msg string
}

func (e *ErrTransactionFailed) Error() string {
	return "transaction failed: " + e.msg
}

// List of transactions.
type TransferList []*Transfer

// Transfer hold details about relocating e-money from a wallet
// to another one.
//
// See http://docs.mangopay.com/api-references/transfers/.
type Transfer struct {
	Id               string
	Tag              string
	AuthorId         string
	CreditedUserId   string
	DebitedFunds     Money
	Fees             Money
	DebitedWalletId  string
	CreditedWalletId string
	CreationDate     int
	CreditedFunds    Money
	Status           string
	ResultCode       string
	ResultMessage    string
	ExecutionDate    int
	service          *MangoPay
}

func (t *Transfer) String() string {
	return fmt.Sprintf(`
Id               : %s
Tag              : %s
CreditedUserId   : %s
DebitedFunds     : %s
Fees             : %s
DebitedWalletId  : %s
CreditedWalletId : %s
Creation date    : %d
CreditedFunds    : %s
Status           : %s
ResultCode       : %s
ResultMessage    : %s
ExecutionDate    : %d
`, t.Id, t.Tag, t.CreditedUserId, t.DebitedFunds.String(), t.Fees.String(), t.DebitedWalletId, t.CreditedWalletId, t.CreationDate, t.CreditedFunds.String(), t.Status, t.ResultCode, t.ResultMessage, t.ExecutionDate)
}

// NewTransfer creates a new tranfer (or transaction).
func (m *MangoPay) NewTransfer(author Buyer, amount Money, fees Money, from, to *Wallet) (*Transfer, error) {
	msg := "new tranfer: "
	if author == nil {
		return nil, errors.New(msg + "nil author")
	}
	if from == nil {
		return nil, errors.New(msg + "nil source wallet")
	}
	if to == nil {
		return nil, errors.New(msg + "nil dest wallet")
	}
	if from.Id == "" {
		return nil, errors.New(msg + "source wallet has empty Id")
	}
	if to.Id == "" {
		return nil, errors.New(msg + "dest wallet has empty Id")
	}
	id := ""
	switch author.(type) {
	case *LegalUser:
		id = author.(*LegalUser).Id
	case *NaturalUser:
		id = author.(*NaturalUser).Id
	}
	if id == "" {
		return nil, errors.New(msg + "author has empty Id")
	}
	t := &Transfer{
		AuthorId:         id,
		DebitedFunds:     amount,
		Fees:             fees,
		DebitedWalletId:  from.Id,
		CreditedWalletId: to.Id,
	}
	t.service = m
	return t, nil
}

// Save creates a transfer.
func (t *Transfer) Save() error {
	data := JsonObject{}
	j, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	// Force float64 to int conversion after unmarshalling.
	for _, field := range []string{"CreationDate", "ExecutionDate"} {
		data[field] = int(data[field].(float64))
	}

	// Fields not allowed when creating a tranfer.
	for _, field := range []string{"Id", "CreationDate", "ExecutionDate", "CreditedFunds", "CreditedUserId", "ResultCode", "ResultMessage", "Status"} {
		delete(data, field)
	}

	tr, err := t.service.transferRequest(actionCreateTransfer, data)
	if err != nil {
		return err
	}
	*t = *tr

	if t.Status == "FAILED" {
		return &ErrTransactionFailed{t.ResultMessage}
	}
	return nil
}

// Transfer finds a legal user using the user_id attribute.
func (m *MangoPay) Transfer(id string) (*Transfer, error) {
	w, err := m.transferRequest(actionFetchTransfer, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (m *MangoPay) transferRequest(action mangoAction, data JsonObject) (*Transfer, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}
	u := new(Transfer)
	if err := m.unMarshalJSONResponse(resp, u); err != nil {
		return nil, err
	}
	return u, nil
}
