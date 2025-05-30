package sportradar

import (
	"context"
	"fmt"
	"time"

	"github.com/abaika-abay/live_sports_project/match-service/repository"
	"sync" // Import sync for mutex
)

// MockSportradarClient implements SportradarClientI for testing purposes.
type MockSportradarClient struct {
	liveData map[string]*repository.Match // Mock live data storage
	mu       sync.RWMutex
}

// NewMockSportradarClient creates a new mock Sportradar client.
func NewMockSportradarClient() *MockSportradarClient {
	return &MockSportradarClient{
		liveData: make(map[string]*repository.Match),
	}
}

// FetchMatchData simulates fetching real-time data for a specific match.
func (c *MockSportradarClient) FetchMatchData(ctx context.Context, matchID string) (*repository.Match, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	match, ok := c.liveData[matchID]
	if !ok {
		// If not found in mock, return a default/error, or simulate "not found"
		return nil, fmt.Errorf("match %s not found in mock Sportradar data", matchID)
	}

	// Simulate dynamic score/stat changes for live testing
	match.HomeScore = (match.HomeScore + int32(time.Now().Second()%2)) % 5 // Score changes slowly
	match.AwayScore = (match.AwayScore + int32(time.Now().Second()%2)) % 4
	match.Possession = int32(50 + (time.Now().Minute()%10 - 5)) // Possession fluctuates
	match.LastEvent = fmt.Sprintf("Mock update at %s", time.Now().Format("15:04:05"))

	// Return a copy to prevent external modification of internal mock state
	copiedMatch := *match
	return &copiedMatch, nil
}

// AddInitialMatchData is a helper for the mock to pre-populate some data.
// This is used by main.go for initial setup.
func (c *MockSportradarClient) AddInitialMatchData(match *repository.Match) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.liveData[match.MatchID] = match
	fmt.Printf("Mock Sportradar: Added initial data for %s\n", match.MatchID)
}

// UpdateData allows external simulation of Sportradar pushing updates.
// (Not part of the SportradarClientI interface, but useful for the mock)
func (c *MockSportradarClient) UpdateData(ctx context.Context, match *repository.Match) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.liveData[match.MatchID] = match // Overwrite with new state
	fmt.Printf("Mock Sportradar: Manually updated data for match %s\n", match.MatchID)
	return nil
}