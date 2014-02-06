// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"fmt"
)

// NaturalUser describes all the properties of a MangoPay legal user object.
type NaturalUser struct {
	Tag                 string
	Email               string
	FirstName, LastName string
	Address             string
	Birthday            int
	Nationality         string
	CountryOfResidence  string
	Occupation          string
	IncomeRange         string
	ProofOfIdentity     string
	ProofOfAddress      string
	CreationDate        int
	PersonType          string
	Id                  string
	service             *MangoPay // Current service
}

func (u *NaturalUser) String() string {
	return fmt.Sprintf(`Id                      : %s
Tag                     : %s
First name              : %s
Last name               : %s
Email                   : %s
Address                 : %s
Birthday                : %d
Nationality             : %s
Country of residence    : %s
Occupation              : %s
Income range            : %s
Proof of identity       : %s
Proof of address        : %s
Creation date           : %d
Person type             : %s`, u.Id, u.Tag, u.FirstName, u.LastName, u.Email, u.Address, u.Birthday, u.Nationality, u.CountryOfResidence, u.Occupation, u.IncomeRange, u.ProofOfIdentity, u.ProofOfAddress, u.CreationDate, u.PersonType)
}

func (m *MangoPay) userRequest(action MangoAction, data JsonObject) (*NaturalUser, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}
	u := new(NaturalUser)
	if err := unMarshalJSONResponse(resp, u); err != nil {
		return nil, err
	}
	return u, nil
}

// NewNaturalUser creates a new natural user.
func (m *MangoPay) NewNaturalUser() *NaturalUser {
	u := new(NaturalUser)
	u.service = m
	return u
}

// Save creates or updates a natural user. The Create API is used
// if the user's Id is an empty string. The Edit API is used when
// the Id is a non-empty string.
func (u *NaturalUser) Save() error {
	var action MangoAction
	if u.Id == "" {
		action = ActionCreateNaturalUser
	} else {
		action = ActionEditNaturalUser
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
	if action == ActionCreateNaturalUser {
		delete(data, "Id")
		delete(data, "CreationDate")
	}

	user, err := u.service.userRequest(action, data)
	if err != nil {
		return err
	}
	*u = *user
	return nil
}

// NaturalUser finds a natural user using the user_id attribute.
func (m *MangoPay) NaturalUser(uid int) (*NaturalUser, error) {
	/*
		u, err := s.userRequest(ActionFetchNaturalUser, JsonObject{"user_id": uid})
		if err != nil {
			return nil, err
		}
		u.service = s
	*/
	return nil, nil
}
