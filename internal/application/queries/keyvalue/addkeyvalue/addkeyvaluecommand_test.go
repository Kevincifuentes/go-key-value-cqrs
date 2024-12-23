package addkeyvalue

import (
	"errors"
	"github.com/jaswdr/faker/v2"
	"github.com/labstack/gommon/log"
	"go-key-value-cqrs/application/queries/cqrs/commandbus"
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

const maxLengthKey = 200

var mock testInterfaceMock
var fakerInstance = faker.New()
var keyValueObjectMother objectmothers.KeyValueObjectMother = objectmothers.KeyValueObjectMother{FakerInstance: &fakerInstance}

func registerCommand(expectedValue *domain.KeyValue, expectedError error) {
	mock = testInterfaceMock{numberOfCalls: 0, expectedKeyValue: expectedValue, expectedError: expectedError}
	handler := CommandHandler{&mock}
	err := commandbus.Load(handler)
	if err != nil {
		log.Warnf("Error loading command handler %v. Error=%v", handler, err)
	}
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
