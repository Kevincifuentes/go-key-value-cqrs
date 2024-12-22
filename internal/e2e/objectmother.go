package e2e

import (
	"context"
	"github.com/jaswdr/faker/v2"
	"go-key-value-cqrs/e2e/client"
	"net/http"
)

type KeyValueObjectMother struct {
	fakerInstance *faker.Faker
	client        *client.ClientWithResponses
}

func (objectMother *KeyValueObjectMother) createRandom() (string, string) {
	expectedKey := objectMother.fakerInstance.UUID().V4()
	expectedValue := objectMother.fakerInstance.Person().Name()
	request := client.AddKeyRequest{expectedKey: expectedValue}
	response, err := objectMother.client.PostKeyWithResponse(context.TODO(), request)
	if response != nil && response.StatusCode() == http.StatusNoContent {
		return expectedKey, expectedValue
	}
	panic(err)
}
