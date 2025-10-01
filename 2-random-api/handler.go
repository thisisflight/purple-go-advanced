package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
)

type RandomHandler struct{}

func (h *RandomHandler) GetRandomNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		randomNumber := rand.IntN(6) + 1
		fmt.Fprintf(w, "%d", randomNumber)
	}
}

func NewRandomHandler(router *http.ServeMux) *RandomHandler {
	handler := &RandomHandler{}
	router.HandleFunc("/random", handler.GetRandomNumber())
	return handler
}
