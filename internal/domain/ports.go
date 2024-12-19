package domain

type KeyValueReader interface {
	Get(key string) (KeyValueView, error)
}
