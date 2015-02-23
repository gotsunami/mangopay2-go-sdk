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
	service          *MangoPay
}

func (p *PayIn) String() string {
	return struct2string(p)
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
	return struct2string(p)
}

// WebPayIn hold details about making a payment through a web interface.
//
// See http://docs.mangopay.com/api-references/payins/payins-card-web/
type WebPayIn struct {
	PayIn
	ReturnUrl          string
	TemplateUrlOptions struct {
		Payline string
	} `json:",omitempty"`
	Culture     string
	CardType    string
	RedirectUrl string
	service     *MangoPay
}

func (p *WebPayIn) String() string {
	return struct2string(p)
}

// NewWebPayIn creates a new payment.
func (m *MangoPay) NewWebPayIn(author Consumer, amount Money, fees Money, credit *Wallet, returnUrl string, culture string, templateUrl string) (*WebPayIn, error) {
	msg := "new web payIn: "
	if author == nil {
		return nil, errors.New(msg + "nil author")
	}
	if credit == nil {
		return nil, errors.New(msg + "nil dest wallet")
	}
	id := consumerId(author)
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
			service:          m,
		},
		ReturnUrl: u.String(),
		CardType:  "CB_VISA_MASTERCARD",
		Culture:   culture,
		service:   m,
	}

	// Optional parameter
	if len(templateUrl) > 0 {
		t, err := url.Parse(templateUrl)
		if err != nil {
			return nil, errors.New(msg + err.Error())
		}
		p.TemplateUrlOptions.Payline = t.String()
	}

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

	tr, err := t.service.anyRequest(new(WebPayIn), actionCreateWebPayIn, data)
	if err != nil {
		return err
	}
	serv := t.service
	*t = *(tr.(*WebPayIn))
	t.service = serv
	t.PayIn.service = serv

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
		id := consumerId(con)
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
			service:          m,
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

	tr, err := p.service.anyRequest(new(DirectPayIn), actionCreateDirectPayIn, data)
	if err != nil {
		return err
	}
	serv := p.service
	*p = *(tr.(*DirectPayIn))
	p.service = serv
	p.PayIn.service = serv

	if p.Status == "FAILED" {
		return &ErrPayInFailed{p.Id, p.ResultMessage}
	}
	return nil
}

// Refund allows to refund a pay-in. Call the Refund's Save() method
// to make a request to reimburse a user on his payment card.
func (p *PayIn) Refund() (*Refund, error) {
	r := &Refund{
		ProcessReply: ProcessReply{},
		payIn:        p,
		kind:         payInRefund,
	}
	if err := r.save(); err != nil {
		return nil, err
	}
	return r, nil
}

// PayIn finds a payment.
func (m *MangoPay) PayIn(id string) (*WebPayIn, error) {
	p, err := m.anyRequest(new(WebPayIn), actionFetchPayIn, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return p.(*WebPayIn), nil
}
