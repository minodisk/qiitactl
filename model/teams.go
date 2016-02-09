package model

import (
	"encoding/json"

	"github.com/minodisk/qiitactl/api"
)

// Teams is a collection of team.
type Teams []Team

// FetchTeams fetches teams that the authenticated user belongs.
func FetchTeams(client api.Client) (teams Teams, err error) {
	body, _, err := client.Get("", "/teams", nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &teams)
	return
}
