package ops

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"time"
)

// https://www.linode.com/docs/email/postfix/email-with-postfix-dovecot-and-mysql
// http://wiki.dovecot.org/Authentication/PasswordSchemes
// https://mad9scientist.com/dovecot-password-creation-php/

// Domain domain
type Domain struct {
	tableName struct{} `sql:"mail_domains"`
	ID        uint
	Name      string
	UpdatedAt time.Time
	CreatedAt time.Time
}

// User user
type User struct {
	tableName struct{} `sql:"mail_users"`
	ID        uint
	FullName  string
	Email     string
	Password  string
	Enable    bool
	Domain    Domain
	Begin     time.Time
	End       time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (p *User) sum(password string, salt []byte) string {
	buf := sha512.Sum512(append([]byte(password), salt...))
	return base64.StdEncoding.EncodeToString(append(buf[:], salt...))
}

// SetPassword set  password (SSHA512-CRYPT)
func (p *User) SetPassword(password string) error {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	p.Password = p.sum(password, salt)
	return nil
}

// ChkPassword check password
func (p *User) ChkPassword(password string) bool {
	buf, err := base64.StdEncoding.DecodeString(p.Password)
	if err != nil {
		return false
	}

	return len(buf) > sha512.Size && p.Password == p.sum(password, buf[sha512.Size:])
}

// Alias alias
type Alias struct {
	tableName   struct{} `sql:"mail_aliases"`
	ID          uint
	Source      string
	Destination string
	Domain      Domain
	UpdatedAt   time.Time
	CreatedAt   time.Time
}
