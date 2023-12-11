package user

import (
	"context"
	"crumbs/internal/database"
	"crumbs/internal/model"
	"crumbs/internal/util"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type LoginSuccesssResponse struct {
	Status       int     `json:"status"`
	Message      string  `json:"message,omitempty"`
	Token        *string `json:"token"`
	RefreshToken *string `json:"freshToken"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user model.User
	var foundUser model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Xác định phương thức đăng nhập và truy vấn dựa trên đó
	var query bson.M

	switch r.URL.Path {
	case "/login/email":
		query = bson.M{"email": user.Email}
		// fmt.Println(*user.Email)
	case "/login/phone":
		query = bson.M{"phone": user.Phone}
	case "/login/user_name":
		query = bson.M{"user_name": user.User_name}
	default:
		http.Error(w, "Invalid login method", http.StatusBadRequest)
		return
	}
	// fmt.Println(r.URL.Path)
	// Thực hiện truy vấn để kiểm tra tính đúng đắn của thông tin đăng nhập
	err := database.UserCollection.FindOne(ctx, query).Decode(&foundUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println(*user.Password)

	passwordIsValid, msg := util.VerifyPassword(*user.Password, *foundUser.Password)
	if !passwordIsValid {
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	if foundUser.Email == nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	token, refreshToken, _ := util.GenerateAllTokens(*foundUser.Email, *foundUser.Phone, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
	util.UpdateAllTokens(token, refreshToken, foundUser.User_id)

	// Decode and return founduser from database
	err = database.UserCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := LoginSuccesssResponse{
		Status:       http.StatusOK,
		Message:      "Login successful",
		Token:        foundUser.Token,
		RefreshToken: foundUser.Refresh_token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
