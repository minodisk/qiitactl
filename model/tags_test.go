package model_test

import (
	"testing"

	"github.com/minodisk/qiitactl/model"
	"github.com/minodisk/qiitactl/testutil"

	"gopkg.in/yaml.v2"
)

func TestTagsMarshalYAMLWithEmptyTags(t *testing.T) {
	tags := model.Tags{}
	b, err := yaml.Marshal(tags)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := `[]
`
	if actual != expected {
		t.Errorf("wrong string:\n%s", testutil.Diff(expected, actual))
	}
}

func TestTagsMarshalYAMLWithATagAndNoVersion(t *testing.T) {
	tags := model.Tags{
		model.Tag{
			Name: "Go",
		},
	}
	b, err := yaml.Marshal(tags)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := `- Go
`
	if actual != expected {
		t.Errorf("wrong string:\n%s", testutil.Diff(expected, actual))
	}
}

func TestTagsMarshalYAMLWithATagAndAVersion(t *testing.T) {
	tags := model.Tags{
		model.Tag{
			Name: "Go",
			Versions: []string{
				"1.5.3",
			},
		},
	}
	b, err := yaml.Marshal(tags)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := `- Go:
  - 1.5.3
`
	if actual != expected {
		t.Errorf("wrong string:\n%s", testutil.Diff(expected, actual))
	}
}

func TestTagsMarshalYAMLWithATagAndVersions(t *testing.T) {
	tags := model.Tags{
		model.Tag{
			Name: "Go",
			Versions: []string{
				"1.5.3",
				"1.6.0",
			},
		},
	}
	b, err := yaml.Marshal(tags)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := `- Go:
  - 1.5.3
  - 1.6.0
`
	if actual != expected {
		t.Errorf("wrong string:\n%s", testutil.Diff(expected, actual))
	}
}

func TestTagsMarshalYAMLWithTagsAndVersions(t *testing.T) {
	tags := model.Tags{
		model.Tag{
			Name: "Go",
			Versions: []string{
				"1.5.3",
				"1.6.0",
			},
		},
		model.Tag{
			Name: "Ruby",
			Versions: []string{
				"2.3.0",
			},
		},
	}
	b, err := yaml.Marshal(tags)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := `- Go:
  - 1.5.3
  - 1.6.0
- Ruby:
  - 2.3.0
`
	if actual != expected {
		t.Errorf("wrong string:\n%s", testutil.Diff(expected, actual))
	}
}

func TestTagsUnmarshalYAMLWithEmptyTags(t *testing.T) {
	var tags model.Tags
	err := yaml.Unmarshal([]byte(`[]
`), &tags)
	if err != nil {
		t.Fatal(err)
	}

	if len(tags) != 0 {
		t.Fatalf("wrong length: expected %d, but actual %d", 0, len(tags))
	}
}

func TestTagsUnmrshalYAMLWithATagAndNoVersion(t *testing.T) {
	var tags model.Tags
	err := yaml.Unmarshal([]byte(`- Go
`), &tags)
	if err != nil {
		t.Fatal(err)
	}

	if len(tags) != 1 {
		t.Fatalf("wrong length: expected %d, but actual %d", 1, len(tags))
	}
	if tags[0].Name != "Go" {
		t.Errorf("wrong name: expected %d, but actual %d", "Go", tags[0].Name)
	}
	if len(tags[0].Versions) != 0 {
		t.Fatalf("wrong length: expected %d, but actual %d", 0, len(tags[0].Versions))
	}
}

func TestTagsUnmarshalYAMLWithATagAndAVersion(t *testing.T) {
	var tags model.Tags
	err := yaml.Unmarshal([]byte(`- Go:
  - 1.5.3
`), &tags)
	if err != nil {
		t.Fatal(err)
	}

	if len(tags) != 1 {
		t.Fatalf("wrong length: expected %d, but actual %d", 1, len(tags))
	}
	if tags[0].Name != "Go" {
		t.Errorf("wrong name: expected %d, but actual %d", "Go", tags[0].Name)
	}
	if len(tags[0].Versions) != 1 {
		t.Fatalf("wrong length: expected %d, but actual %d", 1, len(tags[0].Versions))
	}
	if tags[0].Versions[0] != "1.5.3" {
		t.Errorf("wrong version: expected %s, but actual %s", "1.5.3", tags[0].Versions[0])
	}
}

func TestTagsUnmarshalYAMLWithATagAndVersions(t *testing.T) {
	var tags model.Tags
	err := yaml.Unmarshal([]byte(`- Go:
  - 1.5.3
  - 1.6.0
`), &tags)
	if err != nil {
		t.Fatal(err)
	}

	if len(tags) != 1 {
		t.Fatalf("wrong length: expected %d, but actual %d", 1, len(tags))
	}
	if tags[0].Name != "Go" {
		t.Errorf("wrong name: expected %d, but actual %d", "Go", tags[0].Name)
	}
	if len(tags[0].Versions) != 2 {
		t.Fatalf("wrong length: expected %d, but actual %d", 2, len(tags[0].Versions))
	}
	if tags[0].Versions[0] != "1.5.3" {
		t.Errorf("wrong version: expected %s, but actual %s", "1.5.3", tags[0].Versions[0])
	}
	if tags[0].Versions[1] != "1.6.0" {
		t.Errorf("wrong version: expected %s, but actual %s", "1.6.0", tags[0].Versions[0])
	}
}

func TestTagsUnmarshalYAMLWithTagsAndVersions(t *testing.T) {
	var tags model.Tags
	err := yaml.Unmarshal([]byte(`- Go:
  - 1.5.3
  - 1.6.0
- Ruby:
  - 2.3.0
`), &tags)
	if err != nil {
		t.Fatal(err)
	}

	if len(tags) != 2 {
		t.Fatalf("wrong length: expected %d, but actual %d", 2, len(tags))
	}
	if tags[0].Name != "Go" {
		t.Errorf("wrong name: expected %d, but actual %d", "Go", tags[0].Name)
	}
	if len(tags[0].Versions) != 2 {
		t.Fatalf("wrong length: expected %d, but actual %d", 2, len(tags[0].Versions))
	}
	if tags[0].Versions[0] != "1.5.3" {
		t.Errorf("wrong version: expected %s, but actual %s", "1.5.3", tags[0].Versions[0])
	}
	if tags[0].Versions[1] != "1.6.0" {
		t.Errorf("wrong version: expected %s, but actual %s", "1.6.0", tags[0].Versions[0])
	}
	if tags[1].Name != "Ruby" {
		t.Errorf("wrong name: expected %d, but actual %d", "Go", tags[1].Name)
	}
	if len(tags[1].Versions) != 1 {
		t.Fatalf("wrong length: expected %d, but actual %d", 1, len(tags[1].Versions))
	}
	if tags[1].Versions[0] != "2.3.0" {
		t.Errorf("wrong version: expected %s, but actual %s", "2.3.0", tags[1].Versions[0])
	}
}
