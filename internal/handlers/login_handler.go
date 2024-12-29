package handlers

import (
	"encoding/json"
	"file_mgmt_system/helper"
	"file_mgmt_system/internal/models"
	"file_mgmt_system/internal/service"
	"file_mgmt_system/middleware"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("KA11EL4943")

type LoginHandler struct {
	service *service.LoginService
}

func NewLoginHandler(service *service.LoginService) *LoginHandler {
	return &LoginHandler{
		service: service,
	}
}

type Response struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	FileExists   bool   `json:"fileExists"`
	RowsAffected int    `json:"rowsAffected"`
}

func (handler *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var req models.Input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	expirationTime := time.Now().Add(1 * time.Hour) // Token valid for 1 hour
	claims := middleware.JWTClaims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	success := false
	rowsAffected, fileExists, err := handler.service.Login(&req)
	if err != nil {
		http.Error(w, "Failed to process login", http.StatusInternalServerError)
		return
	}
	success = true
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true, // Prevent client-side scripts from accessing the cookie
		Path:     "/",  // Make the cookie accessible to the entire site
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Success:      success,
		Message:      "Login successful",
		FileExists:   fileExists,
		RowsAffected: rowsAffected,
	})
}

func (handler *LoginHandler) GetEmail(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
		Email   string `json:"email,omitempty"` // Omits the field if empty
	}
	email, ok := helper.GetEmailFromContext(r.Context())
	if !ok {
		response := Response{
			Success: false,
			Message: "error while fetching email",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{
		Success: true,
		Email:   email,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
