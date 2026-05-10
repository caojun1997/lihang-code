<template>
  <div class="min-h-[calc(100vh-200px)]">
    <div class="container mx-auto px-4 py-8 max-w-6xl">
      <div class="mb-8">
        <h2 class="text-3xl font-bold text-gray-900 mb-2">任务列表</h2>
        <p class="text-gray-600">查看和管理所有解析任务</p>
      </div>

      <div class="bg-white rounded-xl border shadow-sm p-6 mb-6">
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <label class="text-sm font-medium text-gray-700">状态筛选:</label>
            <select
              v-model="statusFilter"
              class="border rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              @change="handleFilterChange"
            >
              <option value="">全部</option>
              <option value="pending">等待中</option>
              <option value="processing">处理中</option>
              <option value="completed">已完成</option>
              <option value="failed">失败</option>
            </select>
          </div>
        </div>
      </div>

      <TaskList
        :tasks="tasks"
        :current-task="currentTask"
        @select="handleSelect"
        @view="handleView"
        @download="handleDownload"
        @copy="handleCopy"
        @delete="handleDelete"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/task'
import TaskList from '@/components/task/TaskList.vue'
import type { Task } from '@/api/types'

const router = useRouter()
const taskStore = useTaskStore()

const tasks = computed(() => taskStore.tasks)
const currentTask = computed(() => taskStore.currentTask)
const statusFilter = ref('')

async function loadTasks() {
  await taskStore.fetchTasks({
    page: 1,
    page_size: 50,
    status: statusFilter.value || undefined,
  })
}

function handleSelect(task: Task) {
  taskStore.currentTask = task
}

function handleView(task: Task) {
  router.push(`/tasks/${task.task_id}`)
}

async function handleDownload(task: Task) {
  try {
    await taskStore.downloadResult(task.task_id)
  } catch (error) {
    console.error('Download failed:', error)
  }
}

async function handleCopy(task: Task) {
  try {
    const result = await taskStore.fetchTaskResult(task.task_id)
    if (result) {
      navigator.clipboard.writeText(result.content)
    }
  } catch (error) {
    console.error('Copy failed:', error)
  }
}

async function handleDelete(task: Task) {
  if (confirm('确定要删除这个任务吗？')) {
    await taskStore.deleteTask(task.task_id)
  }
}

function handleFilterChange() {
  taskStore.setStatusFilter(statusFilter.value as any)
  loadTasks()
}

onMounted(() => {
  loadTasks()
})
</script>
