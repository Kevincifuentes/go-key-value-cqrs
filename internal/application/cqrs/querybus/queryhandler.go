package querybus

type QueryHandler[T Query, K any] interface {
	Ask(query T) (K, error)
}
