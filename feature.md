# Session Restore Feature (cliamp)

## Goal
- Restore last listening session when cliamp is reopened
- Auto-resume playback for the last active provider

## Primary Use Case
- Radio provider
  - Save last selected station (URL)
  - On restart: auto-select and auto-play that station

## Existing System (NOT reused)
- `internal/resume`
  - Handles track position within a song/stream
  - Uses `resume.json`
- New feature is separate:
  - Session = what was being listened to
  - Resume = where in the track you were

## New Persistent State
- Stored in:
  - `~/.config/cliamp/session.json`
- Contains:
  - `provider` (e.g. "radio")
  - `id` (e.g. radio station URL)

## Config Toggle
```toml
restore_last_session = true


# Optional Built-in Radio Playlist (cliamp)

## Goal
- Make the built-in “cliamp radio” station optional via configuration
- Allow users to disable bundled default radio streams

## Current Behavior
- cliamp always injects a built-in radio station:
  - Name: `cliamp radio`
  - URL: `https://radio.cliamp.stream/streams.m3u`
- This happens inside the radio provider at runtime
- It is not configurable today

## Proposed Change
- Add a config flag to enable/disable built-in radio station

## Config Option
```toml
[radio]
enable_builtin = true