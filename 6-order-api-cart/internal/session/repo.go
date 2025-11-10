package session

import (
	"purple/links/internal/user"
	"purple/links/pkg/db"
	"purple/links/pkg/tokens"
)

type SessionRepository struct {
	db *db.DB
}

func NewSessionRepository(db *db.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (s *SessionRepository) Create(user *user.User) (*Session, error) {
	var session Session
	session.Key, _ = tokens.GenerateSecureToken(32)
	session.User = *user
	result := s.db.Create(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (s *SessionRepository) FindByKey(key string) (*Session, error) {
	var session Session
	result := s.db.Preload("User").First(&session, "key = ?", key)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}
