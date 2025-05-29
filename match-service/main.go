package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time" // Added for initial match data

	"google.golang.org/grpc"

	"github.com/abaika-abay/live_sports_project/common/pkg/config"
	"github.com/abaika-abay/live_sports_project/common/pkg/db"
	"github.com/abaika-abay/live_sports_project/match-service/proto"
	"github.com/abaika-abay/live_sports_project/match-service/repository"
	"github.com/abaika-abay/live_sports_project/match-service/service"
	"github.com/abaika-abay/live_sports_project/match-service/sportradar" // Import sportradar
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbHandler, err := db.InitMongoDB(c.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := dbHandler.Client.Disconnect(context.Background()); err != nil {
			log.Fatalf("failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Initialize Sportradar Client
	srClient := sportradar.NewSportradarClient()

	// --- Optional: Add some initial mock match data to Sportradar and MongoDB ---
	// You might automate this with a migration script or a dedicated admin tool
	repo := repository.NewMatchRepository(dbHandler)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	initialMatchID := "match-123"
	initialMatch := &repository.Match{
		MatchID:    initialMatchID,
		HomeTeam:   "Real Madrid",
		AwayTeam:   "Barcelona",
		StartTime:  time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		Status:     "Scheduled",
		HomeScore:  0,
		AwayScore:  0,
		LastEvent:  "Match scheduled",
		Possession: 50,
		Shots:      0,
		Fouls:      0,
		Cards:      []string{},
	}
	err = repo.CreateMatch(ctx, initialMatch)
	if err != nil {
		fmt.Printf("Warning: Failed to create initial match %s in DB (might already exist): %v\n", initialMatchID, err)
	}
	srClient.AddInitialMatchData(initialMatch) // Add to mock Sportradar
	fmt.Printf("Initialized mock match data for: %s\n", initialMatchID)
	// --- End of Optional Initial Data ---

	matchService := service.NewMatchService(dbHandler, srClient) // Pass Sportradar client here

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterMatchServiceServer(s, matchService)

	log.Printf("Match Service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
