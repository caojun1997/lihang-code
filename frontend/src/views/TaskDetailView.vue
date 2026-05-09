<template>
  <div class="min-h-[calc(100vh-200px)]">
    <div class="container mx-auto px-4 py-8 max-w-6xl">
      <button
        type="button"
        class="flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-6"
        @click="goBack"
      >
        <ArrowLeft class="w-5 h-5" />
        返回
      </button>

      <div v-if="loading" class="flex items-center justify-center py-24">
        <Loader2 class="w-12 h-12 animate-spin text-blue-500" />
      </div>

      <div v-else-if="task" class="space-y-6">
        <div class="bg-white rounded-xl border shadow-sm p-6">
          <div class="flex items-start justify-between mb-6">
            <div class="flex items-center gap-4">
              <div class="w-16 h-16 bg-red-100 rounded-xl flex items-center justify-center">
                <FileText class="w-8 h-8 text-red-500" />
              </div>
              <div>
                <h2 class="text-xl font-bold text-gray-900">{{ task.file_name }}</h2>
                <p class="text-sm text-gray-500 mt-1">
                  创建于 {{ formatDate(task.created_at) }}
                </p>
              </div>
            </div>
            <TaskStatus :status="task.status" />
          </div>

          <div class="grid grid-cols-4 gap-6 mb-6">
            <div>
              <p class="text-sm text-gray-500">文件大小</p>
              <p class="text-base font-medium text-gray-900">{{ formatFileSize(task.file_size) }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">输出格式</p>
              <p class="text-base font-medium text-gray-900">
                {{ task.output_format === 'markdown' ? 'Markdown' : '纯文本' }}
              </p>
            </div>
            <div>
              <p class="text-sm text-gray-500">任务状态</p>
              <p class="text-base font-medium text-gray-900">{{ statusText[task.status] }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">任务ID</p>
              <p class="text-base font-medium text-gray-900 font-mono">{{ task.task_id }}</p>
            </div>
          </div>

          <div v-if="task.status === 'processing'" class="mb-6">
            <div class="flex items-center justify-between text-sm mb-2">
              <span class="text-gray-500">处理进度</span>
              <span class="font-medium text-blue-600">{{ progress }}%</span>
            </div>
            <div class="h-3 bg-gray-100 rounded-full overflow-hidden">
              <div
                class="h-full bg-blue-500 rounded-full transition-all duration-500"
                :style="{ width: `${progress}%` }"
              />
            </div>
          </div>

          <div v-if="task.error_message" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p class="text-sm text-red-600">错误信息: {{ task.error_message }}</p>
          </div>

          <div class="flex items-center gap-4">
            <button
              v-if="task.status === 'completed'"
              type="button"
              class="flex-1 bg-blue-600 text-white px-6 py-3 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors flex items-center justify-center gap-2"
              @click="handleDownload"
            >
              <Download class="w-5 h-5" />
              下载解析结果
            </button>
            <button
              type="button"
              class="text-red-600 px-6 py-3 rounded-lg text-sm font-medium hover:bg-red-50 transition-colors flex items-center gap-2"
              @click="handleDelete"
            >
              <Trash2 class="w-5 h-5" />
              删除任务
            </button>
          </div>
        </div>

        <ResultViewer v-if="task.status === 'completed'" :task-id="task.task_id" />
      </div>

      <div v-else class="text-center py-24">
        <p class="text-gray-500">任务不存在</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, FileText, Download, Trash2, Loader2 } from 'lucide-vue-next'
import { useTaskStore } from '@/stores/task'
import { formatFileSize, formatDate } from '@/lib/utils'
import TaskStatus from '@/components/task/TaskStatus.vue'
import ResultViewer from '@/components/result/ResultViewer.vue'

const route = useRoute()
const router = useRouter()
const taskStore = useTaskStore()

const task = ref(taskStore.currentTask)
const loading = ref(true)
const progress = ref(0)

const statusText: Record<string, string> = {
  pending: '等待中',
  processing: '处理中',
  completed: '已完成',
  failed: '失败',
}

let pollInterval: number | null = null

async function loadTask() {
  loading.value = true
  const taskId = route.params.id as string
  await taskStore.fetchTask(taskId)
  task.value = taskStore.currentTask

  if (task.value?.status === 'processing') {
    startPolling()
  }
  loading.value = false
}

function startPolling() {
  if (pollInterval) return

  pollInterval = window.setInterval(async () => {
    if (!task.value) return

    try {
      const status = await taskStore.fetchTaskStatus(task.value.task_id)
      progress.value = status.progress

      if (task.value && status.status) {
        task.value.status = status.status
      }

      if (status.status === 'completed' || status.status === 'failed') {
        stopPolling()
      }
    } catch (e) {
      console.error('Failed to poll status:', e)
    }
  }, 2000)
}

function stopPolling() {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
}

async function handleDownload() {
  if (!task.value) return

  try {
    const blob = await fetch(`/api/v1/tasks/${task.value.task_id}/download`).then((r) => r.blob())
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = task.value.file_name.replace(
      '.pdf',
      task.value.output_format === 'markdown' ? '_parsed.md' : '_parsed.txt'
    )
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) {
    console.error('Download failed:', e)
  }
}

async function handleDelete() {
  if (!task.value) return

  if (confirm('确定要删除这个任务吗？')) {
    await taskStore.deleteTask(task.value.task_id)
    router.push('/tasks')
  }
}

function goBack() {
  router.back()
}

onMounted(() => {
  loadTask()
})

onUnmounted(() => {
  stopPolling()
})
</script>
