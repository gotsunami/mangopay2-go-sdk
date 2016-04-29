// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package mango is a library for the MangoPay service v2.
//
// http://www.mangopay.com
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Adrien-P/mangopay2-go-sdk"
)

func perror(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}

// sendRegistrationData sends the user's card number, expiration date and cvx and
// returns the registration data token at the specified returnUrl.
//
// Note that this function is for __testing__ only as this is NOT the correct
// way to proceed. Indeed, user's card details must be sent directly
// through an HTML form to the CardRegistrationUrl (which is an external banking
// service).
func sendRegistrationData(c *mango.CardRegistration, cardNumber,
	expirationDate, cvx string) (string, error) {
	data := url.Values{
		"data":               []string{c.PreregistrationData},
		"accessKeyRef":       []string{c.AccessKey},
		"cardNumber":         []string{cardNumber},
		"cardExpirationDate": []string{expirationDate},
		"cardCvx":            []string{cvx},
	}

	req, err := http.NewRequest("POST", c.CardRegistrationUrl,
		strings.NewReader(string([]byte(data.Encode()))))
	if err != nil {
		return "", err
	}
	resp, err := mango.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Error code %d: %v", resp.StatusCode, err))
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	c.CardRegistrationData = string(b)
	return c.CardRegistrationData, nil
}

// Parse config file
func parseConfig(configfile string) (*mango.Config, error) {
	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		return nil, err
	}
	conf := new(mango.Config)
	if err := json.Unmarshal(data, conf); err != nil {
		return nil, err
	}
	c, err := mango.NewConfig(conf.ClientId, conf.Name, conf.Email,
		conf.Passphrase, conf.Env)
	return c, err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Usage: %s [options] action configfile\n", os.Args[0]))
		fmt.Fprintf(os.Stderr, fmt.Sprintf(` 
where action is one of: 
  conf              show config
  events            list all events (PayIns, PayOuts, Transfers)
  users             list all users
  user*             fetch a user (natural or legal)

  addnatuser*       create a natural user
  editnatuser*      update natural user info
  natuser*          fetch natural user info

  addlegaluser*     create a legal user
  editlegaluser*    update legal user info
  legaluser*        fetch legal user info

  addwallet*        create a new wallet
  editwallet*       update wallet info
  trwallet*         fetch all wallet's transactions
  wallet*           fetch wallet info
  wallets*          list all user's wallet

  addtransfer*      create a new tranfer
  transfer*         fetch transfer info
  transfers*        list all user's transactions

  addwebpayin*      create a payIn through web interface
  adddirectpayin*   create a direct payIn (with tokenized card)
  payin*            fetch a payIn

  addcard*          register a credit card
  card*             fetch a credit card
  cards*            list all user's cards

  addrefund*        refund a payment (provide TransferId or PayInId)
  refund*           fetch a refund (transfer or payin)

  addaccount*       create an IBAN bank account
  account*          fetch a user's bank account
  accounts*         list all user's bank accounts

  addpayout*        create a bank wire
  payout*           fetch a bank wire

Actions with an asterisk(*) require input JSON data (-d).

Options:
`))
		flag.PrintDefaults()
		os.Exit(2)
	}

	post := flag.String("d", "", "JSON data part of the HTTP request")
	verbose := flag.Int("v", 0, "Verbosity level (1 for debug)")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Fprint(os.Stderr, "action or config file missing.\n")
		flag.Usage()
	}

	conf, err := parseConfig(flag.Arg(1))
	if err != nil {
		perror(fmt.Sprintf("config parsing error: %s\n", err.Error()))
	}

	service, err := mango.NewMangoPay(conf, mango.OAuth)
	if err != nil {
		perror(fmt.Sprintf("can't use service: %s\n", err.Error()))
	}

	if *verbose == 1 {
		service.Option(mango.Verbosity(mango.Debug))
	}

	action := flag.Arg(0)
	switch action {
	case "conf":
		fmt.Println(conf)
	case "events":
		evs, err := service.Events()
		if err != nil {
			perror(err.Error())
		}
		if len(evs) == 0 {
			fmt.Println("No event.")
		} else {
			for _, ev := range evs {
				fmt.Println(ev)
			}
		}
	case "addnatuser":
		u := &mango.NaturalUser{}
		if err := json.Unmarshal([]byte(*post), u); err != nil {
			perror(err.Error())
		}
		n := service.NewNaturalUser(u.FirstName, u.LastName, u.Email, u.Birthday, u.Nationality, u.CountryOfResidence)
		if err := n.Save(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Natural user created:")
		fmt.Println(n)
	case "editnatuser":
		u := &mango.NaturalUser{}
		if err := json.Unmarshal([]byte(*post), u); err != nil {
			perror(err.Error())
		}
		n := service.NewNaturalUser(u.FirstName, u.LastName, u.Email, u.Birthday, u.Nationality, u.CountryOfResidence)
		if err := n.Save(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Natural user updated:")
		fmt.Println(n)
	case "users":
		users, err := service.Users()
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(users)
		for _, u := range users {
			fmt.Println(u)
		}
	case "natuser":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		u, err := service.NaturalUser(data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(u)
	case "addlegaluser":
		u := &mango.LegalUser{}
		if err := json.Unmarshal([]byte(*post), u); err != nil {
			perror(err.Error())
		}
		n := service.NewLegalUser(u.Name, u.Email, u.LegalPersonType, u.LegalRepresentativeFirstName, u.LegalRepresentativeLastName, u.LegalRepresentativeBirthday, u.LegalRepresentativeNationality, u.LegalRepresentativeCountryOfResidence)
		if err := n.Save(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Legal user created:")
		fmt.Println(n)
	case "editlegaluser":
		u := &mango.LegalUser{}
		if err := json.Unmarshal([]byte(*post), u); err != nil {
			perror(err.Error())
		}
		n := service.NewLegalUser(u.Name, u.Email, u.LegalPersonType, u.LegalRepresentativeFirstName, u.LegalRepresentativeLastName, u.LegalRepresentativeBirthday, u.LegalRepresentativeNationality, u.LegalRepresentativeCountryOfResidence)
		if err := n.Save(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Legal user updated:")
		fmt.Println(n)
	case "legaluser":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		u, err := service.LegalUser(data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(u)
	case "user":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		u, err := service.User(data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(u)
	case "addwallet":
		w := &mango.Wallet{}
		if err := json.Unmarshal([]byte(*post), w); err != nil {
			perror(err.Error())
		}
		ows := mango.ConsumerList{}
		for _, o := range w.Owners {
			u := new(mango.LegalUser)
			u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: o}}
			ows = append(ows, u)
		}
		n, err := service.NewWallet(ows, w.Description, w.Currency)
		if err != nil {
			perror(err.Error())
		}
		if err := n.Save(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Wallet created:")
		fmt.Println(n)
	case "wallet":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		w, err := service.Wallet(data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(w)
	case "editwallet":
		w := &mango.Wallet{}
		if err := json.Unmarshal([]byte(*post), w); err != nil {
			perror(err.Error())
		}
		ows := mango.ConsumerList{}
		for _, o := range w.Owners {
			u := new(mango.LegalUser)
			u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: o}}
			ows = append(ows, u)
		}
		n, err := service.NewWallet(ows, w.Description, w.Currency)
		if err != nil {
			perror(err.Error())
		}
		if err := n.Save(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Wallet updated:")
		fmt.Println(n)
	case "trwallet":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		w, err := service.Wallet(data.Id)
		if err != nil {
			perror(err.Error())
		}
		trs, err := w.Transactions()
		if err != nil {
			perror(err.Error())
		}
		for _, t := range trs {
			fmt.Println(t)
		}
	case "wallets":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: data.Id}}
		ws, err := service.Wallets(u)
		if err != nil {
			perror(err.Error())
		}
		if len(ws) == 0 {
			fmt.Println("No wallets.")
		} else {
			for _, tr := range ws {
				fmt.Println(tr)
			}
		}
	case "addtransfer":
		t := &mango.Transfer{}
		if err := json.Unmarshal([]byte(*post), t); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: t.AuthorId}}
		k, err := service.NewTransfer(u, t.DebitedFunds, t.Fees,
			&mango.Wallet{
				ProcessIdent: mango.ProcessIdent{Id: t.DebitedWalletId},
			},
			&mango.Wallet{
				ProcessIdent: mango.ProcessIdent{Id: t.CreditedWalletId},
			})
		if err != nil {
			perror(err.Error())
		}
		if err := k.Save(); err != nil {
			if _, ok := err.(*mango.ErrTransferFailed); ok {
				fmt.Println(k)
			}
			perror(err.Error())
		}
		fmt.Println("Transfer created:")
		fmt.Println(k)
	case "transfer":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		t, err := service.Transfer(data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(t)
	case "transfers":
		var data struct {
			Id   string
			Type string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: data.Id}}
		trs, err := service.Transfers(u, data.Type)
		if err != nil {
			perror(err.Error())
		}
		if len(trs) == 0 {
			fmt.Println("No transfers.")
		} else {
			for _, tr := range trs {
				fmt.Println(tr)
			}
		}
	case "addwebpayin":
		w := &mango.WebPayIn{}
		if err := json.Unmarshal([]byte(*post), w); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: w.AuthorId}}

		k, err := service.NewWebPayIn(u, w.DebitedFunds, w.Fees,
			&mango.Wallet{
				ProcessIdent: mango.ProcessIdent{Id: w.CreditedWalletId},
			},
			w.ReturnUrl, w.Culture, w.TemplateURLOptions)
		if err != nil {
			perror(err.Error())
		}
		if err := k.Save(); err != nil {
			if _, ok := err.(*mango.ErrPayInFailed); ok {
				fmt.Println(k)
			}
			perror(err.Error())
		}
		fmt.Println("Web PayIn created:")
		fmt.Println(k)
	case "payin":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		t, err := service.PayIn(data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(t)
	case "addcard":
		c := &mango.CardRegistration{}
		if err := json.Unmarshal([]byte(*post), c); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: c.Id}}
		card, err := service.NewCardRegistration(u, c.Currency)
		if err != nil {
			perror(err.Error())
		}
		if err := card.Init(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Card registration:")
		fmt.Println(card)

		fmt.Println("Sending credit card info for registering this test card:")
		num, date, ccv := "4929683808277688", "0217", "184"

		fmt.Println("Okay! now registering the card to MangoPay...")
		// Simulates a user-supplied HTML form POST to the external
		// bank service.
		rdata, err := sendRegistrationData(card, num, date, ccv)
		if err != nil {
			perror(err.Error())
		}
		if err := card.Register(rdata); err != nil {
			perror(err.Error())
		}
		fmt.Println(card)
	case "card":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		c, err := service.Card(data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(c)
	case "adddirectpayin":
		w := &mango.DirectPayIn{}
		if err := json.Unmarshal([]byte(*post), w); err != nil {
			perror(err.Error())
		}
		from := new(mango.LegalUser)
		from.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: w.AuthorId}}
		to := new(mango.LegalUser)
		to.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: w.CreditedUserId}}
		card := mango.Card{ProcessIdent: mango.ProcessIdent{Id: w.CardId}}
		wallet := mango.Wallet{ProcessIdent: mango.ProcessIdent{Id: w.CreditedWalletId}}
		k, err := service.NewDirectPayIn(from, to, &card, &wallet,
			w.DebitedFunds, w.Fees, w.SecureModeReturnUrl)
		if err != nil {
			perror(err.Error())
		}
		if err := k.Save(); err != nil {
			if _, ok := err.(*mango.ErrPayInFailed); ok {
				fmt.Println(k)
			}
			perror(err.Error())
		}
		fmt.Println("Direct PayIn created:")
		fmt.Println(k)
	case "cards":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: data.Id}}
		cs, err := service.Cards(u)
		if err != nil {
			perror(err.Error())
		}
		if len(cs) == 0 {
			fmt.Println("No card.")
		} else {
			for _, c := range cs {
				fmt.Println(c)
			}
		}
	case "addrefund":
		r := &mango.Refund{}
		if err := json.Unmarshal([]byte(*post), r); err != nil {
			perror(err.Error())
		}

		t, err := service.Transfer(r.InitialTransactionId)
		if err != nil {
			perror(err.Error())
		}
		r, err = t.Refund()
		if err != nil {
			perror(err.Error())
		}
		fmt.Println("Transfer refund:")
		fmt.Println(r)
	case "refund":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		r, err := service.Refund(data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(r)
	case "addaccount":
		r := &mango.BankAccount{}
		if err := json.Unmarshal([]byte(*post), r); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: r.UserId}}
		acc, err := service.NewBankAccount(u, r.OwnerName, r.OwnerAddress, mango.IBAN)
		if err != nil {
			perror(err.Error())
		}
		acc.IBAN = r.IBAN
		acc.BIC = r.BIC
		if err := acc.Save(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Bank account:")
		fmt.Println(acc)
	case "account":
		var data struct {
			Id     string
			UserId string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: data.UserId}}
		a, err := service.BankAccount(u, data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println(a)
	case "accounts":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: data.Id}}
		accs, err := service.BankAccounts(u)
		if err != nil {
			perror(err.Error())
		}
		if len(accs) == 0 {
			fmt.Println("No bank account.")
		} else {
			for _, acc := range accs {
				fmt.Println(acc)
			}
		}
	case "addpayout":
		r := &mango.PayOut{}
		if err := json.Unmarshal([]byte(*post), r); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{ProcessIdent: mango.ProcessIdent{Id: r.AuthorId}}
		wallet := mango.Wallet{ProcessIdent: mango.ProcessIdent{Id: r.DebitedWalletId}}
		acc := mango.BankAccount{ProcessIdent: mango.ProcessIdent{Id: r.BankAccountId}}
		pay, err := service.NewPayOut(u, r.DebitedFunds, r.Fees, &wallet, &acc)
		if err != nil {
			perror(err.Error())
		}
		if err := pay.Save(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Bank wire:")
		fmt.Println(pay)
	case "payout":
		var data struct {
			Id string
		}
		if err := json.Unmarshal([]byte(*post), &data); err != nil {
			perror(err.Error())
		}
		p, err := service.PayOut(data.Id)
		if err != nil {
			perror(err.Error())
		}
		fmt.Println("Bank wire:")
		fmt.Println(p)
	default:
		flag.Usage()
		perror(fmt.Sprintf("No such action '%s'.", action))
	}
}
