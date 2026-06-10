package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mastermou8/goProject/internal/store"
	"github.com/mastermou8/goProject/internal/tokens"
	"github.com/mastermou8/goProject/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}
type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	// get the user
	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("Error fetching user: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid credentials"})
		return
	}
	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("Error comparing passwords: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if !passwordsDoMatch {
		h.logger.Printf("Invalid credentials for user: %v", req.Username)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid credentials"})
		return
	}
	// Generate token
	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("Error generating token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
