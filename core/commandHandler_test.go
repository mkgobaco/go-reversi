package core_test

import (
	"reversi/core"
	"testing"
)

type TestEventConsumer struct {
	events []core.Event
}

func (consumer *TestEventConsumer) SendEvent(event core.Event) {
	existingEventCount := len(consumer.events)
	nextEvents := make([]core.Event, existingEventCount+1)

	for i, event := range consumer.events {
		nextEvents[i] = event
	}
	nextEvents[existingEventCount] = event

	consumer.events = nextEvents
}

func NewTestEventConsumer() TestEventConsumer {
	return TestEventConsumer{events: []core.Event{}}
}

type TestCommandRejectHandler struct {
	rejectWasCalled bool
}

func (rejectHandler *TestCommandRejectHandler) InvalidCommand(command core.Command) {
	rejectHandler.rejectWasCalled = true
}
func NewTestCommandRejectHandler() TestCommandRejectHandler {
	return TestCommandRejectHandler{
		rejectWasCalled: false,
	}
}

func Test_InitializeCommand_triggersGameStateUpdate(t *testing.T) {
	testEventConsumer := NewTestEventConsumer()
	commandHandler := core.NewCommandHandler(&testEventConsumer, core.NewGameEventAggregator())

	initializeCommand := core.NewInitializeCommand()

	testRejectHandler := NewTestCommandRejectHandler()
	commandHandler.AttemptCommand(initializeCommand, &testRejectHandler)

	eventCount := len(testEventConsumer.events)
	if !(eventCount == 1) {
		t.Errorf("Expected 1 event, instead got %d", eventCount)
	}

	event := testEventConsumer.events[0]
	if !(event.EventType == core.INITILIZED) {
		t.Error("Should have been an INITIALIZED event")
	}
}

func Test_SecondInitializeCommand_isRejected(t *testing.T) {
	testEventConsumer := NewTestEventConsumer()

	commandHandler := core.NewCommandHandler(&testEventConsumer, core.NewGameEventAggregator())
	initializeCommand := core.NewInitializeCommand()

	testRejectHandler := NewTestCommandRejectHandler()
	commandHandler.AttemptCommand(initializeCommand, &testRejectHandler)
	if testRejectHandler.rejectWasCalled {
		t.Error("First Initialize should not have been rejected")
	}

	commandHandler.AttemptCommand(initializeCommand, &testRejectHandler)
	if !testRejectHandler.rejectWasCalled {
		t.Error("Second Initialize should have been rejected")
	}
}
