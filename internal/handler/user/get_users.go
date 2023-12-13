package user

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

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Respond with the user data (you might want to exclude sensitive information)
	vars := mux.Vars(r)
	if err := util.CheckUserType(r, "ADMIN"); err != nil {
		util.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(vars["startIndex"])

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}
	result, err := database.UserCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage})
	if err != nil {
		util.HandleError(w, "error occurred while listing user items", http.StatusInternalServerError)
		return
	}

	var allusers []bson.M
	if err = result.All(ctx, &allusers); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allusers[0])
}
