// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

type option func(*MangoPay)

type Level int

const (
	Info Level = iota
	Debug
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
