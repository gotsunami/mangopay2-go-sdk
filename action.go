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
	actionCreateBankwireDirectPayIn
	actionCreateDirectDebitWebPayIn
	actionCreateDirectDebitDirectPayIn

	actionCreateCardRegistration
	actionSendCardRegistrationData

	actionFetchCard

	actionFetchMandate

	actionCreateTransferRefund
	actionCreatePayInRefund
	actionFetchRefund

	actionCreateBankAccount
	actionFetchBankAccount

	actionCreatePayOut
	actionFetchPayOut

	actionCreateKYCDocument
	actionFetchKYCDocument
	actionSubmitKYCDocument
	actionCreateKYCPage
	actionFetchUserKYCDocuments
	actionFetchAllKYCDocuments

	actionCreateHook
	actionUpdateHook
	actionFetchHook
	actionFetchAllHooks
)

// JsonObject is used to manage JSON data.
type JsonObject map[string]interface{}

type mangoRequest struct {
	Method, Path string
	PathValues   JsonObject
}

// Defines mango requests metadata.
var mangoRequests = map[mangoAction]mangoRequest{
	actionEvents: {
		"GET",
		"/events",
		nil,
	},
	actionCreateNaturalUser: {
		"POST",
		"/users/natural",
		nil,
	},
	actionEditNaturalUser: {
		"PUT",
		"/users/natural/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionAllUsers: {
		"GET",
		"/users",
		nil,
	},
	actionFetchNaturalUser: {
		"GET",
		"/users/natural/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionCreateLegalUser: {
		"POST",
		"/users/legal",
		nil,
	},
	actionEditLegalUser: {
		"PUT",
		"/users/legal/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchLegalUser: {
		"GET",
		"/users/legal/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchUser: {
		"GET",
		"/users/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchUserTransfers: {
		"GET",
		"/users/{{Id}}/transactions",
		JsonObject{"Id": ""},
	},
	actionFetchUserWallets: {
		"GET",
		"/users/{{Id}}/wallets",
		JsonObject{"Id": ""},
	},
	actionFetchUserCards: {
		"GET",
		"/users/{{Id}}/cards",
		JsonObject{"Id": ""},
	},
	actionCreateWallet: {
		"POST",
		"/wallets",
		nil,
	},
	actionEditWallet: {
		"PUT",
		"/wallets/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchWallet: {
		"GET",
		"/wallets/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchWalletTransactions: {
		"GET",
		"/wallets/{{Id}}/transactions",
		JsonObject{"Id": ""},
	},
	actionCreateTransfer: {
		"POST",
		"/transfers",
		nil,
	},
	actionFetchTransfer: {
		"GET",
		"/transfers/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchPayIn: {
		"GET",
		"/payins/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionCreateWebPayIn: {
		"POST",
		"/payins/card/web",
		nil,
	},
	actionCreateDirectPayIn: {
		"POST",
		"/payins/card/direct",
		nil,
	},
	actionCreateBankwireDirectPayIn: {
		"POST",
		"/payins/bankwire/direct",
		nil,
	},
	actionCreateDirectDebitWebPayIn: {
		"POST",
		"/payins/directdebit/web",
		nil,
	},
	actionCreateDirectDebitDirectPayIn: {
		"POST",
		"/payins/directdebit/direct",
		nil,
	},
	actionCreateCardRegistration: {
		"POST",
		"/cardregistrations",
		nil,
	},
	actionSendCardRegistrationData: {
		"PUT",
		"/CardRegistrations/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchCard: {
		"GET",
		"/cards/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchMandate: {
		"GET",
		"/mandates/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionCreateTransferRefund: {
		"POST",
		"/transfers/{{TransferId}}/refunds",
		JsonObject{"TransferId": ""},
	},
	actionCreatePayInRefund: {
		"POST",
		"/payins/{{PayInId}}/refunds",
		JsonObject{"PayInId": ""},
	},
	actionFetchRefund: {
		"GET",
		"/refunds/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionCreateBankAccount: {
		"POST",
		"/users/{{UserId}}/bankaccounts/{{Type}}",
		JsonObject{"UserId": "", "Type": ""},
	},
	actionFetchBankAccount: {
		"GET",
		"/users/{{UserId}}/bankaccounts/{{Id}}",
		JsonObject{"UserId": "", "Id": ""},
	},
	actionFetchUserBankAccounts: {
		"GET",
		"/users/{{Id}}/bankaccounts",
		JsonObject{"Id": ""},
	},
	actionCreatePayOut: {
		"POST",
		"/payouts/bankwire",
		nil,
	},
	actionFetchPayOut: {
		"GET",
		"/payouts/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionCreateKYCDocument: {
		"POST",
		"/users/{{UserId}}/kyc/documents",
		JsonObject{"UserId": ""},
	},
	actionFetchKYCDocument: {
		"GET",
		"/kyc/documents/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionSubmitKYCDocument: {
		"PUT",
		"/users/{{UserId}}/kyc/documents/{{Id}}",
		JsonObject{"UserId": "", "Id": ""},
	},
	actionCreateKYCPage: {
		"POST",
		"/users/{{UserId}}/kyc/documents/{{Id}}/pages",
		JsonObject{"UserId": "", "Id": ""},
	},
	actionFetchUserKYCDocuments: {
		"GET",
		"/users/{{UserId}}/kyc/documents?status={{Status}}",
		JsonObject{"UserId": "", "Status": ""},
	},
	actionFetchAllKYCDocuments: {
		"GET",
		"/kyc/documents",
		nil,
	},

	actionCreateHook: {
		"POST",
		"/hooks",
		nil,
	},
	actionUpdateHook: {
		"PUT",
		"/hooks/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchHook: {
		"GET",
		"/hooks/{{Id}}",
		JsonObject{"Id": ""},
	},
	actionFetchAllHooks: {
		"GET",
		"/hooks/",
		nil,
	},
}
