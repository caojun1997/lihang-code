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
        :loading="isLoading"
        :total="total"
        :page="page"
        :page-size="pageSize"
        @page-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useTaskStore } from '@/stores/task'
import TaskList from '@/components/task/TaskList.vue'

const taskStore = useTaskStore()

const tasks = ref(taskStore.tasks)
const isLoading = ref(false)
const total = ref(taskStore.total)
const page = ref(1)
const pageSize = ref(10)
const statusFilter = ref('')

let pollInterval: number | null = null

async function loadTasks() {
  isLoading.value = true
  await taskStore.fetchTasks({
    page: page.value,
    page_size: pageSize.value,
    status: statusFilter.value,
  })
  tasks.value = taskStore.tasks
  total.value = taskStore.total
  isLoading.value = false
}

function handlePageChange(newPage: number) {
  page.value = newPage
  loadTasks()
}

function handleFilterChange() {
  page.value = 1
  loadTasks()
}

function startPolling() {
  pollInterval = window.setInterval(() => {
    const processingTasks = tasks.value.filter((t) => t.status === 'processing')
    if (processingTasks.length > 0) {
      processingTasks.forEach((task) => {
        taskStore.fetchTaskStatus(task.task_id)
      })
    }
  }, 5000)
}

onMounted(() => {
  loadTasks()
  startPolling()
})

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval)
  }
})
</script>
