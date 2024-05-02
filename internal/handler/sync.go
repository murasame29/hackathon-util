package handler

import (
	"encoding/json"
	"net/http"

	"github.com/murasame29/hackathon-util/internal/application"
)

type SyncRequest struct {
	GuildID       string `json:"guild_id"`
	SpreadSheetID string `json:"spread_sheet_id"`
	SpreadRange   string `json:"spread_range"`
}

type SyncResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	var req SyncRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.app.DeleteRole(r.Context(), application.DeleteRoleParam{
		GuildID:       req.GuildID,
		SpreadSheetID: req.SpreadSheetID,
		Range:         req.SpreadRange,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.app.CreateRole(r.Context(), application.CreateRoleParam{
		GuildID:       req.GuildID,
		SpreadSheetID: req.SpreadSheetID,
		Range:         req.SpreadRange,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.app.BindRole(r.Context(), application.BindRoleParam{
		GuildID:       req.GuildID,
		SpreadSheetID: req.SpreadSheetID,
		Range:         req.SpreadRange,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&SyncResponse{Message: "success"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
