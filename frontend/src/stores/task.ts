import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Task, TaskStatus, TaskResult, TaskStats, UploadProgress } from '@/api/types'
import { taskApi } from '@/api'

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([])
  const currentTask = ref<Task | null>(null)
  const currentResult = ref<TaskResult | null>(null)
  const isLoading = ref(false)
  const isUploading = ref(false)
  const uploadProgress = ref<UploadProgress | null>(null)
  const error = ref<string | null>(null)
  
  const pagination = ref({
    page: 1,
    pageSize: 20,
    total: 0,
    totalPages: 0,
  })
  
  const stats = ref<TaskStats>({
    total: 0,
    pending: 0,
    processing: 0,
    completed: 0,
    failed: 0,
  })
  
  const selectedOutputFormat = ref<'markdown' | 'txt' | 'json'>('markdown')
  const statusFilter = ref<TaskStatus | 'all'>('all')
  const searchQuery = ref('')
  const selectedTaskIds = ref<string[]>([])

  const completedTasks = computed(() => 
    tasks.value.filter(t => t.status === 'completed')
  )

  const processingTasks = computed(() => 
    tasks.value.filter(t => ['pending', 'uploading', 'processing'].includes(t.status))
  )

  const failedTasks = computed(() => 
    tasks.value.filter(t => t.status === 'failed')
  )

  const filteredTasks = computed(() => {
    let result = tasks.value
    
    if (statusFilter.value !== 'all') {
      result = result.filter(t => t.status === statusFilter.value)
    }
    
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      result = result.filter(t => 
        t.file_name.toLowerCase().includes(query)
      )
    }
    
    return result
  })

  const hasSelection = computed(() => selectedTaskIds.value.length > 0)

  const allSelected = computed(() => 
    tasks.value.length > 0 && selectedTaskIds.value.length === tasks.value.length
  )

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
        created_at: response.data?.created_at || new Date().toISOString(),
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

  async function uploadBatchFiles(
    files: File[],
    outputFormat: 'markdown' | 'txt' | 'json',
    onProgress?: (progress: UploadProgress) => void
  ) {
    isUploading.value = true
    error.value = null
    uploadProgress.value = null
    
    try {
      const response = await taskApi.uploadBatchFiles(files, outputFormat, (progress) => {
        uploadProgress.value = progress
        onProgress?.(progress)
      })
      
      const newTasks: Task[] = (response.data?.tasks || []).map((t, i) => ({
        id: Date.now() + i,
        task_id: t.task_id,
        file_name: t.file_name || files[i]?.name || 'unknown',
        file_path: '',
        file_size: t.file_size || files[i]?.size || 0,
        output_format: outputFormat,
        status: 'pending' as TaskStatus,
        progress: 0,
        created_at: t.created_at || new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }))
      
      tasks.value = [...newTasks, ...tasks.value]
      
      return response.data
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isUploading.value = false
    }
  }

  async function fetchTasks(params?: { 
    page?: number
    page_size?: number
    status?: string
    search?: string
  }) {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await taskApi.listTasks({
        page: params?.page || pagination.value.page,
        page_size: params?.page_size || pagination.value.pageSize,
        status: params?.status && params.status !== 'all' ? params.status : undefined,
        search: params?.search || searchQuery.value || undefined,
      })
      
      tasks.value = response.data?.tasks || []
      
      if (response.data?.pagination) {
        pagination.value = {
          page: response.data.pagination.page,
          pageSize: response.data.pagination.page_size,
          total: response.data.pagination.total,
          totalPages: response.data.pagination.total_page,
        }
      }
      
      if (response.data?.stats) {
        stats.value = response.data.stats
      }
    } catch (e: any) {
      error.value = e.message
    } finally {
      isLoading.value = false
    }
  }

  async function fetchStats() {
    try {
      const response = await taskApi.getTaskStats()
      if (response.data) {
        stats.value = response.data
      }
    } catch (e: any) {
      console.error('Failed to fetch stats:', e)
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

  async function pollTaskStatus(taskId: string): Promise<void> {
    try {
      const response = await taskApi.getTaskStatus(taskId)
      if (response.data) {
        await updateTaskStatus(taskId, response.data.status as TaskStatus, response.data.progress)
      }
    } catch (e) {
      console.error('Poll task status failed:', e)
    }
  }

  async function pollBatchStatus(taskIds: string[]): Promise<void> {
    try {
      const response = await taskApi.getBatchStatus(taskIds)
      if (response.data?.tasks) {
        for (const taskStatus of response.data.tasks) {
          await updateTaskStatus(
            taskStatus.task_id, 
            taskStatus.status as TaskStatus, 
            taskStatus.progress
          )
        }
      }
    } catch (e) {
      console.error('Poll batch status failed:', e)
    }
  }

  async function deleteTask(taskId: string) {
    isLoading.value = true
    error.value = null
    
    try {
      await taskApi.deleteTask(taskId)
      tasks.value = tasks.value.filter((t) => t.task_id !== taskId)
      selectedTaskIds.value = selectedTaskIds.value.filter(id => id !== taskId)
      
      if (currentTask.value?.task_id === taskId) {
        currentTask.value = null
        currentResult.value = null
      }
      
      await fetchStats()
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function deleteBatchTasks(taskIds?: string[]) {
    const idsToDelete = taskIds || selectedTaskIds.value
    
    if (idsToDelete.length === 0) return
    
    isLoading.value = true
    error.value = null
    
    try {
      await taskApi.deleteBatchTasks(idsToDelete)
      tasks.value = tasks.value.filter(t => !idsToDelete.includes(t.task_id))
      selectedTaskIds.value = []
      
      if (currentTask.value && idsToDelete.includes(currentTask.value.task_id)) {
        currentTask.value = null
        currentResult.value = null
      }
      
      await fetchStats()
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
      const task = tasks.value.find(t => t.task_id === taskId)
      const fileName = task?.file_name?.replace(/\.[^.]+$/, '') || 'result'
      
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${fileName}.${format === 'md' ? 'md' : format}`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(url)
    } catch (e: any) {
      error.value = e.message
      throw e
    }
  }

  async function downloadBatchResults(taskIds?: string[]) {
    const idsToDownload = taskIds || selectedTaskIds.value
    
    if (idsToDownload.length === 0) return
    
    try {
      const blob = await taskApi.downloadBatchResults(idsToDownload)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `results_${Date.now()}.zip`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(url)
    } catch (e: any) {
      error.value = e.message
      throw e
    }
  }

  function toggleTaskSelection(taskId: string) {
    const index = selectedTaskIds.value.indexOf(taskId)
    if (index > -1) {
      selectedTaskIds.value.splice(index, 1)
    } else {
      selectedTaskIds.value.push(taskId)
    }
  }

  function toggleAllSelection() {
    if (allSelected.value) {
      selectedTaskIds.value = []
    } else {
      selectedTaskIds.value = tasks.value.map(t => t.task_id)
    }
  }

  function clearSelection() {
    selectedTaskIds.value = []
  }

  function setStatusFilter(status: TaskStatus | 'all') {
    statusFilter.value = status
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  function setOutputFormat(format: 'markdown' | 'txt' | 'json') {
    selectedOutputFormat.value = format
  }

  function setPage(page: number) {
    pagination.value.page = page
  }

  function clearError() {
    error.value = null
  }

  function reset() {
    tasks.value = []
    currentTask.value = null
    currentResult.value = null
    selectedTaskIds.value = []
    error.value = null
    pagination.value = { page: 1, pageSize: 20, total: 0, totalPages: 0 }
    stats.value = { total: 0, pending: 0, processing: 0, completed: 0, failed: 0 }
  }

  return {
    tasks,
    currentTask,
    currentResult,
    isLoading,
    isUploading,
    uploadProgress,
    error,
    pagination,
    stats,
    selectedOutputFormat,
    statusFilter,
    searchQuery,
    selectedTaskIds,
    completedTasks,
    processingTasks,
    failedTasks,
    filteredTasks,
    hasSelection,
    allSelected,
    createTask,
    uploadBatchFiles,
    fetchTasks,
    fetchStats,
    fetchTask,
    fetchTaskResult,
    updateTaskStatus,
    pollTaskStatus,
    pollBatchStatus,
    deleteTask,
    deleteBatchTasks,
    downloadResult,
    downloadBatchResults,
    toggleTaskSelection,
    toggleAllSelection,
    clearSelection,
    setStatusFilter,
    setSearchQuery,
    setOutputFormat,
    setPage,
    clearError,
    reset,
  }
})
