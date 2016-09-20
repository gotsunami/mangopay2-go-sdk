package mango

import (
	"testing"
)

func TestBankwireDirectPayIn_Save(test *testing.T) {
	serv, _ := newTestService()
	user := createTestUser(serv)
	if err := user.Save(); err != nil {
		test.Fatal("Unable to store user", err)
	}
	wallet := createTestWallet(test, serv, user)

	test.Log("Creating pay-in...")
	amount := Money{
		Currency: "EUR",
		Amount:   10000,
	}
	fees := Money{
		Currency: "EUR",
		Amount:   0,
	}
	payIn, err := serv.NewBankwireDirectPayIn(user, wallet, amount, fees)
	if err != nil {
		test.Fatal("Unable to create pay-in:", err)
	}
	if err = payIn.Save(); err != nil {
		test.Fatal("Unable to store pay-in:", err)
	}
}

func TestDirectDebitWebPayIn_Save(test *testing.T) {
	serv, _ := newTestService()
	user := createTestUser(serv)
	if err := user.Save(); err != nil {
		test.Fatal("Unable to store user", err)
	}
	wallet := createTestWallet(test, serv, user)

	test.Log("Creating pay-in...")
	amount := Money{
		Currency: "EUR",
		Amount:   10000,
	}
	fees := Money{
		Currency: "EUR",
		Amount:   0,
	}
	payIn, err := serv.NewDirectDebitWebPayIn(user, wallet, amount, fees, "https://google.com", DirectDebitTypeSofort, "DE")
	if err != nil {
		test.Fatal("Unable to create pay-in:", err)
	}
	if err = payIn.Save(); err != nil {
		test.Fatal("Unable to store pay-in:", err)
	}
}
