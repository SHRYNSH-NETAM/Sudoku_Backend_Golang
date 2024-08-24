package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/models"
	"github.com/golang-jwt/jwt/v5"
)
var secretKey = []byte("secret-key")

const jwtPayloadKey models.Key = "jwtPayload"

func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), jwtPayloadKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println("1")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}

func CreateToken (email string, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
				jwt.MapClaims{
					"email": email,
					"username": username,
					"exp": time.Now().Add(time.Hour * 48).Unix(),
				})

	tokenstring, err := token.SignedString(secretKey)
	if err!=nil{
		return "",err
	}
	return tokenstring, nil
}