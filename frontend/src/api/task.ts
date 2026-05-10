import axios from 'axios'
import type {
  ApiResponse,
  Task,
  TaskStatus,
  TaskResult,
  CreateTaskResponse,
  TaskListResponse,
  TaskStats,
  BatchTaskResponse,
  UploadProgress,
} from './types'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 120000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token')
      window.location.href = '/login'
    }
    const message = error.response?.data?.msg || error.message || '请求失败'
    return Promise.reject(new Error(message))
  }
)

export const taskApi = {
  createTask: async (
    file: File,
    outputFormat: 'markdown' | 'txt' | 'json',
    onProgress?: (progress: UploadProgress) => void
  ): Promise<ApiResponse<CreateTaskResponse>> => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('output_format', outputFormat)

    return api.post('/tasks', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress: UploadProgress = {
            loaded: progressEvent.loaded,
            total: progressEvent.total,
            percentage: Math.round((progressEvent.loaded * 100) / progressEvent.total),
          }
          onProgress(progress)
        }
      },
    })
  },

  uploadBatchFiles: async (
    files: File[],
    outputFormat: 'markdown' | 'txt' | 'json',
    onProgress?: (progress: UploadProgress) => void
  ): Promise<ApiResponse<BatchTaskResponse>> => {
    const formData = new FormData()
    files.forEach((file) => {
      formData.append('files', file)
    })
    formData.append('output_format', outputFormat)

    return api.post('/tasks/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress: UploadProgress = {
            loaded: progressEvent.loaded,
            total: progressEvent.total,
            percentage: Math.round((progressEvent.loaded * 100) / progressEvent.total),
          }
          onProgress(progress)
        }
      },
    })
  },

  uploadFromUrl: async (
    url: string,
    fileName?: string,
    outputFormat: 'markdown' | 'txt' | 'json' = 'markdown'
  ): Promise<ApiResponse<CreateTaskResponse>> => {
    return api.post('/tasks/url', {
      url,
      file_name: fileName,
      output_format: outputFormat,
    })
  },

  getTask: async (taskId: string): Promise<ApiResponse<Task>> => {
    return api.get(`/tasks/${taskId}`)
  },

  getTaskStatus: async (taskId: string): Promise<ApiResponse<{
    task_id: string
    file_name?: string
    status: TaskStatus
    progress: number
    error_msg?: string
    created_at?: string
    updated_at?: string
  }>> => {
    return api.get(`/tasks/${taskId}/status`)
  },

  getBatchStatus: async (taskIds: string[]): Promise<ApiResponse<{
    tasks: Array<{
      task_id: string
      file_name?: string
      status: TaskStatus
      progress: number
      error_msg?: string
    }>
    total: number
  }>> => {
    return api.get('/tasks/batch-status', {
      params: { ids: taskIds.join(',') },
    })
  },

  getTaskResult: async (taskId: string): Promise<ApiResponse<TaskResult>> => {
    return api.get(`/tasks/${taskId}/result`)
  },

  downloadResult: async (taskId: string, format: string = 'md'): Promise<Blob> => {
    const response = await api.get(`/tasks/${taskId}/download?format=${format}`, {
      responseType: 'blob',
    }) as unknown as Blob
    return response
  },

  downloadBatchResults: async (taskIds: string[]): Promise<Blob> => {
    const response = await api.get('/tasks/download-all', {
      params: { ids: taskIds.join(',') },
      responseType: 'blob',
    }) as unknown as Blob
    return response
  },

  listTasks: async (params: {
    page?: number
    page_size?: number
    status?: string
    search?: string
  }): Promise<ApiResponse<TaskListResponse>> => {
    return api.get('/tasks', { params })
  },

  getTaskStats: async (): Promise<ApiResponse<TaskStats>> => {
    return api.get('/tasks/stats')
  },

  deleteTask: async (taskId: string): Promise<ApiResponse<null>> => {
    return api.delete(`/tasks/${taskId}`)
  },

  deleteBatchTasks: async (taskIds: string[]): Promise<ApiResponse<{ deleted: number }>> => {
    return api.delete('/tasks/batch', { data: { task_ids: taskIds } })
  },
}

export const authApi = {
  getAuthUrl: async (): Promise<ApiResponse<{ auth_url: string; state: string }>> => {
    return api.get('/auth/url')
  },

  handleCallback: async (code: string, state: string): Promise<ApiResponse<{
    session_id: string
    user: any
    expires_in: number
  }>> => {
    return api.get('/auth/callback', { params: { code, state } })
  },

  getCurrentUser: async (): Promise<ApiResponse<any>> => {
    return api.get('/auth/me')
  },

  refreshToken: async (): Promise<ApiResponse<{
    access_token: string
    refresh_token: string
    expires_in: number
  }>> => {
    return api.post('/auth/refresh')
  },

  logout: async (): Promise<ApiResponse<{ logout_url: string }>> => {
    return api.post('/auth/logout')
  },

  getAPIKey: async (): Promise<ApiResponse<{
    api_key: string
    api_key_secret: string
  }>> => {
    return api.get('/auth/apikey')
  },
}

export default api
