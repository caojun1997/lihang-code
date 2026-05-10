<template>
  <div class="min-h-screen bg-background">
    <header class="sticky top-0 z-50 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div class="container flex h-16 items-center justify-between px-4 md:px-6">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-lg shadow-primary/30">
            <FileText class="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 class="text-xl font-bold text-foreground">PDF解析助手</h1>
            <p class="text-xs text-muted-foreground">智能识别 · 高效转换</p>
          </div>
        </div>

        <div class="flex items-center gap-4">
          <button
            class="p-2 rounded-lg hover:bg-muted transition-colors"
            @click="toggleTheme"
          >
            <Sun v-if="isDark" class="w-5 h-5 text-muted-foreground" />
            <Moon v-else class="w-5 h-5 text-muted-foreground" />
          </button>
          <div class="w-9 h-9 bg-primary/10 rounded-full flex items-center justify-center">
            <User class="w-5 h-5 text-primary" />
          </div>
        </div>
      </div>
    </header>

    <main class="container px-4 md:px-6 py-8">
      <div class="grid lg:grid-cols-3 gap-8">
        <div class="lg:col-span-2 space-y-8">
          <div class="bg-card rounded-2xl border border-border p-8 shadow-sm">
            <div class="flex items-center gap-2 mb-6">
              <div class="w-1 h-6 bg-primary rounded-full" />
              <h2 class="text-xl font-semibold text-foreground">上传文件</h2>
            </div>

            <FileUploader
              ref="uploaderRef"
              @files-selected="handleFilesSelected"
              @error="handleError"
              @upload-progress="handleUploadProgress"
            />

            <div class="mt-6">
              <FormatSelector v-model="outputFormat" />
            </div>

            <div class="mt-6 flex items-center gap-4">
              <button
                :disabled="selectedFiles.length === 0 || isUploading"
                class="flex-1 h-12 bg-primary text-primary-foreground rounded-xl font-medium hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center gap-2"
                @click="startParsing"
              >
                <Loader2 v-if="isUploading" class="w-5 h-5 animate-spin" />
                <Rocket v-else class="w-5 h-5" />
                {{ isUploading ? '正在解析...' : '开始解析' }}
              </button>
            </div>
          </div>

          <div v-if="currentResult" class="bg-card rounded-2xl border border-border overflow-hidden shadow-sm">
            <ResultPreview
              :content="currentResult.content"
              :file-name="currentTask?.file_name"
              :format="outputFormat"
              :created-at="currentResult.created_at"
              @copy="handleCopy"
              @download="handleDownload"
            />
          </div>
        </div>

        <div class="space-y-6">
          <div class="bg-card rounded-2xl border border-border p-6 shadow-sm">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-2">
                <div class="w-1 h-6 bg-primary rounded-full" />
                <h2 class="text-lg font-semibold text-foreground">任务列表</h2>
              </div>
              <span class="text-sm text-muted-foreground">{{ filteredTasks.length }} 个任务</span>
            </div>

            <div class="space-y-3 max-h-[500px] overflow-y-auto">
              <TransitionGroup
                enter-active-class="transition duration-300 ease-out"
                enter-from-class="opacity-0 -translate-x-4"
                enter-to-class="opacity-100 translate-x-0"
                leave-active-class="transition duration-200 ease-in"
                leave-from-class="opacity-100 translate-x-0"
                leave-to-class="opacity-0 translate-x-4"
              >
                <TaskCard
                  v-for="task in filteredTasks"
                  :key="task.task_id"
                  :task="task"
                  :is-active="currentTask?.task_id === task.task_id"
                  @click="selectTask(task)"
                  @view="viewTask(task)"
                  @download="downloadTask(task)"
                  @copy="copyTask(task)"
                  @delete="deleteTask(task)"
                />
              </TransitionGroup>

              <div v-if="filteredTasks.length === 0" class="text-center py-8">
                <div class="w-16 h-16 bg-muted rounded-2xl flex items-center justify-center mx-auto mb-4">
                  <FileText class="w-8 h-8 text-muted-foreground" />
                </div>
                <p class="text-sm text-muted-foreground">暂无任务</p>
              </div>
            </div>
          </div>

          <div class="bg-card rounded-2xl border border-border p-6 shadow-sm">
            <h3 class="text-sm font-medium text-foreground mb-4">使用统计</h3>
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">总解析次数</span>
                <span class="font-semibold text-foreground">{{ totalTasks }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">成功次数</span>
                <span class="font-semibold text-green-600">{{ completedTasks.length }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">处理中</span>
                <span class="font-semibold text-blue-600">{{ processingTasks.length }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <footer class="border-t border-border mt-12">
      <div class="container px-4 md:px-6 py-6">
        <div class="flex items-center justify-between text-sm text-muted-foreground">
          <p>© 2026 PDF解析助手. All rights reserved.</p>
          <div class="flex items-center gap-4">
            <button class="hover:text-foreground transition-colors">帮助</button>
            <button class="hover:text-foreground transition-colors">隐私</button>
            <button class="hover:text-foreground transition-colors">使用条款</button>
          </div>
        </div>
      </div>
    </footer>

    <Toast v-if="toast.show" :message="toast.message" :type="toast.type" @close="toast.show = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTaskStore } from '@/stores/task'
import FileUploader from '@/components/upload/FileUploader.vue'
import FormatSelector from '@/components/upload/FormatSelector.vue'
import TaskCard from '@/components/task/TaskCard.vue'
import ResultPreview from '@/components/result/ResultPreview.vue'
import Toast from '@/components/common/Toast.vue'
import { FileText, Sun, Moon, User, Loader2, Rocket } from 'lucide-vue-next'
import type { Task } from '@/api/types'

const taskStore = useTaskStore()
const uploaderRef = ref<InstanceType<typeof FileUploader> | null>(null)
const selectedFiles = ref<File[]>([])
const outputFormat = ref<'markdown' | 'txt' | 'json'>('markdown')
const isDark = ref(false)
const toast = ref({
  show: false,
  message: '',
  type: 'success' as 'success' | 'error' | 'info'
})

const tasks = computed(() => taskStore.tasks)
const filteredTasks = computed(() => taskStore.filteredTasks)
const currentTask = computed(() => taskStore.currentTask)
const currentResult = computed(() => taskStore.currentResult)
const isUploading = computed(() => taskStore.isUploading)
const completedTasks = computed(() => taskStore.completedTasks)
const processingTasks = computed(() => taskStore.processingTasks)
const totalTasks = computed(() => tasks.value.length)

onMounted(async () => {
  await taskStore.fetchTasks()
  checkDarkMode()
})

function checkDarkMode() {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

function toggleTheme() {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}

function handleFilesSelected(files: File[]) {
  selectedFiles.value = files
}

function handleError(message: string) {
  showToast(message, 'error')
}

function handleUploadProgress(index: number, progress: number) {
  uploaderRef.value?.updateProgress(index, progress)
}

async function startParsing() {
  if (selectedFiles.value.length === 0) return

  try {
    for (const file of selectedFiles.value) {
      await taskStore.createTask(file, outputFormat.value)
    }
    
    showToast(`已提交 ${selectedFiles.value.length} 个任务`, 'success')
    uploaderRef.value?.clearFiles()
    selectedFiles.value = []

    for (const task of taskStore.processingTasks) {
      await pollTaskStatus(task.task_id)
    }
  } catch (error: any) {
    showToast(error.message || '解析失败', 'error')
  }
}

async function pollTaskStatus(taskId: string) {
  const poll = async () => {
    try {
      const task = tasks.value.find(t => t.task_id === taskId)
      if (!task) return

      if (task.status === 'completed' || task.status === 'failed') {
        if (task.status === 'completed') {
          await taskStore.fetchTaskResult(taskId)
        }
        return
      }

      await taskStore.updateTaskStatus(taskId, task.status, (task.progress || 0) + 10)
      setTimeout(poll, 2000)
    } catch (error) {
      console.error('轮询任务状态失败:', error)
    }
  }

  poll()
}

function selectTask(task: Task) {
  taskStore.currentTask = task
}

async function viewTask(task: Task) {
  taskStore.currentTask = task
  try {
    await taskStore.fetchTaskResult(task.task_id)
  } catch (error) {
    showToast('获取结果失败', 'error')
  }
}

async function downloadTask(task: Task) {
  try {
    await taskStore.downloadResult(task.task_id, outputFormat.value === 'markdown' ? 'md' : outputFormat.value)
    showToast('下载成功', 'success')
  } catch (error) {
    showToast('下载失败', 'error')
  }
}

async function copyTask(task: Task) {
  try {
    const result = await taskStore.fetchTaskResult(task.task_id)
    if (result) {
      navigator.clipboard.writeText(result.content)
      showToast('已复制到剪贴板', 'success')
    }
  } catch (error) {
    showToast('复制失败', 'error')
  }
}

async function deleteTask(task: Task) {
  try {
    await taskStore.deleteTask(task.task_id)
    showToast('删除成功', 'success')
  } catch (error) {
    showToast('删除失败', 'error')
  }
}

function handleCopy() {
  showToast('已复制到剪贴板', 'success')
}

function handleDownload(_format: string) {
  if (currentTask.value) {
    downloadTask(currentTask.value)
  }
}

function showToast(message: string, type: 'success' | 'error' | 'info' = 'success') {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}
</script>
