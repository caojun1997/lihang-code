<template>
  <MainLayout>
    <template #header-actions>
      <div class="hidden md:flex items-center gap-2 px-3 py-1.5 bg-muted rounded-full text-sm">
        <span class="text-muted-foreground">总解析</span>
        <span class="font-semibold text-foreground">{{ stats.total }}</span>
        <span class="text-muted-foreground">次</span>
      </div>
    </template>

    <div class="container px-4 md:px-6 py-6">
      <div class="mb-8">
        <StepIndicator
          :steps="['上传文件', '配置选项', '下载结果']"
          :current-step="currentStep"
        />
      </div>

      <div class="grid lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 space-y-6">
          <Transition
            mode="out-in"
            enter-active-class="transition duration-300 ease-out"
            enter-from-class="opacity-0 translate-x-4"
            enter-to-class="opacity-100 translate-x-0"
            leave-active-class="transition duration-200 ease-in"
            leave-from-class="opacity-100 translate-x-0"
            leave-to-class="opacity-0 -translate-x-4"
          >
            <div
              :key="currentStep"
              class="bg-card rounded-xl border border-border p-6 shadow-sm"
            >
              <div v-if="currentStep === 0" class="space-y-6">
                <div class="flex items-center gap-2">
                  <div class="w-1 h-6 bg-primary rounded-full" />
                  <h2 class="text-lg font-semibold text-foreground">上传文件</h2>
                  <span class="text-xs text-muted-foreground ml-auto">支持批量上传</span>
                </div>

                <FileUploader
                  ref="uploaderRef"
                  v-model="selectedFiles"
                  :max-size-m-b="100"
                  @error="handleError"
                />

                <div class="flex justify-end">
                  <button
                    :disabled="selectedFiles.length === 0"
                    class="h-10 px-6 bg-primary text-primary-foreground rounded-lg font-medium hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center gap-2"
                    @click="nextStep"
                  >
                    下一步
                    <ArrowRight class="w-4 h-4" />
                  </button>
                </div>
              </div>

              <div v-else-if="currentStep === 1" class="space-y-6">
                <div class="flex items-center gap-2">
                  <div class="w-1 h-6 bg-primary rounded-full" />
                  <h2 class="text-lg font-semibold text-foreground">配置解析选项</h2>
                  <span class="text-xs text-muted-foreground ml-auto">已选择 {{ selectedFiles.length }} 个文件</span>
                </div>

                <FormatSelector
                  v-model="outputFormat"
                  :options="parseOptions"
                  @update:options="parseOptions = $event"
                />

                <div class="flex items-center justify-between">
                  <button
                    class="h-10 px-6 border border-border rounded-lg font-medium hover:bg-muted transition-all flex items-center gap-2"
                    @click="prevStep"
                  >
                    <ArrowLeft class="w-4 h-4" />
                    上一步
                  </button>

                  <button
                    :disabled="isUploading"
                    class="h-10 px-6 bg-primary text-primary-foreground rounded-lg font-medium hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center gap-2"
                    @click="startParsing"
                  >
                    <Loader2 v-if="isUploading" class="w-4 h-4 animate-spin" />
                    <Rocket v-else class="w-4 h-4" />
                    {{ isUploading ? '上传中...' : '开始解析' }}
                  </button>
                </div>
              </div>

              <div v-else class="space-y-6">
                <div class="flex items-center gap-2">
                  <div class="w-1 h-6 bg-primary rounded-full" />
                  <h2 class="text-lg font-semibold text-foreground">解析结果</h2>
                </div>

                <div v-if="hasCompletedTasks" class="space-y-4">
                  <div class="flex items-center justify-between p-4 bg-emerald-50 dark:bg-emerald-950/20 rounded-lg border border-emerald-200 dark:border-emerald-800">
                    <div class="flex items-center gap-3">
                      <CheckCircle class="w-5 h-5 text-emerald-600" />
                      <span class="text-sm text-emerald-700 dark:text-emerald-400">
                        解析完成！共 {{ completedTasks.length }} 个文件
                      </span>
                    </div>
                  </div>

                  <div class="flex items-center gap-3">
                    <button
                      class="flex-1 h-10 bg-primary text-primary-foreground rounded-lg font-medium hover:bg-primary/90 transition-all flex items-center justify-center gap-2"
                      @click="downloadAllResults"
                    >
                      <Download class="w-4 h-4" />
                      下载全部
                    </button>
                    <button
                      class="h-10 px-4 border border-border rounded-lg font-medium hover:bg-muted transition-all flex items-center gap-2"
                      @click="copyAllResults"
                    >
                      <Copy class="w-4 h-4" />
                      复制全部
                    </button>
                    <router-link
                      to="/history"
                      class="h-10 px-4 border border-border rounded-lg font-medium hover:bg-muted transition-all flex items-center gap-2"
                    >
                      <History class="w-4 h-4" />
                      查看历史
                    </router-link>
                  </div>
                </div>

                <div v-else-if="hasProcessingTasks" class="text-center py-8">
                  <div class="w-16 h-16 bg-primary/10 rounded-2xl flex items-center justify-center mx-auto mb-4">
                    <Loader2 class="w-8 h-8 text-primary animate-spin" />
                  </div>
                  <p class="text-muted-foreground">正在解析中，请稍候...</p>
                  <p class="text-sm text-muted-foreground mt-1">{{ processingTasks.length }} 个任务处理中</p>
                </div>

                <div class="flex items-center justify-between">
                  <button
                    class="h-10 px-6 border border-border rounded-lg font-medium hover:bg-muted transition-all flex items-center gap-2"
                    @click="resetUpload"
                  >
                    <ArrowLeft class="w-4 h-4" />
                    继续上传
                  </button>
                </div>
              </div>
            </div>
          </Transition>

          <div v-if="currentStep >= 1 && tasks.length > 0" class="bg-card rounded-xl border border-border p-6 shadow-sm">
            <div class="flex items-center gap-2 mb-4">
              <div class="w-1 h-6 bg-primary rounded-full" />
              <h2 class="text-lg font-semibold text-foreground">当前任务</h2>
            </div>

            <div class="space-y-3 max-h-[300px] overflow-y-auto scrollbar-thin">
              <TaskCard
                v-for="task in tasks.slice(0, 10)"
                :key="task.task_id"
                :task="task"
                :is-active="currentTask?.task_id === task.task_id"
                :show-actions="task.status === 'completed'"
                @click="selectTask(task)"
                @view="viewTask(task)"
                @download="downloadTask(task)"
                @copy="copyTask(task)"
                @delete="deleteTask(task)"
              />
            </div>

            <div v-if="tasks.length > 10" class="mt-4 text-center">
              <router-link
                to="/history"
                class="text-sm text-primary hover:underline"
              >
                查看全部 {{ tasks.length }} 个任务 →
              </router-link>
            </div>
          </div>
        </div>

        <div class="space-y-6">
          <ResultPreview
            :content="currentResult?.content"
            :file-name="currentTask?.file_name"
            :format="outputFormat"
            :created-at="currentResult?.created_at"
            :loading="isLoadingResult"
            @copy="handleCopy"
            @download="handleDownload"
          />

          <div class="bg-card rounded-xl border border-border p-5 shadow-sm">
            <h3 class="text-sm font-medium text-foreground mb-4">使用统计</h3>
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">总解析次数</span>
                <span class="font-semibold text-foreground">{{ stats.total }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">成功次数</span>
                <span class="font-semibold text-emerald-600">{{ stats.completed }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">处理中</span>
                <span class="font-semibold text-blue-600">{{ stats.processing }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">失败</span>
                <span class="font-semibold text-red-600">{{ stats.failed }}</span>
              </div>
            </div>
            
            <div class="mt-4 pt-4 border-t border-border">
              <router-link
                to="/history"
                class="flex items-center justify-center gap-2 w-full py-2 text-sm text-primary hover:bg-muted rounded-lg transition-colors"
              >
                <History class="w-4 h-4" />
                查看完整历史
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ToastNotification
      :visible="toast.visible"
      :message="toast.message"
      :type="toast.type"
      @close="toast.visible = false"
    />
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useTaskStore } from '@/stores/task'
import MainLayout from '@/layouts/MainLayout.vue'
import StepIndicator from '@/components/common/StepIndicator.vue'
import FileUploader from '@/components/upload/FileUploader.vue'
import FormatSelector from '@/components/upload/FormatSelector.vue'
import TaskCard from '@/components/task/TaskCard.vue'
import ResultPreview from '@/components/result/ResultPreview.vue'
import ToastNotification from '@/components/common/ToastNotification.vue'
import {
  ArrowRight,
  ArrowLeft,
  Download,
  Copy,
  Rocket,
  Loader2,
  CheckCircle,
  History
} from 'lucide-vue-next'
import type { Task } from '@/api/types'

const taskStore = useTaskStore()
const uploaderRef = ref<InstanceType<typeof FileUploader> | null>(null)

const currentStep = ref(0)
const selectedFiles = ref<File[]>([])
const outputFormat = ref<'markdown' | 'txt' | 'json'>('markdown')
const parseOptions = ref({
  preserve_tables: true,
  extract_images: true,
  latex_formulas: true,
  preserve_headings: true,
  language: 'auto'
})
const isLoadingResult = ref(false)
const toast = ref({
  visible: false,
  message: '',
  type: 'success' as 'success' | 'error' | 'info' | 'warning'
})

const tasks = computed(() => taskStore.tasks)
const currentTask = computed(() => taskStore.currentTask)
const currentResult = computed(() => taskStore.currentResult)
const isUploading = computed(() => taskStore.isUploading)
const completedTasks = computed(() => taskStore.completedTasks)
const processingTasks = computed(() => taskStore.processingTasks)
const stats = computed(() => taskStore.stats)
const hasCompletedTasks = computed(() => completedTasks.value.length > 0)
const hasProcessingTasks = computed(() => processingTasks.value.length > 0)

let pollingInterval: number | null = null

onMounted(async () => {
  await taskStore.fetchTasks()
  await taskStore.fetchStats()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})

function startPolling() {
  pollingInterval = window.setInterval(async () => {
    const processing = tasks.value.filter(t => 
      ['pending', 'uploading', 'processing'].includes(t.status)
    )
    
    if (processing.length > 0) {
      await taskStore.pollBatchStatus(processing.map(t => t.task_id))
    }
  }, 3000)
}

function stopPolling() {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

watch(
  () => processingTasks.value.length,
  (newVal, oldVal) => {
    if (oldVal > 0 && newVal === 0 && currentStep.value === 1) {
      currentStep.value = 2
    }
  }
)

function nextStep() {
  if (currentStep.value < 2) {
    currentStep.value++
  }
}

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

function handleError(message: string) {
  showToast(message, 'error')
}

async function startParsing() {
  if (selectedFiles.value.length === 0) return

  try {
    if (selectedFiles.value.length === 1) {
      await taskStore.createTask(selectedFiles.value[0], outputFormat.value)
    } else {
      await taskStore.uploadBatchFiles(selectedFiles.value, outputFormat.value)
    }

    showToast(`已提交 ${selectedFiles.value.length} 个任务`, 'success')
    uploaderRef.value?.clearFiles()
    selectedFiles.value = []
    nextStep()
  } catch (error: any) {
    showToast(error.message || '解析失败', 'error')
  }
}

function resetUpload() {
  currentStep.value = 0
  selectedFiles.value = []
  taskStore.currentTask = null
  taskStore.currentResult = null
}

function selectTask(task: Task) {
  taskStore.currentTask = task
}

async function viewTask(task: Task) {
  taskStore.currentTask = task
  isLoadingResult.value = true
  try {
    await taskStore.fetchTaskResult(task.task_id)
  } catch {
    showToast('获取结果失败', 'error')
  } finally {
    isLoadingResult.value = false
  }
}

async function downloadTask(task: Task) {
  try {
    await taskStore.downloadResult(task.task_id, outputFormat.value === 'markdown' ? 'md' : outputFormat.value)
    showToast('下载成功', 'success')
  } catch {
    showToast('下载失败', 'error')
  }
}

async function copyTask(task: Task) {
  try {
    const result = await taskStore.fetchTaskResult(task.task_id)
    if (result) {
      await navigator.clipboard.writeText(result.content)
      showToast('已复制到剪贴板', 'success')
    }
  } catch {
    showToast('复制失败', 'error')
  }
}

async function deleteTask(task: Task) {
  try {
    await taskStore.deleteTask(task.task_id)
    showToast('删除成功', 'success')
  } catch {
    showToast('删除失败', 'error')
  }
}

async function downloadAllResults() {
  const completedIds = completedTasks.value.map(t => t.task_id)
  if (completedIds.length === 0) return
  
  try {
    await taskStore.downloadBatchResults(completedIds)
    showToast('下载成功', 'success')
  } catch {
    showToast('下载失败', 'error')
  }
}

async function copyAllResults() {
  for (const task of completedTasks.value.slice(0, 5)) {
    await copyTask(task)
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

function showToast(message: string, type: 'success' | 'error' | 'info' | 'warning' = 'success') {
  toast.value = { visible: true, message, type }
  setTimeout(() => {
    toast.value.visible = false
  }, 3000)
}
</script>
