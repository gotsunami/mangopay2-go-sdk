// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

// Config hold environment credentials required for using the API.
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
	c.env = env
	if env == Sandbox {
		c.Env = "sandbox"
	} else {
		c.Env = "production"
	}
	return c, nil
}
