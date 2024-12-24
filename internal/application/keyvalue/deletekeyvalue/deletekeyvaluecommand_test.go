package deletekeyvalue

import (
	"errors"
	"github.com/jaswdr/faker/v2"
	"go-key-value-cqrs/application/cqrs/commandbus"
	"go-key-value-cqrs/domain"
	"go-key-value-cqrs/objectmothers"
	"testing"
)

type testInterfaceMock struct {
	numberOfCalls int
	expectedKey   string
	expectedError error
}

func (mock *testInterfaceMock) Add(_ domain.KeyValue) error {
	return nil
}

func (mock *testInterfaceMock) Delete(key string) error {
	mock.numberOfCalls++
	mock.expectedKey = key
	return mock.expectedError
}

var mock testInterfaceMock
var fakerInstance = faker.New()
var keyValueObjectMother = objectmothers.KeyValueObjectMother{FakerInstance: &fakerInstance}

func registerCommand(expectedError error) {
	mock = testInterfaceMock{numberOfCalls: 0, expectedKey: "", expectedError: expectedError}
	handler := CommandHandler{&mock}
	commandbus.Load(handler)
}

func TestDeleteKeyValueCommandResolvedCorrectly(t *testing.T) {
	// given
	expectedKeyValue := keyValueObjectMother.CreateRandom()
	expectedKey := expectedKeyValue.Key.Key
	expectedErr := errors.New("any")
	registerCommand(expectedErr)
	command := Command{Key: expectedKey}

	// when
	err := commandbus.Execute(command)

	// then
	if err == nil || !errors.Is(err, expectedErr) {
		t.Errorf("Test failed! Expected to be successful call for %v command. Actual: Error=%v",
			command, err)
	}
	expectedNumberOfCalls := 1
	if mock.numberOfCalls != expectedNumberOfCalls {
		t.Errorf("Expected number of calls not meet on testInterfaceMock. Expected=%v Actual=%v",
			expectedNumberOfCalls, mock.numberOfCalls)
	}
}
