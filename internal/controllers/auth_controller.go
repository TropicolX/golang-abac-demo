package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"golang-abac-demo/internal/models"
	"golang-abac-demo/internal/utils"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

var JwtKey = []byte("my_secret_key")

const UserKey contextKey = "user"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// This should be replaced with a call to your user repository
	user, err := models.GetUserByUsername(creds.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if user.Password != creds.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.InfoLogger.Printf("User '%s' logged in", user.Username)

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "token": tokenString})
}
