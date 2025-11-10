package auth

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

type RussianPhone string

func (p *RussianPhone) UnmarshalJSON(data []byte) error {
	phoneStr := strings.Trim(string(data), `"`)
	num, err := phonenumbers.Parse(phoneStr, "RU")
	if err != nil {
		return fmt.Errorf("неверный формат телефонного номера: %v", err)
	}
	if !phonenumbers.IsValidNumberForRegion(num, "RU") {
		return fmt.Errorf("похоже, введен не российский телефонный номер")
	}
	normalized := p.normalizePhone(num)
	*p = RussianPhone(normalized)
	return nil
}

func (p RussianPhone) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(p) + `"`), nil
}

func (p *RussianPhone) normalizePhone(num *phonenumbers.PhoneNumber) string {
	nationalFormat := phonenumbers.Format(num, phonenumbers.NATIONAL)
	digitsOnly := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, nationalFormat)

	return digitsOnly
}

func (p RussianPhone) String() string {
	return string(p)
}

type AuthRequest struct {
	PhoneNumber RussianPhone `json:"phone_number" validate:"required"`
}

type AuthResponse struct {
	SessionID string `json:"session_id"`
}

type VerifyRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	Code      uint16 `json:"code" validate:"required"`
}

type VerifyResponse struct {
	Token string `json:"token"`
}
