import request from './request'

export interface Config {
  max_concurrent_tasks: number
  max_segment_workers: number
  download_dir: string
  max_consecutive_errors: number
  default_headers: Record<string, string>
  interceptor_enabled: boolean
  agent_address: string
    vpn_address: string
}

export const configApi = {
  get: () => request.get<Config>('/config'),
  update: (data: Partial<Config>) => request.put<Config>('/config', data),
}

export interface FfmpegStatus {
  exists: boolean
  supported: boolean
}

export const ffmpegApi = {
  status: () => request.get<FfmpegStatus>('/ffmpeg/status'),
  // 下载耗时较长（数十 MB），放宽超时到 5 分钟。
  download: () => request.post<{ exists: boolean }>('/ffmpeg/download', null, { timeout: 300000 }),
}
