package mango

import (
	"testing"
)

func TestHook_Save(test *testing.T) {
	service, _ := newTestService()
	test.Log("Hook creatind...")
	hook, err := service.NewHook(EventDisputeClosed, "http://sdfkjkdjf.com/hook/fdjsf")
	if err != nil {
		test.Fatal("Unable to create hook:", err)
	}
	if err = hook.Save(); err != nil {
		test.Fatal("Unable to store hook:", err)
	}
	test.Log("Hook updating...")
	hook.Url = "http://sdfkjkdjf.com/hook/12345"
	if err = hook.Save(); err != nil {
		test.Fatal("Unable to update hook", err)
	}
}
