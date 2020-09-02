package mango

import (
	"strings"
	"testing"
)

var EUR0 = Money{"EUR", 0}
var EUR10 = Money{"EUR", 1000}
var EUR100 = Money{"EUR", 10000}

func TestPayout_Save(test *testing.T) {
	serv := newTestService(test)
	user := createTestUser(serv)
	if err := user.Save(); err != nil {
		test.Fatal("Unable to store user:", err)
	}
	wallet := createTestWallet(test, serv, user)
	account := createTestBankAccount(test, serv, user)
	createTestDirectDebitWebPayIn(test, serv, user, EUR100, EUR0, wallet)

	payout, err := serv.NewPayOut(user, EUR10, EUR0, wallet, account)
	if err != nil {
		test.Fatal("Unable to create PayOut", err)
	}

	if err := payout.Save(); err != nil {
		// TODO: register real bank account to test payout
		// see https://mangopay.desk.com/customer/en/portal/articles/2339775-can-i-test-the-payout-process-in-sandbox-
		// Now ignore 'Unsufficient wallet balance' error
		if !strings.Contains(err.Error(), "Unsufficient wallet balance") {
			test.Fatal("Unable to save PayOut", err)
		}
	}
}
