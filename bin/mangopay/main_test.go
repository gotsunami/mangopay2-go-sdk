// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package mango is a library for the MangoPay service v2.
//
// http://www.mangopay.com
package main

import (
	"encoding/json"
	"github.com/github.com/Adrien-P/mangopay2-go-sdk"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	env      = "sandbox"
	testconf = "testing.conf" // Credentials used for testing
)

var (
	service        *mango.MangoPay
	birth1, birth2 int64
	users          []*mango.NaturalUser
	usersinfo      []user
	noFees         = mango.Money{currency, 0}
	transfer       *mango.Transfer
	payin          *mango.PayIn
)

type user struct {
	first, last    string
	email, country string
	birthday       int64
	ccn, cvv, exp  string // Credit card number, CVV, exp. date (MMYY)
	wallet         *mango.Wallet
	card           *mango.Card
}

const (
	firstName1 = "Alice"
	lastName1  = "Doe"
	email1     = "alice@doe.org"

	firstName2 = "Bob"
	lastName2  = "Doe"
	email2     = "bob@doe.org"

	country  = "FR"
	currency = "EUR"
)

func init() {
	var f *os.File
	var err error
	var c *mango.Config

	f, err = os.Open(testconf)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(testconf)
			ti := strconv.FormatInt(time.Now().Unix(), 10)
			c, err = mango.RegisterClient("testclient"+ti, "A name",
				"m"+ti+"@gmail.com", mango.Sandbox)
			m, _ := json.Marshal(c)
			f.Write(m)
			f.Close()
		} else {
			panic(err)
		}
	} else {
		c, err = parseConfig(testconf)
		if err != nil {
			panic(err)
		}
	}

	log.Printf("Running tests in sandbox as user %s", c.ClientId)

	service, err = mango.NewMangoPay(c, mango.OAuth)
	if err != nil {
		log.Fatalf("can't use service: %s\n", err.Error())
	}
	birth1 = time.Now().Add(-20 * 24 * time.Hour * 365).Unix()
	birth2 = time.Now().Add(-25 * 24 * time.Hour * 365).Unix()

	usersinfo = []user{
		user{firstName1, lastName1, email1, country, birth1,
			"4970101122334463", "123", "0219", nil, nil},
		user{firstName2, lastName2, email2, country, birth2,
			"4970101122334471", "123", "0919", nil, nil},
	}
	users = make([]*mango.NaturalUser, 2)
}

func TestNewNaturalUser(t *testing.T) {
	for k, u := range usersinfo {
		log.Printf("Creating user %s ...", u.first)
		users[k] = service.NewNaturalUser(u.first, u.last, u.email, u.birthday,
			u.country, u.country)
		if err := users[k].Save(); err != nil {
			t.Fatalf("can't create user: " + err.Error())
		}
	}
}

func TestFetchNaturalUser(t *testing.T) {
	for k, u := range usersinfo {
		log.Printf("Fetching user %s ...", u.first)
		if _, err := service.NaturalUser(users[k].Id); err != nil {
			t.Errorf("can't find user %s", u.first)
		}
		log.Printf("%s has Id %s", u.first, users[k].Id)
	}
}

func TestNewWallet(t *testing.T) {
	for k, _ := range usersinfo {
		u := users[k]
		log.Printf("Creating wallet for %s ...", u.FirstName)
		w, err := service.NewWallet(mango.ConsumerList{u}, u.FirstName+"'s wallet", currency)
		if err != nil {
			t.Errorf("can't create wallet for %s: %s", u.FirstName, err.Error())
		}
		if err := w.Save(); err != nil {
			t.Errorf("can't save wallet for %s: %s", u.FirstName, err.Error())
		}
		usersinfo[k].wallet = w
		log.Printf("%s has wallet Id %s", u.FirstName, w.Id)
	}
}

func TestRegisterCreditCard(t *testing.T) {
	for k, u := range users {
		log.Printf("New credit card for %s ... ", u.FirstName)
		card, err := service.NewCardRegistration(u, currency)
		if err != nil {
			t.Fatal(err.Error())
		}
		if err := card.Init(); err != nil {
			t.Fatal(err.Error())
		}
		log.Printf("Using fake credit card number for %s: %s", u.FirstName, usersinfo[k].ccn)
		// Simulates a user-supplied HTML form POST to the external
		// bank service.
		rdata, err := sendRegistrationData(card, usersinfo[k].ccn,
			usersinfo[k].exp, usersinfo[k].cvv)
		if err != nil {
			t.Fatal(err.Error())
		}
		if err := card.Register(rdata); err != nil {
			t.Fatal(err.Error())
		}
		log.Printf("%s has card Id %s", u.FirstName, card.CardId)
		c, err := service.Card(card.CardId)
		if err != nil {
			t.Fatal(err.Error())
		}
		usersinfo[k].card = c
	}
}

func TestDirectPayin(t *testing.T) {
	amount := 100
	for k, u := range users {
		log.Printf("Sending %d EUR to %s's wallet ... ", amount, u.FirstName)
		p, err := service.NewDirectPayIn(u, u, usersinfo[k].card,
			usersinfo[k].wallet, mango.Money{currency, amount * 100}, noFees,
			"http://myreturnurl")
		if err != nil {
			t.Fatal(err.Error)
		}
		if err := p.Save(); err != nil {
			t.Fatal(err.Error())
		}
		// Save Bob's direct paying for later refund
		if k == 1 {
			payin = &p.PayIn
		}
	}
	for k, _ := range users {
		w, err := service.Wallet(usersinfo[k].wallet.Id)
		if err != nil {
			t.Fatal(err.Error())
		}
		if w.Balance.Amount != amount*100 {
			t.Errorf("expect 100 EUR in wallet, got %d", w.Balance.Amount/100)
		}
	}
}

func TestTransferBetweenWallets(t *testing.T) {
	amount := 30
	fees := 2
	log.Printf("Alice pays %d EUR to Bob (%d EUR fees) ...", amount, fees)
	tr, err := service.NewTransfer(users[0], mango.Money{currency, amount * 100},
		mango.Money{currency, int(fees * 100)},
		usersinfo[0].wallet, usersinfo[1].wallet)
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := tr.Save(); err != nil {
		t.Fatal(err.Error())
	}
	// Hold current successful transfer for later refund
	transfer = tr
	log.Printf("Checking Bob's wallet balance ...")
	w, err := service.Wallet(usersinfo[1].wallet.Id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if w.Balance.Amount != (10000 + 2800) {
		t.Errorf("wrong Bob's wallet balance. Expect %.2f, got %.2f", 128, w.Balance.Amount)
	}
}

func TestTransferRefund(t *testing.T) {
	log.Printf("Try refunding Alice ...")
	if _, err := transfer.Refund(); err != nil {
		t.Fatal(err.Error())
	}
	log.Printf("Checking Alice's wallet balance ...")
	w, err := service.Wallet(usersinfo[0].wallet.Id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if w.Balance.Amount != 10000 {
		t.Errorf("wrong Alice's wallet balance. Expect 100 EUR, got %.2f", w.Balance.Amount)
	}
}

func TestPayInRefund(t *testing.T) {
	log.Printf("Try refunding Bob after direct payin ...")
	_, err := payin.Refund()
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Printf("Checking Bob's wallet balance ...")
	w, err := service.Wallet(usersinfo[1].wallet.Id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if w.Balance.Amount != 0 {
		t.Errorf("wrong Bob's wallet balance. Expect 0 EUR, got %.2f", w.Balance.Amount)
	}
}

func TestCreateBankAccounts(t *testing.T) {
	banks := []struct{ IBAN, BIC string }{
		// Using different IBAN/BIC numbers seems not to be supported by the
		// testing API.
		{"FR3020041010124530725S03383", "CRLYFRPP"},
		{"FR3020041010124530725S03383", "CRLYFRPP"},
	}
	for k, u := range users {
		log.Printf("Creating IBAN account for %s ...", u.FirstName)
		acc, err := service.NewBankAccount(u, u.FirstName, "one great place", mango.IBAN)
		if err != nil {
			t.Fatal(err.Error())
		}
		acc.IBAN, acc.BIC = banks[k].IBAN, banks[k].BIC
		if err := acc.Save(); err != nil {
			t.Fatal(err.Error())
		}
	}
}
