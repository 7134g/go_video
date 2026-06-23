# Progress Push

## Purpose

向 WebSocket 客户端实时推送单个任务的进度更新，替代定时轮询全量快照，降低延迟和带宽消耗。

## Requirements

### Requirement: Progress push on counter change

When `Progress.AddDone()` or `Progress.IncrementDone()` is called and updates the `Done` counter, the system SHALL immediately broadcast a `ProgressInfo` message containing the updated progress values of that specific task.

#### Scenario: AddDone triggers progress push
- **WHEN** `Progress.AddDone(n)` is called on a task with `Done=0, Total=100`
- **THEN** a `ProgressInfo` JSON object SHALL be broadcast with `percent=50` (after adding 50)
- **AND** the push SHALL contain `id`, `name`, `type`, `done`, `total`, and `percent` fields matching the task's current state

#### Scenario: IncrementDone triggers progress push
- **WHEN** `Progress.IncrementDone()` is called on a task with `Done=5, Total=10`
- **THEN** a `ProgressInfo` SHALL be broadcast with `done=6, percent=60`
- **AND** the push SHALL contain the updated values from the increment

#### Scenario: No push on SetTotal
- **WHEN** `Progress.SetTotal(total)` is called
- **THEN** no `ProgressInfo` SHALL be broadcast (only Done changes trigger pushes)

#### Scenario: Zero total does not panic
- **WHEN** `Progress.AddDone(1)` is called on a task with `Total=0`
- **THEN** the broadcast SHALL contain `percent=0` (no division by zero)
- **AND** no panic SHALL occur

### Requirement: Progress push on task start

When a task is started via `AddAndStart`, `StartTask`, or `StartAll`, the system SHALL broadcast a `ProgressInfo` for that task with its initial state (Done=0, Total=0, Percent=0).

#### Scenario: AddAndStart pushes initial ProgressInfo
- **WHEN** `AddAndStart()` is called for a new task
- **THEN** a `ProgressInfo` SHALL be broadcast with `done=0`, `total=0`, `percent=0`
- **AND** the push SHALL contain the task's `id`, `name`, and `type`

#### Scenario: StartTask pushes ProgressInfo
- **WHEN** `StartTask(id)` is called for an existing task
- **THEN** a `ProgressInfo` SHALL be broadcast for that task

#### Scenario: StartAll pushes ProgressInfo for each task
- **WHEN** `StartAll()` is called with 3 tasks
- **THEN** 3 `ProgressInfo` SHALL be broadcast, one per task

### Requirement: WebSocket receives progress pushes

The WebSocket handler SHALL receive individual `ProgressInfo` pushes on a dedicated channel and write them as single JSON object messages.

#### Scenario: Connection receives initial snapshot
- **WHEN** a WebSocket client connects to `/api/tasks/progress`
- **THEN** the server SHALL immediately write a JSON array of all task `ProgressInfo` (one-shot initial state)
- **AND** the server SHALL NOT use a periodic ticker for progress

#### Scenario: Progress push delivered to connected client
- **WHEN** a progress push is broadcast while a WebSocket client is connected
- **THEN** the client SHALL receive the `ProgressInfo` as a single JSON object message

#### Scenario: Multiple connected clients receive pushes
- **WHEN** a progress push is broadcast and multiple WebSocket clients are connected
- **THEN** all connected clients SHALL receive the push

#### Scenario: Non-blocking drop protects broadcaster
- **WHEN** a client's progress channel buffer is full
- **THEN** the push SHALL be dropped for that client (non-blocking)
- **AND** other clients SHALL still receive the push

### Requirement: Frontend handles single progress object

The Vue frontend SHALL distinguish between the initial array snapshot and individual `ProgressInfo` objects, updating only the affected task.

#### Scenario: Single ProgressInfo updates one task
- **WHEN** a single `ProgressInfo` JSON object is received via WebSocket (not an array)
- **THEN** the frontend SHALL update only the `progress[data.id]` value for that specific task
- **AND** SHALL NOT replace the entire `taskProgressList`

#### Scenario: Array snapshot replaces full list
- **WHEN** a JSON array is received via WebSocket
- **THEN** the frontend SHALL replace `taskProgressList` entirely (existing behavior unchanged)
- **AND** SHALL update all `progress[item.id]` values
