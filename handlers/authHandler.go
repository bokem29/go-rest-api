package handlers

import (
	"encoding/json"
	"go-rest/utils"
	"net/http"
	"time"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

// LoginHandler handles /api/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if !utils.Authenticate(req.Username, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	Token, err := utils.CreateToken(req.Username)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}
	refresh, err := utils.CreateRefreshToken(req.Username)
	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}
	// set cookies so browser requests (no custom headers) can access protected endpoints
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: Token, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: refresh, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse{Token: Token, Refresh: refresh})
}

// LogoutHandler handles /api/logout, invalidates token
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	token, err := utils.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	utils.InvalidateToken(token)
	// clear auth cookies
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "", Path: "/", Expires: time.Unix(0, 0), MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Path: "/", Expires: time.Unix(0, 0), MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	w.WriteHeader(http.StatusNoContent)
}

// RefreshHandler rotates refresh token and returns new access + refresh token
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Refresh string `json:"refresh"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Refresh == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	access, refresh, err := utils.ValidateAndRotateRefresh(body.Refresh)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// rotate cookies so browser stays authenticated
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: access, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: refresh, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse{Token: access, Refresh: refresh})
}
