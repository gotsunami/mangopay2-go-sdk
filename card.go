// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	CardId string

	service *MangoPay
	// true after Init() is successful.
	isInitialized bool
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
		UserId:        id,
		Currency:      currency,
		ProcessReply:  ProcessReply{},
		isInitialized: false,
	}
	cr.service = m
	return cr, nil
}

// Init initiates the process of getting pre-registration data and access
// key from MangoPay to allow a user to post his credit card info to the
// c.CardRegistrationUrl (which is an external banking service).
//
// User's card details must be sent directly through an HTML form to the
// c.CardRegistrationUrl.
//
// The HTML form must have the following input fields:
//  - "data" (hidden) equals to c.PreregistrationData
//  - "accessKeyRef" (hidden) equals to c.AccessKey
//  - "cardNumber" equals to the user's credit card number
//  - "cardExpirationDate" equals to the user's card expiration date (format: MMYY)
//  - "cardCvx" equals to user's 3-digits cvx code
//  - "returnURL" so we can retrieve the final registration data token
//
// A successful call to Init() will fill in the PreregistrationData and
// AccessKey fields of the current CardRegistration object automatically.
func (c *CardRegistration) Init() error {
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
	// Backup private service
	service := c.service
	*c = *cr
	c.service = service

	// Okay for step 2.
	c.isInitialized = true
	return nil
}

// Register effectively registers the credit card against the MangoPay service. The
// registrationData value is returned by the external banking service that deals with
// the credit card information, and is obtained by submitting an HTML form to
// the external banking service.
func (c *CardRegistration) Register(registrationData string) error {
	if !strings.HasPrefix(registrationData, "data=") {
		return errors.New("invalid registration data. Must start with data=")
	}
	if !c.isInitialized {
		return errors.New("card registration process not initialized. Did you call Init() first?")
	}
	cr, err := c.service.cardRegistartionRequest(actionSendCardRegistrationData,
		JsonObject{"Id": c.Id, "RegistrationData": registrationData})
	if err != nil {
		return err
	}
	// Backup private members
	serv := c.service
	isr := c.isInitialized
	*c = *cr
	c.CardRegistrationData = registrationData
	c.service = serv
	c.isInitialized = isr
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
