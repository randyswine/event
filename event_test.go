package event

import (
	"testing"
)

func TestFireEventWithNilCallback(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			t.Errorf("FireEvent() no handle nil callback")
			t.Fail()
		}
	}()
	wantEventName := "test_fire_event_with_nil_callback"
	On(wantEventName, nil)
	FireEvent(NewEvent(wantEventName, nil))
}

func TestFireEvent(t *testing.T) {
	wantEventName := "test"
	var gotEventName string
	On(wantEventName, func(e Event) error {
		gotEventName = e.Name()
		return nil
	})
	FireEvent(NewEvent("test", nil))
	if gotEventName != wantEventName {
		t.Errorf("FireEvent() got event name \"%s\"; want \"%s\"", gotEventName, wantEventName)
		t.Fail()
	}
}

func TestSubscibe(t *testing.T) {
	wantEventName := "test_listener"
	listener := &testListener{}
	Subscribe(wantEventName, listener)
	FireEvent(NewEvent("test_listener", nil))
	if listener.takedEventName != wantEventName {
		t.Errorf("FireEvent() got event name \"%s\"; want \"%s\"", listener.takedEventName, wantEventName)
		t.Fail()
	}
}

func TestUnsub(t *testing.T) {
	wantEventName := "test_unsub"
	testData := "foo"
	beforeUnsubGot := "bar"
	afterUnsubGot := "fizzbazz"
	handler := On(wantEventName, func(e Event) error {
		testData = afterUnsubGot
		return nil
	})
	On(wantEventName, func(e Event) error {
		testData = beforeUnsubGot
		return nil
	})
	handler.Unsub()
	FireEvent(NewEvent(wantEventName, nil))
	if testData != beforeUnsubGot {
		t.Errorf("FireEvent() got event name \"%s\"; want \"%s\"", afterUnsubGot, beforeUnsubGot)
		t.Fail()
	}
}

func TestRepeatedUnsub(t *testing.T) {
	wantEventName := "test_unsub"
	handler := On(wantEventName, func(e Event) error {
		return nil
	})
	handler.Unsub()
	handler.Unsub()
}

type testListener struct {
	takedEventName string
}

func (tl *testListener) Handle(e Event) error {
	tl.takedEventName = e.Name()
	return nil
}
