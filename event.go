// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"time"
)

// See http://docs.mangopay.com/api-references/events/.

// An event ressource.
type Event struct {
	RessourceId string
	Type        EventType `json:"EventType"`
	Date        *time.Time
}

type EventList []*Event

type EventType int

const (
	PAYIN_NORMAL_CREATED EventType = iota
	PAYIN_NORMAL_SUCCEEDED
	PAYIN_NORMAL_FAILED
	PAYOUT_NORMAL_CREATED
	PAYOUT_NORMAL_SUCCEEDED
	PAYOUT_NORMAL_FAILED
	TRANSFER_NORMAL_CREATED
	TRANSFER_NORMAL_SUCCEEDED
	TRANSFER_NORMAL_FAILED
	PAYIN_REFUND_CREATED
	PAYIN_REFUND_SUCCEEDED
	PAYIN_REFUND_FAILED
	PAYOUT_REFUND_CREATED
	PAYOUT_REFUND_SUCCEEDED
	PAYOUT_REFUND_FAILED
	TRANSFER_REFUND_CREATED
	TRANSFER_REFUND_SUCCEEDED
	TRANSFER_REFUND_FAILED
)

// Events returns a list of all financial events. This include PayIns, PayOuts and
// transfers.
//
// TODO: add support for pagination and date range.
func (m *MangoPay) Events() (EventList, error) {
	es := EventList{}
	resp, err := m.request(actionEvents, nil)
	if err != nil {
		return nil, err
	}
	if err := m.unMarshalJSONResponse(resp, &es); err != nil {
		return nil, err
	}
	return es, err
}
