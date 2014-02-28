// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AuthMode defines authentication methods for communicating with
// the service.
type AuthMode int

const (
	// Basic Access Authentication
	BasicAuth AuthMode = iota
	// OAuth 2.0, token based authentication
	OAuth
)

// Holds OAuth2.0 token info.
type oAuth2 struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
	created   int64
}

func newToken(m *MangoPay) (*oAuth2, error) {
	if m == nil {
		return nil, errors.New("newToken: nil service")
	}

	gentoken := func(m *MangoPay) (*oAuth2, error) {
		// Let's get a new access token
		auth := basicAuthorization(m.clientId, m.password)
		u, err := url.Parse(rootURLs[m.env] + "oauth/token")
		if err != nil {
			return nil, err
		}

		body := "grant_type=client_credentials"
		req, err := http.NewRequest("POST", u.String(), strings.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", auth)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		access := new(oAuth2)
		if err := json.Unmarshal(b, access); err != nil {
			return nil, err
		}
		access.created = time.Now().Unix()
		// Now we can track the token during its lifetime
		m.oauth = access
		return access, nil
	}

	if m.oauth != nil {
		if time.Now().Unix()-m.oauth.created < m.oauth.ExpiresIn-60 {
			return m.oauth, nil // Reuse token
		}
	}
	access, err := gentoken(m)
	return access, err
}

// Returns authorization string for basic auth.
func basicAuthorization(clientId, passwd string) string {
	credential := clientId + ":" + passwd
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credential))
}

// Returns authorization string for OAuth2.0 auth.
func oAuthAuthorization(m *MangoPay) (string, error) {
	if m == nil {
		return "", errors.New("oAuth2.0: nil service")
	}
	o, err := newToken(m)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", o.Type, o.Token), nil
}

// Config hold environment credentials required for using the API.
//
// See http://docs.mangopay.com/api-references/sandbox-credentials/
type Config struct {
	ClientId   string
	Name       string
	Email      string
	Passphrase string
	Env        string
	env        ExecEnvironment
}

func (c *Config) String() string {
	return struct2string(c)
}

// NewConfig creates a config suitable for NewMangoPay().
func NewConfig(clientId, name, email, passwd, env string) *Config {
	c := &Config{ClientId: clientId, Name: name, Email: email, Passphrase: passwd, Env: env}
	if env == "sandbox" {
		c.env = Sandbox
	} else {
		c.env = Production
	}
	return c
}

// RegisterClient asks MangoPay to create a new client account.
func RegisterClient(clientId, name, email string, env ExecEnvironment) (*Config, error) {
	c := &Config{ClientId: clientId, Name: name, Email: email}
	body, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(fmt.Sprintf("%sclients/", rootURLs[env]))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	if env == Sandbox {
		c.Env = "sandbox"
	} else {
		c.Env = "production"
	}
	c.env = env
	return c, nil
}
