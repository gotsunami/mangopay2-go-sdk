package mango

import (
	"testing"
	"time"
)

func TestNaturalUserSave(test *testing.T) {
	serv := newTestService(test)
	user := createTestUser(serv)
	test.Log("Storing user...")
	if err := user.Save(); err != nil {
		test.Fatal("Unable to save user:", err)
	}
	test.Log("Fetching user...", user.Id)
	if _, err := serv.User(user.Id); err != nil {
		test.Fatal("Unable to fetch user:", err)
	}
}

func createTestUser(serv *MangoPay) *NaturalUser {
	user := serv.NewNaturalUser("Firstname", "Lastname", "example@example.com",
		time.Date(2000, time.January, 01, 0, 0, 0, 0, time.UTC).Unix(), "EN", "EN")
	user.IncomeRange = 3
	return user
}
