package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigRestoreLastSession(t *testing.T) {
	cfg := defaultConfig()
	if cfg.RestoreLastSession {
		t.Error("RestoreLastSession should be false by default")
	}
}

func TestLoadRestoreLastSessionTrue(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	path := filepath.Join(os.Getenv("HOME"), ".config", "cliamp", "config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte("restore_last_session = true\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.RestoreLastSession {
		t.Error("RestoreLastSession should be true")
	}
}

func TestLoadRadioEnableBuiltinTrue(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	path := filepath.Join(os.Getenv("HOME"), ".config", "cliamp", "config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte("[radio]\nenable_builtin = true\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.Radio.EnableBuiltin {
		t.Error("Radio.EnableBuiltin should be true")
	}
}

func TestLoadRadioEnableBuiltinFalse(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	path := filepath.Join(os.Getenv("HOME"), ".config", "cliamp", "config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte("[radio]\nenable_builtin = false\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Radio.EnableBuiltin {
		t.Error("Radio.EnableBuiltin should be false")
	}
}

func TestLoadRadioSectionAfterSpotify(t *testing.T) {
	// Reproduce the user's config: [radio] at the end after [spotify].
	t.Setenv("HOME", t.TempDir())

	path := filepath.Join(os.Getenv("HOME"), ".config", "cliamp", "config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	content := "restore_last_session = true\n" +
		"speed = 1.00\n" +
		"eq_preset = \"Electronic\"\n" +
		"[spotify]\n" +
		"client_id = \"abc123\"\n" +
		"bitrate = 320\n" +
		"[radio]\n" +
		"enable_builtin = false\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.RestoreLastSession {
		t.Error("RestoreLastSession should be true")
	}
	if cfg.Spotify.ClientID != "abc123" {
		t.Errorf("Spotify.ClientID = %q, want abc123", cfg.Spotify.ClientID)
	}
	if cfg.Radio.EnableBuiltin {
		t.Error("Radio.EnableBuiltin should be false (set after [spotify] section)")
	}
}

func TestLoadRestoreLastSessionFalseAndDefault(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{name: "explicit false", content: "restore_last_session = false\n", expected: false},
		{name: "not set", content: "volume = -5\n", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HOME", t.TempDir())

			path := filepath.Join(os.Getenv("HOME"), ".config", "cliamp", "config.toml")
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatalf("MkdirAll: %v", err)
			}
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("WriteFile: %v", err)
			}

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load: %v", err)
			}
			if cfg.RestoreLastSession != tt.expected {
				t.Errorf("RestoreLastSession = %v, want %v", cfg.RestoreLastSession, tt.expected)
			}
		})
	}
}
