package info_test

import (
	"testing"

	"github.com/minodisk/qiitactl/info"
)

func TestInfo(t *testing.T) {
	i, err := info.New([]byte(`{
	"ArtifactsDest": "build",
	"Tasks": [
		"compile",
		"package"
	],
	"BuildConstraints": "linux darwin windows",
	"PackageVersion": "0.1.1",
	"TaskSettings": {
		"publish-github": {
			"body": "[Changes](https://github.com/minodisk/qiitactl/blob/master/CHANGELOG.md)",
			"owner": "minodisk",
			"repository": "qiitactl"
		}
	},
	"ConfigVersion": "0.9"
}`))
	if err != nil {
		t.Fatal(err)
	}
	if i.Version != "0.1.1" {
		t.Errorf("wrong version: %s", i.Version)
	}
	if i.Name != "qiitactl" {
		t.Errorf("wrong name: %s", i.Name)
	}
	if i.Author != "minodisk" {
		t.Errorf("wrong author: %s", i.Author)
	}
}

func TestInfoError(t *testing.T) {
	_, err := info.New([]byte("non JSON format"))
	if err == nil {
		t.Fatal("error should occur")
	}
}
