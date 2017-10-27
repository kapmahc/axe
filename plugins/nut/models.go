package nut

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// RoleAdmin admin role
	RoleAdmin = "admin"
	// RoleRoot root role
	RoleRoot = "root"
	// UserTypeEmail email user
	UserTypeEmail = "email"

	// DefaultResourceType default resource type
	DefaultResourceType = ""
	// DefaultResourceID default resourc id
	DefaultResourceID = 0
)

// User user
type User struct {
	tableName       struct{} `sql:"users"`
	ID              uint
	Name            string
	Email           string
	UID             string
	Password        []byte
	ProviderID      string
	ProviderType    string
	Logo            string
	SignInCount     uint
	LastSignInAt    *time.Time
	LastSignInIP    string
	CurrentSignInAt *time.Time
	CurrentSignInIP string
	ConfirmedAt     *time.Time
	LockedAt        *time.Time
	Logs            []Log
	UpdatedAt       time.Time
	CreatedAt       time.Time
}

// IsConfirm is confirm?
func (p *User) IsConfirm() bool {
	return p.ConfirmedAt != nil
}

// IsLock is lock?
func (p *User) IsLock() bool {
	return p.LockedAt != nil
}

// SetGravatarLogo set logo by gravatar
func (p *User) SetGravatarLogo() {
	// https: //en.gravatar.com/site/implement/
	buf := md5.Sum([]byte(strings.ToLower(p.Email)))
	p.Logo = fmt.Sprintf("https://gravatar.com/avatar/%s.png", hex.EncodeToString(buf[:]))
}

//SetUID generate uid
func (p *User) SetUID() {
	p.UID = uuid.New().String()
}

func (p User) String() string {
	return fmt.Sprintf("%s<%s>", p.Name, p.Email)
}

// Attachment attachment
type Attachment struct {
	tableName    struct{} `sql:"attachments"`
	ID           uint
	Title        string
	URL          string
	Length       int64
	MediaType    string
	ResourceID   uint
	ResourceType string
	User         User
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

// IsPicture is picture?
func (p *Attachment) IsPicture() bool {
	return strings.HasPrefix(p.MediaType, "image/")
}

// Log log
type Log struct {
	tableName struct{}  `sql:"logs"`
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	IP        string    `json:"ip"`
	User      *User     `json:"user"`
	UserID    uint      `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (p Log) String() string {
	return fmt.Sprintf("%s: [%s]\t %s", p.CreatedAt.Format(time.ANSIC), p.IP, p.Message)
}

// Policy policy
type Policy struct {
	tableName struct{} `sql:"policies"`
	ID        uint
	Begin     time.Time `sql:"_begin"`
	End       time.Time `sql:"_end"`
	User      *User
	UserID    uint
	Role      *Role
	RoleID    uint
	UpdatedAt time.Time
	CreatedAt time.Time
}

//Enable is enable?
func (p *Policy) Enable() bool {
	now := time.Now()
	return now.After(p.Begin) && now.Before(p.End)
}

// Role role
type Role struct {
	tableName    struct{} `sql:"roles"`
	ID           uint
	Name         string
	ResourceID   uint   `sql:",notnull"`
	ResourceType string `sql:",notnull"`
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func (p Role) String() string {
	return fmt.Sprintf("%s@%s://%d", p.Name, p.ResourceType, p.ResourceID)
}

// Vote vote
type Vote struct {
	tableName    struct{} `sql:"votes"`
	ID           uint
	Point        int
	ResourceID   uint
	ResourceType string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

// LeaveWord leave-word
type LeaveWord struct {
	tableName struct{} `sql:"leave_words"`
	ID        uint
	Body      string
	Type      string
	CreatedAt time.Time
}

// Link link
type Link struct {
	tableName struct{} `sql:"links"`
	ID        uint
	Loc       string
	Href      string
	Label     string
	SortOrder int
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Card card
type Card struct {
	tableName struct{} `sql:"cards"`
	ID        uint
	Loc       string
	Title     string
	Summary   string
	Type      string
	Href      string
	Logo      string
	SortOrder int
	Action    string
	UpdatedAt time.Time
	CreatedAt time.Time
}

// FriendLink friend_links
type FriendLink struct {
	tableName struct{} `sql:"friend_links"`
	ID        uint
	Title     string
	Home      string
	Logo      string
	SortOrder int
	UpdatedAt time.Time
	CreatedAt time.Time
}
