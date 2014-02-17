// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"errors"
	"fmt"
)

// CardRegistration is used to register a credit card.
//
// http://docs.mangopay.com/api-references/card-registration/
type CardRegistration struct {
	ProcessReply
	// Id of the author.
	UserId string
	// Currency of the registered card.
	Currency string
	// Key sent with the card details and the PreregistrationData.
	AccessKey string
	// This passphrase is sent with the card details and the AccessKey.
	PreregistrationData string
	// The actual URL to POST the card details, the access key and the
	// PreregistrationData.
	CardRegistrationUrl string
	// Part of the reply, once the card details, the AccessKey and the
	// PreregistrationData has been sent.
	CardRegistrationData string
	CardType             string
	// CardId if part of the reply, once the CardRegistration has been
	// edited with the CardRegistrationData.
	CardId  string
	service *MangoPay
}

func (c *CardRegistration) String() string {
	return fmt.Sprintf(`
Id                      : %s
Tag                     : %s
CreationDate            : %s
ResultCode              : %s
ResultMessage           : %s
Status                  : %s
UserId                  : %s
Currency                : %s
AccessKey               : %s
PreregistrationData     : %s
CardRegistrationUrl     : %s
CardRegistrationData    : %s
CardType                : %s
CardId                  : %s
`, c.Id, c.Tag, unixTimeToString(c.CreationDate), c.ResultCode, c.ResultMessage, c.Status, c.UserId, c.Currency, c.AccessKey, c.PreregistrationData, c.CardRegistrationUrl, c.CardRegistrationData, c.CardType, c.CardId)
}

// NewCardRegistration creates a new credit card registration process.
func (m *MangoPay) NewCardRegistration(user Consumer, currency string) (*CardRegistration, error) {
	id := ""
	switch user.(type) {
	case *LegalUser:
		id = user.(*LegalUser).Id
	case *NaturalUser:
		id = user.(*NaturalUser).Id
	}
	if id == "" {
		return nil, errors.New("empty user ID. Unable to create card registration.")
	}
	cr := &CardRegistration{
		UserId:       id,
		Currency:     currency,
		ProcessReply: ProcessReply{},
	}
	cr.service = m
	return cr, nil
}

// Register initiates the process of getting pre-registration data and access
// key to allow a user to post his credit cart info the CardRegistrationUrl.
func (c *CardRegistration) Register() error {
	data := JsonObject{}
	j, err := json.Marshal(c)
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

	// Fields not allowed when initiating a card registration.
	for _, field := range []string{"Id", "CreationDate", "ExecutionDate",
		"ResultCode", "ResultMessage", "Status", "AccessKey", "CardId",
		"CardRegistrationData", "CardRegistrationUrl", "CardType",
		"PreregistrationData", "Tag"} {
		delete(data, field)
	}

	cr, err := c.service.cardRegistartionRequest(actionCreateCardRegistration, data)
	if err != nil {
		return err
	}
	*c = *cr
	return nil
}

// SendRegistrationData sends a request with rdata to get the CardID
// ready to be used for a direct payIn.
func (c *CardRegistration) SendRegistrationData(rdata string) error {
	if c.Id == "" {
		return errors.New("empty card registration Id. Can't send registration data.")
	}
	if rdata == "" {
		return errors.New("empty registration data supplied.")
	}

	cr, err := c.service.cardRegistartionRequest(actionSendCardRegistrationData, JsonObject{"Id": c.Id, "CardRegistrationData": rdata})
	if err != nil {
		return err
	}
	*c = *cr
	return nil
}

func (m *MangoPay) cardRegistartionRequest(action mangoAction, data JsonObject) (*CardRegistration, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}
	u := new(CardRegistration)
	if err := m.unMarshalJSONResponse(resp, u); err != nil {
		return nil, err
	}
	return u, nil
}
