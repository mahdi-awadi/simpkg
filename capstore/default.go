package capstore

const DefaultStoreKey = "defaultStore"

type defaultStore struct {
	pool IPool
}

func newDefaultStore() *defaultStore {
	return &defaultStore{
		pool: NewPool(),
	}
}

func (store *defaultStore) Pool() IPool {
	return store.pool
}
