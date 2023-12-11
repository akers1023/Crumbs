package util

import (
	"errors"
	"net/http"
)

func CheckUserType(r *http.Request, role string) error {
	userType := r.Header.Get("user_type")
	if userType != role {
		return errors.New("Unauthorized to access this resource")
	}
	return nil
}

func MatchUserTypeToUid(r *http.Request, userId string) error {
	userType := r.Header.Get("user_type")
	uid := r.Header.Get("uid")

	if userType == "USER" && uid != userId {
		return errors.New("Unauthorized to access this resource")
	}

	return CheckUserType(r, userType)
}
