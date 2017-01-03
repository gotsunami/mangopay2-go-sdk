package mango

import (
	"encoding/json"
	"fmt"
)

func (m *MangoPay) NewHook(eventType EventType, url string) (*Hook, error) {
	hook := &Hook{
		EventType: eventType,
		Url:       url,
		service:   m,
	}
	return hook, nil
}

type Hook struct {
	ProcessIdent
	Url       string
	EventType EventType
	Status    string
	Validity  string

	service *MangoPay
}

func (h *Hook) String() string {
	return struct2string(h)
}

func (h *Hook) Save() error {
	var action mangoAction
	if h.Id == "" {
		action = actionCreateHook
	} else {
		action = actionUpdateHook
	}

	data := JsonObject{}
	j, err := json.Marshal(h)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	// Fields not allowed when creating an object.
	if action == actionCreateWallet {
		delete(data, "Id")
	}
	delete(data, "Status")
	delete(data, "Validity")

	if action == actionUpdateHook {
		// Delete empty values so that existing ones don't get
		// overwritten with empty values.
		for k, v := range data {
			switch v.(type) {
			case string:
				if v.(string) == "" {
					delete(data, k)
				}
			case int:
				if v.(int) == 0 {
					delete(data, k)
				}
			}
		}
	}

	hook, err := h.service.anyRequest(new(Hook), action, data)
	if err != nil {
		return err
	}
	serv := h.service
	*h = *(hook.(*Hook))
	h.service = serv
	return nil
}

func (m *MangoPay) Hook(id string) (*Hook, error) {
	h, err := m.anyRequest(new(Hook), actionFetchHook, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	hook := h.(*Hook)
	hook.service = m
	return hook, nil
}

// HookByEventType returns the only hook for given event type.
// see https://docs.mangopay.com/endpoints/v2.01/hooks#e247_create-a-hook
func (m *MangoPay) HookByEventType(eventType EventType) (*Hook, error) {
	hooks, err := m.Hooks()
	if err != nil {
		return nil, err
	}

	var result *Hook
	for _, hook := range hooks {
		if hook.EventType == eventType {
			result = hook
			break
		}
	}

	if result == nil {
		return nil, fmt.Errorf("No hook found for event type %s", eventType)
	}

	return result, nil
}

func (m *MangoPay) Hooks() (HookList, error) {
	list, err := m.anyRequest(new(HookList), actionFetchAllHooks, nil)
	if err != nil {
		return nil, err
	}
	casted := *(list.(*HookList))
	for _, hook := range casted {
		hook.service = m
	}
	return casted, nil
}

type HookList []*Hook
