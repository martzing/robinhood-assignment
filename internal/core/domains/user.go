package domains

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserParams struct {
	Name     string
	Email    string
	Username string
	Password string
	ImageUrl string
	Role     string
}

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	ImageUrl string             `bson:"imageUrl"`
	Role     string             `bson:"role"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
