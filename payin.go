// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	_ "errors"
	"fmt"
)

// Custom error returned in case of failed payIn.
type ErrPayInFailed struct {
	payinId string
	msg     string
}

func (e *ErrPayInFailed) Error() string {
	return fmt.Sprintf("payIn %s failed: %s ", e.payinId, e.msg)
}

// WebPayIn hold details about making a payment through a web interface.
//
// See http://docs.mangopay.com/api-references/payins/payins-card-web/
type WebPayIn struct {
	Id               string
	Tag              string
	AuthorId         string
	DebitedFunds     Money
	Fees             Money
	CreditedWalletId string
	ReturnUrl        string
	Culture          string
	CardType         string
	SecureMode       string
	CreationDate     int
	CreditedFunds    Money
	CreditedUserId   string
	Status           string
	ResultCode       string
	ResultMessage    string
	ExecutionDate    int
	Type             string // PAY_IN, PAY_OUT or TRANSFER
	Nature           string // REGULAR, REFUND or REPUDIATION
	PaymentType      string
	ExecutionType    string // WEB or DIRECT (with tokenized card)
	RedirectUrl      string
	service          *MangoPay
}

func (p *WebPayIn) String() string {
	return fmt.Sprintf(`
Id               : %s
Tag              : %s
AuthorId         : %s
DebitedFunds     : %s
Fees             : %s
CreditedWalletId : %s
ReturnUrl        : %s
Culture          : %s
CardType         : %s
SecureMode       : %s
CreationDate     : %d
CreditedFunds    : %s
CreditedUserId   : %s
Status           : %s
ResultCode       : %s
ResultMessage    : %s
ExecutionDate    : %d
Type             : %s 
Nature           : %s
PaymentType      : %s
ExecutionType    : %s
RedirectUrl      : %s
`, p.Id, p.Tag, p.AuthorId, p.DebitedFunds, p.Fees, p.CreditedWalletId, p.ReturnUrl, p.Culture, p.CardType, p.SecureMode, p.CreationDate, p.CreditedFunds, p.CreditedUserId, p.Status, p.ResultCode, p.ResultMessage, p.ExecutionDate, p.Type, p.Nature, p.PaymentType, p.ExecutionType, p.RedirectUrl)
}

// NewWebPayIn creates a new payment.
func (m *MangoPay) NewWebPayIn() (*WebPayIn, error) {
	p := &WebPayIn{}
	p.service = m
	return p, nil
}

// Save sends an HTTP query to create a payIn. Upon successful creation,
// it may return an ErrPayInFailed error if the payment has failed.
func (t *WebPayIn) Save() error {
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
	/* FIXME
	for _, field := range []string{"Id", "CreationDate", "ExecutionDate", "CreditedFunds", "CreditedUserId", "ResultCode", "ResultMessage", "Status"} {
		delete(data, field)
	}

	tr, err := t.service.transferRequest(actionCreateWebPayIn, data)
	if err != nil {
		return err
	}
	*t = *tr

	if t.Status == "FAILED" {
		return &ErrWebPayInFailed{t.Id, t.ResultMessage}
	}
	*/
	return nil
}

// PayIn finds a payment.
func (m *MangoPay) PayIn(id string) (*WebPayIn, error) {
	w, err := m.payinRequest(actionFetchPayIn, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (m *MangoPay) payinRequest(action mangoAction, data JsonObject) (*WebPayIn, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}
	u := new(WebPayIn)
	if err := m.unMarshalJSONResponse(resp, u); err != nil {
		return nil, err
	}
	return u, nil
}
