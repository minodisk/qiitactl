package model

// Team is a team data of Qiita:Team.
type Team struct {
	Active bool   `json:"active"` // チームが利用可能な状態かどうか
	ID     string `json:"id"`     // チームの一意なID
	Name   string `json:"name"`   // チームに設定されている名前を表します。
}
