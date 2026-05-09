import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Task, TaskStatus } from '@/api/types'
import { taskApi } from '@/api'

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([])
  const currentTask = ref<Task | null>(null)
  const currentStatus = ref<TaskStatus | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)

  async function createTask(file: File, outputFormat: 'markdown' | 'txt') {
    isLoading.value = true
    error.value = null
    try {
      const response = await taskApi.createTask(file, outputFormat)
      return response.data
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchTasks(params?: { page?: number; page_size?: number; status?: string }) {
    isLoading.value = true
    error.value = null
    try {
      const response = await taskApi.listTasks({
        page: params?.page || page.value,
        page_size: params?.page_size || pageSize.value,
        status: params?.status,
      })
      tasks.value = response.data?.list || []
      total.value = response.data?.total || 0
      page.value = response.data?.page || 1
      pageSize.value = response.data?.page_size || 10
    } catch (e: any) {
      error.value = e.message
    } finally {
      isLoading.value = false
    }
  }

  async function fetchTask(taskId: string) {
    isLoading.value = true
    error.value = null
    try {
      const response = await taskApi.getTask(taskId)
      currentTask.value = response.data
    } catch (e: any) {
      error.value = e.message
    } finally {
      isLoading.value = false
    }
  }

  async function fetchTaskStatus(taskId: string) {
    try {
      const response = await taskApi.getTaskStatus(taskId)
      currentStatus.value = response.data
      if (currentTask.value && currentTask.value.task_id === taskId) {
        currentTask.value.status = response.data.status
      }
      return response.data
    } catch (e: any) {
      error.value = e.message
      throw e
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
      }
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      isLoading.value = false
    }
  }

  return {
    tasks,
    currentTask,
    currentStatus,
    isLoading,
    error,
    total,
    page,
    pageSize,
    createTask,
    fetchTasks,
    fetchTask,
    fetchTaskStatus,
    deleteTask,
  }
})
