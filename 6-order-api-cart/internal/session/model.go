package session

import (
	"purple/links/internal/user"
	"time"
)

type Session struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Key       string    `json:"key" gorm:"uniqueIndex;size:32;not null"`
	UserID    uint      `json:"user_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;index"`
	User      user.User `json:"user" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime;index"`
}

func (Session) TableName() string {
	return "session"
}
