package session

import (
	"os"
	"path/filepath"
	"testing"
)

// withTempHome sets HOME so appdir.Dir() points inside a temp directory,
// restoring the original on cleanup.
func withTempHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	return dir
}

func TestSaveLoadRoundTrip(t *testing.T) {
	withTempHome(t)

	Save("radio", "http://radio.example.com/stream")

	got := Load()
	if got.Provider != "radio" {
		t.Errorf("Provider = %q, want radio", got.Provider)
	}
	if got.ID != "http://radio.example.com/stream" {
		t.Errorf("ID = %q, want http://radio.example.com/stream", got.ID)
	}
}

func TestSaveIgnoresEmptyProvider(t *testing.T) {
	home := withTempHome(t)
	Save("", "http://radio.example.com/stream")

	f := filepath.Join(home, ".config", "cliamp", "session.json")
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Errorf("session.json should not exist for empty provider, got err=%v", err)
	}
}

func TestSaveIgnoresEmptyID(t *testing.T) {
	home := withTempHome(t)
	Save("radio", "")

	f := filepath.Join(home, ".config", "cliamp", "session.json")
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Errorf("session.json should not exist for empty id, got err=%v", err)
	}
}

func TestLoadMissingFileReturnsZero(t *testing.T) {
	withTempHome(t)
	got := Load()
	if got != (State{}) {
		t.Errorf("Load() = %+v, want zero State", got)
	}
}

func TestLoadCorruptFileReturnsZero(t *testing.T) {
	home := withTempHome(t)
	dir := filepath.Join(home, ".config", "cliamp")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "session.json"), []byte("not json {{"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got := Load()
	if got != (State{}) {
		t.Errorf("Load() = %+v, want zero State for corrupt file", got)
	}
}

func TestSaveCreatesParentDirectory(t *testing.T) {
	home := withTempHome(t)

	parent := filepath.Join(home, ".config", "cliamp")
	if _, err := os.Stat(parent); !os.IsNotExist(err) {
		t.Fatalf("precondition: parent should not exist, got err=%v", err)
	}

	Save("radio", "http://radio.example.com/stream")

	if info, err := os.Stat(parent); err != nil || !info.IsDir() {
		t.Errorf("Save should create parent directory, err=%v", err)
	}
}

func TestSaveOverwritesPrevious(t *testing.T) {
	withTempHome(t)

	Save("radio", "http://first.example.com/stream")
	Save("radio", "http://second.example.com/stream")

	got := Load()
	if got.Provider != "radio" || got.ID != "http://second.example.com/stream" {
		t.Errorf("Load() = %+v, want Provider=radio ID=http://second.example.com/stream", got)
	}
}
