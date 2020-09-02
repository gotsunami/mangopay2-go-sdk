package mango

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestHookByEventType(test *testing.T) {
	service := newTestService(test)

	hook, err := service.HookByEventType(EventDisputeClosed)
	if err != nil {
		test.Fatal("Unable to get hook for event type:", err)
	}

	if hook.EventType != EventDisputeClosed {
		test.Fatal("Mismatched event type:", hook.EventType)
	}
}

func TestHook_Save(test *testing.T) {
	service := newTestService(test)
	test.Log("Hook creating...")
	hook, err := service.NewHook(EventDisputeClosed, "http://sdfkjkdjf.com/hook/fdjsf")
	if err != nil {
		test.Fatal("Unable to create hook:", err)
	}

	if err = hook.Save(); err != nil {
		test.Log("Unable to store hook:", err)

		hook, err = service.HookByEventType(EventDisputeClosed)
		if err != nil {
			test.Fatal("Unable to get hook: ", err)
		}
	}

	test.Log("Hook updating...")
	rand.Seed(time.Now().UTC().UnixNano())
	hook.Url = "http://sdfkjkdjf.com/hook/" + strconv.Itoa(rand.Int())
	if err = hook.Save(); err != nil {
		test.Fatal("Unable to update hook", err)
	}
}
