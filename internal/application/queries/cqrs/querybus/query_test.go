package querybus

import (
	"errors"
	"github.com/labstack/gommon/log"
	"strings"
	"testing"
)

type testQuery struct {
	id string
}

func (query testQuery) Config() QueryConfig {
	return QueryConfig{Name: "TestQuery"}
}

type anotherTestQuery struct {
	id string
}

func (query anotherTestQuery) Config() QueryConfig {
	return QueryConfig{Name: "anotherTestQuery"}
}

type testInterface interface {
	Get(id string) (any, error)
}

type testInterfaceMock struct {
	numberOfCalls int
}

func (mock *testInterfaceMock) Get(id string) (any, error) {
	mock.numberOfCalls++
	return id, nil
}

type TestQueryHandler struct {
	repository testInterface
}

func (handler TestQueryHandler) Ask(query testQuery) (any, error) {
	return handler.repository.Get(query.id)
}

type TestAnotherQueryHandler struct {
}

func (handler TestAnotherQueryHandler) Ask(_ anotherTestQuery) (float64, error) {
	return 0, nil
}

func TestMain(m *testing.M) {
	registerQueries()

	m.Run()
}

var mock *testInterfaceMock

func registerQueries() {
	mock = &testInterfaceMock{numberOfCalls: 0}
	handler := TestQueryHandler{mock}
	err := Load(handler)
	if err != nil {
		log.Warnf("Error loading query handler %v. Error=%v", handler, err)
	}
}

func TestQueryResolvedCorrectly(t *testing.T) {
	// given
	expectedId := "TestQueryResolvedCorrectly"
	query := testQuery{expectedId}

	// when
	response, err := Asks[any](query)

	// then
	if err != nil || response != expectedId {
		t.Errorf("Test failed! Expected to be successful call for %v query. Actual: Response=%v Error=%v",
			query, response, err)
	}
	expectedNumberOfCalls := 1
	if mock.numberOfCalls != expectedNumberOfCalls {
		t.Errorf("Expected number of calls not meet on testInterfaceMock. Expected=%v Actual=%v",
			expectedNumberOfCalls, mock.numberOfCalls)
	}
}

func TestQueryNotFound(t *testing.T) {
	// given
	expectedId := "TestQueryNotFound"
	query := anotherTestQuery{expectedId}

	// when
	response, err := Asks[any](query)

	// then
	var errorQueryNotFound *ErrorQueryNotFound
	isErrorQueryNotFound := errors.As(err, &errorQueryNotFound)
	if err == nil || !isErrorQueryNotFound || !strings.Contains(errorQueryNotFound.Error(), "No query found") {
		t.Errorf("Test failed! Expected error on call for %v query. Actual: Response=%v Error=%v",
			query, response, err)
	}
}

func TestHandlerTypeInvalid(t *testing.T) {
	// given
	err := Load(TestAnotherQueryHandler{})
	expectedId := "TestQueryNotFound"
	query := anotherTestQuery{expectedId}

	// when
	response, err := Asks[string](query)

	// then
	var errorTypeNotValid *ErrorQueryHandlerTypeNotValid
	isErrorTypeInvalid := errors.As(err, &errorTypeNotValid)
	if err == nil || !isErrorTypeInvalid || !strings.Contains(errorTypeNotValid.Error(), "incorrect type") {
		t.Errorf("Test failed! Expected error on call for %v query. Actual: Response=%v Error=%v",
			query, response, err)
	}
}
