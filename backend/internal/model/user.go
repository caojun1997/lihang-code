package model

import (
	"time"
)

type User struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     string    `gorm:"column:user_id;type:varchar(64);uniqueIndex;not null" json:"user_id"`
	Username   string    `gorm:"column:username;type:varchar(128)" json:"username"`
	Email      string    `gorm:"column:email;type:varchar(256)" json:"email"`
	Avatar     string    `gorm:"column:avatar;type:varchar(512)" json:"avatar"`
	Nickname   string    `gorm:"column:nickname;type:varchar(128)" json:"nickname"`
	APIKey     string    `gorm:"column:api_key;type:varchar(256)" json:"api_key"`
	APIKeySecret string  `gorm:"column:api_key_secret;type:varchar(256)" json:"api_key_secret"`
	LastLogin  time.Time `gorm:"column:last_login;autoUpdateTime" json:"last_login"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Session struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID    string    `gorm:"column:session_id;type:varchar(128);uniqueIndex;not null" json:"session_id"`
	UserID       string    `gorm:"column:user_id;type:varchar(64);index;not null" json:"user_id"`
	AccessToken  string    `gorm:"column:access_token;type:text;not null" json:"access_token"`
	RefreshToken string    `gorm:"column:refresh_token;type:text" json:"refresh_token"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Session) TableName() string {
	return "sessions"
}

type CodeExchange struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	State       string    `gorm:"column:state;type:varchar(128);uniqueIndex;not null" json:"state"`
	CodeVerifier string   `gorm:"column:code_verifier;type:varchar(256);not null" json:"code_verifier"`
	RedirectURI string    `gorm:"column:redirect_uri;type:varchar(512);not null" json:"redirect_uri"`
	UserID      string    `gorm:"column:user_id;type:varchar(64)" json:"user_id"`
	Used        bool      `gorm:"column:used;default:false" json:"used"`
	ExpiresAt   time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (CodeExchange) TableName() string {
	return "code_exchanges"
}
