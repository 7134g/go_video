## Purpose

TBD — 任务执行队列与并发控制。

## Requirements

### Requirement: 任务执行队列
系统 SHALL 维护一个基于 channel 的任务执行队列，所有待执行的任务通过该队列调度执行。

#### Scenario: 任务入队
- **WHEN** 调用 StartTask、StartAll 或 AddAndStart 启动任务
- **THEN** 系统将任务及其回调发送到 taskQueue channel，由后台 dispatcher 统一调度

#### Scenario: 任务完成自动出队
- **WHEN** 任务执行完成（成功、失败或取消）
- **THEN** 系统释放该任务占用的并发槽位，下一个排队的任务自动开始执行

### Requirement: 并发限制
系统 SHALL 通过信号量 `taskSem` 限制同时运行的任务数不超过 `MaxConcurrentTasks`。

#### Scenario: 槽位未满
- **WHEN** dispatcher 从 taskQueue 取出任务，且当前运行任务数小于 MaxConcurrentTasks
- **THEN** dispatcher 获取槽位并立即启动任务 goroutine

#### Scenario: 槽位已满
- **WHEN** dispatcher 从 taskQueue 取出任务，但当前运行任务数已达到 MaxConcurrentTasks
- **THEN** dispatcher 阻塞等待，直到有运行中的任务完成释放槽位

### Requirement: 统一入队入口
系统 SHALL 确保所有任务启动路径（AddAndStart、StartTask、StartAll）都通过 taskQueue 入队，禁止直接启动 goroutine。

#### Scenario: AddAndStart 通过队列
- **WHEN** 调用 AddAndStart 一步完成添加和启动
- **THEN** 系统存储回调到 DTask 并将任务发送到 taskQueue，不绕过并发限制

#### Scenario: StartTask 通过队列
- **WHEN** 用户手动启动单个任务
- **THEN** 系统存储回调到 DTask 并将任务发送到 taskQueue

#### Scenario: StartAll 通过队列
- **WHEN** 用户批量启动所有任务
- **THEN** 系统依次将每个任务发送到 taskQueue

### Requirement: 已取消任务跳过执行
系统 SHALL 在启动任务下载前检查任务上下文是否已被取消。

#### Scenario: 队列中的任务被取消后出队
- **WHEN** dispatcher 取出一个已被 StopAll 取消的任务
- **THEN** dispatcher 检测到 ctx.Done() 已关闭，直接调用回调返回 context.Canceled 错误，不执行实际下载
