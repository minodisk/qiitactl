package model

import "encoding/json"

type Teams []Team

func (teams *Teams) Fetch(client Client) (err error) {
	body, err := client.Get("", "/teams", nil)
	if err != nil {
		return
	}
	json.Unmarshal(body, teams)
	return
}
