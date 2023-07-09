package domains

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetFortuneTellingsParams struct {
	SearchText *string
	Prefix     *string
	Offset     *uint32
	Limit      *uint32
}

type CountFortuneTellingsParams struct {
	SearchText *string
	Prefix     *string
}

type CreateFortuneTellingParams struct {
	Title        string
	Description  string
	LuckyNumbers []string
	Prefix       string
}

type UpdateFortuneTellingParams struct {
	ID           string
	Title        *string
	Description  *string
	LuckyNumbers []string
	Prefix       *string
}

type FortuneTelling struct {
	ID           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	Description  string             `bson:"description"`
	LuckyNumbers []string           `bson:"luckyNumbers"`
	Prefix       string             `bson:"prefix"`
}
