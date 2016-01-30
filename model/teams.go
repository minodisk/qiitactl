package model

import (
	"encoding/json"

	"github.com/minodisk/qiitactl/api"
)

type Teams []Team

func (teams *Teams) Fetch(client api.Client) (err error) {
	body, err := client.Get("", "/teams", nil)
	if err != nil {
		return
	}
	json.Unmarshal(body, teams)
	return
}
