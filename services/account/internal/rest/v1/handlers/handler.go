package handlers

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/Molov30/go-microservices/services/account/internal/rest/v1/generated"
)

type Handler struct {
	generated.ServerInterface
	logger *zerolog.Logger
}

func NewHandler(logger *zerolog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
