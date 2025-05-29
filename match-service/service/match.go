package service

import (
	"context"
	"fmt"
	"time"

	"github.com/abaika-abay/live_sports_project/common/pkg/db"
	"github.com/abaika-abay/live_sports_project/match-service/proto"
	"github.com/abaika-abay/live_sports_project/match-service/repository"
	"github.com/abaika-abay/live_sports_project/match-service/sportradar" // Import the package
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MatchService struct {
	proto.UnimplementedMatchServiceServer
	repo             *repository.MatchRepository
	sportradarClient sportradar.SportradarClientI
	websocketHub     *WebSocketHub // Added WebSocket hub
	// notificationProducer *kafka.Producer // Placeholder for Kafka/NATS
}

// NewMatchService initializes the MatchService.
func NewMatchService(database *db.MongoDB, srClient sportradar.SportradarClientI, wsHub *WebSocketHub) *MatchService {
	return &MatchService{
		repo:             repository.NewMatchRepository(database),
		sportradarClient: srClient,
		websocketHub:     wsHub, // Pass the hub
	}
}

// GetMatchUpdates fetches real-time match data, combining internal and Sportradar sources.
func (s *MatchService) GetMatchUpdates(ctx context.Context, req *proto.MatchRequest) (*proto.MatchResponse, error) {
	// 1. Get base match data from our internal DB
	match, err := s.repo.GetMatch(ctx, req.MatchId)
	if err != nil {
		// If match not found in DB, try fetching from Sportradar first, then creating
		srMatch, srErr := s.sportradarClient.FetchMatchData(ctx, req.MatchId)
		if srErr != nil {
			return nil, status.Errorf(codes.NotFound, "match not found in internal DB and failed to fetch from Sportradar: %v", srErr)
		}
		// If found in SR but not local, create it locally
		match = &repository.Match{
			MatchID:    srMatch.MatchID,
			HomeTeam:   srMatch.HomeTeam,
			AwayTeam:   srMatch.AwayTeam,
			Status:     srMatch.Status,
			HomeScore:  srMatch.HomeScore,
			AwayScore:  srMatch.AwayScore,
			LastEvent:  srMatch.LastEvent,
			Possession: srMatch.Possession,
			Shots:      srMatch.Shots,
			Fouls:      srMatch.Fouls,
			Cards:      srMatch.Cards,
			StartTime:  time.Now().Format(time.RFC3339), // Placeholder if not in SR initial fetch
		}
		if err := s.repo.CreateMatch(ctx, match); err != nil {
			fmt.Printf("Warning: Failed to create match %s in DB after fetching from Sportradar: %v\n", req.MatchId, err)
			// Proceed with SR data even if DB create fails, but log
		}
	}

	// 2. Fetch real-time data from Sportradar
	srMatch, err := s.sportradarClient.FetchMatchData(ctx, req.MatchId)
	if err != nil {
		fmt.Printf("Warning: Could not fetch real-time data from Sportradar for match %s: %v\n", req.MatchId, err)
		// Return internal data if Sportradar call fails
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

	// 3. Combine/Merge data: Prioritize Sportradar for live scores/events/stats
	// You might have more complex merging logic based on data freshness/completeness
	match.Status = srMatch.Status
	match.HomeScore = srMatch.HomeScore
	match.AwayScore = srMatch.AwayScore
	match.LastEvent = srMatch.LastEvent
	match.Possession = srMatch.Possession
	match.Shots = srMatch.Shots
	match.Fouls = srMatch.Fouls
	match.Cards = srMatch.Cards // Assuming SR also provides card info, or merge

	// Optional: Update internal DB with latest Sportradar data
	// This keeps your internal data fresh but adds a write operation
	if err := s.repo.UpdateMatch(ctx, match); err != nil {
		fmt.Printf("Warning: Failed to update internal match data from Sportradar for match %s: %v\n", req.MatchId, err)
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

// UpdateMatchEvent handles admin-submitted events.
func (s *MatchService) UpdateMatchEvent(ctx context.Context, req *proto.UpdateMatchEventRequest) (*proto.MatchResponse, error) {
	match, err := s.repo.GetMatch(ctx, req.MatchId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "match not found for event update: %v", err)
	}

	eventDescription := fmt.Sprintf("%s: %s", req.EventType, req.Description)
	switch req.EventType {
	case "goal":
		match.HomeScore += req.HomeScoreChange
		match.AwayScore += req.AwayScoreChange
		match.LastEvent = fmt.Sprintf("GOAL! %s (%s)", req.Description, time.Now().Format("15:04:05")) // Add timestamp for clarity
	case "foul":
		match.Fouls++
		match.LastEvent = fmt.Sprintf("FOUL: %s (%s)", req.Description, time.Now().Format("15:04:05"))
	case "card":
		match.Cards = append(match.Cards, fmt.Sprintf("%s_%s", req.CardColor, req.Description))
		match.LastEvent = fmt.Sprintf("%s CARD: %s (%s)", req.CardColor, req.Description, time.Now().Format("15:04:05"))
	case "status_change":
		match.Status = req.Description
		match.LastEvent = fmt.Sprintf("STATUS: %s (%s)", req.Description, time.Now().Format("15:04:05"))
	default:
		match.LastEvent = eventDescription
	}

	if err := s.repo.UpdateMatch(ctx, match); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update match after event: %v", err)
	}

	event := &repository.Event{
		EventID:     fmt.Sprintf("evt-%s-%d", req.MatchId, time.Now().UnixNano()),
		MatchID:     req.MatchId,
		EventType:   req.EventType,
		Description: req.Description,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	if err := s.repo.AddEvent(ctx, event); err != nil {
		fmt.Printf("Warning: Failed to add event record for match %s: %v\n", req.MatchId, err)
	}

	// Trigger WebSocket update here!
	if s.websocketHub != nil {
		// Convert repository.Match to proto.MatchResponse for broadcasting
		protoMatch := &proto.MatchResponse{
			MatchId:    match.MatchID,
			Status:     match.Status,
			HomeScore:  match.HomeScore,
			AwayScore:  match.AwayScore,
			LastEvent:  match.LastEvent,
			Possession: match.Possession,
			Shots:      match.Shots,
			Fouls:      match.Fouls,
			Cards:      match.Cards,
		}
		s.websocketHub.BroadcastMatchUpdate(match.MatchID, protoMatch)
		fmt.Printf("WebSocket: Broadcasted update for match %s (admin event).\n", match.MatchID)
	}

	// ... (notification producer placeholder) ...

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
