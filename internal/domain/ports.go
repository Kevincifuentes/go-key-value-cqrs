package domain

type KeyValueReader interface {
	Get(key string) (KeyValueView, error)
}

type KeyValueWriter interface {
	Add(keyValue KeyValue) error
}
