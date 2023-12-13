package util

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

func CheckUserType(r *http.Request, role string) error {
	tokenClient := r.Header.Get("Token")
	token, _ := jwt.ParseWithClaims(tokenClient, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	claims, _ := token.Claims.(*SignedDetails)
	userType := claims.User_type

	if userType != role {
		return errors.New("Unauthorized to access this resource")
	}
	return nil
}

func MatchUserTypeToUid(r *http.Request, userId string) error {
	tokenClient := r.Header.Get("Token")
	token, _ := jwt.ParseWithClaims(tokenClient, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	claims, _ := token.Claims.(*SignedDetails)
	userType := claims.User_type
	uid := claims.Uid

	if userType == "USER" && uid != userId {
		return errors.New("Unauthorized to access this resource")
	}

	return CheckUserType(r, userType)
}
