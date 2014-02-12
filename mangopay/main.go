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
	"github.com/matm/mangopay2-go-sdk"
	"io/ioutil"
	"os"
	_ "strconv"
)

// JSON config read from config file
type config struct {
	ClientId   string
	Name       string
	Email      string
	Passphrase string
	EnvStr     string `json:"Env"`
	env        mango.ExecEnvironment
}

func (c *config) String() string {
	return fmt.Sprintf("Name: %s\nClientId: %s\nEmail: %s", c.Name, c.ClientId, c.Email)
}

func perror(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(2)
}

// Parse config file
func parseConfig(configfile string) (*config, error) {
	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		return nil, err
	}
	conf := new(config)
	if err := json.Unmarshal(data, conf); err != nil {
		return nil, err
	}
	if conf.EnvStr == "production" {
		conf.env = mango.Production
	} else {
		if conf.EnvStr != "sandbox" {
			return nil, errors.New(fmt.Sprintf("unknown exec environment '%s'. "+
				"Must be one of production or sandbox.", conf.EnvStr))
		}
		conf.env = mango.Sandbox
	}
	return conf, nil
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
  wallet*           fetch wallet info
  editwallet*       update wallet info
  trwallet*         fetch all wallet's transactions

  addtransfer*      create a new tranfer
  transfer*         fetch transfer info

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

	service, err := mango.NewMangoPay(conf.ClientId, conf.Passphrase, conf.env)
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
		ows := mango.BuyerList{}
		for _, o := range w.Owners {
			u := new(mango.LegalUser)
			u.User = mango.User{Id: o}
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
		ows := mango.BuyerList{}
		for _, o := range w.Owners {
			u := new(mango.LegalUser)
			u.User = mango.User{Id: o}
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
		trs := w.Transactions()
		for _, t := range trs {
			fmt.Println(t)
		}
	case "addtransfer":
		t := &mango.Transfer{}
		if err := json.Unmarshal([]byte(*post), t); err != nil {
			perror(err.Error())
		}
		u := new(mango.LegalUser)
		u.User = mango.User{Id: t.AuthorId}
		k, err := service.NewTransfer(u, t.DebitedFunds, t.Fees, &mango.Wallet{Id: t.DebitedWalletId}, &mango.Wallet{Id: t.CreditedWalletId})
		if err != nil {
			perror(err.Error())
		}
		if err := k.Save(); err != nil {
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
	default:
		flag.Usage()
		perror(fmt.Sprintf("No such action '%s'.", action))
	}
}
