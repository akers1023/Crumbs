package post

import (
	"context"
	"crumbs/internal/database"
	val "crumbs/internal/handler/user"
	"crumbs/internal/model"
	"crumbs/internal/util"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var post model.Post
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		util.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	validationErr := val.Validate.Struct(post)
	if validationErr != nil {
		util.HandleError(w, validationErr.Error(), http.StatusBadRequest)
		return
	}

	if post.User_id != nil {
		err := database.UserCollection.FindOne(ctx, bson.M{"user_id": post.User_id}).Decode(&user)

		if err != nil {
			fmt.Print(err)
			msg := fmt.Sprintf("message: Table was not found")
			util.HandleError(w, msg, http.StatusInternalServerError)
			return
		}
	}

	post.Created_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
	post.Updated_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))

	post.ID = primitive.NewObjectID()
	post.Post_id = post.ID.Hex()

	resultInsertion, insertErr := database.PostCollection.InsertOne(ctx, post)
	if insertErr != nil {
		msg := fmt.Sprint("Post was not created")
		util.HandleError(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultInsertion)
}
