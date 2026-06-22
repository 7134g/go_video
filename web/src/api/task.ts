import request from './request'

export interface Task {
  id: number
  name: string
  url: string
  header: string
  type: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateTaskReq {
  name: string
  url: string
  header?: string
  type: string
}

export interface UpdateTaskReq extends Partial<CreateTaskReq> {
  id: number
}

export const taskApi = {
  list: (status?: number) => request.get<Task[]>('/tasks', { params: status !== undefined ? { status } : {} }),
  create: (data: CreateTaskReq) => request.post<Task>('/tasks', data),
  update: (data: UpdateTaskReq) => request.post<Task>('/tasks/update', data),
  delete: (id: number) => request.post('/tasks/delete', { id }),
  start: () => request.post<{ started: number }>('/tasks/start'),
  pause: (id: number) => request.post('/tasks/pause', { id }),
  retry: (id: number) => request.post('/tasks/retry', { id }),
  startOne: (id: number) => request.post('/tasks/start-one', { id }),
  stopAll: () => request.post('/tasks/stop-all'),
  updateTitle: (id: number) => request.post<Task>('/tasks/update-title', { id }),
  redownload: (id: number) => request.post('/tasks/redownload', { id }),
}
