package auth

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

type RussianPhone string

func (p *RussianPhone) UnmarshalJSON(data []byte) error {
	// Убираем кавычки JSON
	phoneStr := strings.Trim(string(data), `"`)

	// Парсим номер
	num, err := phonenumbers.Parse(phoneStr, "RU")
	if err != nil {
		return fmt.Errorf("неверный формат телефонного номера: %v", err)
	}

	// Проверяем валидность для России
	if !phonenumbers.IsValidNumberForRegion(num, "RU") {
		return fmt.Errorf("похоже, введен не российский телефонный номер")
	}

	// Нормализуем к формату 8XXXXXXXXXX
	normalized := p.normalizePhone(num)
	*p = RussianPhone(normalized)

	return nil
}

func (p RussianPhone) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(p) + `"`), nil
}

func (p *RussianPhone) normalizePhone(num *phonenumbers.PhoneNumber) string {
	// Преобразуем в национальный формат (8 XXX XXX-XX-XX)
	nationalFormat := phonenumbers.Format(num, phonenumbers.NATIONAL)

	// Убираем всё, кроме цифр
	digitsOnly := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, nationalFormat)

	return digitsOnly
}

// String возвращает номер в формате 8XXXXXXXXXX
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
