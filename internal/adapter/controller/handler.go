package controller

import "github.com/murasame29/hackathon-util/internal/application"

type Handler struct {
	app *application.ApplicationService
}

func NewHandler(app *application.ApplicationService) *Handler {
	return &Handler{app: app}
}

type ActionType string

const (
	ActionTypeCreate ActionType = "create"
	ActionTypeDelete ActionType = "delete"
	ActionTypeBind   ActionType = "bind"
)
