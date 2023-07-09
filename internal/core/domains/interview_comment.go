package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddInterviewComment struct {
	ID        primitive.ObjectID `bson:"_id"`
	Comment   string             `bson:"comment"`
	UserID    primitive.ObjectID `bson:"userId"`
	CreatedAt time.Time          `bson:"CreatedAt"`
	UpdatedAt time.Time          `bson:"UpdatedAt"`
}

type InterviewComment struct {
	ID        primitive.ObjectID `bson:"_id"`
	Comment   string             `bson:"comment"`
	User      User               `bson:"user,omitempty"`
	CreatedAt time.Time          `bson:"CreatedAt"`
	UpdatedAt time.Time          `bson:"UpdatedAt"`
}

type AddInterviewCommentParams struct {
	ID      primitive.ObjectID `bson:"_id"`
	Comment string             `bson:"comment"`
	UserID  primitive.ObjectID `bson:"userId"`
}
