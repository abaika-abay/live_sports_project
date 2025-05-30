package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": update})
	return err
}
