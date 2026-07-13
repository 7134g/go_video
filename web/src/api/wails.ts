// Wails 绑定调用包装层。
// 将 window.go.main.App.<Method>() 调用包装为 axios 响应格式 { data }，
// 保持与 http taskApi/configApi 相同的调用签名。

interface WailsApp {
  AddTask(name: string, url: string, header: string, taskType: string): Promise<any>
  GetTasks(status: number): Promise<any[]>
  DeleteTask(id: number): Promise<void>
  UpdateTask(id: number, name: string, url: string, header: string, taskType: string): Promise<any>
  StartOneTask(id: number): Promise<void>
  PauseTask(id: number): Promise<void>
  RetryTask(id: number): Promise<void>
  StopAllTasks(): Promise<void>
  StartAllTasks(): Promise<number>
  RedownloadTask(id: number): Promise<void>
  GetConfig(): Promise<any>
  UpdateConfig(updates: Record<string, any>): Promise<any>
  GetFfmpegStatus(): Promise<{ exists: boolean; supported: boolean }>
  DownloadFfmpeg(): Promise<{ exists: boolean }>
  UpdateTaskTitle(id: number): Promise<any>
  GetAllProgress(): Promise<any[]>
}

function app(): WailsApp {
  return (window as any).go.main.App
}

export function isWails() {
  return typeof window !== 'undefined' && !!(window as any).go
}

import type { CreateTaskReq, UpdateTaskReq } from './task'
import type { Config } from './config'

export const taskApi = {
  list: async (status?: number) => {
    const data = await app().GetTasks(status ?? -1)
    return { data }
  },
  create: async (req: CreateTaskReq) => {
    const data = await app().AddTask(req.name, req.url, req.header ?? '', req.type)
    return { data }
  },
  update: async (req: UpdateTaskReq) => {
    const data = await app().UpdateTask(req.id, req.name ?? '', req.url ?? '', req.header ?? '', req.type ?? '')
    return { data }
  },
  delete: async (id: number) => {
    await app().DeleteTask(id)
    return {}
  },
  start: async () => {
    const started = await app().StartAllTasks()
    return { data: { started } }
  },
  pause: async (id: number) => {
    await app().PauseTask(id)
    return {}
  },
  retry: async (id: number) => {
    await app().RetryTask(id)
    return {}
  },
  startOne: async (id: number) => {
    await app().StartOneTask(id)
    return {}
  },
  stopAll: async () => {
    await app().StopAllTasks()
    return {}
  },
  updateTitle: async (id: number) => {
    const data = await app().UpdateTaskTitle(id)
    return { data }
  },
  redownload: async (id: number) => {
    await app().RedownloadTask(id)
    return {}
  },
}

export const configApi = {
  get: async () => {
    const data = await app().GetConfig()
    return { data }
  },
  update: async (data: Partial<Config>) => {
    const result = await app().UpdateConfig(data)
    return { data: result }
  },
}

export const ffmpegApi = {
  status: async () => {
    const data = await app().GetFfmpegStatus()
    return { data }
  },
  download: async () => {
    const data = await app().DownloadFfmpeg()
    return { data }
  },
}
