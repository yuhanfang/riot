package esports

import (
	"context"
	"fmt"
	"net/http"
)

func (c Client) GetHighlanderMatchDetails(ctx context.Context, tournamentID string, matchID string) (*HighlanderMatchDetails, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.lolesports.com/api/v2/highlanderMatchDetails?tournamentId=%s&matchId=%s", tournamentID, matchID), nil)
	if err != nil {
		return nil, err
	}
	var res HighlanderMatchDetails
	_, err = c.doJSON(ctx, req, &res)
	return &res, err
}

type HighlanderMatchDetails struct {
	GameIDMappings []HighlanderMatchDetails_GameIDMapping
}

type HighlanderMatchDetails_GameIDMapping struct {
	ID       string
	GameHash string
}
