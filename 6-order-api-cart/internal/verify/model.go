package verify

import (
	"purple/links/internal/session"
	"time"
)

type VerifyCode struct {
	ID        uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      uint16          `json:"code" gorm:"not null"`
	SessionID uint            `json:"session_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;index"`
	Session   session.Session `json:"session" gorm:"foreignKey:SessionID;references:ID"`
	CreatedAt time.Time       `json:"created_at" gorm:"not null;autoCreateTime;index"`
}

func (VerifyCode) TableName() string {
	return "verify_code"
}
