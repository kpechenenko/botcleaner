package cache

type Cache[K comparable, V any] interface {
	Get(k K) (V, bool)
	Set(k K, v V)
}
