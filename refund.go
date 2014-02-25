// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"fmt"

//	"errors"
)

// List of refunds.
type RefundList []*Refund

type refundKind int

const (
	// Pay a wallet back
	transferRefund refundKind = iota
	// Pay a user (bank account) back
	payInRefund
)

// A refund is a request to pay a wallet back.
//
// http://docs.mangopay.com/api-references/refund/%E2%80%A2-refund-a-transfer/
type Refund struct {
	ProcessReply
	AuthorId               string
	DebitedFunds           Money
	Fees                   Money
	CreditedFunds          Money
	Type                   string // PAY_IN, PAY_OUT or TRANSFER
	Nature                 string
	CreditedUserId         string
	InitialTransactionId   string
	InitialTransactionType string
	DebitedWalletId        string
	CreditedWalletId       string

	transfer *Transfer
	payIn    *PayIn
	kind     refundKind
}

func (r *Refund) String() string {
	return struct2string(r)
}

// Save creates a refund.
func (r *Refund) save() error {
	data := JsonObject{}
	j, err := json.Marshal(r)
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
	for _, field := range []string{"Id", "CreationDate", "ExecutionDate", "CreditedFunds", "CreditedUserId", "ResultCode", "ResultMessage", "Status", "Fees", "InitialTransactionType", "InitialTransactionId", "DebitedFunds", "Nature", "DebitedWalletId", "CreditedWalletId", "Type"} {
		delete(data, field)
	}

	data["AuthorId"] = r.transfer.AuthorId
	var action mangoAction
	r.String()
	switch r.kind {
	case transferRefund:
		action = actionCreateTransferRefund
		data["TransferId"] = r.transfer.Id
	case payInRefund:
		action = actionCreatePayInRefund
		data["PayInId"] = r.payIn.Id
	}
	ins, err := r.transfer.service.anyRequest(new(Refund), action, data)
	if err != nil {
		return err
	}
	t, p, k := r.transfer, r.payIn, r.kind
	*r = *(ins.(*Refund))
	r.transfer, r.payIn, r.kind = t, p, k
	return nil
}
