package event

import (
	"testing"
)

func TestDispatcher(t *testing.T) {
	wantEventName := "firstInitalization"
	firstDispatcherInitialization := Dispatcher()
	var gotEventName string
	firstDispatcherInitialization.On(wantEventName, func(e Event) error {
		gotEventName = e.Name()
		return nil
	})
	secondDispatcherInitalization := Dispatcher()
	secondDispatcherInitalization.FireEvent(NewEvent("firstInitalization", nil))
	if gotEventName != wantEventName {
		t.Errorf("FireEvent() got event name \"%s\"; want \"%s\"", gotEventName, wantEventName)
		t.Fail()
	}
}

func TestFireEvent(t *testing.T) {
	wantEventName := "test"
	dispatcher := Dispatcher()
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
	dispatcher := Dispatcher()
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
