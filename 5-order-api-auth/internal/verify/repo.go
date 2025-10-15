package verify

import (
	"fmt"
	"purple/links/internal/session"
	"purple/links/pkg/db"
	"purple/links/pkg/utils"
)

type VerifyRepository struct {
	db *db.DB
}

func NewVerifyRepository(db *db.DB) *VerifyRepository {
	return &VerifyRepository{
		db: db,
	}
}

func (v *VerifyRepository) Create(session *session.Session) (*VerifyCode, error) {
	var verifyСode VerifyCode
	verifyСode.Code = utils.GenerateRandomCode()
	verifyСode.Session = *session
	result := v.db.Create(&verifyСode)
	if result.Error != nil {
		return nil, result.Error
	}
	return &verifyСode, nil
}

func (v *VerifyRepository) FindBySessionKey(sessionKey string) (*VerifyCode, error) {
	var verifyCode VerifyCode
	result := v.db.Joins("JOIN session ON verify_code.session_id = session.id").
		Where("session.key = ?", sessionKey).
		First(&verifyCode)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("не должно быть выведено")
	return &verifyCode, nil
}

func (v *VerifyRepository) Delete(verifyCode *VerifyCode) error {
	result := v.db.Delete(verifyCode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
