// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// Custom error returned in case of failed payIn.
type ErrPayInFailed struct {
	payinId string
	msg     string
}

func (e *ErrPayInFailed) Error() string {
	return fmt.Sprintf("payIn %s failed: %s ", e.payinId, e.msg)
}

// PayIn holds common fields to all MangoPay's supported payment means
// (through web, direct, preauthorized, bank wire).
type PayIn struct {
	ProcessReply
	AuthorId         string
	CreditedUserId   string
	DebitedFunds     Money
	Fees             Money
	CreditedWalletId string
	SecureMode       string
	CreditedFunds    Money
	Type             string // PAY_IN, PAY_OUT or TRANSFER
	Nature           string // REGULAR, REFUND or REPUDIATION
	PaymentType      string
	ExecutionType    string // WEB or DIRECT (with tokenized card)
}

// DirectPayIn is used to process a payment with registered (tokenized) cards.
type DirectPayIn struct {
	PayIn
	SecureModeReturnUrl string
	CardId              string
	DebitedWalletId     string
	service             *MangoPay
}

func (p *DirectPayIn) String() string {
	return fmt.Sprintf(`
Id                  : %s
Tag                 : %s
AuthorId            : %s
DebitedFunds        : %s
Fees                : %s
CreditedWalletId    : %s
SecureMode          : %s
CreationDate        : %s
CreditedFunds       : %s
CreditedUserId      : %s
Status              : %s
ResultCode          : %s
ResultMessage       : %s
ExecutionDate       : %s
Type                : %s
Nature              : %s
PaymentType         : %s
ExecutionType       : %s
SecureModeReturnUrl : %s
CardId              : %s
DebitedWalletId     : %s
`, p.Id, p.Tag, p.AuthorId, p.DebitedFunds.String(), p.Fees.String(), p.CreditedWalletId, p.SecureMode, unixTimeToString(p.CreationDate), p.CreditedFunds.String(), p.CreditedUserId, p.Status, p.ResultCode, p.ResultMessage, unixTimeToString(p.ExecutionDate), p.Type, p.Nature, p.PaymentType, p.ExecutionType, p.SecureModeReturnUrl, p.CardId, p.DebitedWalletId)
}

// WebPayIn hold details about making a payment through a web interface.
//
// See http://docs.mangopay.com/api-references/payins/payins-card-web/
type WebPayIn struct {
	PayIn
	ReturnUrl   string
	Culture     string
	CardType    string
	RedirectUrl string
	service     *MangoPay
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
CreationDate     : %s
CreditedFunds    : %s
CreditedUserId   : %s
Status           : %s
ResultCode       : %s
ResultMessage    : %s
ExecutionDate    : %s
Type             : %s 
Nature           : %s
PaymentType      : %s
ExecutionType    : %s
RedirectUrl      : %s
`, p.Id, p.Tag, p.AuthorId, p.DebitedFunds.String(), p.Fees.String(), p.CreditedWalletId, p.ReturnUrl, p.Culture, p.CardType, p.SecureMode, unixTimeToString(p.CreationDate), p.CreditedFunds.String(), p.CreditedUserId, p.Status, p.ResultCode, p.ResultMessage, unixTimeToString(p.ExecutionDate), p.Type, p.Nature, p.PaymentType, p.ExecutionType, p.RedirectUrl)
}

// NewWebPayIn creates a new payment.
func (m *MangoPay) NewWebPayIn(author Consumer, amount Money, fees Money, credit *Wallet, returnUrl string, culture string) (*WebPayIn, error) {
	msg := "new web payIn: "
	if author == nil {
		return nil, errors.New(msg + "nil author")
	}
	if credit == nil {
		return nil, errors.New(msg + "nil dest wallet")
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
	u, err := url.Parse(returnUrl)
	if err != nil {
		return nil, errors.New(msg + err.Error())
	}
	p := &WebPayIn{
		PayIn: PayIn{
			AuthorId:         id,
			DebitedFunds:     amount,
			Fees:             fees,
			CreditedWalletId: credit.Id,
		},
		ReturnUrl: u.String(),
		CardType:  "CB_VISA_MASTERCARD",
		Culture:   culture,
	}
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
	for _, field := range []string{"Id", "CreationDate", "ExecutionDate", "CreditedFunds", "CreditedUserId", "ResultCode", "ResultMessage", "Status", "ExecutionType", "PaymentType", "SecureMode", "Type", "Nature"} {
		delete(data, field)
	}

	tr, err := t.service.webPayinRequest(actionCreateWebPayIn, data)
	if err != nil {
		return err
	}
	*t = *tr

	if t.Status == "FAILED" {
		return &ErrPayInFailed{t.Id, t.ResultMessage}
	}
	return nil
}

// NewDirectPayIn creates a direct payment from a tokenized credit card.
//
//  - from     : AuthorId value
//  - to       : CreditedUserId value
//  - src      : CardId value
//  - dst      : CreditedWalletId value
//  - amount   : DebitedFunds value
//  - fees     : Fees value
//  - returnUrl: SecureModeReturnUrl value
//
// See http://docs.mangopay.com/api-references/payins/payindirectcard/
func (m *MangoPay) NewDirectPayIn(from, to Consumer, src *Card, dst *Wallet, amount, fees Money, returnUrl string) (*DirectPayIn, error) {
	msg := "new direct payIn: "
	ps := []struct {
		i   interface{}
		msg string
	}{
		{from, "from parameter"},
		{to, "to parameter"},
		{src, "card"},
		{dst, "wallet"},
	}
	for _, p := range ps {
		if p.i == nil {
			return nil, errors.New(msg + p.msg)
		}
	}
	if returnUrl == "" {
		return nil, errors.New(msg + "empty return url")
	}

	cons := make([]string, 2)
	for k, con := range []Consumer{from, to} {
		id := ""
		switch con.(type) {
		case *LegalUser:
			id = con.(*LegalUser).Id
		case *NaturalUser:
			id = con.(*NaturalUser).Id
		}
		cons[k] = id
	}

	// Check Ids
	for _, i := range []struct{ v, msg string }{
		{cons[0], "from consumer"},
		{cons[1], "to consumer"},
		{dst.Id, "wallet"},
		{src.Id, "card"},
	} {
		if i.v == "" {
			return nil, errors.New(fmt.Sprintf("empty %s id", i.msg))
		}
	}

	u, err := url.Parse(returnUrl)
	if err != nil {
		return nil, errors.New(msg + err.Error())
	}
	p := &DirectPayIn{
		PayIn: PayIn{
			AuthorId:         cons[0],
			CreditedUserId:   cons[1],
			DebitedFunds:     amount,
			Fees:             fees,
			CreditedWalletId: dst.Id,
		},
		SecureModeReturnUrl: u.String(),
		CardId:              src.Id,
	}
	p.service = m
	return p, nil
}

// Save sends an HTTP query to create a direct payIn. Upon successful creation,
// it may return an ErrPayInFailed error if the payment has failed.
func (p *DirectPayIn) Save() error {
	data := JsonObject{}
	j, err := json.Marshal(p)
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
	for _, field := range []string{"Id", "CreationDate", "ExecutionDate", "CreditedFunds", "ResultCode", "ResultMessage", "Status", "ExecutionType", "PaymentType", "SecureMode", "DebitedWalletId", "Type", "Nature"} {
		delete(data, field)
	}

	tr, err := p.service.directPayinRequest(actionCreateDirectPayIn, data)
	if err != nil {
		return err
	}
	serv := p.service
	*p = *tr
	p.service = serv

	if p.Status == "FAILED" {
		return &ErrPayInFailed{p.Id, p.ResultMessage}
	}
	return nil
}

// PayIn finds a payment.
func (m *MangoPay) PayIn(id string) (*WebPayIn, error) {
	p, err := m.webPayinRequest(actionFetchPayIn, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (m *MangoPay) webPayinRequest(action mangoAction, data JsonObject) (*WebPayIn, error) {
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

func (m *MangoPay) directPayinRequest(action mangoAction, data JsonObject) (*DirectPayIn, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}
	u := new(DirectPayIn)
	if err := m.unMarshalJSONResponse(resp, u); err != nil {
		return nil, err
	}
	return u, nil
}
