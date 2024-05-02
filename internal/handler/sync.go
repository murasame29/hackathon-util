package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/murasame29/hackathon-util/internal/application"
)

type SyncRequest struct {
	GuildID       string `json:"guild_id"`
	SpreadSheetID string `json:"spread_sheet_id"`
	SpreadRange   string `json:"spread_range"`
}

type SyncResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	var req SyncRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messageDeleteRole, err := h.app.DeleteRole(r.Context(), application.DeleteRoleParam{
		GuildID:       req.GuildID,
		SpreadSheetID: req.SpreadSheetID,
		Range:         req.SpreadRange,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messageCreateRole, err := h.app.CreateRole(r.Context(), application.CreateRoleParam{
		GuildID:       req.GuildID,
		SpreadSheetID: req.SpreadSheetID,
		Range:         req.SpreadRange,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messageBindRole, err := h.app.BindRole(r.Context(), application.BindRoleParam{
		GuildID:       req.GuildID,
		SpreadSheetID: req.SpreadSheetID,
		Range:         req.SpreadRange,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&SyncResponse{Message: strings.Join(append(messageDeleteRole, append(messageCreateRole, messageBindRole...)...), "\n")}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
