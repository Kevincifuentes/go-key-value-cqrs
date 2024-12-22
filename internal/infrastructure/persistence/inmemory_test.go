package persistence

import (
	"github.com/jaswdr/faker/v2"
	"go-key-value-cqrs/domain"
	"testing"
)

var fakerInstance = faker.New()

const defaultTestKey = "anyKey"
const defaultTestValue = "anyValue"

var repository domain.KeyValueReader

func TestMain(testing *testing.M) {
	initialMap := map[string]string{defaultTestKey: defaultTestValue}
	repository = &InMemoryKeyValueRepository{KeyValueMap: KeyValueMap{keyToValueMap: initialMap}}

	testing.Run()
}

func TestKeyValueFoundWithRepository(t *testing.T) {
	// when
	keyValueView, err := repository.Get(defaultTestKey)

	// then
	if err != nil || keyValueView.Value != defaultTestValue {
		t.Logf(
			"expected no error and expected value %s: actual value=%v, actual error=%v",
			defaultTestValue, keyValueView, err)
		t.Fail()
	}
}

func TestKeyValueNotFoundWithRepository(t *testing.T) {
	//given
	unknownKey := fakerInstance.Person().Name()

	// when
	keyValueView, err := repository.Get(unknownKey)

	// then
	if err == nil || keyValueView != (domain.KeyValueView{}) {
		t.Logf(
			"expected no error and expected value %s: actual value=%v, actual error=%v",
			defaultTestValue, keyValueView, err)
		t.Fail()
	}
}
