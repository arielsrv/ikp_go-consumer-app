package kvs

type Client[TValue any] interface {
	Get(key string) (*TValue, error)
	Save(key string, value *TValue) error
}
