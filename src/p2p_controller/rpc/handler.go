package rpc

import (
	"fmt"
	"github.com/stdi0/archer-problem/src/models"
)

func NewHandler(c models.Controller) *Handler {
	return &Handler{c}
}

type Handler struct {
	Controller models.Controller
}

func (h *Handler) Handle(cmd Command) {
	switch cmd.Cmd {
	case "define_your_number":
		fmt.Println("[RUN] DEFINE_YOUR_NUMBER")
	case "fire":
		fmt.Println("[RUN] FIRE")
	default:
		break
	}
}
