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

type keyValueKey struct {
	key string
}

func newKeyValueKey(key string) (*keyValueKey, error) {
	err := validateLength(key, "key", MinLengthKey, MaxLengthKey)
	if err != nil {
		return nil, err
	}
	return &keyValueKey{key: key}, nil
}

type keyValueValue struct {
	value string
}

func newKeyValueValue(value string) (*keyValueValue, error) {
	err := validateLength(value, "value", MinLengthValue, MaxLengthValue)
	if err != nil {
		return nil, err
	}
	return &keyValueValue{value: value}, nil
}

// AggregateRoot

type KeyValue struct {
	key   *keyValueKey
	value *keyValueValue
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
		key:   keyValue,
		value: valueValue,
	}, nil
}

// View

type KeyValueView struct {
	key   string
	value string
}
