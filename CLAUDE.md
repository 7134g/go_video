# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build / Run

```bash
go build . && ./go_video          # main app (serves web UI on :8080)
go build ./cmd/proxy && ./proxy   # CA cert installer (requires admin)
go test ./...                     # run all tests
```

The frontend is built separately (`web/` — Vue + Element Plus + Vite) and embedded via `//go:embed web/dist`. Build the frontend with `npm run build` from `web/` before building the Go binary.

## Architecture

This is a video downloader that handles **m3u8/HLS streams** (including AES-128 encrypted) and **direct MP4 downloads**. It includes an HTTPS MITM proxy that intercepts browser traffic to auto-capture video URLs.

### Layer flow

```
HTTP API (internal/api) → Service (internal/service) → Controller (internal/controller) → Downloader (internal/downloader)
```

- **`internal/api/`** — Gin HTTP handlers: task CRUD, config get/update, WebSocket progress push.
- **`internal/service/`** — `TaskService` bridges API and controller (lifecycle + header merging). `ConfigService` manages app config and the MITM proxy lifecycle.
- **`internal/controller/`** — Singleton `DownloadController` owns the in-memory task map, dispatches m3u8/mp4 downloads, tracks progress, and fans out broadcast messages to WebSocket listeners.
- **`internal/repository/`** — Task persistence via GORM/SQLite, config persistence via JSON file (`config.json`).
- **`internal/downloader/`** — `Pool` limits concurrent downloads per domain using a semaphore-based `Group`.
- **`pkg/m3u8/`** — M3U8 playlist parser (supports master playlists, AES-128 keys, byte ranges), AES-128-CBC decrypt, IV derivation, key cache, and ffmpeg concat merge.
- **`pkg/proxy/`** — HTTPS MITM proxy using `github.com/google/martian`. `Server.ModifyRequest` intercepts every request; `VideoDetector` matches `.m3u8`/`.mp4` URLs; captured tasks flow through a channel to `ConfigService.doTask` which persists them.

### Startup sequence

1. `repository.InitDB()` — opens SQLite and auto-migrates `Task` table
2. `InitCa()` — verifies the CA cert is installed in the system trust store (panics if not)
3. Config loaded from `config.json`, defaults applied to `DownloadController`
4. If `interceptor_enabled`, starts the MITM proxy in a goroutine
5. Gin server starts on `:8080`, serving the embedded Vue SPA and API routes

### Key behaviors

- **Header merging**: task-specific headers override config defaults; both are JSON-marshaled as `http.Header` (map[string][]string format).
- **MP4 resume**: uses `Range` header to resume incomplete downloads based on local file size.
- **M3U8 resume**: skips segment files that already exist on disk.
- **Concurrency**: task-level concurrency controlled by `MaxConcurrentTasks`; per-segment concurrency limited per-domain by `MaxSegmentWorkers`.
- **Proxy task dedup**: `doTask` checks if a URL already exists before inserting; updates name/header if the existing task is newer.
- **The `pkg/m3u8/m3u8.go` file** contains an older local `parse()` function; the one actually used is `ParseM3u8Data` in `method.go`.
