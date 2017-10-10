package nut

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var (
	_security    *Security
	securityOnce sync.Once
)

// SECURITY security handle
func SECURITY() *Security {
	securityOnce.Do(func() {
		buf, err := base64.StdEncoding.DecodeString(viper.GetString("secret"))
		if err != nil {
			log.Error(err)
			return
		}
		cip, err := aes.NewCipher(buf)
		if err != nil {
			log.Error(err)
			return
		}
		_security = &Security{cip: cip}
	})
	return _security
}

// Security security helper
type Security struct {
	cip cipher.Block
}

// Hash ont-way encrypt
func (p *Security) Hash(plain []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(plain, 16)
}

// Check check hash
func (p *Security) Check(encode, plain []byte) bool {
	return bcrypt.CompareHashAndPassword(encode, plain) == nil
}

// Encrypt encrypt
func (p *Security) Encrypt(buf []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(p.cip, iv)
	val := make([]byte, len(buf))
	cfb.XORKeyStream(val, buf)

	return append(val, iv...), nil
}

// Decrypt decrypt
func (p *Security) Decrypt(buf []byte) ([]byte, error) {
	bln := len(buf)
	cln := bln - aes.BlockSize
	ct := buf[0:cln]
	iv := buf[cln:bln]

	cfb := cipher.NewCFBDecrypter(p.cip, iv)
	val := make([]byte, cln)
	cfb.XORKeyStream(val, ct)
	return val, nil
}
