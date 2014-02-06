// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import ()

type option func(*MangoPay)

type Level int

const (
	Info Level = iota
	Debug
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

// Sets verbosity level. Default verbosity level is Info.
func Verbosity(v Level) option {
	return func(m *MangoPay) {
		m.verbosity = v
	}
}

// AuthMethod sets the preferred method for authenticating against
// the service.
func AuthMethod(auth AuthMode) option {
	return func(m *MangoPay) {
		m.authMethod = auth
	}
}
