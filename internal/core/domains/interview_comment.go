package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InterviewComment struct {
	ID        primitive.ObjectID `bson:"_id"`
	Comment   string             `bson:"comment"`
	User      User               `bson:"user"`
	CreatedAt time.Time          `bson:"CreatedAt"`
	UpdatedAt time.Time          `bson:"UpdatedAt"`
}
