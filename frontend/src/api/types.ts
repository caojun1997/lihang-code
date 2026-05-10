export interface ApiResponse<T = any> {
  code: number
  msg: string
  data?: T
}

export type TaskStatus = 'pending' | 'uploading' | 'processing' | 'completed' | 'partial_failed' | 'failed' | 'cancelled'

export interface Task {
  id: number
  task_id: string
  user_id?: string
  file_name: string
  file_path: string
  file_size: number
  page_count?: number
  output_format: 'markdown' | 'txt' | 'json'
  status: TaskStatus
  progress?: number
  current_page?: number
  total_pages?: number
  error_msg?: string
  created_at: string
  updated_at: string
}

export interface TaskStatusResponse {
  task_id: string
  status: TaskStatus
  progress: number
  current_page?: number
  total_pages?: number
  error_msg?: string
  updated_at?: string
}

export interface TaskResult {
  task_id: string
  content: string
  file_path?: string
  word_count: number
  created_at: string
}

export interface CreateTaskResponse {
  task_id: string
  status: string
  created_at: string
}

export interface TaskListResponse {
  tasks: Task[]
  total: number
  page: number
  size: number
}

export interface User {
  id: string
  username: string
  email?: string
  avatar?: string
  nickname?: string
}

export interface AuthResponse {
  session_id: string
  user: User
  expires_in: number
}

export interface UploadProgress {
  loaded: number
  total: number
  percentage: number
}
