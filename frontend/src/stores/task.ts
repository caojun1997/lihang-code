import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Task, TaskStatus, TaskResult, UploadProgress } from '@/api/types'
import { taskApi } from '@/api'

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([])
  const currentTask = ref<Task | null>(null)
  const currentResult = ref<TaskResult | null>(null)
  const isLoading = ref(false)
  const isUploading = ref(false)
  const uploadProgress = ref<UploadProgress | null>(null)
  const error = ref<string | null>(null)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const selectedOutputFormat = ref<'markdown' | 'txt' | 'json'>('markdown')
  const statusFilter = ref<TaskStatus | 'all'>('all')

  const completedTasks = computed(() => 
    tasks.value.filter(t => t.status === 'completed')
  )

  const processingTasks = computed(() => 
    tasks.value.filter(t => ['pending', 'uploading', 'processing'].includes(t.status))
  )

  const filteredTasks = computed(() => {
    if (statusFilter.value === 'all') {
      return tasks.value
    }
    return tasks.value.filter(t => t.status === statusFilter.value)
  })

  async function createTask(
    file: File, 
    outputFormat: 'markdown' | 'txt' | 'json',
    onProgress?: (progress: UploadProgress) => void
  ) {
    isUploading.value = true
    error.value = null
    uploadProgress.value = null
    
    try {
      const response = await taskApi.createTask(file, outputFormat, (progress) => {
        uploadProgress.value = progress
        onProgress?.(progress)
      })
      
      const newTask: Task = {
        id: Date.now(),
        task_id: response.data?.task_id || '',
        file_name: file.name,
        file_path: '',
        file_size: file.size,
        output_format: outputFormat,
        status: 'pending',
        progress: 0,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      
      tasks.value.unshift(newTask)
      currentTask.value = newTask
      
      return response.data
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isUploading.value = false
    }
  }

  async function fetchTasks(params?: { page?: number; page_size?: number; status?: string }) {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await taskApi.listTasks({
        page: params?.page || page.value,
        page_size: params?.page_size || pageSize.value,
        status: params?.status !== 'all' ? params?.status : undefined,
      })
      
      tasks.value = response.data?.tasks || []
      total.value = response.data?.total || 0
      page.value = response.data?.page || 1
      pageSize.value = response.data?.size || 10
    } catch (e: any) {
      error.value = e.message
    } finally {
      isLoading.value = false
    }
  }

  async function fetchTask(taskId: string): Promise<Task | undefined> {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await taskApi.getTask(taskId)
      currentTask.value = response.data || null
      return response.data
    } catch (e: any) {
      error.value = e.message
      return undefined
    } finally {
      isLoading.value = false
    }
  }

  async function fetchTaskResult(taskId: string): Promise<TaskResult | undefined> {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await taskApi.getTaskResult(taskId)
      currentResult.value = response.data || null
      return response.data
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function updateTaskStatus(taskId: string, status: TaskStatus, progress: number = 0) {
    const task = tasks.value.find(t => t.task_id === taskId)
    if (task) {
      task.status = status
      task.progress = progress
      task.updated_at = new Date().toISOString()
    }
    if (currentTask.value?.task_id === taskId) {
      currentTask.value.status = status
      currentTask.value.progress = progress
    }
  }

  async function deleteTask(taskId: string) {
    isLoading.value = true
    error.value = null
    
    try {
      await taskApi.deleteTask(taskId)
      tasks.value = tasks.value.filter((t) => t.task_id !== taskId)
      
      if (currentTask.value?.task_id === taskId) {
        currentTask.value = null
        currentResult.value = null
      }
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function downloadResult(taskId: string, format: string = 'md') {
    try {
      const blob = await taskApi.downloadResult(taskId, format)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `result.${format === 'md' ? 'md' : format}`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(url)
    } catch (e: any) {
      error.value = e.message
      throw e
    }
  }

  function setStatusFilter(status: TaskStatus | 'all') {
    statusFilter.value = status
  }

  function setOutputFormat(format: 'markdown' | 'txt' | 'json') {
    selectedOutputFormat.value = format
  }

  function clearError() {
    error.value = null
  }

  return {
    tasks,
    currentTask,
    currentResult,
    isLoading,
    isUploading,
    uploadProgress,
    error,
    total,
    page,
    pageSize,
    selectedOutputFormat,
    statusFilter,
    completedTasks,
    processingTasks,
    filteredTasks,
    createTask,
    fetchTasks,
    fetchTask,
    fetchTaskResult,
    updateTaskStatus,
    deleteTask,
    downloadResult,
    setStatusFilter,
    setOutputFormat,
    clearError,
  }
})
