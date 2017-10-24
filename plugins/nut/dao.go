package nut

import (
	"errors"
	"time"

	"github.com/go-pg/pg"
)

// Allow allow
func Allow(tx *pg.Tx, user uint, name, rty string, rid uint, years, months, days int) error {
	now := time.Now()
	end := now.AddDate(years, months, days)
	role, err := getRole(tx, name, rty, rid)
	if err != nil {
		if err != pg.ErrNoRows {
			return err
		}
		role = &Role{
			Name:         name,
			ResourceType: rty,
			ResourceID:   rid,
			UpdatedAt:    now,
		}
		err = tx.Insert(role)
	}
	if err != nil {
		return err
	}

	var p Policy
	err = tx.Model(&p).
		Where("role_id = ? AND user_id = ?", role.ID, user).Limit(1).Select()

	if err == nil {
		p.UpdatedAt = now
		p.Begin = now
		p.End = end
		_, err = tx.Model(&p).Column("_begin", "_end", "updated_at").Update()
	} else if err == pg.ErrNoRows {
		p.RoleID = role.ID
		p.UserID = user
		p.UpdatedAt = now
		p.Begin = now
		p.End = end
		err = tx.Insert(&p)
	}
	return err
}

func getRole(tx *pg.Tx, name, rty string, rid uint) (*Role, error) {
	var it Role
	if err := tx.Model(&it).Column("id").
		Where("name = ? AND resource_type = ? AND resource_id = ?", name, rty, rid).
		Limit(1).Select(); err != nil {
		return nil, err
	}
	return &it, nil
}

// Deny deny
func Deny(tx *pg.Tx, user uint, name, rty string, rid uint) error {
	role, err := getRole(tx, name, rty, rid)
	if err != nil {
		return err
	}
	if _, err := tx.Model(&Policy{}).
		Where("role_id = ? AND user_id = ?", role.ID, user).
		Delete(); err != nil {
		return err
	}

	return nil
}

// Can can?
func Can(user uint, name, rty string, rid uint) bool {
	db := DB()
	var role Role
	if err := db.Model(&role).
		Column("id").
		Where("name = ? AND resource_type = ? AND resource_id = ?", name, rty, rid).
		Limit(1).Select(); err != nil {
		return false
	}
	var it Policy
	if err := db.Model(&it).
		Column("_begin", "_end").
		Where("user_id = ? AND role_id = ?", user, role.ID).
		Limit(1).Select(); err != nil {
		return false
	}

	return it.Enable()
}

// Is is role?
func Is(user uint, role string) bool {
	return Can(user, role, DefaultResourceType, DefaultResourceID)
}

// AddLog add log
func AddLog(tx *pg.Tx, user uint, ip, message string) error {
	return tx.Insert(&Log{
		UserID:  user,
		IP:      ip,
		Message: message,
	})
}

// AddEmailUser add email user
func AddEmailUser(tx *pg.Tx, name, email, password string) (*User, error) {
	now := time.Now()
	cnt, err := tx.Model(&User{}).
		Where("provider_type = ? AND provider_id = ?", UserTypeEmail, email).
		Count()
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("email already exists")
	}
	user := User{
		Email:           email,
		Name:            name,
		ProviderType:    UserTypeEmail,
		ProviderID:      email,
		LastSignInIP:    "0.0.0.0",
		CurrentSignInIP: "0.0.0.0",
		UpdatedAt:       now,
	}
	user.SetUID()
	user.SetGravatarLogo()
	if user.Password, err = SECURITY().Hash([]byte(password)); err != nil {
		return nil, err
	}

	if err = tx.Insert(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
