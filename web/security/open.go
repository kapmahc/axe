package security

import "crypto/aes"

// Open open
func Open(secret string) error {
	cip, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return err
	}
	_cip = cip
	return nil
}
