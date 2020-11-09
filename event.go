// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

func NewEventFromRequest(req *http.Request) (*Event, error) {
	q := req.URL.Query()
	eventType := q.Get("EventType")
	resourceId := q.Get("RessourceId") // ReSSource - Mistake in MangoPay API
	if resourceId == "" {
		// Fallback
		resourceId = q.Get("ResourceId")
	}
	timestampStr := q.Get("Date")
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil, errors.New("Unable parse value of 'Date' query option:" + err.Error())
	}
	return &Event{resourceId, EventType(eventType), time.Unix(timestamp, 0)}, nil
}

// See http://docs.mangopay.com/api-references/events/.

// An event resource.
type Event struct {
	ResourceId string    `json:"RessourceId"` // ReSSource - Mistake in MangoPay API
	Type       EventType `json:"EventType"`
	Date       time.Time
}

type EventList []Event

type EventType string

const (
	EventPayinNormalCreated             EventType = "PAYIN_NORMAL_CREATED"
	EventPayinNormalSucceeded                     = "PAYIN_NORMAL_SUCCEEDED"
	EventPayinNormalFailed                        = "PAYIN_NORMAL_FAILED"
	EventPayoutNormalCreated                      = "PAYOUT_NORMAL_CREATED"
	EventPayoutNormalSucceeded                    = "PAYOUT_NORMAL_SUCCEEDED"
	EventPayoutNormalFailed                       = "PAYOUT_NORMAL_FAILED"
	EventTransferNormalCreated                    = "TRANSFER_NORMAL_CREATED"
	EventTransferNormalSucceeded                  = "TRANSFER_NORMAL_SUCCEEDED"
	EventTransferNormalFailed                     = "TRANSFER_NORMAL_FAILED"
	EventPayinRefundCreated                       = "PAYIN_REFUND_CREATED"
	EventPayinRefundSucceeded                     = "PAYIN_REFUND_SUCCEEDED"
	EventPayinRefundFailed                        = "PAYIN_REFUND_FAILED"
	EventPayoutRefundCreated                      = "PAYOUT_REFUND_CREATED"
	EventPayoutRefundSucceeded                    = "PAYOUT_REFUND_SUCCEEDED"
	EventPayoutRefundFailed                       = "PAYOUT_REFUND_FAILED"
	EventTransferRefundCreated                    = "TRANSFER_REFUND_CREATED"
	EventTransferRefundSucceeded                  = "TRANSFER_REFUND_SUCCEEDED"
	EventTransferRefundFailed                     = "TRANSFER_REFUND_FAILED"
	EventPayinRepudiationCreated                  = "PAYIN_REPUDIATION_CREATED"
	EventPayinRepudiationSucceeded                = "PAYIN_REPUDIATION_SUCCEEDED"
	EventPayinRepudiationFailed                   = "PAYIN_REPUDIATION_FAILED"
	EventKycCreated                               = "KYC_CREATED"
	EventKycSucceeded                             = "KYC_SUCCEEDED"
	EventKycFailed                                = "KYC_FAILED"
	EventKycValidationAsked                       = "KYC_VALIDATION_ASKED"
	EventKycOutdated                              = "KYC_OUTDATED"
	EventDisputeDocumentCreated                   = "DISPUTE_DOCUMENT_CREATED"
	EventDisputeDocumentValidationAsked           = "DISPUTE_DOCUMENT_VALIDATION_ASKED"
	EventDisputeDocumentSucceeded                 = "DISPUTE_DOCUMENT_SUCCEEDED"
	EventDisputeDocumentFailed                    = "DISPUTE_DOCUMENT_FAILED"
	EventDisputeCreated                           = "DISPUTE_CREATED"
	EventDisputeSubmitted                         = "DISPUTE_SUBMITTED"
	EventDisputeActionRequired                    = "DISPUTE_ACTION_REQUIRED"
	EventDisputeFurtherActionRequired             = "DISPUTE_FURTHER_ACTION_REQUIRED"
	EventDisputeClosed                            = "DISPUTE_CLOSED"
	EventDisputeSentToBank                        = "DISPUTE_SENT_TO_BANK"
	EventTransferSettlementCreated                = "TRANSFER_SETTLEMENT_CREATED"
	EventTransferSettlementSucceeded              = "TRANSFER_SETTLEMENT_SUCCEEDED"
	EventTransferSettlementFailed                 = "TRANSFER_SETTLEMENT_FAILED"
	EventMandateCreated                           = "MANDATE_CREATED"
	EventMandatedFailed                           = "MANDATED_FAILED"
	EventMandateActivated                         = "MANDATE_ACTIVATED"
	EventMandateSubmitted                         = "MANDATE_SUBMITTED"
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
