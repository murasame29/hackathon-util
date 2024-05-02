package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/murasame29/hackathon-util/internal/application"
)

type ChannelRequest struct {
	Action        ActionType `json:"action"`
	GuildID       string     `json:"guild_id"`
	SpreadSheetID string     `json:"spread_sheet_id"`
	SpreadRange   string     `json:"spread_range"`
}

type ChannelResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Channel(w http.ResponseWriter, r *http.Request) {
	var req ChannelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		err     error
		message []string
	)
	switch req.Action {
	case ActionTypeCreate:
		message, err = h.app.CraeteChannel(r.Context(), application.CreateChannelParam{
			GuildID:       req.GuildID,
			SpreadSheetID: req.SpreadSheetID,
			Range:         req.SpreadRange,
		})
	case ActionTypeDelete:
		message, err = h.app.DeleteChannel(r.Context(), application.DeleteChannelParam{
			GuildID:       req.GuildID,
			SpreadSheetID: req.SpreadSheetID,
			Range:         req.SpreadRange,
		})
	default:
		http.Error(w, "invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&ChannelResponse{Message: strings.Join(message, "\n")}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
