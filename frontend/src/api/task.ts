import axios from 'axios'
import type {
  ApiResponse,
  Task,
  TaskStatus,
  TaskResult,
  CreateTaskResponse,
  TaskListResponse,
} from './types'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
})

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const message = error.response?.data?.msg || error.message || '请求失败'
    return Promise.reject(new Error(message))
  }
)

export const taskApi = {
  createTask: async (file: File, outputFormat: 'markdown' | 'txt'): Promise<ApiResponse<CreateTaskResponse>> => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('output_format', outputFormat)
    return api.post('/tasks', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
  },

  getTask: async (taskId: string): Promise<ApiResponse<Task>> => {
    return api.get(`/tasks/${taskId}`)
  },

  getTaskStatus: async (taskId: string): Promise<ApiResponse<TaskStatus>> => {
    return api.get(`/tasks/${taskId}/status`)
  },

  getTaskResult: async (taskId: string): Promise<ApiResponse<TaskResult>> => {
    return api.get(`/tasks/${taskId}/result`)
  },

  downloadResult: async (taskId: string): Promise<Blob> => {
    const response = await api.get(`/tasks/${taskId}/download`, {
      responseType: 'blob',
    })
    return response
  },

  listTasks: async (params: {
    page?: number
    page_size?: number
    status?: string
  }): Promise<ApiResponse<TaskListResponse>> => {
    return api.get('/tasks', { params })
  },

  deleteTask: async (taskId: string): Promise<ApiResponse<null>> => {
    return api.delete(`/tasks/${taskId}`)
  },
}

export default api
