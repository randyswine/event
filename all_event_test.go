package event

import (
	"testing"
)

func TestFireEvent(t *testing.T) {
	wantEventName := "test"
	dispatcher := New()
	var gotEventName string
	dispatcher.On(wantEventName, func(e Event) error {
		gotEventName = e.Name()
		return nil
	})
	dispatcher.FireEvent(NewEvent("test", nil))
	if gotEventName != wantEventName {
		t.Errorf("FireEvent() got event name \"%s\"; want \"%s\"", gotEventName, wantEventName)
		t.Fail()
	}
}

func TestSubscibe(t *testing.T) {
	wantEventName := "test_listener"
	dispatcher := New()
	listener := &testListener{}
	dispatcher.Subscribe(wantEventName, listener)
	dispatcher.FireEvent(NewEvent("test_listener", nil))
	if listener.takedEventName != wantEventName {
		t.Errorf("FireEvent() got event name \"%s\"; want \"%s\"", listener.takedEventName, wantEventName)
		t.Fail()
	}
}

type testListener struct {
	takedEventName string
}

func (tl *testListener) Handle(e Event) error {
	tl.takedEventName = e.Name()
	return nil
}
