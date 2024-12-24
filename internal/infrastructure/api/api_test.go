package api

import (
	"encoding/json"
	"errors"
	"go-key-value-cqrs/infrastructure/api/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAddKeyValueRequestReturnsDecodeError(t *testing.T) {
	// given
	invalidRequest := http.Request{Body: http.NoBody}

	// when
	addKeyRequest, keyPresent, err := validateAddKeyValueRequest(&invalidRequest)

	// then
	require.Nil(t, addKeyRequest)
	require.Equal(t, "", keyPresent)
	require.NotNil(t, err)
	expectedError := errors.New("EOF")
	require.Equal(t, expectedError, err)
}

func TestHandleErrorUnknown(t *testing.T) {
	// given
	expectedErrorMessage := "unknown error"
	err := errors.New(expectedErrorMessage)
	writer := httptest.NewRecorder()

	// when
	handleError(writer, err)

	// then
	actualResponse := writer.Result()
	defer actualResponse.Body.Close()

	data, _ := io.ReadAll(actualResponse.Body)
	expectedJsonBytes, _ := json.Marshal(model.ErrorResponse{Message: expectedErrorMessage})
	require.Contains(t, string(data), string(expectedJsonBytes))
}
