export interface ApiResponse<T = any> {
  code: number
  msg: string
  data?: T
}

export interface Task {
  id: number
  task_id: string
  file_name: string
  file_path: string
  file_size: number
  page_count: number
  output_format: 'markdown' | 'txt'
  status: 'pending' | 'processing' | 'completed' | 'failed'
  mineru_task_id?: string
  error_message?: string
  created_at: string
  updated_at: string
}

export interface TaskStatus {
  task_id: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  progress: number
  updated_at?: string
}

export interface TaskResult {
  task_id: string
  content: string
  word_count: number
}

export interface CreateTaskResponse {
  task_id: string
  status: string
  created_at: string
}

export interface TaskListResponse {
  total: number
  page: number
  page_size: number
  list: Task[]
}
