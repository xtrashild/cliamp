// Package session persists the last listening session (provider + track ID)
// so the user can resume on the next launch, separate from the per-track
// position resume handled by internal/resume.
package session

import (
	"encoding/json"
	"os"
	"path/filepath"

	"cliamp/internal/appdir"
)

// State holds enough information to restore the last listening session.
// Provider identifies the provider key (e.g. "radio", "spotify").
// ID is a provider-specific identifier: for radio it is the station URL,
// for other providers it is the playlist ID.
type State struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
}

func stateFile() (string, error) {
	dir, err := appdir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "session.json"), nil
}

// Save writes the session state to disk. No-ops for empty provider or id
// to avoid overwriting a valid file with useless data.
// Errors are silently ignored so a failed write never disrupts normal exit.
func Save(provider, id string) {
	if provider == "" || id == "" {
		return
	}
	f, err := stateFile()
	if err != nil {
		return
	}
	data, err := json.Marshal(State{Provider: provider, ID: id})
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(f), 0o755)
	_ = os.WriteFile(f, data, 0o600)
}

// Load reads the session state from disk. Returns a zero State if the file
// does not exist or cannot be parsed.
func Load() State {
	f, err := stateFile()
	if err != nil {
		return State{}
	}
	data, err := os.ReadFile(f)
	if err != nil {
		return State{}
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}
	}
	return s
}
