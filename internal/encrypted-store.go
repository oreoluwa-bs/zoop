package internal

type EncryptedStore struct {
	store  Store
	cipher Cipher
}

func NewEncryptedStore(store Store, cipher Cipher) *EncryptedStore {
	return &EncryptedStore{
		store:  store,
		cipher: cipher,
	}
}

func (e *EncryptedStore) Set(key, value string) error {
	encrypted, err := e.cipher.Encrypt(value)
	if err != nil {
		return err
	}

	return e.store.Set(key, encrypted)
}

func (e *EncryptedStore) Get(key string) (string, error) {
	encrypted, err := e.store.Get(key)
	if err != nil {
		return "", err
	}

	return e.cipher.Decrypt(encrypted)
}

func (e *EncryptedStore) Delete(key string) error {
	return e.store.Delete(key)
}

func (e *EncryptedStore) GetAllKeys() ([]string, error) {
	return e.store.GetAllKeys()
}
