package auth

import (
	"purple/links/configs"
	errs "purple/links/errs/verify"
	"purple/links/internal/session"
	"purple/links/internal/user"
	"purple/links/internal/verify"
	"purple/links/pkg/jwt"
)

type AuthServiceDeps struct {
	Conf              *configs.Config
	UserRepository    *user.UserRepository
	SessionRepository *session.SessionRepository
	VerifyRepository  *verify.VerifyRepository
	JWT               *jwt.JWT
}

type AuthService struct {
	Conf              *configs.Config
	UserRepository    *user.UserRepository
	SessionRepository *session.SessionRepository
	VerifyRepository  *verify.VerifyRepository
	JWT               *jwt.JWT
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{
		Conf:              deps.Conf,
		UserRepository:    deps.UserRepository,
		SessionRepository: deps.SessionRepository,
		VerifyRepository:  deps.VerifyRepository,
		JWT:               deps.JWT,
	}
}

func (s *AuthService) Auth(phoneNumber string) (string, error) {
	existedUser, _ := s.UserRepository.GetByPhoneNumber(phoneNumber)
	authUser := existedUser
	if authUser == nil {
		authUser = &user.User{PhoneNumber: phoneNumber}
		s.UserRepository.Create(authUser)
	}
	session, err := s.SessionRepository.Create(authUser)
	if err != nil {
		return "", err
	}
	_, err = s.VerifyRepository.Create(session)
	if err != nil {
		return "", err
	}
	return session.Key, nil
}

func (s *AuthService) Verify(sessionKey string, code uint16) (string, error) {
	session, err := s.SessionRepository.FindByKey(sessionKey)
	if err != nil {
		return "", err
	}
	verifyCode, err := s.VerifyRepository.FindBySessionKey(sessionKey)
	if err != nil {
		return "", err
	}
	if verifyCode.Code != code {
		return "", &errs.CodeDoesNotMatchError{}
	}
	user := session.User
	token, err := s.JWT.Create(user.PhoneNumber)
	if err != nil {
		return "", err
	}
	return token, nil
}
