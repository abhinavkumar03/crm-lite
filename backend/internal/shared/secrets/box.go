// Package secrets provides AES-GCM encryption for communication provider credentials.
package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var ErrMissingKey = errors.New("communication secrets key not configured")

// Box encrypts/decrypts JSON secret maps with AES-256-GCM.
type Box struct {
	key []byte
}

// NewBox derives a 32-byte key from the given secret string (base64 or raw).
func NewBox(secret string) (*Box, error) {
	if secret == "" {
		return nil, ErrMissingKey
	}
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil || len(key) != 32 {
		// Accept raw 32-byte strings or pad/hash via repeating.
		raw := []byte(secret)
		key = make([]byte, 32)
		copy(key, raw)
		if len(raw) < 32 {
			for i := len(raw); i < 32; i++ {
				key[i] = byte(i)
			}
		}
	}
	return &Box{key: key}, nil
}

// EncryptJSON marshals v and encrypts it.
func (b *Box) EncryptJSON(v any) ([]byte, error) {
	if b == nil || len(b.key) == 0 {
		return nil, ErrMissingKey
	}
	plain, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(b.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plain, nil), nil
}

// DecryptJSON decrypts ciphertext into dest.
func (b *Box) DecryptJSON(ciphertext []byte, dest any) error {
	if b == nil || len(b.key) == 0 {
		return ErrMissingKey
	}
	if len(ciphertext) == 0 {
		return fmt.Errorf("empty ciphertext")
	}
	block, err := aes.NewCipher(b.key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}
	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(plain, dest)
}
