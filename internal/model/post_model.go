package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID          primitive.ObjectID `bson:"_id"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	Tittle      string             `json:"tittle" validate:"required"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	Post_image  *[]string          `json:"post_image"`
	// Comment_lid *[]string          `json:"comment"`
}
