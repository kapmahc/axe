package security

import "golang.org/x/crypto/bcrypt"

// Password encode password
func Password(plain []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(plain, 16)
}

// Check check password
func Check(encode, plain []byte) bool {
	return bcrypt.CompareHashAndPassword(encode, plain) == nil
}
