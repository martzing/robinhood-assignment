package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddInterviewComment struct {
	ID        primitive.ObjectID `bson:"_id"`
	Comment   string             `bson:"comment"`
	UserID    primitive.ObjectID `bson:"userId"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type InterviewComment struct {
	ID        primitive.ObjectID `bson:"_id"`
	Comment   string             `bson:"comment"`
	User      User               `bson:"user"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type AddInterviewCommentParams struct {
	ID      primitive.ObjectID
	Comment string
	UserID  primitive.ObjectID
}

type UpdateInterviewCommentParams struct {
	ID        primitive.ObjectID
	CommentID primitive.ObjectID
	Comment   string
}
