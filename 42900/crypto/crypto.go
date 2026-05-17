package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	KeySize  = 32
	NonceSize = 12
	SaltSize  = 16

	Argon2Time    = 3
	Argon2Memory  = 64 * 1024
	Argon2Threads = 4
)

type EncryptedData struct {
	Ciphertext []byte
	Nonce      []byte
	Salt        []byte
}

func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, Argon2Time, Argon2Memory, Argon2Threads, KeySize)
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func Encrypt(plaintext []byte, key []byte) (*EncryptedData, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("invalid key size: expected %d, got %d", KeySize, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, NonceSize)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return &EncryptedData{
		Ciphertext: ciphertext,
		Nonce:      nonce,
	}, nil
}

func Decrypt(data *EncryptedData, key []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("invalid key size: expected %d, got %d", KeySize, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, data.Nonce, data.Ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed: invalid key or corrupted data")
	}

	return plaintext, nil
}

func (ed *EncryptedData) Encode() string {
	data := make([]byte, 0, len(ed.Nonce)+len(ed.Ciphertext))
	data = append(data, ed.Nonce...)
	data = append(data, ed.Ciphertext...)
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeEncodedData(encoded string) (*EncryptedData, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	if len(data) < NonceSize {
		return nil, errors.New("invalid encrypted data")
	}

	return &EncryptedData{
		Nonce:      data[:NonceSize],
		Ciphertext: data[NonceSize:],
	}, nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
