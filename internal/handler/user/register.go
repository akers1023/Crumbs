package user

import (
	"context"
	"crumbs/internal/database"
	"crumbs/internal/model"
	"crumbs/internal/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Validate = validator.New()

func Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		util.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	validationErr := Validate.Struct(user)
	if validationErr != nil {
		util.HandleError(w, validationErr.Error(), http.StatusBadRequest)
		return
	}

	password := util.HashPassword(*user.Password)
	user.Password = &password

	if util.IsValidPhoneNumber(*user.Phone) {
		util.HandleError(w, "invalid phone number", http.StatusBadRequest)
		return
	}
	// if util.IsValidEmail(*user.Email) {
	// 	http.Error(w, "invalid email", http.StatusBadRequest)
	// 	return
	// }
	count, err := database.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		log.Panic(err)
		util.HandleError(w, "error occurred while checking for the email", http.StatusInternalServerError)
		return
	}
	count, err = database.UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		log.Panic(err)
		util.HandleError(w, "error occurred while checking for the phone number", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		util.HandleError(w, "This email or phone number is already in use", http.StatusInternalServerError)
		return
	}

	user.Created_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
	user.Updated_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	token, refreshToken, _ := util.GenerateAllTokens(*user.Email, *user.Phone, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	resultInsertion, insertErr := database.UserCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := fmt.Sprint("User was not created")
		util.HandleError(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultInsertion)
}
