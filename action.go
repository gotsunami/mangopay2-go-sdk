// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

type MangoAction int

const (
	ActionEvents MangoAction = iota
	ActionCreateNaturalUser

	/*
		ActionCreateUser MangoAction = iota
		ActionFetchUser
		ActionUpdateUser
		ActionFetchUserWallets

		ActionCreateBeneficiary
		ActionFetchBeneficiary

		ActionCreatePaymentCard
		ActionFetchPaymentCard
		ActionFetchUserPaymentCards
		ActionDeletePaymentCard

		ActionCreateWallet
		ActionFetchWallet
		ActionUpdateWallet
		ActionFetchUsersOfWallet
	*/
)

// JsonObject is used to manage JSON data.
type JsonObject map[string]interface{}

type mangoRequest struct {
	Method, Path string
	PathValues   JsonObject
}

// Defines mango requests metadata.
var mangoRequests = map[MangoAction]mangoRequest{
	ActionEvents: mangoRequest{
		"GET",
		"/events",
		nil,
	},
	ActionCreateNaturalUser: mangoRequest{
		"POST",
		"/users/natural",
		nil,
	},
	/*
		// User
		ActionFetchUser: mangoRequest{
			"GET",
			"/users/{{user_id}}",
			JsonObject{"user_id": ""},
		},
		ActionUpdateUser: mangoRequest{
			"PUT",
			"/users/{{user_id}}",
			JsonObject{"user_id": ""},
		},
		ActionFetchUserWallets: mangoRequest{
			"GET",
			"/users/{{user_id}}/wallets",
			JsonObject{"user_id": ""},
		},
		// Beneficiary
		ActionCreateBeneficiary: mangoRequest{
			"POST",
			"/beneficiaries",
			nil,
		},
		ActionFetchBeneficiary: mangoRequest{
			"GET",
			"/beneficiaries/{{beneficiary_id}}",
			JsonObject{"beneficiary_id": ""},
		},
		// Payment card
		ActionCreatePaymentCard: mangoRequest{
			"POST",
			"/cards",
			nil,
		},
		ActionFetchPaymentCard: mangoRequest{
			"GET",
			"/cards/{{paymentcard_id}}",
			JsonObject{"paymentcard_id": ""},
		},
		ActionFetchUserPaymentCards: mangoRequest{
			"GET",
			"/users/{{user_id}}/cards",
			JsonObject{"user_id": ""},
		},
		ActionDeletePaymentCard: mangoRequest{
			"DELETE",
			"/cards/{{paymentcard_id}}",
			JsonObject{"paymentcard_id": ""},
		},
		// Wallet
		ActionCreateWallet: mangoRequest{
			"POST",
			"/wallets",
			nil,
		},
		ActionFetchWallet: mangoRequest{
			"GET",
			"/wallets/{{wallet_id}}",
			JsonObject{"wallet_id": ""},
		},
		ActionUpdateWallet: mangoRequest{
			"PUT",
			"/wallets/{{wallet_id}}",
			JsonObject{"wallet_id": ""},
		},
		ActionFetchUsersOfWallet: mangoRequest{
			"GET",
			"/wallets/{{wallet_id}}/users",
			JsonObject{"wallet_id": ""},
		},
	*/
}
