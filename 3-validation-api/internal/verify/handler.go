package verify

import (
	"net/http"
	"purple/links/configs"
)

type VerifyHandler struct {
	*configs.Config
}

func NewVerifyHandler(router *http.ServeMux, conf *configs.Config) {
	handler := VerifyHandler{Config: conf}
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("/verify/{hash}", handler.Verify())
}

func (h *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}

func (h *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}
