package nut

import (
	"errors"
	"time"

	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
)

// Dao dao
type Dao struct {
	DB       *pg.DB        `inject:""`
	I18n     *web.I18n     `inject:""`
	Security *web.Security `inject:""`
}

// Allow allow
func (p *Dao) Allow(tx *pg.Tx, user uint, name, rty string, rid uint, years, months, days int) error {
	now := time.Now()
	end := now.AddDate(years, months, days)
	role, err := p.getRole(tx, name, rty, rid)
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

	var it Policy
	err = tx.Model(&it).
		Where("role_id = ? AND user_id = ?", role.ID, user).Limit(1).Select()

	if err == nil {
		it.UpdatedAt = now
		it.Begin = now
		it.End = end
		_, err = tx.Model(&it).Column("_begin", "_end", "updated_at").Update()
	} else if err == pg.ErrNoRows {
		it.RoleID = role.ID
		it.UserID = user
		it.UpdatedAt = now
		it.Begin = now
		it.End = end
		err = tx.Insert(&it)
	}
	return err
}

func (p *Dao) getRole(tx *pg.Tx, name, rty string, rid uint) (*Role, error) {
	var it Role
	if err := tx.Model(&it).Column("id").
		Where("name = ? AND resource_type = ? AND resource_id = ?", name, rty, rid).
		Limit(1).Select(); err != nil {
		return nil, err
	}
	return &it, nil
}

// Deny deny
func (p *Dao) Deny(tx *pg.Tx, user uint, name, rty string, rid uint) error {
	role, err := p.getRole(tx, name, rty, rid)
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
func (p *Dao) Can(user uint, name, rty string, rid uint) bool {
	var role Role
	if err := p.DB.Model(&role).
		Column("id").
		Where("name = ? AND resource_type = ? AND resource_id = ?", name, rty, rid).
		Limit(1).Select(); err != nil {
		return false
	}
	var it Policy
	if err := p.DB.Model(&it).
		Column("_begin", "_end").
		Where("user_id = ? AND role_id = ?", user, role.ID).
		Limit(1).Select(); err != nil {
		return false
	}

	return it.Enable()
}

// Is is role?
func (p *Dao) Is(user uint, role string) bool {
	return p.Can(user, role, DefaultResourceType, DefaultResourceID)
}

// AddLog add log
func (p *Dao) AddLog(tx *pg.Tx, user uint, ip, lang, format string, args ...interface{}) error {
	return tx.Insert(&Log{
		UserID:  user,
		IP:      ip,
		Message: p.I18n.T(lang, format, args...),
	})
}

// AddEmailUser add email user
func (p *Dao) AddEmailUser(tx *pg.Tx, name, email, password string) (*User, error) {
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
	if user.Password, err = p.Security.Hash([]byte(password)); err != nil {
		return nil, err
	}

	if err = tx.Insert(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
