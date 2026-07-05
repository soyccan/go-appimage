package helpers_test

import (
	"testing"

	"github.com/probonopd/go-appimage/internal/helpers"
)

func TestValidStartupWMClass(t *testing.T) {
	valid := []string{"Heptabase", "org.gnome.Terminal", " DB Browser "}
	for _, value := range valid {
		if !helpers.ValidStartupWMClass(value) {
			t.Fatalf("expected %q to be valid", value)
		}
	}

	invalid := []string{"", "   ", "undefined", "UNDEFINED"}
	for _, value := range invalid {
		if helpers.ValidStartupWMClass(value) {
			t.Fatalf("expected %q to be invalid", value)
		}
	}
}

func TestExtractJSONStringField(t *testing.T) {
	data := []byte(`{"name":"foo","desktopName":"Heptabase","productName":"Ignored"}`)
	if got := helpers.ExtractJSONStringField(data, "desktopName"); got != "Heptabase" {
		t.Fatalf("desktopName = %q, want Heptabase", got)
	}
	if got := helpers.ExtractJSONStringField(data, "missing"); got != "" {
		t.Fatalf("missing = %q, want empty", got)
	}
}

func TestParseAppRunBinary(t *testing.T) {
	content := `#!/bin/bash
BIN="$APPDIR/project-meta"
exec "$BIN"
`
	if got := helpers.ParseAppRunBinary(content); got != "project-meta" {
		t.Fatalf("ParseAppRunBinary() = %q, want project-meta", got)
	}
}
