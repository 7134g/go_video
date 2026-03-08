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

export const taskApi = {
  list: (status?: number) => request.get<Task[]>('/tasks', { params: status !== undefined ? { status } : {} }),
  create: (data: CreateTaskReq) => request.post<Task>('/tasks', data),
  update: (id: number, data: Partial<CreateTaskReq>) => request.put<Task>(`/tasks/${id}`, data),
  delete: (id: number) => request.delete(`/tasks/${id}`),
  start: () => request.post<{ started: number }>('/tasks/start'),
  pause: (id: number) => request.post(`/tasks/${id}/pause`),
  retry: (id: number) => request.post(`/tasks/${id}/retry`),
}
