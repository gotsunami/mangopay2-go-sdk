// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"errors"
	"strings"
)

// List of cards.
type CardList []*Card

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
	return struct2string(c)
}

// Card holds all credit card details.
type Card struct {
	ProcessIdent
	ExpirationDate string // MMYY
	Alias          string // Obfuscated card number, i.e 497010XXXXXX4414
	CardProvider   string // CB, VISA, MASTERCARD etc.
	CardType       string // CB_VISA_MASTERCARD
	Product        string
	BankCode       string
	Active         bool
	Currency       string // Currency accepted in the waller, i.e EUR, USD etc.
	Validity       string // UNKNOWN, VALID, INVALID
}

func (c *Card) String() string {
	return struct2string(c)
}

// Card fetches a registered credit card.
func (m *MangoPay) Card(id string) (*Card, error) {
	any, err := m.anyRequest(new(Card), actionFetchCard, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return any.(*Card), nil
}

// Card finds all user's cards.
func (m *MangoPay) Cards(user Consumer) (CardList, error) {
	id := consumerId(user)
	if id == "" {
		return nil, errors.New("user has empty Id")
	}
	resp, err := m.request(actionFetchUserCards, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	cl := CardList{}
	if err := m.unMarshalJSONResponse(resp, &cl); err != nil {
		return nil, err
	}
	return cl, nil
}

// NewCardRegistration creates a new credit card registration object that can
// be used to register a new credit card for a given user.
//
// Registering a new credit card involves the following workflow:
//
//  1. Create a new CardRegistration object
//  2. Call .Init() to pre-register the card against MangoPay services and get
//     access tokens required to register the credit card againts an external
//     banking service
//  3. Insert those tokens in an HTML form submitted by the user directly to
//     the external banking service
//  4. Get the final token from the external banking service
//  5. Call .Register() with this token to commit credit card registration at
//     MangoPay
//
// See http://docs.mangopay.com/api-references/card-registration/
//
// Example:
//
//  user := NewNaturalUser(...)
//  cr, err := NewCardRegistration(user, "EUR")
//  if err != nil {
//      log.Fatal(err)
//  }
//  if err := cr.Init(); err != nil {
//      log.Fatal(err)
//  }}
//
// Now render an HTML form for user card details (see Init()). Once submitted,
// you get the final token as a string starting with "data=". Use this token to
// finally register the card:
//
//  if err := cr.Register(token); err != nil {
//      log.Fatal(err)
//  }}
func (m *MangoPay) NewCardRegistration(user Consumer, currency string) (*CardRegistration, error) {
	id := consumerId(user)
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

	cr, err := c.service.anyRequest(new(CardRegistration), actionCreateCardRegistration, data)
	if err != nil {
		return err
	}
	// Backup private service
	service := c.service
	*c = *(cr.(*CardRegistration))
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
	cr, err := c.service.anyRequest(new(CardRegistration),
		actionSendCardRegistrationData, JsonObject{"Id": c.Id,
			"RegistrationData": registrationData})
	if err != nil {
		return err
	}
	// Backup private members
	serv := c.service
	isr := c.isInitialized
	*c = *(cr.(*CardRegistration))
	c.CardRegistrationData = registrationData
	c.service = serv
	c.isInitialized = isr
	return nil
}
