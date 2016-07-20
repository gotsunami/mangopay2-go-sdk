package mango

import (
	"testing"
	"time"
)

func TestNaturalUserSave(test *testing.T) {
	serv, _ := newTestService()
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
	user := serv.NewNaturalUser("Sergey", "Yarmonov", "sergey.yarmonov@gmail.com",
		time.Date(1988, time.January, 18, 0, 0, 0, 0, time.UTC).Unix(), "DE", "DE")
	user.IncomeRange = 3
	return user
}
