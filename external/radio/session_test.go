package radio

import (
	"testing"
)

func TestStationIDByURL_LocalStation(t *testing.T) {
	p := &Provider{
		stations: []station{
			{name: "Built-in", url: "http://builtin.example.com/stream"},
			{name: "User Station", url: "http://user.example.com/stream"},
		},
		favorites: &Favorites{byURL: make(map[string]struct{})},
	}

	t.Run("finds first station", func(t *testing.T) {
		id, ok := p.StationIDByURL("http://builtin.example.com/stream")
		if !ok {
			t.Fatal("expected to find station")
		}
		if id != "l:0" {
			t.Errorf("ID = %q, want l:0", id)
		}
	})

	t.Run("finds second station", func(t *testing.T) {
		id, ok := p.StationIDByURL("http://user.example.com/stream")
		if !ok {
			t.Fatal("expected to find station")
		}
		if id != "l:1" {
			t.Errorf("ID = %q, want l:1", id)
		}
	})

	t.Run("returns false for unknown URL", func(t *testing.T) {
		_, ok := p.StationIDByURL("http://unknown.example.com/stream")
		if ok {
			t.Error("expected not found for unknown URL")
		}
	})
}

func TestStationIDByURL_CatalogStation(t *testing.T) {
	p := &Provider{
		favorites: &Favorites{byURL: make(map[string]struct{})},
		catalog: []CatalogStation{
			{Name: "Jazz FM", URL: "http://jazz.example.com/stream"},
			{Name: "Rock Radio", URL: "http://rock.example.com/stream"},
		},
	}

	id, ok := p.StationIDByURL("http://rock.example.com/stream")
	if !ok {
		t.Fatal("expected to find catalog station")
	}
	if id != "c:1" {
		t.Errorf("ID = %q, want c:1", id)
	}
}

func TestStationIDByURL_FavoriteStation(t *testing.T) {
	p := &Provider{
		favorites: &Favorites{
			stations: []CatalogStation{
				{Name: "Fav Station", URL: "http://fav.example.com/stream"},
			},
			byURL: map[string]struct{}{
				"http://fav.example.com/stream": {},
			},
		},
	}

	id, ok := p.StationIDByURL("http://fav.example.com/stream")
	if !ok {
		t.Fatal("expected to find favorite station")
	}
	if id != "f:0" {
		t.Errorf("ID = %q, want f:0", id)
	}
}

func TestStationIDByURL_LocalTakesPriority(t *testing.T) {
	// A URL that exists in both local stations and catalog should return the local ID first.
	p := &Provider{
		stations: []station{
			{name: "Local Override", url: "http://shared.example.com/stream"},
		},
		favorites: &Favorites{byURL: make(map[string]struct{})},
		catalog: []CatalogStation{
			{Name: "Catalog Duplicate", URL: "http://shared.example.com/stream"},
		},
	}

	id, ok := p.StationIDByURL("http://shared.example.com/stream")
	if !ok {
		t.Fatal("expected to find station")
	}
	if id != "l:0" {
		t.Errorf("ID = %q, want l:0 (local takes priority)", id)
	}
}
