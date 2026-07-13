// API 层统一入口。
// 运行时检测 window.go 自动选择 Wails 绑定或 HTTP（axios）实现。
// 使用 Proxy 延迟检测，确保 Wails runtime 注入后再判断。

import { isWails, taskApi as wailsTaskApi, configApi as wailsConfigApi, ffmpegApi as wailsFfmpegApi, caApi as wailsCaApi } from './wails'
import { taskApi as httpTaskApi } from './task'
import { configApi as httpConfigApi, ffmpegApi as httpFfmpegApi, caApi as httpCaApi } from './config'

function createApi<T extends object>(wailsImpl: T, httpImpl: T): T {
  return new Proxy({} as T, {
    get(_, prop) {
      const api = isWails() ? wailsImpl : httpImpl
      const val = (api as any)[prop]
      return typeof val === 'function' ? val.bind(api) : val
    },
  })
}

export const taskApi = createApi(wailsTaskApi, httpTaskApi)
export const configApi = createApi(wailsConfigApi, httpConfigApi)
export const ffmpegApi = createApi(wailsFfmpegApi, httpFfmpegApi)
export const caApi = createApi(wailsCaApi, httpCaApi)

// 从原始模块重导出类型，保持类型导入路径兼容。
export type { Task, CreateTaskReq, UpdateTaskReq } from './task'
export type { Config, FfmpegStatus, CaStatus } from './config'
