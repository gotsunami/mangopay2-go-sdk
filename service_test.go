package mango

import (
	"os"
	"testing"
)

func TestNewMangoPay(test *testing.T) {
	newTestService(test)
}

func newTestService(test *testing.T) *MangoPay {
	clientId := os.Getenv("MANGOPAY_CLIENT_ID")
	name := os.Getenv("MANGOPAY_NAME")
	email := os.Getenv("MANGOPAY_EMAIL")
	passwd := os.Getenv("MANGOPAY_PASSWD")
	conf, err := NewConfig(clientId, name, email, passwd, "sandbox")
	if err != nil {
		test.Fatal("Unable to create service:", err)
	}
	service, err := NewMangoPay(conf, OAuth)
	if err != nil {
		test.Fatal("Unable to create service:", err)
	}
	Verbosity(Debug)(service)
	return service
}
