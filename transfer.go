// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"fmt"
)

// List of transactions.
type TransferList []*Transfer

// Transfer hold details about relocating e-money from a wallet
// to another one.
//
// See http://docs.mangopay.com/api-references/transfers/.
type Transfer struct {
	Id                 string
	Tag                string
	AuthorId           string
	CreditedUserId     string
	DebitedFunds       Money
	Fees               Money
	DebitedTransferID  string
	CreditedTransferID string
	CreationDate       int
	CreditedFunds      Money
	Status             string
	ResultCode         string
	ResultMessage      string
	ExecutionDate      int
	service            *MangoPay
}

func (t *Transfer) String() string {
	return fmt.Sprintf(`
Id               : %s
Tag              : %s
CreditedUserId   : %s
DebitedFunds     : %s
Fees             : %s
DebitedTransferID  : %s
CreditedTransferID : %s
Creation date    : %d
CreditedFunds    : %s
Status           : %s
ResultCode       : %s
ResultMessage    : %s
ExecutionDate    time
`, t.Id, t.Tag, t.CreditedUserId, t.DebitedFunds, t.Fees, t.DebitedTransferID, t.CreditedTransferID, t.CreationDate, t.CreditedFunds, t.Status, t.ResultCode, t.ResultMessage, t.ExecutionDate)
}

// NewTransfer creates a new tranfer (or transaction).
func (m *MangoPay) NewTransfer() *Transfer {
	t := new(Transfer)
	t.service = m
	return t
}

// Save creates a transfer.
func (t *Transfer) Save() error {
	var action mangoAction
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
	delete(data, "Id")
	delete(data, "CreationDate")
	delete(data, "ExecutionDate")

	tr, err := t.service.transferRequest(action, data)
	if err != nil {
		return err
	}
	*t = *tr
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
