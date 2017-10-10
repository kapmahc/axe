package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

var (
	_cip cipher.Block
)

// Encrypt encrypt
func Encrypt(buf []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(_cip, iv)
	val := make([]byte, len(buf))
	cfb.XORKeyStream(val, buf)

	return append(val, iv...), nil
}

// Decrypt decrypt
func Decrypt(buf []byte) ([]byte, error) {
	bln := len(buf)
	cln := bln - aes.BlockSize
	ct := buf[0:cln]
	iv := buf[cln:bln]

	cfb := cipher.NewCFBDecrypter(_cip, iv)
	val := make([]byte, cln)
	cfb.XORKeyStream(val, ct)
	return val, nil
}
