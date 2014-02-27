// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	_ "errors"
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
}

func (c *Config) String() string {
	return struct2string(c)
}

// RegisterClient asks MangoPay to create a new client account.
func RegisterClient(clientId, name, email string) (*Config, error) {
	return nil, nil
}
