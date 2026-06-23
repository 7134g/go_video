## 1. Backend — Progress broadcast mechanism

- [x] 1.1 Add `taskID` and `taskName` fields to `Progress` struct in `internal/controller/dtask.go`
- [x] 1.2 Add `BroadcastProgress(info ProgressInfo)` / `AddProgressListener` / `RemoveProgressListener` in `internal/controller/broadcast.go`
- [x] 1.3 In `AddDone()` and `IncrementDone()`, release lock then call `BroadcastProgress` with computed `ProgressInfo`
- [x] 1.4 Set `Progress.taskID` and `Progress.taskName` when creating `DTask` in `controller.go` (`AddTask` / `AddAndStart`)
- [x] 1.5 In `AddAndStart()`, `StartTask()`, and `StartAll()`, call `BroadcastProgress` with initial `ProgressInfo` after the task starts

## 2. Backend — WebSocket handler changes

- [x] 2.1 Add `progressCh` listener channel in `ProgressWS`, register via `AddProgressListener`
- [x] 2.2 Send `GetAllProgress()` full snapshot once on connect
- [x] 2.3 Remove the 1s ticker; listen on `progressCh` for single `ProgressInfo` push

## 3. Frontend — Handle single progress push

- [x] 3.1 In `ws.onmessage` handler, detect single `ProgressInfo` object (non-array with `percent` property) and update only that task's progress

## 4. Cleanup and verify

- [x] 4.1 Build frontend and Go binary, verify no compilation errors (PASS)
- [x] 4.2 Test WebSocket connection: verify initial snapshot array is delivered
- [x] 4.3 Test progress push: verify single ProgressInfo is delivered when download progresses
- [x] 4.4 Verify frontend properly updates progress bar from single push
