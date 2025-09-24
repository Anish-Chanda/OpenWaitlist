package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/anish-chanda/openwaitlist/backend/internal/db"
	"github.com/anish-chanda/openwaitlist/backend/internal/logger"
	"github.com/anish-chanda/openwaitlist/backend/internal/models"
	"github.com/anish-chanda/openwaitlist/backend/internal/utils"
)

// HandleLogin validates user credentials for the local auth provider
func HandleLogin(database db.Database, email, password string) (bool, error) {
	// Get user by email
	user, err := database.GetUserByEmail(context.Background(), email)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user uses local auth provider
	if user.AuthProvider != "local" {
		return false, fmt.Errorf("user does not use local authentication")
	}

	// Check if password hash exists
	if user.PasswordHash == nil {
		return false, fmt.Errorf("user has no password set")
	}

	// Verify password
	isValid, err := utils.VerifyPassword(password, *user.PasswordHash)
	if err != nil {
		return false, fmt.Errorf("password verification error: %w", err)
	}

	return isValid, nil
}

type SignupRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name,omitempty"`
}

type SignupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int64  `json:"user_id,omitempty"`
}

// SignupHandler handles user signup requests
func SignupHandler(database db.Database, log logger.ServiceLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Failed to parse signup request: ", err)
			writeErrorResponse(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Validate input
		if req.Email == "" || req.Password == "" {
			writeErrorResponse(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Validate email format (basic validation)
		if !strings.Contains(req.Email, "@") {
			writeErrorResponse(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		// Check if user already exists
		existingUser, err := database.GetUserByEmail(context.Background(), req.Email)
		if err == nil && existingUser != nil {
			writeErrorResponse(w, "User with this email already exists", http.StatusConflict)
			return
		}

		// Hash password
		hashedPasswordStr, err := utils.HashPassword(req.Password)
		if err != nil {
			log.Error("Failed to hash password: ", err)
			writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Create user
		user := &models.User{
			Email:        req.Email,
			AuthProvider: "local",
			PasswordHash: &hashedPasswordStr,
		}

		if req.DisplayName != "" {
			user.DisplayName = &req.DisplayName
		}

		if err := database.CreateUser(context.Background(), user); err != nil {
			log.Error("Failed to create user: ", err)
			writeErrorResponse(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Return success response
		response := SignupResponse{
			Success: true,
			Message: "User created successfully",
			UserID:  user.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
		
		log.Info("User created successfully", map[string]interface{}{
			"email": req.Email,
			"user_id": user.ID,
		})
	}
}

// writeErrorResponse writes an error response in JSON format
func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := SignupResponse{
		Success: false,
		Message: message,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
