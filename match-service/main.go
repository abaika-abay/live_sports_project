package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http" // New import for HTTP server
	"slices"
	"time"

	"google.golang.org/grpc"

	"github.com/abaika-abay/live_sports_project/common/pkg/config"
	"github.com/abaika-abay/live_sports_project/common/pkg/db"
	"github.com/abaika-abay/live_sports_project/match-service/proto"
	"github.com/abaika-abay/live_sports_project/match-service/repository"
	"github.com/abaika-abay/live_sports_project/match-service/service"
	"github.com/abaika-abay/live_sports_project/match-service/sportradar"
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
	srClient := sportradar.NewMockSportradarClient() // srClient is now *sportradar.MockSportradarClient

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
	if err := repo.UpdateMatch(ctx, initialMatch); err != nil {
		fmt.Printf("Warning: Failed to ensure initial match %s exists in DB: %v\n", initialMatchID, err)
	} else {
		fmt.Printf("Ensured initial match data for: %s in DB\n", initialMatchID)
	}

	// Directly call AddInitialMatchData on srClient
	// No type assertion needed because srClient is already *sportradar.MockSportradarClient
	srClient.AddInitialMatchData(initialMatch) // <--- FIXED LINE
	fmt.Printf("Initialized mock Sportradar data for: %s\n", initialMatchID)

	// --- Start WebSocket setup (Part 2) ---
	websocketHub := service.NewWebSocketHub()
	go websocketHub.Run() // Run the hub's goroutine

	// Create a new HTTP server for WebSockets (often on a different port)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", websocketHub.HandleConnections)
	go func() {
		log.Printf("WebSocket server starting on %s", c.WebSocketPort)
		if err := http.ListenAndServe(c.WebSocketPort, mux); err != nil {
			log.Fatalf("WebSocket server failed to start: %v", err)
		}
	}()
	// --- End WebSocket setup ---

	matchService := service.NewMatchService(dbHandler, srClient, websocketHub) // Pass WebSocket hub here

	// --- Start Background Polling (Part 3) ---
	go func() {
		pollInterval := 5 * time.Second // Poll Sportradar every 5 seconds
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		log.Printf("Starting Sportradar background polling every %s for %s", pollInterval, initialMatchID) // Polling for fixed match for demo
		for range ticker.C {
			pollCtx, cancelPoll := context.WithTimeout(context.Background(), 3*time.Second) // Shorter timeout for polling
			matchIDToPoll := initialMatchID                                                 // Or iterate over all live matches in DB
			currentMatch, err := repo.GetMatch(pollCtx, matchIDToPoll)                      // Get current state from DB
			if err != nil {
				log.Printf("Polling: Could not get match %s from DB: %v", matchIDToPoll, err)
				cancelPoll()
				continue
			}

			srUpdate, err := srClient.FetchMatchData(pollCtx, matchIDToPoll)
			if err != nil {
				log.Printf("Polling: Failed to fetch real-time data for %s from Sportradar: %v", matchIDToPoll, err)
				cancelPoll()
				continue
			}

			// Check if there are actual changes before updating and broadcasting
			if srUpdate.HomeScore != currentMatch.HomeScore || srUpdate.AwayScore != currentMatch.AwayScore ||
				srUpdate.Status != currentMatch.Status || srUpdate.LastEvent != currentMatch.LastEvent ||
				srUpdate.Possession != currentMatch.Possession || srUpdate.Shots != currentMatch.Shots ||
				srUpdate.Fouls != currentMatch.Fouls || !slices.Equal(srUpdate.Cards, currentMatch.Cards) { // Use slices.Equal for slices
				// Update current match object with new SR data
				currentMatch.HomeScore = srUpdate.HomeScore
				currentMatch.AwayScore = srUpdate.AwayScore
				currentMatch.Status = srUpdate.Status
				currentMatch.LastEvent = srUpdate.LastEvent
				currentMatch.Possession = srUpdate.Possession
				currentMatch.Shots = srUpdate.Shots
				currentMatch.Fouls = srUpdate.Fouls
				currentMatch.Cards = srUpdate.Cards // Deep copy if needed

				// Update database
				if err := repo.UpdateMatch(pollCtx, currentMatch); err != nil {
					log.Printf("Polling: Failed to update DB for match %s: %v", matchIDToPoll, err)
					cancelPoll()
					continue
				}

				// Broadcast update via WebSockets
				protoMatch := &proto.MatchResponse{
					MatchId:    currentMatch.MatchID,
					Status:     currentMatch.Status,
					HomeScore:  currentMatch.HomeScore,
					AwayScore:  currentMatch.AwayScore,
					LastEvent:  currentMatch.LastEvent,
					Possession: currentMatch.Possession,
					Shots:      currentMatch.Shots,
					Fouls:      currentMatch.Fouls,
					Cards:      currentMatch.Cards,
				}
				websocketHub.BroadcastMatchUpdate(currentMatch.MatchID, protoMatch)
				log.Printf("Polling: Broadcasted live update for match %s (score %d-%d).\n", currentMatch.MatchID, currentMatch.HomeScore, currentMatch.AwayScore)
			}
			cancelPoll()
		}
	}()
	// --- End Background Polling ---

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterMatchServiceServer(s, matchService)

	log.Printf("Match Service gRPC listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
