package objectmothers

import (
	"github.com/jaswdr/faker/v2"
	"go-key-value-cqrs/domain"
)

type KeyValueObjectMother struct {
	FakerInstance *faker.Faker
}

func (objectMother *KeyValueObjectMother) CreateRandom() domain.KeyValue {
	expectedKey := objectMother.FakerInstance.UUID().V4()
	expectedValue := objectMother.FakerInstance.Person().Name()
	return domain.KeyValue{
		Key:   &domain.KeyValueKey{Key: expectedKey},
		Value: &domain.KeyValueValue{Value: expectedValue},
	}
}
