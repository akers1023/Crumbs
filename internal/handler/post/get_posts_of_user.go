package post

import (
	"context"
	"crumbs/internal/database"
	"crumbs/internal/util"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get posts by user
func GetPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	recordPerPage, err := strconv.Atoi(vars["recordPerPage"])
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err1 := strconv.Atoi(vars["page"])
	if err1 != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(vars["startIndex"])

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}
	result, err := database.PostCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage})
	if err != nil {
		util.HandleError(w, "error occurred while listing user items", http.StatusInternalServerError)
		return
	}

	var allPostOfUser []bson.M
	if err = result.All(ctx, &allPostOfUser); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allPostOfUser[0])

}
