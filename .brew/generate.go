package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Release struct {
	Meta
	Assets []Asset `json:"assets"`
}

type Meta struct {
	TagName string `json:"tag_name"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type DarwinAMD64 struct {
	Meta
	Asset  Asset
	Sha256 string
}

func main() {
	res, err := http.Get("https://api.github.com/repos/minodisk/qiitactl/releases/latest")
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var release Release
	err = json.Unmarshal(b, &release)
	if err != nil {
		log.Fatal(err)
	}

	darwin, err := findDarwinAMD64(release)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", darwin)

	tmpl, err := template.ParseFiles(".brew/brew.rb.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("homebrew-qiitactl/qiitactl.rb")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(f, darwin)
	if err != nil {
		log.Fatal(err)
	}

}

func findDarwinAMD64(release Release) (darwin DarwinAMD64, err error) {
	darwin.Meta = release.Meta
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, "_darwin_amd64.zip") {
			darwin.Asset = asset
			resp, err := http.Get(asset.BrowserDownloadURL)
			if err != nil {
				return darwin, err
			}
			defer resp.Body.Close()
			hasher := sha256.New()
			_, err = io.Copy(hasher, resp.Body)
			if err != nil {
				return darwin, err
			}
			darwin.Sha256 = hex.EncodeToString(hasher.Sum(nil))
			return darwin, nil
		}
	}
	err = errors.New("not found")
	return
}
