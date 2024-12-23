package getvalue

import (
	"errors"
	"github.com/jaswdr/faker/v2"
	"github.com/labstack/gommon/log"
	"go-key-value-cqrs/application/queries/cqrs/querybus"
	"go-key-value-cqrs/domain"
	"testing"
)

type testInterfaceMock struct {
	numberOfCalls int
	expectedValue domain.KeyValueView
	expectedError error
}

func (mock *testInterfaceMock) Get(key string) (domain.KeyValueView, error) {
	mock.numberOfCalls++
	if key != mock.expectedValue.Key {
		log.Errorf("Expected key %s but got %s", mock.expectedValue.Key, key)
	}
	return mock.expectedValue, mock.expectedError
}

var mock testInterfaceMock
var fakerInstance = faker.New()

func registerQuery(expectedValue domain.KeyValueView, expectedError error) {
	mock = testInterfaceMock{numberOfCalls: 0, expectedValue: expectedValue, expectedError: expectedError}
	handler := QueryHandler{&mock}
	err := querybus.Load(handler)
	if err != nil {
		log.Warnf("Error loading query handler %v. Error=%v", handler, err)
	}
}

func TestGetValueQueryResolvedCorrectly(t *testing.T) {
	// given
	expectedKey := fakerInstance.Person().Name()
	expectedValue := fakerInstance.UUID().V4()
	expectedKeyValueView := domain.KeyValueView{Key: expectedKey, Value: expectedValue}
	registerQuery(expectedKeyValueView, nil)
	query := Query{expectedKey}

	// when
	response, err := querybus.Asks[domain.KeyValueView](query)

	// then
	if err != nil || response != expectedKeyValueView {
		t.Errorf("Test failed! Expected to be successful call for %v query. Actual: Response=%v Error=%v",
			query, response, err)
	}
	expectedNumberOfCalls := 1
	if mock.numberOfCalls != expectedNumberOfCalls {
		t.Errorf("Expected number of calls not meet on testInterfaceMock. Expected=%v Actual=%v",
			expectedNumberOfCalls, mock.numberOfCalls)
	}
}

func TestGetValueQueryReturnsError(t *testing.T) {
	// given
	expectedError := errors.New("test error")
	registerQuery(domain.KeyValueView{}, expectedError)
	query := Query{fakerInstance.UUID().V4()}

	// when
	response, err := querybus.Asks[domain.KeyValueView](query)

	// then
	if err == nil || !errors.Is(err, expectedError) {
		t.Errorf("Test failed! Expected to fail call for %v query. "+
			"Expected: Error=%v, Actual: Response=%v Error=%v",
			query, expectedError, response, err)
	}
	expectedNumberOfCalls := 1
	if mock.numberOfCalls != expectedNumberOfCalls {
		t.Errorf("Expected number of calls not meet on testInterfaceMock. Expected=%v Actual=%v",
			expectedNumberOfCalls, mock.numberOfCalls)
	}
}
