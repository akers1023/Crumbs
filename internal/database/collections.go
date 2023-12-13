package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection *mongo.Collection = OpenCollection(Client, "user")
var PostCollection *mongo.Collection = OpenCollection(Client, "post")
