package domain

const MaxLengthKey = 200
const MinLengthKey = 1
const MaxLengthValue = 200
const MinLengthValue = 1

// Value Objects

func validateLength(value string, valueName string, minLength int, maxLength int) error {
	if len(value) < minLength || len(value) > maxLength {
		return NewInvalidLengthError(value, valueName, MinLengthKey, MaxLengthKey)
	}
	return nil
}

type KeyValueKey struct {
	Key string
}

func newKeyValueKey(key string) (*KeyValueKey, error) {
	err := validateLength(key, "key", MinLengthKey, MaxLengthKey)
	if err != nil {
		return nil, err
	}
	return &KeyValueKey{Key: key}, nil
}

type KeyValueValue struct {
	Value string
}

func newKeyValueValue(value string) (*KeyValueValue, error) {
	err := validateLength(value, "value", MinLengthValue, MaxLengthValue)
	if err != nil {
		return nil, err
	}
	return &KeyValueValue{Value: value}, nil
}

// AggregateRoot

type KeyValue struct {
	Key   *KeyValueKey
	Value *KeyValueValue
}

func newKeyValue(key string, value string) (*KeyValue, error) {
	keyValue, err := newKeyValueKey(key)
	if err != nil {
		return nil, err
	}

	valueValue, err := newKeyValueValue(value)
	if err != nil {
		return nil, err
	}
	return &KeyValue{
		Key:   keyValue,
		Value: valueValue,
	}, nil
}

// View

type KeyValueView struct {
	Key   string
	Value string
}
