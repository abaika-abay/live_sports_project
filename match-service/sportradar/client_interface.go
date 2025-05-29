package sportradar

import (
	"context"

	"github.com/abaika-abay/live_sports_project/match-service/repository"
)

// SportradarClientI defines the interface for interacting with the Sportradar API.
type SportradarClientI interface {
	FetchMatchData(ctx context.Context, matchID string) (*repository.Match, error)
	// You might add other methods like FetchLiveMatchesList, etc.
}
