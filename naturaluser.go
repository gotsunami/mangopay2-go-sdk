// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
)

// NaturalUser describes all the properties of a MangoPay natural user object.
type NaturalUser struct {
	User
	FirstName, LastName string
	Address             string
	Birthday            int64
	Nationality         string
	CountryOfResidence  string
	Occupation          string
	IncomeRange         string
	ProofOfIdentity     string
	ProofOfAddress      string
	service             *MangoPay // Current service
	wallets             WalletList
}

func (u *NaturalUser) String() string {
	return struct2string(u)
}

// NewNaturalUser creates a new natural user.
func (m *MangoPay) NewNaturalUser(first, last string, email string, birthday int64, nationality, country string) *NaturalUser {
	u := &NaturalUser{
		FirstName:          first,
		LastName:           last,
		Birthday:           birthday,
		Nationality:        nationality,
		CountryOfResidence: country,
	}
	u.User = User{Email: email}
	u.service = m
	return u
}

// Wallets returns user's wallets.
func (u *NaturalUser) Wallets() (WalletList, error) {
	ws, err := u.service.wallets(u)
	return ws, err
}

// Transfer gets all user's transaction.
func (u *NaturalUser) Transfers(t string) (TransferList, error) {
	trs, err := u.service.transfers(u, t)
	return trs, err
}

// Save creates or updates a natural user. The Create API is used
// if the user's Id is an empty string. The Edit API is used when
// the Id is a non-empty string.
func (u *NaturalUser) Save() error {
	var action mangoAction
	if u.Id == "" {
		action = actionCreateNaturalUser
	} else {
		action = actionEditNaturalUser
	}

	data := JsonObject{}
	j, err := json.Marshal(u)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	// Force float64 to int conversion after unmarshalling.
	for _, field := range []string{"Birthday", "CreationDate"} {
		data[field] = int(data[field].(float64))
	}

	// Fields not allowed when creating a user
	if action == actionCreateNaturalUser {
		delete(data, "Id")
	}
	delete(data, "CreationDate")

	if action == actionEditNaturalUser {
		// Delete empty values so that existing ones don't get
		// overwritten with empty values.
		for k, v := range data {
			switch v.(type) {
			case string:
				if v.(string) == "" {
					delete(data, k)
				}
			case int:
				if v.(int) == 0 {
					delete(data, k)
				}
			}
		}
	}

	user, err := u.service.anyRequest(new(NaturalUser), action, data)
	if err != nil {
		return err
	}
	serv := u.service
	*u = *(user.(*NaturalUser))
	u.service = serv
	return nil
}

// NaturalUser finds a natural user using the user_id attribute.
func (m *MangoPay) NaturalUser(id string) (*NaturalUser, error) {
	u, err := m.anyRequest(new(NaturalUser), actionFetchNaturalUser, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return u.(*NaturalUser), nil
}
