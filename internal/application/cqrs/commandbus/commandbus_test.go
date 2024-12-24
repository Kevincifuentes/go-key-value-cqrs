package commandbus

import (
	"errors"
	"strings"
	"testing"
)

type testCommand struct {
	id string
}

func (command testCommand) Config() CommandConfig {
	return CommandConfig{Name: "TestCommand"}
}

type anotherTestCommand struct {
	id string
}

func (command anotherTestCommand) Config() CommandConfig {
	return CommandConfig{Name: "anotherTestCommand"}
}

type anotherTestWithSameNameCommand struct {
	id string
}

func (command anotherTestWithSameNameCommand) Config() CommandConfig {
	return CommandConfig{Name: "TestCommand"}
}

type testInterface interface {
	increment(id string) error
}

type testInterfaceMock struct {
	numberOfCalls int
	expectedId    string
}

func (mock *testInterfaceMock) increment(id string) error {
	mock.numberOfCalls++
	mock.expectedId = id
	return nil
}

type TestCommandHandler struct {
	repository testInterface
}

func (handler TestCommandHandler) Execute(command testCommand) error {
	return handler.repository.increment(command.id)
}

type TestAnotherCommandHandler struct {
}

func (handler TestAnotherCommandHandler) Execute(_ anotherTestCommand) error {
	return nil
}

func TestMain(m *testing.M) {
	registerQueries()

	m.Run()
}

var mock *testInterfaceMock

func registerQueries() {
	mock = &testInterfaceMock{numberOfCalls: 0}
	handler := TestCommandHandler{mock}
	Load(handler)
}

func TestCommandResolvedCorrectly(t *testing.T) {
	// given
	expectedId := "TestCommandResolvedCorrectly"
	command := testCommand{expectedId}

	// when
	err := Execute(command)

	// then
	if err != nil {
		t.Errorf("Test failed! Expected to be successful call for %v command. Actual: Error=%v",
			command, err)
	}
	expectedNumberOfCalls := 1
	if mock.numberOfCalls != expectedNumberOfCalls {
		t.Errorf("Expected number of calls not meet on testInterfaceMock. Expected=%v Actual=%v",
			expectedNumberOfCalls, mock.numberOfCalls)
	}
	if mock.expectedId != expectedId {
		t.Errorf("Expected id not meet on testInterfaceMock. Expected=%v Actual=%v", expectedId, mock.expectedId)
	}
}

func TestCommandNotFound(t *testing.T) {
	// given
	expectedId := "TestCommandNotFound"
	command := anotherTestCommand{expectedId}

	// when
	err := Execute(command)

	// then
	var errorCommandNotFound *ErrorCommandNotFound
	isErrorCommandNotFound := errors.As(err, &errorCommandNotFound)
	if err == nil || !isErrorCommandNotFound || !strings.Contains(errorCommandNotFound.Error(), "No command found") {
		t.Errorf("Test failed! Expected error on call for %v command. Actual: Error=%v",
			command, err)
	}
}

func TestHandlerTypeInvalid(t *testing.T) {
	// given
	Load(TestAnotherCommandHandler{})
	expectedId := "TestCommandNotFound"
	command := anotherTestWithSameNameCommand{expectedId}

	// when
	err := Execute(command)

	// then
	var errorTypeNotValid *ErrorCommandHandlerTypeNotValid
	isErrorTypeInvalid := errors.As(err, &errorTypeNotValid)
	if err == nil || !isErrorTypeInvalid || !strings.Contains(errorTypeNotValid.Error(), "incorrect type") {
		t.Errorf("Test failed! Expected error on call for %v command. Actual: Error=%v",
			command, err)
	}
}
