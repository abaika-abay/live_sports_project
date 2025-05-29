package service

import (
	"context"

	"github.com/abaika-abay/live_sports_project/common/pkg/db"
	"github.com/abaika-abay/live_sports_project/match-service/proto"
	"github.com/abaika-abay/live_sports_project/match-service/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MatchService struct {
	proto.UnimplementedMatchServiceServer
	repo *repository.MatchRepository
}

func NewMatchService(db *db.MongoDB) *MatchService {
	return &MatchService{
		repo: repository.NewMatchRepository(db),
	}
}

func (s *MatchService) GetMatchUpdates(ctx context.Context, req *proto.MatchRequest) (*proto.MatchResponse, error) {
	match, err := s.repo.GetMatch(ctx, req.MatchId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "match not found")
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

func (s *MatchService) CreateMatch(ctx context.Context, req *proto.CreateMatchRequest) (*proto.MatchResponse, error) {
	match := &repository.Match{
		MatchID:    req.MatchId,
		Status:     "Scheduled",
		HomeScore:  0,
		AwayScore:  0,
		LastEvent:  "Match Created",
		Possession: 50,
		Shots:      0,
		Fouls:      0,
		Cards:      []string{},
	}

	if err := s.repo.CreateMatch(ctx, match); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create match")
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
