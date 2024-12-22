package model

import "go-key-value-cqrs/domain"

type KeyValueResponse map[string]string

func ToKeyValueResponse(view domain.KeyValueView) KeyValueResponse {
	return KeyValueResponse{
		view.Key: view.Value,
	}
}
