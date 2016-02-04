package cli_test

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/cli"
)

var (
	client api.Client
)

func TestMain(m *testing.M) {
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client, err = api.NewClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestNewAppNoCommand(t *testing.T) {
	outBuf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, outBuf, errBuf)
	err := app.Run([]string{"qiitactl"})
	if err != nil {
		t.Error(err)
	}
	if len(errBuf.Bytes()) != 0 {
		t.Fatalf("error shouldn't occur: %s", errBuf.Bytes())
	}
	out := string(outBuf.Bytes())
	if !strings.HasPrefix(out, "NAME:") {
		t.Errorf("wrong output: %s", out)
	}
}

func TestNewAppUndefinedCommand(t *testing.T) {
	outBuf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, outBuf, errBuf)
	err := app.Run([]string{"qiitactl", "nop"})
	if err != nil {
		t.Error(err)
	}
	if string(errBuf.Bytes()) != "qiitactl: 'nop' is not a qiitactl command. See 'qiitactl --help'." {
		t.Fatalf("wrong error: %s", errBuf.Bytes())
	}
	if len(outBuf.Bytes()) != 0 {
		t.Errorf("shouldn't output: %s", outBuf.Bytes())
	}
}
