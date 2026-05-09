<template>
  <div class="space-y-4">
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Loader2 class="w-8 h-8 animate-spin text-blue-500" />
    </div>

    <div v-else-if="tasks.length === 0" class="text-center py-12">
      <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <FileText class="w-8 h-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">暂无任务</h3>
      <p class="text-sm text-gray-500">上传一个PDF文件开始解析</p>
    </div>

    <div v-else class="space-y-4">
      <TaskCard
        v-for="task in tasks"
        :key="task.task_id"
        :task="task"
        :progress="progressMap[task.task_id]"
        @view="handleView(task)"
        @download="handleDownload(task)"
        @delete="handleDelete(task)"
      />
    </div>

    <div v-if="total > 0" class="flex items-center justify-between pt-4">
      <p class="text-sm text-gray-500">
        共 {{ total }} 条任务
      </p>
      <div class="flex items-center gap-2">
        <button
          :disabled="page <= 1"
          class="px-3 py-1 rounded-lg text-sm border hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          @click="handlePageChange(page - 1)"
        >
          上一页
        </button>
        <span class="text-sm text-gray-600">
          第 {{ page }} / {{ Math.ceil(total / pageSize) }} 页
        </span>
        <button
          :disabled="page >= Math.ceil(total / pageSize)"
          class="px-3 py-1 rounded-lg text-sm border hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          @click="handlePageChange(page + 1)"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { FileText, Loader2 } from 'lucide-vue-next'
import { useTaskStore } from '@/stores/task'
import TaskCard from './TaskCard.vue'
import type { Task } from '@/api/types'

defineProps<{
  tasks: Task[]
  loading: boolean
  total: number
  page: number
  pageSize: number
}>()

const emit = defineEmits<{
  (e: 'page-change', page: number): void
}>()

const router = useRouter()
const taskStore = useTaskStore()
const progressMap = ref<Record<string, number>>({})
let pollInterval: number | null = null

async function pollTaskStatus(taskId: string) {
  try {
    const status = await taskStore.fetchTaskStatus(taskId)
    progressMap.value[taskId] = status.progress

    if (status.status === 'completed' || status.status === 'failed') {
      if (pollInterval) {
        clearInterval(pollInterval)
        pollInterval = null
      }
    }
  } catch (e) {
    console.error('Failed to poll task status:', e)
  }
}

function handleView(task: Task) {
  router.push(`/tasks/${task.task_id}`)
}

async function handleDownload(task: Task) {
  try {
    const blob = await fetch(`/api/v1/tasks/${task.task_id}/download`).then((r) => r.blob())
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = task.file_name.replace('.pdf', task.output_format === 'markdown' ? '_parsed.md' : '_parsed.txt')
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) {
    console.error('Download failed:', e)
  }
}

async function handleDelete(task: Task) {
  if (confirm('确定要删除这个任务吗？')) {
    await taskStore.deleteTask(task.task_id)
    emit('page-change', 1)
  }
}

function handlePageChange(page: number) {
  emit('page-change', page)
}

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval)
  }
})
</script>
