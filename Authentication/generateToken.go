package authentication

import (
	"os"
	"time"
     "net/http"
	"github.com/golang-jwt/jwt/v5"
	"online_bookStore/Handlers"
	"strings"
	"context" 
	
	"online_bookStore/Interfaces"
	"io"
	"encoding/json"
	
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if len(jwtSecret) == 0 {
		panic("JWT_SECRET is not set")
	}
}

func GenerateToken(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handlers.WriteError(w, http.StatusUnauthorized, "missing token")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			handlers.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		// store user info in context
		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		ctx = context.WithValue(ctx, "role", claims["role"])

		next(w, r.WithContext(ctx))
	}
}



type AuthHandler struct {
	userStore interfaces.UserStore
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	// 1. verify user 
	user, err := h.userStore.GetByEmail(ctx,req.Email)
	if err != nil || !CheckPassword(req.Password, user.Password) {
		handlers.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// 2. generate token
	token, err := GenerateToken(user.ID, user.Role)
	if err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// 3. return token
	resp, _ := json.Marshal(map[string]string{
		"token": token,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}


func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
	return err == nil
}


