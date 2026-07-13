// API 层统一入口。
// 运行时检测 window.go 自动选择 Wails 绑定或 HTTP（axios）实现。

import { isWails, taskApi as wailsTaskApi, configApi as wailsConfigApi, ffmpegApi as wailsFfmpegApi } from './wails'
import { taskApi as httpTaskApi } from './task'
import { configApi as httpConfigApi, ffmpegApi as httpFfmpegApi } from './config'

export const taskApi = isWails() ? wailsTaskApi : httpTaskApi
export const configApi = isWails() ? wailsConfigApi : httpConfigApi
export const ffmpegApi = isWails() ? wailsFfmpegApi : httpFfmpegApi

// 从原始模块重导出类型，保持类型导入路径兼容。
export type { Task, CreateTaskReq, UpdateTaskReq } from './task'
export type { Config, FfmpegStatus } from './config'
