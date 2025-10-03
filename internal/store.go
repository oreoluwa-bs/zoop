package internal

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	GetAllKeys() ([]string, error)
}

type StoreManager struct {
	Store Store `json:"-"`
}

func NewStoreManager(store Store) *StoreManager {
	return &StoreManager{
		Store: store,
	}
}
