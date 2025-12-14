package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/crutchm/elite/internal/auth"
	"github.com/crutchm/elite/internal/service"
)

type AuthHandler struct {
	authService *auth.TelegramAuth
	userService *service.UserService
}

func NewAuthHandler(authService *auth.TelegramAuth, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// AuthRequest содержит данные от Telegram Login Widget
type AuthRequest struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Error("invalid http method")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("invalid request body", slog.Any("err", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Преобразуем в LoginWidgetData
	loginData := &auth.LoginWidgetData{
		ID:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		PhotoURL:  req.PhotoURL,
		AuthDate:  req.AuthDate,
		Hash:      req.Hash,
	}

	telegramUser, err := h.authService.ValidateLoginWidgetData(loginData)
	if err != nil {
		slog.Error("invalid request body", slog.Any("err", err))
		http.Error(w, "Invalid login data: "+err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = h.userService.GetOrCreateUser(r.Context(), telegramUser.ID)
	if err != nil {
		slog.Error("failed to get or create user", slog.Any("err", err))
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := h.authService.GenerateToken(telegramUser.ID)
	if err != nil {
		slog.Error("failed to generate token", slog.Any("err", err))
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
