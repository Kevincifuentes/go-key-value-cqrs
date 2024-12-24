package persistence

import (
	"errors"
	"fmt"
	"github.com/jaswdr/faker/v2"
	"go-key-value-cqrs/domain"
	"go-key-value-cqrs/objectmothers"
	"testing"
)

const defaultTestValue = "anyValue"

var repositoryReader domain.KeyValueReader
var repositoryWriter domain.KeyValueWriter
var keyValueObjectMother objectmothers.KeyValueObjectMother
var fakerInstance faker.Faker

func TestMain(testing *testing.M) {
	fakerInstance = faker.New()
	inMemoryKeyValueRepository := NewInMemoryKeyValueRepository()
	repositoryReader = inMemoryKeyValueRepository
	repositoryWriter = inMemoryKeyValueRepository
	keyValueObjectMother = objectmothers.KeyValueObjectMother{FakerInstance: &fakerInstance}

	testing.Run()
}

func TestKeyValueFoundWithRepository(t *testing.T) {
	// given
	expectedKeyValue := keyValueObjectMother.CreateRandom()
	_ = repositoryWriter.Add(expectedKeyValue)

	// when
	keyValueView, err := repositoryReader.Get(expectedKeyValue.Key.Key)

	// then
	expectedValue := expectedKeyValue.Value.Value
	if err != nil || keyValueView.Value != expectedValue {
		t.Logf(
			"expected no error and expected value %s: actual value=%v, actualError=%v",
			expectedValue, keyValueView, err)
		t.Fail()
	}
}

func TestKeyValueNotFoundWithRepository(t *testing.T) {
	//given
	unknownKey := fakerInstance.Person().Name()

	// when
	keyValueView, err := repositoryReader.Get(unknownKey)

	// then
	if err == nil || keyValueView != (domain.KeyValueView{}) {
		t.Logf(
			"expected no error and expected value %s: actual value=%v, actualError=%v",
			defaultTestValue, keyValueView, err)
		t.Fail()
	}
}

func TestAddNewKeySuccessfully(t *testing.T) {
	// given
	expectedKeyValue := keyValueObjectMother.CreateRandom()
	expectedValue := expectedKeyValue.Value.Value
	expectedKey := expectedKeyValue.Key.Key

	// when
	err := repositoryWriter.Add(expectedKeyValue)

	// then
	keyValueView, err := repositoryReader.Get(expectedKey)
	if err != nil || keyValueView.Value != expectedValue || keyValueView.Key != expectedKey {
		t.Logf(
			"expected no error and receive (expectedKey=%v, expectedValue=%v): "+
				"actualError=%v, actualKey=%v, actualValue=%v",
			expectedKey, expectedValue, err, keyValueView.Key, keyValueView.Value)
		t.Fail()
	}
}

func TestAddNewKeyFailsOnExistingKey(t *testing.T) {
	// given
	existingKeyValue := keyValueObjectMother.CreateRandom()
	existingKey := existingKeyValue.Key.Key
	_ = repositoryWriter.Add(existingKeyValue)
	// and
	duplicatedKeyKeyValue := keyValueObjectMother.WithKey(existingKey)

	// when
	err := repositoryWriter.Add(duplicatedKeyKeyValue)

	// then
	var keyExistsError *domain.KeyExistsError
	isKeyExistsError := errors.As(err, &keyExistsError)
	expectedStringError := fmt.Sprintf("Key '%v' already exists", existingKey)
	if err == nil || !isKeyExistsError || keyExistsError.Error() != expectedStringError {
		t.Logf(
			"expected error and KeyExistsError type with message (expectedMessage=%v): actualError=%v",
			expectedStringError, err)
		t.Fail()
	}
}

func TestDeleteKeySuccessfully(t *testing.T) {
	// given
	expectedKeyValue := keyValueObjectMother.CreateRandom()
	expectedKey := expectedKeyValue.Key.Key
	_ = repositoryWriter.Add(expectedKeyValue)

	// when
	deleteErr := repositoryWriter.Delete(expectedKey)

	// then
	if deleteErr != nil {
		t.Logf(
			"expected no error on DELETE: actualError=%v", deleteErr)
		t.Fail()
	}
	//and
	keyValueView, getErr := repositoryReader.Get(expectedKey)
	var keyNotFoundError *domain.KeyNotFoundError
	isKeyNotFoundError := errors.As(getErr, &keyNotFoundError)
	if getErr == nil || keyValueView != (domain.KeyValueView{}) || !isKeyNotFoundError {
		t.Logf(
			"expected KeyNotFoundError and no KeyValueView: actualKeyValueView=%v, actualError=%v",
			keyValueView, getErr)
		t.Fail()
	}
}

func TestDeleteKeyNotFoundKey(t *testing.T) {
	// given
	unknownKey := fakerInstance.UUID().V4()

	// when
	deleteErr := repositoryWriter.Delete(unknownKey)

	// then
	var keyNotFoundError *domain.KeyNotFoundError
	isKeyNotFoundError := errors.As(deleteErr, &keyNotFoundError)
	if deleteErr == nil || !isKeyNotFoundError {
		t.Logf(
			"expected KeyNotFoundError on DELETE: actualError=%v", deleteErr)
		t.Fail()
	}
}
