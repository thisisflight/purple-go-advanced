package user

import "purple/links/pkg/db"

type UserRepository struct {
	db *db.DB
}

func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(user *User) (*User, error) {
	result := u.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (u *UserRepository) GetByPhoneNumber(phoneNumber string) (*User, error) {
	var user User
	result := u.db.First(&user, "phone_number = ?", phoneNumber)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
