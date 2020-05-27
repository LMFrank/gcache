package sgcache

type Getter interface {
	Get(key string) ([]byte, error)
}

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type GetterFunc func(key string) ([]byte, error)
