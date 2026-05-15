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
