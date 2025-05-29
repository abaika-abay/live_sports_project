package sportradar

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/abaika-abay/live_sports_project/match-service/repository" // Assuming your Match struct is here or define a new one
)

// SportradarClient simulates an external Sportradar API.
// In a real application, this would make HTTP requests.
type SportradarClient struct {
	liveData map[string]*repository.Match // Mock live data storage
	mu       sync.RWMutex
}

// NewSportradarClient creates a new mock Sportradar client.
func NewSportradarClient() *SportradarClient {
	return &SportradarClient{
		liveData: make(map[string]*repository.Match),
	}
}

// FetchMatchData simulates fetching real-time data for a specific match.
func (c *SportradarClient) FetchMatchData(ctx context.Context, matchID string) (*repository.Match, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	match, ok := c.liveData[matchID]
	if !ok {
		return nil, errors.New("match not found in Sportradar simulation")
	}
	return match, nil
}

// UpdateData simulates an internal process updating Sportradar's data.
// In a real scenario, Sportradar pushes updates, or you poll.
// For this mock, we allow manual updates.
func (c *SportradarClient) UpdateData(ctx context.Context, match *repository.Match) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Simulate some dynamic changes
	// This would be much more sophisticated in a real system
	if existingMatch, ok := c.liveData[match.MatchID]; ok {
		existingMatch.Status = match.Status
		existingMatch.HomeScore = match.HomeScore
		existingMatch.AwayScore = match.AwayScore
		existingMatch.LastEvent = match.LastEvent
		existingMatch.Possession = match.Possession
		existingMatch.Shots = match.Shots
		existingMatch.Fouls = match.Fouls
		existingMatch.Cards = match.Cards
		// Simulate slight changes over time if needed for demo
		existingMatch.Possession = (existingMatch.Possession + 1) % 100 // Example
	} else {
		c.liveData[match.MatchID] = match
	}
	fmt.Printf("Sportradar Mock: Updated live data for match %s at %s\n", match.MatchID, time.Now().Format(time.Kitchen))
	return nil
}

// AddInitialMatchData is a helper for the mock to pre-populate some data
func (c *SportradarClient) AddInitialMatchData(match *repository.Match) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.liveData[match.MatchID] = match
}
