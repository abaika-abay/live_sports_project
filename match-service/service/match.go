package service

import (
	"context"
	"fmt"
	"time"

	"github.com/abaika-abay/live_sports_project/common/pkg/db"
	"github.com/abaika-abay/live_sports_project/match-service/proto"
	"github.com/abaika-abay/live_sports_project/match-service/repository"
	"github.com/abaika-abay/live_sports_project/match-service/sportradar" // Import sportradar client
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb" // Added for GetAdminMatchList
)

// MatchService implements the gRPC server for match operations.
type MatchService struct {
	proto.UnimplementedMatchServiceServer
	repo             *repository.MatchRepository
	sportradarClient *sportradar.SportradarClient // New Sportradar client
	// You might add a Kafka/NATS producer here for notifications
	// notificationProducer *kafka.Producer // Placeholder
}

// NewMatchService initializes the MatchService with a repository and Sportradar client.
func NewMatchService(database *db.MongoDB, srClient *sportradar.SportradarClient) *MatchService {
	return &MatchService{
		repo:             repository.NewMatchRepository(database),
		sportradarClient: srClient,
		// notificationProducer: kafka.NewProducer(), // Initialize your producer
	}
}

// GetMatchUpdates fetches real-time match data, combining internal and Sportradar sources.
func (s *MatchService) GetMatchUpdates(ctx context.Context, req *proto.MatchRequest) (*proto.MatchResponse, error) {
	// 1. Get base match data from our internal DB
	match, err := s.repo.GetMatch(ctx, req.MatchId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "match not found in internal DB: %v", err)
	}

	// 2. Fetch real-time data from Sportradar (mocked here)
	srMatch, err := s.sportradarClient.FetchMatchData(ctx, req.MatchId)
	if err != nil {
		// Log the error but don't fail the request if Sportradar is down or has no data
		fmt.Printf("Warning: Could not fetch real-time data from Sportradar for match %s: %v\n", req.MatchId, err)
		// Optionally, return only internal data if Sportradar data is critical
	} else {
		// 3. Combine/Merge data: Prioritize Sportradar for live scores/events/stats
		// You might have more complex merging logic based on data freshness/completeness
		match.Status = srMatch.Status
		match.HomeScore = srMatch.HomeScore
		match.AwayScore = srMatch.AwayScore
		match.LastEvent = srMatch.LastEvent
		match.Possession = srMatch.Possession
		match.Shots = srMatch.Shots
		match.Fouls = srMatch.Fouls
		match.Cards = srMatch.Cards

		// Optional: Update internal DB with latest Sportradar data
		// This keeps your internal data fresh but adds a write operation
		if err := s.repo.UpdateMatch(ctx, match); err != nil {
			fmt.Printf("Warning: Failed to update internal match data from Sportradar for match %s: %v\n", req.MatchId, err)
		}
	}

	return &proto.MatchResponse{
		MatchId:    match.MatchID,
		Status:     match.Status,
		HomeScore:  match.HomeScore,
		AwayScore:  match.AwayScore,
		LastEvent:  match.LastEvent,
		Possession: match.Possession,
		Shots:      match.Shots,
		Fouls:      match.Fouls,
		Cards:      match.Cards,
	}, nil
}

// CreateMatch creates a new match in the database, potentially also notifying Sportradar mock.
func (s *MatchService) CreateMatch(ctx context.Context, req *proto.CreateMatchRequest) (*proto.MatchResponse, error) {
	match := &repository.Match{
		MatchID:    req.MatchId,
		HomeTeam:   req.HomeTeam,
		AwayTeam:   req.AwayTeam,
		StartTime:  req.StartTime,
		Status:     "Scheduled",
		HomeScore:  0,
		AwayScore:  0,
		LastEvent:  "Match Created",
		Possession: 50, // Default or initial value
		Shots:      0,
		Fouls:      0,
		Cards:      []string{},
	}

	if err := s.repo.CreateMatch(ctx, match); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create match in DB: %v", err)
	}

	// Also add to Sportradar mock for consistency in simulation
	s.sportradarClient.AddInitialMatchData(match)

	return &proto.MatchResponse{
		MatchId:    match.MatchID,
		Status:     match.Status,
		HomeScore:  match.HomeScore,
		AwayScore:  match.AwayScore,
		LastEvent:  match.LastEvent,
		Possession: match.Possession,
		Shots:      match.Shots,
		Fouls:      match.Fouls,
		Cards:      match.Cards,
	}, nil
}

// UpdateMatchEvent handles admin-submitted events like goals, fouls, cards.
func (s *MatchService) UpdateMatchEvent(ctx context.Context, req *proto.UpdateMatchEventRequest) (*proto.MatchResponse, error) {
	// 1. Get the current match state from DB
	match, err := s.repo.GetMatch(ctx, req.MatchId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "match not found for event update: %v", err)
	}

	// 2. Apply event specific logic
	eventDescription := fmt.Sprintf("%s: %s", req.EventType, req.Description)
	switch req.EventType {
	case "goal":
		match.HomeScore += req.HomeScoreChange
		match.AwayScore += req.AwayScoreChange
		match.LastEvent = fmt.Sprintf("GOAL! %s", req.Description)
	case "foul":
		match.Fouls++
		match.LastEvent = fmt.Sprintf("FOUL: %s", req.Description)
	case "card":
		match.Cards = append(match.Cards, fmt.Sprintf("%s_%s", req.CardColor, req.Description)) // e.g., "yellow_Ronaldo"
		match.LastEvent = fmt.Sprintf("%s CARD: %s", req.CardColor, req.Description)
	case "status_change": // Example for changing match status
		match.Status = req.Description // e.g., "Halftime", "Fulltime"
		match.LastEvent = fmt.Sprintf("STATUS: %s", req.Description)
	default:
		// Generic update for other event types
		match.LastEvent = eventDescription
	}

	// 3. Update the match in the database
	if err := s.repo.UpdateMatch(ctx, match); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update match after event: %v", err)
	}

	// 4. Record the event itself in the events collection
	event := &repository.Event{
		EventID:     fmt.Sprintf("evt-%s-%d", req.MatchId, time.Now().UnixNano()), // Unique ID for event
		MatchID:     req.MatchId,
		EventType:   req.EventType,
		Description: req.Description,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	if err := s.repo.AddEvent(ctx, event); err != nil {
		fmt.Printf("Warning: Failed to add event record for match %s: %v\n", req.MatchId, err)
		// Don't fail the main request if event record fails, but log it.
	}

	// 5. Notify Sportradar mock (optional, for simulation consistency)
	s.sportradarClient.UpdateData(ctx, match)

	// 6. Trigger notification (via Kafka/NATS to Notification Service)
	// Example:
	// if s.notificationProducer != nil {
	// 	notificationMsg := &proto.NotificationMessage{
	// 		MatchId:     req.MatchId,
	// 		EventType:   req.EventType,
	// 		Description: match.LastEvent, // Send the formatted last event
	// 		Timestamp:   time.Now().Format(time.RFC3339),
	// 	}
	// 	// s.notificationProducer.Publish("match_events", notificationMsg) // Assuming a publish method
	// }

	return &proto.MatchResponse{
		MatchId:    match.MatchID,
		Status:     match.Status,
		HomeScore:  match.HomeScore,
		AwayScore:  match.AwayScore,
		LastEvent:  match.LastEvent,
		Possession: match.Possession,
		Shots:      match.Shots,
		Fouls:      match.Fouls,
		Cards:      match.Cards,
	}, nil
}

// GetAdminMatchList retrieves all matches for the admin panel.
func (s *MatchService) GetAdminMatchList(ctx context.Context, req *emptypb.Empty) (*proto.MatchListResponse, error) {
	matches, err := s.repo.GetMatchListForAdmin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve admin match list: %v", err)
	}

	var protoMatches []*proto.MatchResponse
	for _, match := range matches {
		protoMatches = append(protoMatches, &proto.MatchResponse{
			MatchId:    match.MatchID,
			Status:     match.Status,
			HomeScore:  match.HomeScore,
			AwayScore:  match.AwayScore,
			LastEvent:  match.LastEvent,
			Possession: match.Possession,
			Shots:      match.Shots,
			Fouls:      match.Fouls,
			Cards:      match.Cards,
		})
	}

	return &proto.MatchListResponse{Matches: protoMatches}, nil
}
