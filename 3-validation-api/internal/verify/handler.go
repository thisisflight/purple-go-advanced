package verify

import (
	"fmt"
	"net/http"
	"purple/links/configs"
	"purple/links/pkg/mail"
	"purple/links/pkg/req"
	res "purple/links/pkg/res"
	"purple/links/pkg/tokens"
	"purple/links/storage"
)

type VerifyHandler struct {
	*configs.Config
	repo *storage.TokenRepository
}

func NewVerifyHandler(router *http.ServeMux, conf *configs.Config, repo *storage.TokenRepository) {
	handler := VerifyHandler{Config: conf, repo: repo}
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("/verify/{hash}", handler.Verify())
}

func (h *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[SendRequest](w, r)
		if err != nil {
			return
		}
		token, err := tokens.GenerateSecureToken(32)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		email_text := fmt.Sprintf(
			"Перейдите по ссылке, чтобы подтвердить почту: %s%s/verify/%s",
			h.Config.ServerConfig.Schema, h.Config.ServerConfig.Addr, token)
		mail.SendVerifyEmail(body.Email, email_text, h.Config)
		h.repo.AddTokenRecord(&storage.TokenRecord{Email: body.Email, Token: token})
		res.Json(w, "Sent", http.StatusOK)
	}
}

func (h *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		isRemoved := h.repo.RemoveRecordByToken(hash)
		if isRemoved {
			res.Json(w, "Deleted", http.StatusNoContent)
			return
		}
		res.Json(w, "Not found", http.StatusNotFound)
	}
}
