// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

type mangoAction int

const (
	actionEvents mangoAction = iota
	actionAllUsers

	actionCreateNaturalUser
	actionEditNaturalUser
	actionFetchNaturalUser

	actionCreateLegalUser
	actionEditLegalUser
	actionFetchLegalUser

	actionFetchUser
)

// JsonObject is used to manage JSON data.
type JsonObject map[string]interface{}

type mangoRequest struct {
	Method, Path string
	PathValues   JsonObject
}

// Defines mango requests metadata.
var mangoRequests = map[mangoAction]mangoRequest{
	actionEvents: mangoRequest{
		"GET",
		"/events",
		nil,
	},
	actionCreateNaturalUser: mangoRequest{
		"POST",
		"/users/natural",
		nil,
	},
	actionEditNaturalUser: mangoRequest{
		"PUT",
		"/users/natural/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionAllUsers: mangoRequest{
		"GET",
		"/users",
		nil,
	},
	actionFetchNaturalUser: mangoRequest{
		"GET",
		"/users/natural/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionCreateLegalUser: mangoRequest{
		"POST",
		"/users/legal",
		nil,
	},
	actionEditLegalUser: mangoRequest{
		"PUT",
		"/users/legal/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchLegalUser: mangoRequest{
		"GET",
		"/users/legal/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchUser: mangoRequest{
		"GET",
		"/users/{{Id}}",
		JsonObject{"Id": ""},
	},
}
