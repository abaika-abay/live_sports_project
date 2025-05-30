package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/abaika-abay/live_sports_project/common/pkg/db"
)

// Match represents the structure of a match document in MongoDB
type Match struct {
	MatchID    string   `bson:"match_id"`
	HomeTeam   string   `bson:"home_team"`  // Added
	AwayTeam   string   `bson:"away_team"`  // Added
	StartTime  string   `bson:"start_time"` // Added
	Status     string   `bson:"status"`
	HomeScore  int32    `bson:"home_score"`
	AwayScore  int32    `bson:"away_score"`
	LastEvent  string   `bson:"last_event"`
	Possession int32    `bson:"possession"`
	Shots      int32    `bson:"shots"`
	Fouls      int32    `bson:"fouls"`
	Cards      []string `bson:"cards"` // e.g., ["home_yellow", "away_red"]
}

// Event represents a match event, to be stored in a separate collection or embedded
type Event struct {
	EventID     string `bson:"event_id"`
	MatchID     string `bson:"match_id"`
	EventType   string `bson:"event_type"` // "goal", "foul", "card", "substitution"
	Description string `bson:"description"`
	Timestamp   string `bson:"timestamp"` // ISO 8601 string
	// Potentially other fields for specific event types (e.g., player_id, team_id, score_change)
}

// MatchRepository handles database operations for matches and events
type MatchRepository struct {
	matchesCollection *mongo.Collection
	eventsCollection  *mongo.Collection
}

func NewMatchRepository(database *db.MongoDB) *MatchRepository {
	// Call GetDatabase(). Since it only returns one value,
	// assign it directly to a local variable 'mongoDBInstance'.
	// This variable holds the *mongo.Database object.
	mongoDBInstance := database.GetDatabase() // Corrected Line 49

	return &MatchRepository{
		// Now use the 'mongoDBInstance' variable to get your collections
		matchesCollection: mongoDBInstance.Collection("matches"), // Corrected Line 52
		eventsCollection:  mongoDBInstance.Collection("events"),  // Corrected Line 53
	}
}
// GetMatch retrieves a match by its ID
func (r *MatchRepository) GetMatch(ctx context.Context, matchID string) (*Match, error) {
	var match Match
	err := r.matchesCollection.FindOne(ctx, bson.M{"match_id": matchID}).Decode(&match)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("match with ID %s not found", matchID)
		}
		return nil, fmt.Errorf("failed to get match: %w", err)
	}
	return &match, nil
}

// CreateMatch inserts a new match into the database
func (r *MatchRepository) CreateMatch(ctx context.Context, match *Match) error {
	_, err := r.matchesCollection.InsertOne(ctx, match)
	if err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}
	return nil
}

// UpdateMatch updates an existing match in the database.
// This is crucial for integrating Sportradar updates and admin event changes.
func (r *MatchRepository) UpdateMatch(ctx context.Context, match *Match) error {
	filter := bson.M{"match_id": match.MatchID}
	update := bson.M{"$set": match}                                                                // Update all fields
	_, err := r.matchesCollection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true)) // Upsert if not exists
	if err != nil {
		return fmt.Errorf("failed to update match: %w", err)
	}
	return nil
}

// AddEvent adds a new event to the events collection
func (r *MatchRepository) AddEvent(ctx context.Context, event *Event) error {
	_, err := r.eventsCollection.InsertOne(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}
	return nil
}

// GetMatchListForAdmin retrieves a list of all matches (for admin panel)
func (r *MatchRepository) GetMatchListForAdmin(ctx context.Context) ([]*Match, error) {
	cursor, err := r.matchesCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get match list: %w", err)
	}
	defer cursor.Close(ctx)

	var matches []*Match
	if err = cursor.All(ctx, &matches); err != nil {
		return nil, fmt.Errorf("failed to decode match list: %w", err)
	}
	return matches, nil
}
