package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"
)

// Encryptor struct
type Encryptor struct {
	key         []byte
	expireInDay bool
}

// IEncryptor is the interface of the Encryptor struct
type IEncryptor interface {
	SetKey(key []byte)
	ExpiryInDay() IEncryptor
	NoExpiryInDay() IEncryptor
	Encrypt(text string) (string, error)
	EncryptByKey(key []byte, text string) (string, error)
	Decrypt(cryptoText string) (string, error)
	DecryptByKey(key []byte, cryptoText string) (string, error)
}

// New create new instance of Encryptor
func New() IEncryptor {
	return &Encryptor{}
}

// SetKey set key
func (e *Encryptor) SetKey(key []byte) {
	e.key = key
}

// ExpiryInDay set expire in day
func (e *Encryptor) ExpiryInDay() IEncryptor {
	e.expireInDay = true
	return e
}

// NoExpiryInDay remove expire in day
func (e *Encryptor) NoExpiryInDay() IEncryptor {
	e.expireInDay = false
	return e
}

// Encrypt string by global key
func (e *Encryptor) Encrypt(text string) (string, error) {
	return e.EncryptByKey(e.getKey(), text)
}

// Decrypt string by global key
func (e *Encryptor) Decrypt(cryptoText string) (string, error) {
	return e.DecryptByKey(e.getKey(), cryptoText)
}

// returns key
func (e *Encryptor) getKey() []byte {
	key := e.key
	if e.expireInDay {
		key = []byte(time.Now().Format("20060102") + fmt.Sprintf("%v", key)) // YYYY-MM-DD
	}
	return key
}

// ensure key size
func (e *Encryptor) ensureKey(key []byte) ([]byte, error) {
	keyLen := len(key)
	keySize := keyLen
	if keyLen < 16 {
		return nil, errors.New("key must be at least 16 characters")
	}
	if keyLen > 16 && keyLen < 24 {
		keySize = 16
	}
	if keyLen > 24 && keyLen < 32 {
		keySize = 24
	}
	if keyLen > 32 {
		keySize = 32
	}

	return key[:keySize], nil
}

// EncryptByKey encrypt string by given key
func (e *Encryptor) EncryptByKey(key []byte, text string) (string, error) {
	b, err := e.ensureKey(key)
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	block, err := aes.NewCipher(b)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore, it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), err
}

// DecryptByKey decrypt string by given key
func (e *Encryptor) DecryptByKey(key []byte, cryptoText string) (string, error) {
	b, err := e.ensureKey(key)
	if err != nil {
		return "", err
	}

	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher(b)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore, it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}
