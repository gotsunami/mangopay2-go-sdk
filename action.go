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
	actionFetchUserTransfers
	actionFetchUserWallets
	actionFetchUserCards
	actionFetchUserBankAccounts

	actionCreateWallet
	actionEditWallet
	actionFetchWallet
	actionFetchWalletTransactions

	actionCreateTransfer
	actionFetchTransfer

	actionFetchPayIn
	actionCreateWebPayIn
	actionCreateDirectPayIn

	actionCreateCardRegistration
	actionSendCardRegistrationData

	actionFetchCard

	actionCreateTransferRefund
	actionCreatePayInRefund
	actionFetchRefund

	actionCreateBankAccount
	actionFetchBankAccount

	actionCreatePayOut
	actionFetchPayOut
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
	actionFetchUserTransfers: mangoRequest{
		"GET",
		"/users/{{Id}}/transactions?Type={{Type}}",
		JsonObject{"Id": "", "Type": ""},
	},
	actionFetchUserWallets: mangoRequest{
		"GET",
		"/users/{{Id}}/wallets",
		JsonObject{"Id": ""},
	},
	actionFetchUserCards: mangoRequest{
		"GET",
		"/users/{{Id}}/cards",
		JsonObject{"Id": ""},
	},
	actionCreateWallet: mangoRequest{
		"POST",
		"/wallets",
		nil,
	},
	actionEditWallet: mangoRequest{
		"PUT",
		"/wallets/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchWallet: mangoRequest{
		"GET",
		"/wallets/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchWalletTransactions: mangoRequest{
		"GET",
		"/wallets/{{Id}}/transactions",
		JsonObject{"Id": ""},
	},
	actionCreateTransfer: mangoRequest{
		"POST",
		"/transfers",
		nil,
	},
	actionFetchTransfer: mangoRequest{
		"GET",
		"/transfers/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchPayIn: mangoRequest{
		"GET",
		"/payins/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionCreateWebPayIn: mangoRequest{
		"POST",
		"/payins/card/web",
		nil,
	},
	actionCreateDirectPayIn: mangoRequest{
		"POST",
		"/payins/card/direct",
		nil,
	},
	actionCreateCardRegistration: mangoRequest{
		"POST",
		"/cardregistrations",
		nil,
	},
	actionSendCardRegistrationData: mangoRequest{
		"PUT",
		"/CardRegistrations/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchCard: mangoRequest{
		"GET",
		"/cards/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionCreateTransferRefund: mangoRequest{
		"POST",
		"/transfers/{{TransferId}}/refunds",
		JsonObject{"TransferId": ""},
	},
	actionCreatePayInRefund: mangoRequest{
		"POST",
		"/payins/{{PayInId}}/refunds",
		JsonObject{"PayInId": ""},
	},
	actionFetchRefund: mangoRequest{
		"GET",
		"/refunds/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionCreateBankAccount: mangoRequest{
		"POST",
		"/users/{{UserId}}/bankaccounts/{{Type}}",
		JsonObject{"UserId": "", "Type": ""},
	},
	actionFetchBankAccount: mangoRequest{
		"GET",
		"/users/{{UserId}}/bankaccounts/{{Id}}",
		JsonObject{"UserId": "", "Id": ""},
	},
	actionFetchUserBankAccounts: mangoRequest{
		"GET",
		"/users/{{Id}}/bankaccounts",
		JsonObject{"Id": ""},
	},
	actionCreatePayOut: mangoRequest{
		"POST",
		"/payouts/bankwire",
		nil,
	},
	actionFetchPayOut: mangoRequest{
		"GET",
		"/payouts/{{Id}}",
		JsonObject{"Id": ""},
	},
}
