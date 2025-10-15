package auth

import (
	"errors"
	"net/http"

	"purple/links/pkg/req"
	res "purple/links/pkg/res"

	"gorm.io/gorm"
)

type AuthHandlerDeps struct {
	*AuthService
}

type AuthHandler struct {
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth", handler.Auth())
	router.HandleFunc("POST /verify", handler.Verify())
}

func (h *AuthHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[AuthRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sessionID, err := h.AuthService.Auth(string(body.PhoneNumber))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		data := AuthResponse{
			SessionID: sessionID,
		}
		res.Json(w, data, http.StatusCreated)
	}
}

func (h *AuthHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[VerifyRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := h.AuthService.Verify(body.SessionID, body.Code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data := VerifyResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusOK)
	}
}
