package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID          primitive.ObjectID `bson:"_id"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	Description string             `json:"description" validate:"required"`
	Post_id     *string            `json:"post_id"`
	User_id     *string            `json:"user_id" validate:"required"`
}
