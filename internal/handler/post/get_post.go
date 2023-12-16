package post

import (
	"context"
	"crumbs/internal/database"
	"crumbs/internal/model"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PostId := vars["post_id"]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var post model.Post
	err := database.PostCollection.FindOne(ctx, bson.M{"post_id": PostId}).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}
