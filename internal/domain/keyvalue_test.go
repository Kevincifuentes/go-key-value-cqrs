package domain

import (
	"errors"
	"fmt"
	"github.com/jaswdr/faker/v2"
	"strings"
	"testing"
)

var Faker = faker.New()

func TestKeyValueConstructorThrowKeyLengthError(t *testing.T) {
	// given
	validValue := Faker.RandomStringWithLength(MaxLengthValue)
	keyInvalidValues := []string{"", Faker.RandomStringWithLength(MaxLengthKey + 1)}

	// when
	for _, actualTestKey := range keyInvalidValues {
		t.Run(fmt.Sprintf("Testing [%v]", actualTestKey), func(t *testing.T) {
			// when
			keyValue, err := newKeyValue(actualTestKey, validValue)

			// then
			var invalidLengthError *InvalidLengthError
			isInvalidLengthError := errors.As(err, &invalidLengthError)
			if err == nil || !isInvalidLengthError || !strings.Contains(invalidLengthError.Error(), "'key'") {
				t.Logf("expected an InvalidLengthError: , actual keyValue: %v, actual error: %v", keyValue, err)
				t.Fail()
			}
		})
	}
}

func TestKeyValueConstructorThrowValueLengthError(t *testing.T) {
	// given
	validKey := Faker.RandomStringWithLength(MaxLengthKey)
	valueInvalidValues := []string{"", Faker.RandomStringWithLength(MaxLengthValue + 1)}

	for _, actualTestValue := range valueInvalidValues {
		t.Run(fmt.Sprintf("Testing [%v]", actualTestValue), func(t *testing.T) {
			// when
			keyValue, err := newKeyValue(validKey, actualTestValue)

			// then
			var invalidLengthError *InvalidLengthError
			isInvalidLengthError := errors.As(err, &invalidLengthError)
			if err == nil || !isInvalidLengthError || !strings.Contains(invalidLengthError.Error(), "'value'") {
				t.Logf("expected an InvalidLengthError: , actual keyValue: %v, actual error: %v", keyValue, err)
				t.Fail()
			}
		})
	}
}

func TestKeyValueConstructorShouldConstructWithoutError(t *testing.T) {
	// given
	validKey := Faker.RandomStringWithLength(MaxLengthKey)
	validValue := Faker.RandomStringWithLength(MaxLengthValue)

	// when
	keyValue, err := newKeyValue(validKey, validValue)

	// then
	if err != nil || keyValue.Key.Key != validKey || keyValue.Value.Value != validValue {
		t.Logf("expected an InvalidLengthError: , actual keyValue: %v, actual error: %v", keyValue, err)
		t.Fail()
	}
}
