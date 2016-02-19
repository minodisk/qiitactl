package info

import "encoding/json"

type Info struct {
	Version      string `json:"PackageVersion"`
	TaskSettings `json:"TaskSettings"`
}

type TaskSettings struct {
	GitHub `json:"publish-github"`
}

type GitHub struct {
	Name   string `json:"repository"`
	Author string `json:"owner"`
}

func New(bindata []byte) (info Info, err error) {
	err = json.Unmarshal(bindata, &info)
	if err != nil {
		return
	}
	return
}
