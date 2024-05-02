package handler

import (
	"encoding/json"
	"net/http"

	"github.com/murasame29/hackathon-util/internal/application"
)

type RoleRequest struct {
	Action        ActionType `json:"action"`
	GuildID       string     `json:"guild_id"`
	SpreadSheetID string     `json:"spread_sheet_id"`
	SpreadRange   string     `json:"spread_range"`
}

type RoleResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Role(w http.ResponseWriter, r *http.Request) {
	var req RoleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	switch req.Action {
	case ActionTypeCreate:
		err = h.app.CreateRole(r.Context(), application.CreateRoleParam{
			GuildID:       req.GuildID,
			SpreadSheetID: req.SpreadSheetID,
			Range:         req.SpreadRange,
		})
	case ActionTypeDelete:
		err = h.app.DeleteRole(r.Context(), application.DeleteRoleParam{
			GuildID:       req.GuildID,
			SpreadSheetID: req.SpreadSheetID,
			Range:         req.SpreadRange,
		})
	case ActionTypeBind:
		err = h.app.BindRole(r.Context(), application.BindRoleParam{
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
	if err := json.NewEncoder(w).Encode(&ChannelResponse{Message: "success"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
