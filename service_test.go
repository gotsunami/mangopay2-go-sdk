package mango

import (
	"os"
	"testing"
)

func TestNewMangoPay(test *testing.T) {
	_, err := newTestService()
	if err != nil {
		test.Fatal("Unable to create service:", err)
	}
}

func newTestService() (*MangoPay, error) {
	clientId := os.Getenv("MANGOPAY_CLIENT_ID")
	name := os.Getenv("MANGOPAY_NAME")
	email := os.Getenv("MANGOPAY_EMAIL")
	passwd := os.Getenv("MANGOPAY_PASSWD")
	conf, err := NewConfig(clientId, name, email, passwd, "sandbox")
	if err != nil {
		return nil, err
	}
	service, err := NewMangoPay(conf, OAuth)
	if err != nil {
		return nil, err
	}
	Verbosity(Debug)(service)
	return service, nil
}
