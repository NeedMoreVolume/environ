package main_test

import (
	"bytes"
	"log/slog"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	cmd := exec.Command("./environ", `-input=C:\Users\pival\Workspace\Go\environ\testdata\config\config.go`, `-output=C:\Users\pival\Workspace\Go\environ\testdata\.env\`)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	slog.Info("output", "out", out.String())
}
