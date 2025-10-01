package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/anish-chanda/openwaitlist/backend/internal/db"
	"github.com/anish-chanda/openwaitlist/backend/internal/logger"
	"github.com/anish-chanda/openwaitlist/backend/internal/models"
	"github.com/anish-chanda/openwaitlist/backend/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-pkgz/auth/v2/token"
)

type WaitlistResponse struct {
	ID                 int64  `json:"id"`
	Slug               string `json:"slug"`
	Name               string `json:"name"`
	OwnerUserID        int64  `json:"owner_user_id"`
	IsPublic           bool   `json:"is_public"`
	ShowVendorBranding bool   `json:"show_vendor_branding"`
	CreatedAt          string `json:"created_at"`
}

type CreateWaitlistRequest struct {
	Name               string `json:"name"`
	IsPublic           bool   `json:"is_public"`
	ShowVendorBranding bool   `json:"show_vendor_branding"`
}

type WaitlistsResponse struct {
	Waitlists []WaitlistResponse `json:"waitlists"`
	Total     int                `json:"total"`
}

// getUserIDFromRequest extracts the database user ID from the authenticated request
func getUserIDFromRequest(r *http.Request, database db.Database, log logger.ServiceLogger) (int64, error) {
	tokenUser, err := token.GetUserInfo(r)
	if err != nil {
		return 0, fmt.Errorf("failed to get user info from token: %w", err)
	}

	// Check if user ID is stored in user attributes
	if userIDStr := tokenUser.StrAttr("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			return userID, nil
		}
	}

	// Get email from token (go-pkgz/auth stores email in Name field)
	email := tokenUser.Name
	if email == "" {
		return 0, fmt.Errorf("no email found in token")
	}

	// Get user from database using email
	dbUser, err := database.GetUserByEmail(context.Background(), email)
	if err != nil {
		return 0, fmt.Errorf("failed to get user from database: %w", err)
	}

	return dbUser.ID, nil
}

// GetWaitlistsHandler returns all waitlists for the authenticated user with optional search
func GetWaitlistsHandler(database db.Database, log logger.ServiceLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user ID
		userID, err := getUserIDFromRequest(r, database, log)
		if err != nil {
			log.Error("Failed to get user ID: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get search parameter from query
		searchName := strings.TrimSpace(r.URL.Query().Get("search"))

		// Get waitlists from database
		waitlists, err := database.GetWaitlistsByUserID(context.Background(), userID, searchName)
		if err != nil {
			log.Error("Failed to get waitlists: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Convert to response format
		var waitlistResponses []WaitlistResponse
		for _, wl := range waitlists {
			waitlistResponses = append(waitlistResponses, WaitlistResponse{
				ID:                 wl.ID,
				Slug:               wl.Slug,
				Name:               wl.Name,
				OwnerUserID:        wl.OwnerUserID,
				IsPublic:           wl.IsPublic,
				ShowVendorBranding: wl.ShowVendorBranding,
				CreatedAt:          wl.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		response := WaitlistsResponse{
			Waitlists: waitlistResponses,
			// total would be usefull for pagination in future
			Total: len(waitlistResponses),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("Failed to encode response: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// CreateWaitlistHandler creates a new waitlist for the authenticated user
func CreateWaitlistHandler(database db.Database, log logger.ServiceLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user ID
		userID, err := getUserIDFromRequest(r, database, log)
		if err != nil {
			log.Error("Failed to get user ID: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse request body
		var req CreateWaitlistRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Failed to decode request body: ", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if strings.TrimSpace(req.Name) == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		// Create waitlist
		waitlist := &models.Waitlist{
			Slug:               utils.GenerateSlugFromName(strings.TrimSpace(req.Name)),
			Name:               strings.TrimSpace(req.Name),
			OwnerUserID:        userID,
			IsPublic:           req.IsPublic,
			ShowVendorBranding: req.ShowVendorBranding,
		}

		if err := database.CreateWaitlist(r.Context(), waitlist); err != nil {
			log.Error("Failed to create waitlist: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		} // Return created waitlist
		response := WaitlistResponse{
			ID:                 waitlist.ID,
			Slug:               waitlist.Slug,
			Name:               waitlist.Name,
			OwnerUserID:        waitlist.OwnerUserID,
			IsPublic:           waitlist.IsPublic,
			ShowVendorBranding: waitlist.ShowVendorBranding,
			CreatedAt:          waitlist.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("Failed to encode response: ", err)
		}
	}
}

// GetWaitlistHandler returns a specific waitlist by slug
func GetWaitlistHandler(database db.Database, log logger.ServiceLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user ID
		userID, err := getUserIDFromRequest(r, database, log)
		if err != nil {
			log.Error("Failed to get user ID: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		slug := chi.URLParam(r, "slug")
		if slug == "" {
			http.Error(w, "Slug is required", http.StatusBadRequest)
			return
		}

		waitlist, err := database.GetWaitlistBySlug(context.Background(), slug)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "Waitlist not found", http.StatusNotFound)
				return
			}
			log.Error("Failed to get waitlist: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check if the user owns this waitlist
		if waitlist.OwnerUserID != userID {
			http.Error(w, "Forbidden: You don't own this waitlist", http.StatusForbidden)
			return
		}

		response := WaitlistResponse{
			ID:                 waitlist.ID,
			Slug:               waitlist.Slug,
			Name:               waitlist.Name,
			OwnerUserID:        waitlist.OwnerUserID,
			IsPublic:           waitlist.IsPublic,
			ShowVendorBranding: waitlist.ShowVendorBranding,
			CreatedAt:          waitlist.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("Failed to encode response: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

// UpdateWaitlistHandler updates a waitlist by slug
func UpdateWaitlistHandler(database db.Database, log logger.ServiceLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r, database, log)
		if err != nil {
			log.Error("Failed to get user ID: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		slug := chi.URLParam(r, "slug")
		if slug == "" {
			http.Error(w, "Slug is required", http.StatusBadRequest)
			return
		}

		// Parse request body
		var req CreateWaitlistRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate input
		if strings.TrimSpace(req.Name) == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		// First get the waitlist to check ownership
		waitlist, err := database.GetWaitlistBySlug(r.Context(), slug)
		if err != nil {
			log.Error("Failed to get waitlist: ", err)
			if err.Error() == "waitlist not found" {
				http.Error(w, "Waitlist not found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Check if the user owns this waitlist
		if waitlist.OwnerUserID != userID {
			http.Error(w, "Forbidden: You don't own this waitlist", http.StatusForbidden)
			return
		}

		// Update the waitlist fields
		waitlist.Name = strings.TrimSpace(req.Name)
		waitlist.Slug = utils.GenerateSlugFromName(strings.TrimSpace(req.Name))
		waitlist.IsPublic = req.IsPublic
		waitlist.ShowVendorBranding = req.ShowVendorBranding

		// Update in database
		if err := database.UpdateWaitlist(r.Context(), waitlist); err != nil {
			log.Error("Failed to update waitlist: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Create response
		response := WaitlistResponse{
			ID:                 waitlist.ID,
			Slug:               waitlist.Slug,
			Name:               waitlist.Name,
			OwnerUserID:        waitlist.OwnerUserID,
			IsPublic:           waitlist.IsPublic,
			ShowVendorBranding: waitlist.ShowVendorBranding,
			CreatedAt:          waitlist.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}

		log.Info(fmt.Sprintf("Waitlist updated successfully: %s by user %d", slug, userID))
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("Failed to encode response: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

// DeleteWaitlistHandler deletes a waitlist by slug
func DeleteWaitlistHandler(database db.Database, log logger.ServiceLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r, database, log)
		if err != nil {
			log.Error("Failed to get user ID: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		slug := chi.URLParam(r, "slug")
		if slug == "" {
			http.Error(w, "Slug is required", http.StatusBadRequest)
			return
		}

		// First get the waitlist to check ownership
		waitlist, err := database.GetWaitlistBySlug(r.Context(), slug)
		if err != nil {
			log.Error("Failed to get waitlist: ", err)
			if err.Error() == "waitlist not found" {
				http.Error(w, "Waitlist not found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Check if the user owns this waitlist
		if waitlist.OwnerUserID != userID {
			http.Error(w, "Forbidden: You don't own this waitlist", http.StatusForbidden)
			return
		}

		// Delete the waitlist
		if err := database.DeleteWaitlist(r.Context(), waitlist.ID); err != nil {
			log.Error("Failed to delete waitlist: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Info(fmt.Sprintf("Waitlist deleted successfully: %s by user %d", slug, userID))
		w.WriteHeader(http.StatusNoContent)
	}
}
