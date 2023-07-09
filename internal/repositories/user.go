package repositories

import (
	"context"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type user struct {
	mc  *mongo.Client
	db  string
	cn  string
	col *mongo.Collection
}

func NewUserRepository(mc *mongo.Client, db string) ports.UserRepository {
	cn := "user"
	return &user{
		mc:  mc,
		db:  db,
		cn:  cn,
		col: mc.Database(db).Collection(cn),
	}
}

func (u *user) Get(ctx context.Context, id primitive.ObjectID) (*domains.User, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	res := domains.User{}
	if err := u.col.FindOne(ctx, filter).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (u *user) GetByUsername(ctx context.Context, username string) (*domains.User, error) {
	filter := bson.D{{Key: "username", Value: username}}
	res := domains.User{}
	if err := u.col.FindOne(ctx, filter).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (u *user) Create(ctx context.Context, params *domains.CreateUserParams) (*domains.User, error) {
	user := domains.User{
		ID:       primitive.NewObjectID(),
		Name:     params.Name,
		Email:    params.Email,
		Username: params.Username,
		Password: params.Password,
		ImageUrl: params.ImageUrl,
		Role:     params.Role,
	}
	if _, err := u.col.InsertOne(ctx, user); err != nil {
		return nil, err
	}
	return &user, nil
}
