package handler

import (
	"encoding/json"
	"net/http"

	"github.com/crutchm/elite/internal/models"
	"github.com/crutchm/elite/internal/service"
)

type VoteHandler struct {
	voteService *service.VoteService
}

func NewVoteHandler(voteService *service.VoteService) *VoteHandler {
	return &VoteHandler{voteService: voteService}
}

func (h *VoteHandler) Vote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tgUserID, ok := r.Context().Value("tg_user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CategoryID == 0 {
		http.Error(w, "category_id is required", http.StatusBadRequest)
		return
	}

	if err := h.voteService.CreateVote(r.Context(), tgUserID, &req); err != nil {
		errMsg := err.Error()
		if errMsg == "vote already exists for this category" {
			http.Error(w, errMsg, http.StatusConflict)
			return
		}
		if errMsg == "nominant not found" || errMsg == "category not found" {
			http.Error(w, errMsg, http.StatusNotFound)
			return
		}
		if errMsg == "nominant is not participating in this category" {
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create vote: "+errMsg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vote created successfully"})
}
