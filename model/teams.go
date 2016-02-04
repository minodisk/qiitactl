package model

import (
	"encoding/json"

	"github.com/minodisk/qiitactl/api"
)

type Teams []Team

func FetchTeams(client api.Client) (teams Teams, err error) {
	body, _, err := client.Get("", "/teams", nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &teams)
	return
}
