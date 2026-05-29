# REASONIX.md

## Stack

- **Go 1.24** — module `go_video`
- **Gin v1.9.1** — HTTP framework
- **GORM + SQLite** — task persistence
- **Vue 3 + Element Plus + Vite** — frontend SPA at `web/`
- **google/martian** — HTTPS MITM proxy for browser interception
- **gorilla/websocket** — real-time download progress push

## Layout

- `main.go` — app entrypoint; embeds `web/dist`, creates Gin routes, starts MITM proxy
- `cmd/proxy/main.go` — standalone CA certificate installer binary
- `internal/api/` — HTTP handlers for tasks CRUD, config, WebSocket
- `internal/service/` — TaskService / ConfigService (orchestration layer)
- `internal/controller/` — DownloadController (m3u8/mp4 download, progress, broadcast)
- `internal/repository/` — DB (GORM/SQLite) for tasks, JSON file for config
- `internal/downloader/` — semaphore-based concurrency pool per domain
- `internal/model/` — Task / Config structs
- `pkg/m3u8/` — M3U8 parser, AES-128-CBC decryption, ffmpeg merge
- `pkg/proxy/` — HTTPS MITM proxy, video URL detection, request channel
- `web/` — Vue 3 + TypeScript + Vite frontend

## Commands

| Action | Command |
|---|---|
| Build all | `cd web && npm run build && go build -o build/go_video.exe` |
| Build cert installer | `go build -o build/install_cert.exe cmd/proxy/main.go` |
| Run | `go build . && ./go_video` |
| Frontend dev | `cd web && npm run dev` |
| Test | `go test ./...` |
| Build (script) | `bash build.sh` |

## Conventions

- Frontend built with `vue-tsc -b && vite build` (type-checked)
- Go tests only in `pkg/m3u8/` (`crypto_test.go`, `m3u8_test.go`)
- No lint config (no `.golangci.yml`, no Makefile)
- Config stored in `config.json` at project root
- Task DB is `video.db` (SQLite), auto-migrated on startup

## Watch out for

- **CA cert required.** `InitCa()` panics if not installed — run `install_cert.exe` first
- **`web/dist/` is generated.** Edit files in `web/src/`, not `web/dist/`
- **`ffmpeg.exe` must be at project root** for m3u8 segment merging
- **`config.json` and `video.db` are runtime files** at project root, not in `.gitignore` for `config.json` (it is gitignored per `.gitignore`)
- **MITM proxy** listens on `127.0.0.1:9999` by default; upstream proxy at `127.0.0.1:7890`
