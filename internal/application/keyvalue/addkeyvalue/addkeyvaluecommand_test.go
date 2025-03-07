package addkeyvalue

import (
	"errors"
	"github.com/jaswdr/faker/v2"
	"go-key-value-cqrs/application/cqrs/commandbus"
	"go-key-value-cqrs/domain"
	"go-key-value-cqrs/objectmothers"
	"testing"
)

type testInterfaceMock struct {
	numberOfCalls    int
	expectedKeyValue *domain.KeyValue
	expectedError    error
}

func (mock *testInterfaceMock) Add(keyValue domain.KeyValue) error {
	mock.numberOfCalls++
	mock.expectedKeyValue = &keyValue
	return mock.expectedError
}

func (mock *testInterfaceMock) Delete(_ string) error {
	return nil
}

const maxLengthKey = 200

var mock testInterfaceMock
var fakerInstance = faker.New()
var keyValueObjectMother = objectmothers.KeyValueObjectMother{FakerInstance: &fakerInstance}

func registerCommand(expectedValue *domain.KeyValue, expectedError error) {
	mock = testInterfaceMock{numberOfCalls: 0, expectedKeyValue: expectedValue, expectedError: expectedError}
	handler := CommandHandler{&mock}
	commandbus.Load(handler)
}

func TestAddKeyValueCommandResolvedCorrectly(t *testing.T) {
	// given
	expectedKeyValue := keyValueObjectMother.CreateRandom()
	expectedKey := expectedKeyValue.Key.Key
	expectedValue := expectedKeyValue.Value.Value
	registerCommand(&expectedKeyValue, nil)
	command := Command{Key: expectedKey, Value: expectedValue}

	// when
	err := commandbus.Execute(command)

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
}

func TestAddKeyValueCommandReturnsDomainError(t *testing.T) {
	// given
	tooLongKey := fakerInstance.RandomStringWithLength(maxLengthKey + 1)
	expectedKeyValue := keyValueObjectMother.WithKey(tooLongKey)
	expectedValue := expectedKeyValue.Value.Value
	registerCommand(&expectedKeyValue, nil)
	command := Command{Key: tooLongKey, Value: expectedValue}

	// when
	err := commandbus.Execute(command)

	// then
	var invalidLengthError *domain.KeyValueDomainError
	isInvalidLengthError := errors.As(err, &invalidLengthError)
	if err == nil || !isInvalidLengthError {
		t.Errorf("Test failed! Expected to fail call for %v command. Expected: Error=InvalidLengthError, Actual: Error=%v",
			command, err)
	}
	expectedNumberOfCalls := 0
	if mock.numberOfCalls != expectedNumberOfCalls {
		t.Errorf("Expected number of calls not meet on testInterfaceMock. Expected=%v Actual=%v",
			expectedNumberOfCalls, mock.numberOfCalls)
	}
}
