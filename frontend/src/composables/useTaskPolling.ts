import { ref, onUnmounted, type Ref } from 'vue'

export type TaskStatus = 'pending' | 'uploading' | 'processing' | 'completed' | 'partial_failed' | 'failed' | 'cancelled'

export interface TaskStatusInfo {
  status: TaskStatus
  progress: number
  currentPage?: number
  totalPages?: number
  message?: string
}

export function useTaskPolling(
  taskId: Ref<string>,
  onStatusChange?: (status: TaskStatusInfo) => void,
  onComplete?: (status: TaskStatusInfo) => void,
  onError?: (error: any) => void
) {
  const status = ref<TaskStatusInfo>({
    status: 'pending',
    progress: 0
  })
  const isPolling = ref(false)
  let intervalId: number | null = null

  async function poll(fetchStatus: () => Promise<TaskStatusInfo>) {
    try {
      const result = await fetchStatus()
      status.value = result

      onStatusChange?.(result)

      if (result.status === 'completed' || result.status === 'partial_failed') {
        stopPolling()
        onComplete?.(result)
      } else if (result.status === 'failed' || result.status === 'cancelled') {
        stopPolling()
        onError?.(new Error(`Task ${result.status}`))
      }
    } catch (error) {
      console.error('轮询失败:', error)
      onError?.(error)
    }
  }

  function startPolling(
    fetchStatus: () => Promise<TaskStatusInfo>,
    interval = 2000
  ) {
    if (isPolling.value) return

    isPolling.value = true

    poll(fetchStatus)

    intervalId = window.setInterval(() => {
      poll(fetchStatus)
    }, interval)
  }

  function stopPolling() {
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }
    isPolling.value = false
  }

  function updateTaskId(newTaskId: string) {
    taskId.value = newTaskId
  }

  onUnmounted(() => {
    stopPolling()
  })

  return {
    status,
    isPolling,
    startPolling,
    stopPolling,
    updateTaskId
  }
}

export function useMultiTaskPolling(
  taskIds: Ref<string[]>,
  onTaskUpdate?: (taskId: string, status: TaskStatusInfo) => void,
  onAllComplete?: (results: Map<string, TaskStatusInfo>) => void
) {
  const taskStatuses = ref<Map<string, TaskStatusInfo>>(new Map())
  const isPolling = ref(false)
  const intervalId = ref<number | null>(null)

  function updateStatus(taskId: string, info: TaskStatusInfo) {
    taskStatuses.value.set(taskId, info)
    onTaskUpdate?.(taskId, info)

    const allDone = Array.from(taskStatuses.value.values()).every(
      s => ['completed', 'partial_failed', 'failed', 'cancelled'].includes(s.status)
    )

    if (allDone) {
      stopPolling()
      onAllComplete?.(taskStatuses.value)
    }
  }

  async function pollAll(fetchStatuses: (taskIds: string[]) => Promise<Map<string, TaskStatusInfo>>) {
    try {
      const results = await fetchStatuses(taskIds.value)

      results.forEach((info, taskId) => {
        updateStatus(taskId, info)
      })
    } catch (error) {
      console.error('批量轮询失败:', error)
    }
  }

  function startPolling(
    fetchStatuses: (taskIds: string[]) => Promise<Map<string, TaskStatusInfo>>,
    interval = 3000
  ) {
    if (isPolling.value) return

    isPolling.value = true

    pollAll(fetchStatuses)

    intervalId.value = window.setInterval(() => {
      pollAll(fetchStatuses)
    }, interval)
  }

  function stopPolling() {
    if (intervalId.value) {
      clearInterval(intervalId.value)
      intervalId.value = null
    }
    isPolling.value = false
  }

  function reset() {
    taskStatuses.value = new Map()
    stopPolling()
  }

  onUnmounted(() => {
    stopPolling()
  })

  return {
    taskStatuses,
    isPolling,
    startPolling,
    stopPolling,
    reset
  }
}
