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
