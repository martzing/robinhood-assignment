package domains

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Seamsi struct {
	ID           primitive.ObjectID `bson:"_id"`
	No           uint32             `bson:"no"`
	Title        string             `bson:"title"`
	Description  string             `bson:"description"`
	LuckyNumbers []string           `bson:"luckyNumbers"`
}
