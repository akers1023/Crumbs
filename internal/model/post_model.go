package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID          primitive.ObjectID `bson:"_id"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	Title       *string            `json:"title" validate:"required"`
	Description *string            `json:"description" validate:"required,min=1,max=5000"`
	Status      *string            `json:"status"`
	Post_image  *[]string          `json:"post_image"`
	Post_id     string             `json:"post_id"`
	User_id     *string            `json:"user_id" validate:"required"`
	// Comment_lid *[]string          `json:"comment"`
}
