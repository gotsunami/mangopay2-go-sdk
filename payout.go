// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Custom error returned in case of failed payOut.
type ErrPayOutFailed struct {
	payinId string
	msg     string
}

func (e *ErrPayOutFailed) Error() string {
	return fmt.Sprintf("payOut %s failed: %s ", e.payinId, e.msg)
}

// A PayOut Bank wire is a request to withdraw money from a wallet to a
// registered bank account.
//
// See http://docs.mangopay.com/api-references/pay-out-bank-wire/
type PayOut struct {
	ProcessReply
	AuthorId          string
	CreditedUserId    string
	DebitedFunds      Money
	Fees              Money
	Type              string // PAY_IN, PAY_OUT or TRANSFER
	Nature            string // REGULAR, REFUND or REPUDIATION
	PaymentType       string
	DebitedWalletId   string
	BankAccountId     string
	CreditedFunds     Money
	MeanOfPaymentType string
	BankWireRef       string
	service           *MangoPay
}

func (p *PayOut) String() string {
	return struct2string(p)
}

// NewPayOut creates a new bank wire.
func (m *MangoPay) NewPayOut(author Consumer, amount Money, fees Money, from *Wallet, to *BankAccount) (*PayOut, error) {
	msg := "new payOut: "
	if author == nil {
		return nil, errors.New(msg + "nil author")
	}
	if from == nil {
		return nil, errors.New(msg + "nil wallet")
	}
	if to == nil {
		return nil, errors.New(msg + "nil bank account")
	}
	id := consumerId(author)
	if id == "" {
		return nil, errors.New(msg + "author has empty Id")
	}
	p := &PayOut{
		AuthorId:        id,
		DebitedWalletId: from.Id,
		Fees:            fees,
		DebitedFunds:    amount,
		BankAccountId:   to.Id,
		service:         m,
	}
	return p, nil
}

// Save sends an HTTP query to create a bank wire. It may return an
// ErrPayOutFailed error if the payment has failed.
func (p *PayOut) Save() error {
	data := JsonObject{}
	j, err := json.Marshal(p)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	// Fields not allowed when creating a payOut.
	for _, field := range []string{"Id", "CreationDate", "ExecutionDate", "ResultCode", "ResultMessage", "Status", "CreditedUserId", "Type", "Nature", "PaymentType", "CreditedFunds", "MeanOfPaymentType"} {
		delete(data, field)
	}

	pay, err := p.service.anyRequest(new(PayOut), actionCreatePayOut, data)
	if err != nil {
		return err
	}
	serv := p.service
	*p = *(pay.(*PayOut))
	p.service = serv

	if p.Status == "FAILED" {
		return &ErrPayOutFailed{p.Id, p.ResultMessage}
	}
	return nil
}

// PayOut finds a bank wire.
func (m *MangoPay) PayOut(id string) (*PayOut, error) {
	p, err := m.anyRequest(new(PayOut), actionFetchPayOut, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return p.(*PayOut), nil
}
