package mango

import (
	"testing"
)

func TestWalletSave(test *testing.T) {
	serv := newTestService(test)
	createTestWallet(test, serv, nil)
}

func createTestWallet(test *testing.T, serv *MangoPay, user *NaturalUser) *Wallet {
	if user == nil {
		user = createTestUser(serv)
		if err := user.Save(); err != nil {
			test.Fatal("Unable to store user", err)
		}
	}
	test.Log("Storing wallet...")
	wallet, err := serv.NewWallet(ConsumerList{user}, "Test wallet", "EUR")
	if err != nil {
		test.Fatal("Unable to create wallet:", err)
	}
	if err := wallet.Save(); err != nil {
		test.Fatal("Unable to store wallet:", err)
	}
	return wallet
}
