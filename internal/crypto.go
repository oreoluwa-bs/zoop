package internal

type Cipher interface {
	Decrypt(ciphertext string) (string, error)
	Encrypt(plaintext string) (string, error)
}
