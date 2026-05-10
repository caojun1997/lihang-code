<template>
  <MainLayout>
    <template #header-actions>
      <div class="flex items-center gap-3">
        <div class="hidden md:flex items-center gap-2 px-3 py-1.5 bg-muted rounded-full text-sm">
          <span class="text-muted-foreground">总任务</span>
          <span class="font-semibold text-foreground">{{ stats.total }}</span>
        </div>
      </div>
    </template>

    <div class="container px-4 md:px-6 py-6">
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-foreground">任务历史</h1>
        <p class="text-muted-foreground mt-1">查看和管理您的所有解析任务</p>
      </div>

      <div class="grid lg:grid-cols-4 gap-6">
        <div class="lg:col-span-3 space-y-4">
          <div class="bg-card rounded-xl border border-border p-4 shadow-sm">
            <div class="flex flex-col sm:flex-row gap-4">
              <div class="relative flex-1">
                <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <input
                  v-model="searchInput"
                  type="text"
                  placeholder="搜索文件名..."
                  class="w-full pl-10 pr-4 py-2.5 border border-border rounded-lg bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
                  @input="handleSearch"
                />
              </div>

              <div class="flex items-center gap-2">
                <select
                  v-model="statusFilter"
                  class="h-10 px-3 border border-border rounded-lg bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
                  @change="handleFilterChange"
                >
                  <option value="all">全部状态</option>
                  <option value="pending">等待中</option>
                  <option value="processing">处理中</option>
                  <option value="completed">已完成</option>
                  <option value="failed">失败</option>
                </select>

                <button
                  v-if="hasSelection"
                  class="h-10 px-4 bg-destructive text-destructive-foreground rounded-lg text-sm font-medium hover:bg-destructive/90 transition-colors flex items-center gap-2"
                  @click="handleBatchDelete"
                >
                  <Trash2 class="w-4 h-4" />
                  删除 ({{ selectedTaskIds.length }})
                </button>

                <button
                  v-if="hasSelection"
                  class="h-10 px-4 bg-primary text-primary-foreground rounded-lg text-sm font-medium hover:bg-primary/90 transition-colors flex items-center gap-2"
                  @click="handleBatchDownload"
                >
                  <Download class="w-4 h-4" />
                  下载全部
                </button>
              </div>
            </div>
          </div>

          <div v-if="isLoading" class="space-y-3">
            <div v-for="i in 5" :key="i" class="h-20 bg-muted rounded-xl animate-shimmer" />
          </div>

          <div v-else-if="tasks.length === 0" class="bg-card rounded-xl border border-border p-12 text-center">
            <div class="w-16 h-16 bg-muted rounded-2xl flex items-center justify-center mx-auto mb-4">
              <FileText class="w-8 h-8 text-muted-foreground/50" />
            </div>
            <p class="text-muted-foreground">暂无任务记录</p>
            <router-link
              to="/"
              class="inline-flex items-center gap-2 mt-4 px-4 py-2 bg-primary text-primary-foreground rounded-lg text-sm font-medium hover:bg-primary/90 transition-colors"
            >
              <Upload class="w-4 h-4" />
              上传文件
            </router-link>
          </div>

          <div v-else class="space-y-3">
            <div class="flex items-center gap-3 px-1">
              <input
                type="checkbox"
                :checked="allSelected"
                class="w-4 h-4 rounded border-input"
                @change="toggleAllSelection"
              />
              <span class="text-sm text-muted-foreground">全选</span>
            </div>

            <TransitionGroup
              enter-active-class="transition duration-300 ease-out"
              enter-from-class="opacity-0 translate-y-2"
              enter-to-class="opacity-100 translate-y-0"
              leave-active-class="transition duration-200 ease-in"
              leave-from-class="opacity-100 translate-y-0"
              leave-to-class="opacity-0 translate-y-2"
            >
              <div
                v-for="task in tasks"
                :key="task.task_id"
                class="bg-card rounded-xl border border-border p-4 hover:shadow-md transition-all"
                :class="{ 'ring-2 ring-primary': selectedTaskIds.includes(task.task_id) }"
              >
                <div class="flex items-start gap-3">
                  <input
                    type="checkbox"
                    :checked="selectedTaskIds.includes(task.task_id)"
                    class="w-4 h-4 rounded border-input mt-1"
                    @change="toggleTaskSelection(task.task_id)"
                  />

                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-1">
                      <p class="text-sm font-medium text-foreground truncate">
                        {{ task.file_name }}
                      </p>
                      <StatusBadge :status="task.status" size="sm" />
                    </div>

                    <div class="flex items-center gap-3 text-xs text-muted-foreground">
                      <span>{{ formatFileSize(task.file_size) }}</span>
                      <span>·</span>
                      <span>{{ formatTime(task.created_at) }}</span>
                      <span v-if="task.output_format !== 'markdown'">·</span>
                      <span v-if="task.output_format !== 'markdown'">{{ task.output_format.toUpperCase() }}</span>
                    </div>

                    <div v-if="task.status === 'processing'" class="mt-2">
                      <div class="h-1.5 bg-muted rounded-full overflow-hidden">
                        <div
                          class="h-full bg-primary rounded-full transition-all duration-300"
                          :style="{ width: `${task.progress}%` }"
                        />
                      </div>
                    </div>

                    <div v-if="task.status === 'failed' && task.error_msg" class="mt-2 text-xs text-destructive">
                      {{ task.error_msg }}
                    </div>
                  </div>

                  <div class="flex items-center gap-1">
                    <button
                      v-if="task.status === 'completed'"
                      class="p-2 text-muted-foreground hover:text-primary hover:bg-muted rounded-lg transition-colors"
                      title="查看结果"
                      @click="viewTask(task)"
                    >
                      <Eye class="w-4 h-4" />
                    </button>
                    <button
                      v-if="task.status === 'completed'"
                      class="p-2 text-muted-foreground hover:text-primary hover:bg-muted rounded-lg transition-colors"
                      title="下载"
                      @click="handleDownload(task.task_id)"
                    >
                      <Download class="w-4 h-4" />
                    </button>
                    <button
                      class="p-2 text-muted-foreground hover:text-destructive hover:bg-muted rounded-lg transition-colors"
                      title="删除"
                      @click="handleDelete(task.task_id)"
                    >
                      <Trash2 class="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            </TransitionGroup>

            <Pagination
              v-if="pagination.totalPages > 1"
              :current-page="pagination.page"
              :total-pages="pagination.totalPages"
              :total="pagination.total"
              :page-size="pagination.pageSize"
              @change="handlePageChange"
            />
          </div>
        </div>

        <div class="space-y-4">
          <div class="bg-card rounded-xl border border-border p-5 shadow-sm">
            <h3 class="text-sm font-medium text-foreground mb-4">任务统计</h3>
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">全部</span>
                <span class="font-semibold text-foreground">{{ stats.total }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">处理中</span>
                <span class="font-semibold text-blue-600">{{ stats.processing }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">已完成</span>
                <span class="font-semibold text-emerald-600">{{ stats.completed }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">失败</span>
                <span class="font-semibold text-red-600">{{ stats.failed }}</span>
              </div>
            </div>
          </div>

          <div class="bg-card rounded-xl border border-border p-5 shadow-sm">
            <h3 class="text-sm font-medium text-foreground mb-3">快捷操作</h3>
            <div class="space-y-2">
              <router-link
                to="/"
                class="flex items-center gap-2 p-2.5 rounded-lg hover:bg-muted transition-colors text-sm text-muted-foreground hover:text-foreground"
              >
                <Upload class="w-4 h-4" />
                上传新文件
              </router-link>
              <button
                class="w-full flex items-center gap-2 p-2.5 rounded-lg hover:bg-muted transition-colors text-sm text-muted-foreground hover:text-foreground text-left"
                @click="refreshTasks"
              >
                <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
                刷新列表
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <TaskDetailModal
      v-if="showDetailModal"
      :task="selectedTask"
      :result="currentResult"
      :loading="isLoadingResult"
      @close="showDetailModal = false"
      @download="handleDownload(selectedTask?.task_id || '')"
    />

    <ToastNotification
      :visible="toast.visible"
      :message="toast.message"
      :type="toast.type"
      @close="toast.visible = false"
    />
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useTaskStore } from '@/stores/task'
import MainLayout from '@/layouts/MainLayout.vue'
import StatusBadge from '@/components/task/StatusBadge.vue'
import Pagination from '@/components/common/Pagination.vue'
import TaskDetailModal from '@/components/task/TaskDetailModal.vue'
import ToastNotification from '@/components/common/ToastNotification.vue'
import { FileText, Search, Upload, Download, Trash2, Eye, RefreshCw } from 'lucide-vue-next'
import type { Task } from '@/api/types'

const taskStore = useTaskStore()

const tasks = computed(() => taskStore.tasks)
const stats = computed(() => taskStore.stats)
const pagination = computed(() => taskStore.pagination)
const isLoading = computed(() => taskStore.isLoading)
const selectedTaskIds = computed(() => taskStore.selectedTaskIds)
const hasSelection = computed(() => taskStore.hasSelection)
const allSelected = computed(() => taskStore.allSelected)
const currentResult = computed(() => taskStore.currentResult)

const searchInput = ref('')
const statusFilter = ref('all')
const showDetailModal = ref(false)
const selectedTask = ref<Task | null>(null)
const isLoadingResult = ref(false)
const toast = ref({
  visible: false,
  message: '',
  type: 'success' as 'success' | 'error' | 'info' | 'warning'
})

let pollingInterval: number | null = null

onMounted(async () => {
  await taskStore.fetchTasks()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})

function startPolling() {
  pollingInterval = window.setInterval(async () => {
    const processingTasks = tasks.value.filter(t => 
      ['pending', 'processing'].includes(t.status)
    )
    
    if (processingTasks.length > 0) {
      await taskStore.pollBatchStatus(processingTasks.map(t => t.task_id))
    }
  }, 3000)
}

function stopPolling() {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

async function refreshTasks() {
  await taskStore.fetchTasks()
  showToast('列表已刷新', 'success')
}

function handleSearch() {
  taskStore.setSearchQuery(searchInput.value)
  taskStore.fetchTasks({ page: 1 })
}

function handleFilterChange() {
  taskStore.setStatusFilter(statusFilter.value as any)
  taskStore.fetchTasks({ page: 1, status: statusFilter.value })
}

function handlePageChange(page: number) {
  taskStore.setPage(page)
  taskStore.fetchTasks({ page })
}

function toggleTaskSelection(taskId: string) {
  taskStore.toggleTaskSelection(taskId)
}

function toggleAllSelection() {
  taskStore.toggleAllSelection()
}

async function viewTask(task: Task) {
  selectedTask.value = task
  showDetailModal.value = true
  isLoadingResult.value = true
  
  try {
    await taskStore.fetchTaskResult(task.task_id)
  } catch {
    showToast('获取结果失败', 'error')
  } finally {
    isLoadingResult.value = false
  }
}

async function handleDownload(taskId: string) {
  try {
    await taskStore.downloadResult(taskId)
    showToast('下载成功', 'success')
  } catch {
    showToast('下载失败', 'error')
  }
}

async function handleDelete(taskId: string) {
  if (!confirm('确定要删除这个任务吗？')) return
  
  try {
    await taskStore.deleteTask(taskId)
    showToast('删除成功', 'success')
  } catch {
    showToast('删除失败', 'error')
  }
}

async function handleBatchDelete() {
  if (!confirm(`确定要删除选中的 ${selectedTaskIds.value.length} 个任务吗？`)) return
  
  try {
    await taskStore.deleteBatchTasks()
    showToast('批量删除成功', 'success')
  } catch {
    showToast('删除失败', 'error')
  }
}

async function handleBatchDownload() {
  try {
    await taskStore.downloadBatchResults()
    showToast('下载成功', 'success')
  } catch {
    showToast('下载失败', 'error')
  }
}

function formatFileSize(bytes?: number): string {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function formatTime(dateString?: string): string {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function showToast(message: string, type: 'success' | 'error' | 'info' | 'warning' = 'success') {
  toast.value = { visible: true, message, type }
  setTimeout(() => {
    toast.value.visible = false
  }, 3000)
}
</script>
