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

func (e *EncryptedStore) Get(key string) (string, error) {
	encryptedValue, err := e.store.Get(key)
	if err != nil {
		return "", err
	}
	return e.cipher.Decrypt(encryptedValue)
}

func (e *EncryptedStore) Set(key, value string) error {
	encryptedValue, err := e.cipher.Encrypt(value)
	if err != nil {
		return err
	}
	return e.store.Set(key, encryptedValue)
}

func (e *EncryptedStore) Delete(key string) error {
	return e.store.Delete(key)
}

func (e *EncryptedStore) GetAllKeys() ([]string, error) {
	return e.store.GetAllKeys()
}
