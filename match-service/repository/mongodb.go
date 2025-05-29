package repository

import (
	"context"

	"github.com/abaika-abay/live_sports_project/common/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Match struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	MatchID    string             `bson:"match_id"`
	Status     string             `bson:"status"`
	HomeScore  int32              `bson:"home_score"`
	AwayScore  int32              `bson:"away_score"`
	LastEvent  string             `bson:"last_event"`
	Possession int32              `bson:"possession"`
	Shots      int32              `bson:"shots"`
	Fouls      int32              `bson:"fouls"`
	Cards      []string           `bson:"cards"`
}

type MatchRepository struct {
	collection *mongo.Collection
}

func NewMatchRepository(db *db.MongoDB) *MatchRepository {
	return &MatchRepository{
		collection: db.Database.Collection("matches"),
	}
}

func (r *MatchRepository) CreateMatch(ctx context.Context, match *Match) error {
	_, err := r.collection.InsertOne(ctx, match)
	return err
}

func (r *MatchRepository) GetMatch(ctx context.Context, matchID string) (*Match, error) {
	var match Match
	err := r.collection.FindOne(ctx, bson.M{"match_id": matchID}).Decode(&match)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) UpdateMatch(ctx context.Context, matchID string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"match_id": matchID}, bson.M{"$set": update})
	return err
}
