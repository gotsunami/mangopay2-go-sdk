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
  addnatuser        create a natural user

Options:
`))
		/* OLD actions
		   fetchuser         get user info
		   updateuser        update user info
		   fetchuserwallets  get wallets of a user

		   createbenef       create a beneficiary
		   fetchbenef        get beneficiary info

		   createpc          create a payment card
		   deletepc          delete a payment card
		   fetchpc           get payment card info
		   fetchuserpc       get user payment cards

		   createwallet      create a wallet
		   fetchwallet       get wallet info
		   updatewallet      update wallet info
		   fetchuserswallet  get users of a wallet
		*/
		flag.PrintDefaults()
		os.Exit(2)
	}

	post := flag.String("d", "", "JSON for POST or PUT data")
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

	// FIXME: set an option
	service.Option(mango.Verbosity(mango.Debug))

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
		u := service.NewNaturalUser()
		if err := json.Unmarshal([]byte(*post), u); err != nil {
			perror(err.Error())
		}
		if err := u.Save(); err != nil {
			perror(err.Error())
		}
		fmt.Println("Natural user created:")
		fmt.Println(u)
		/*
			case "fetchuser":
				var data struct {
					UserId int `json:"user_id"`
				}
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				u, err := mango.FindUser(service, data.UserId)
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(u)
			case "updateuser":
				data := new(mango.JsonObject)
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				b := *data
				if _, ok := b["user_id"]; !ok {
					perror("missing user_id in JSON data")
				}
				var val int
				switch b["user_id"].(type) {
				case float64:
					val = int(b["user_id"].(float64))
				case string:
					i, err := strconv.ParseInt(b["user_id"].(string), 10, 0)
					if err != nil {
						perror(err.Error())
					}
					val = int(i)
				}
				u, err := mango.FindUser(service, val)
				if err != nil {
					perror(err.Error())
				}
				u, err = u.Update(*data)
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(u)
			case "createbenef":
				data := new(mango.JsonObject)
				if err := json.Unmarshal([]byte(*post), data); err != nil {
					perror(err.Error())
				}
				b, err := mango.NewBeneficiary(service, *data)
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(b)
			case "fetchbenef":
				var data struct {
					BeneficiaryId int `json:"beneficiary_id"`
				}
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				b, err := mango.FindBeneficiary(service, data.BeneficiaryId)
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(b)
			case "createpc":
				data := new(mango.JsonObject)
				if err := json.Unmarshal([]byte(*post), data); err != nil {
					perror(err.Error())
				}
				pc, err := mango.NewPaymentCard(service, *data)
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(pc)
			case "fetchpc":
				var data struct {
					PaymentCardId int `json:"paymentcard_id"`
				}
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				pc, err := mango.FindPaymentCard(service, data.PaymentCardId)
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(pc)
			case "fetchuserpc":
				var data struct {
					UserId int `json:"user_id"`
				}
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				u, err := mango.FindUser(service, data.UserId)
				if err != nil {
					perror(err.Error())
				}
				pcs, err := u.PaymentCards()
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(pcs)
			case "deletepc":
				var data struct {
					PaymentCardId int `json:"paymentcard_id"`
				}
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				pc, err := mango.FindPaymentCard(service, data.PaymentCardId)
				if err != nil {
					perror(err.Error())
				}
				if err := pc.Delete(); err != nil {
					perror(err.Error())
				}
			case "createwallet":
				data := new(mango.JsonObject)
				if err := json.Unmarshal([]byte(*post), data); err != nil {
					perror(err.Error())
				}
				w, err := mango.NewWallet(service, *data)
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(w)
			case "fetchwallet":
				var data struct {
					WalletId int `json:"wallet_id"`
				}
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				w, err := mango.FindWallet(service, data.WalletId)
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(w)
			case "updatewallet":
				data := new(mango.JsonObject)
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				b := *data
				if _, ok := b["wallet_id"]; !ok {
					perror("missing wallet_id in JSON data")
				}
				var val int
				switch b["wallet_id"].(type) {
				case float64:
					val = int(b["wallet_id"].(float64))
				case string:
					i, err := strconv.ParseInt(b["wallet_id"].(string), 10, 0)
					if err != nil {
						perror(err.Error())
					}
					val = int(i)
				}
				w, err := mango.FindWallet(service, val)
				if err != nil {
					perror(err.Error())
				}
				r, err := w.Update(*data)
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(r)
			case "fetchuserswallet":
				var data struct {
					WalletId int `json:"wallet_id"`
				}
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				w, err := mango.FindWallet(service, data.WalletId)
				if err != nil {
					perror(err.Error())
				}
				users, err := w.FindUsers()
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(users)
			case "fetchuserwallets":
				var data struct {
					UserId int `json:"user_id"`
				}
				if err := json.Unmarshal([]byte(*post), &data); err != nil {
					perror(err.Error())
				}
				u, err := mango.FindUser(service, data.UserId)
				if err != nil {
					perror(err.Error())
				}
				ws, err := u.FindWallets()
				if err != nil {
					perror(err.Error())
				}
				fmt.Println(ws)
		*/
	default:
		flag.Usage()
		perror(fmt.Sprintf("No such action '%s'.", action))
	}
}
