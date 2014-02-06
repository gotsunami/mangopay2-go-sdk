// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

type mangoAction int

const (
	actionEvents mangoAction = iota
	actionCreateNaturalUser
	actionEditNaturalUser

	/*
		actionCreateUser Mangoaction = iota
		actionFetchUser
		actionUpdateUser
		actionFetchUserWallets

		actionCreateBeneficiary
		actionFetchBeneficiary

		actionCreatePaymentCard
		actionFetchPaymentCard
		actionFetchUserPaymentCards
		actionDeletePaymentCard

		actionCreateWallet
		actionFetchWallet
		actionUpdateWallet
		actionFetchUsersOfWallet
	*/
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
	/*
		// User
		actionFetchUser: mangoRequest{
			"GET",
			"/users/{{user_id}}",
			JsonObject{"user_id": ""},
		},
		actionUpdateUser: mangoRequest{
			"PUT",
			"/users/{{user_id}}",
			JsonObject{"user_id": ""},
		},
		actionFetchUserWallets: mangoRequest{
			"GET",
			"/users/{{user_id}}/wallets",
			JsonObject{"user_id": ""},
		},
		// Beneficiary
		actionCreateBeneficiary: mangoRequest{
			"POST",
			"/beneficiaries",
			nil,
		},
		actionFetchBeneficiary: mangoRequest{
			"GET",
			"/beneficiaries/{{beneficiary_id}}",
			JsonObject{"beneficiary_id": ""},
		},
		// Payment card
		actionCreatePaymentCard: mangoRequest{
			"POST",
			"/cards",
			nil,
		},
		actionFetchPaymentCard: mangoRequest{
			"GET",
			"/cards/{{paymentcard_id}}",
			JsonObject{"paymentcard_id": ""},
		},
		actionFetchUserPaymentCards: mangoRequest{
			"GET",
			"/users/{{user_id}}/cards",
			JsonObject{"user_id": ""},
		},
		actionDeletePaymentCard: mangoRequest{
			"DELETE",
			"/cards/{{paymentcard_id}}",
			JsonObject{"paymentcard_id": ""},
		},
		// Wallet
		actionCreateWallet: mangoRequest{
			"POST",
			"/wallets",
			nil,
		},
		actionFetchWallet: mangoRequest{
			"GET",
			"/wallets/{{wallet_id}}",
			JsonObject{"wallet_id": ""},
		},
		actionUpdateWallet: mangoRequest{
			"PUT",
			"/wallets/{{wallet_id}}",
			JsonObject{"wallet_id": ""},
		},
		actionFetchUsersOfWallet: mangoRequest{
			"GET",
			"/wallets/{{wallet_id}}/users",
			JsonObject{"wallet_id": ""},
		},
	*/
}
