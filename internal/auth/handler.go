package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"server-auth-go/internal/token"
)

// Handler handles auth requests
type Handler struct {
	store        Store
	tokenService *token.Service
}

// NewHandler creates a new auth handler
func NewHandler(store Store, tokenService *token.Service) *Handler {
	return &Handler{
		store:        store,
		tokenService: tokenService,
	}
}

type tokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// HandleToken handles the token exchange request
func (h *Handler) HandleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, err := h.store.VerifyCredentials(req.ClientID, req.ClientSecret)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !valid {
		http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
		return
	}

	duration := time.Hour
	accessToken, err := h.tokenService.GenerateAccessToken(req.ClientID, duration)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := tokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(duration.Seconds()),
		TokenType:   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
