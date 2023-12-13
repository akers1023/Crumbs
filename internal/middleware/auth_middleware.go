package middleware

import (
	"crumbs/internal/util"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Lấy giá trị token từ header
		clientToken := r.Header.Get("token")

		_, err := jwt.ParseWithClaims(clientToken, &util.SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(util.SECRET_KEY), nil
		})
		if err != nil {
			util.HandleError(w, fmt.Sprintf("Error parsing token: %v", err), http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
