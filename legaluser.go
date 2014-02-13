// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"fmt"
)

// LegalUser describes all the properties of a MangoPay legal user object.
type LegalUser struct {
	User
	Name                                  string
	LegalPersonType                       string
	HeadquartersAddress                   string
	LegalRepresentativeFirstName          string
	LegalRepresentativeLastName           string
	LegalRepresentativeAddress            string
	LegalRepresentativeEmail              string
	LegalRepresentativeBirthday           int64
	LegalRepresentativeNationality        string
	LegalRepresentativeCountryOfResidence string
	Statute                               string
	ProofOfRegistration                   string
	ShareholderDeclaration                string
	service                               *MangoPay // Current service
	wallets                               WalletList
}

func (u *LegalUser) String() string {
	return fmt.Sprintf(`
Id                                      : %s
Tag                                     : %s
Email                                   : %s
Creation date                           : %s
Person type                             : %s
Name                                    : %s
Legal Person Type                       : %s
Headquarters Address                    : %s
Legal Representative FirstName          : %s
Legal Representative LastName           : %s
Legal Representative Address            : %s
Legal Representative Email              : %s
Legal Representative Birthday           : %s
Legal Representative Nationality        : %s
Legal Representative CountryOfResidence : %s
Statute                                 : %s
Proof Of Registration                   : %s
Shareholder Declaration                 : %s
`, u.Id, u.Tag, u.Email, unixTimeToString(u.CreationDate), u.PersonType, u.Name, u.LegalPersonType, u.HeadquartersAddress, u.LegalRepresentativeFirstName, u.LegalRepresentativeLastName, u.LegalRepresentativeAddress, u.LegalRepresentativeEmail, unixTimeToString(u.LegalRepresentativeBirthday), u.LegalRepresentativeNationality, u.LegalRepresentativeCountryOfResidence, u.Statute, u.ProofOfRegistration, u.ShareholderDeclaration)
}

// NewLegalUser creates a new legal user.
func (m *MangoPay) NewLegalUser(name string, email string, personType string, legalFirstName, legalLastName string, birthday int64, nationality string, country string) *LegalUser {
	u := &LegalUser{
		Name:                                  name,
		LegalPersonType:                       personType,
		LegalRepresentativeFirstName:          legalFirstName,
		LegalRepresentativeLastName:           legalLastName,
		LegalRepresentativeBirthday:           birthday,
		LegalRepresentativeNationality:        nationality,
		LegalRepresentativeCountryOfResidence: country}
	u.User = User{Email: email}
	u.service = m
	return u
}

// Wallets returns user's wallets.
func (u *LegalUser) Wallets() WalletList {
	u.wallets = nil
	return u.wallets
}

// Save creates or updates a legal user. The Create API is used
// if the user's Id is an empty string. The Edit API is used when
// the Id is a non-empty string.
func (u *LegalUser) Save() error {
	var action mangoAction
	if u.Id == "" {
		action = actionCreateLegalUser
	} else {
		action = actionEditLegalUser
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
	for _, field := range []string{"LegalRepresentativeBirthday", "CreationDate"} {
		data[field] = int(data[field].(float64))
	}

	// Fields not allowed when creating a user
	if action == actionCreateLegalUser {
		delete(data, "Id")
	}
	delete(data, "CreationDate")

	if action == actionEditLegalUser {
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

	user, err := u.service.legalUserRequest(action, data)
	if err != nil {
		return err
	}
	*u = *user
	return nil
}

// LegalUser finds a legal user using the user_id attribute.
func (m *MangoPay) LegalUser(id string) (*LegalUser, error) {
	u, err := m.legalUserRequest(actionFetchLegalUser, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *MangoPay) legalUserRequest(action mangoAction, data JsonObject) (*LegalUser, error) {
	resp, err := m.request(action, data)
	if err != nil {
		return nil, err
	}
	u := new(LegalUser)
	if err := m.unMarshalJSONResponse(resp, u); err != nil {
		return nil, err
	}
	return u, nil
}
