<template>
  <MainLayout>
    <template #header-actions>
      <div class="hidden md:flex items-center gap-2 px-3 py-1.5 bg-muted rounded-full text-sm">
        <span class="text-muted-foreground">总解析</span>
        <span class="font-semibold text-foreground">{{ totalTasks }}</span>
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
                    {{ isUploading ? '解析中...' : '开始解析' }}
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
                  </div>
                </div>

                <div v-else-if="hasProcessingTasks" class="text-center py-8">
                  <div class="w-16 h-16 bg-primary/10 rounded-2xl flex items-center justify-center mx-auto mb-4">
                    <Loader2 class="w-8 h-8 text-primary animate-spin" />
                  </div>
                  <p class="text-muted-foreground">正在解析中，请稍候...</p>
                </div>

                <div class="flex items-center justify-between">
                  <button
                    class="h-10 px-6 border border-border rounded-lg font-medium hover:bg-muted transition-all flex items-center gap-2"
                    @click="prevStep"
                  >
                    <ArrowLeft class="w-4 h-4" />
                    重新解析
                  </button>
                </div>
              </div>
            </div>
          </Transition>

          <div v-if="currentStep === 2 && hasProcessingTasks" class="bg-card rounded-xl border border-border p-6 shadow-sm">
            <div class="flex items-center gap-2 mb-4">
              <div class="w-1 h-6 bg-primary rounded-full" />
              <h2 class="text-lg font-semibold text-foreground">处理中</h2>
            </div>

            <div class="space-y-3">
              <TaskCard
                v-for="task in processingTasks"
                :key="task.task_id"
                :task="task"
                :is-active="currentTask?.task_id === task.task_id"
                @click="selectTask(task)"
              />
            </div>
          </div>
        </div>

        <div class="space-y-6">
          <div class="bg-card rounded-xl border border-border p-6 shadow-sm">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-2">
                <div class="w-1 h-6 bg-primary rounded-full" />
                <h2 class="text-base font-semibold text-foreground">任务列表</h2>
              </div>
              <span class="text-xs text-muted-foreground">{{ tasks.length }} 个任务</span>
            </div>

            <div class="space-y-3 max-h-[400px] overflow-y-auto scrollbar-thin">
              <TransitionGroup
                enter-active-class="transition duration-300 ease-out"
                enter-from-class="opacity-0 -translate-x-2"
                enter-to-class="opacity-100 translate-x-0"
                leave-active-class="transition duration-200 ease-in"
                leave-from-class="opacity-100 translate-x-0"
                leave-to-class="opacity-0 translate-x-2"
              >
                <TaskCard
                  v-for="task in tasks"
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

              <div v-if="tasks.length === 0" class="text-center py-8">
                <div class="w-12 h-12 bg-muted rounded-xl flex items-center justify-center mx-auto mb-3">
                  <FileText class="w-6 h-6 text-muted-foreground" />
                </div>
                <p class="text-sm text-muted-foreground">暂无任务</p>
              </div>
            </div>
          </div>

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
                <span class="font-semibold text-foreground">{{ totalTasks }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">成功次数</span>
                <span class="font-semibold text-emerald-600">{{ completedTasks.length }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-muted-foreground">处理中</span>
                <span class="font-semibold text-blue-600">{{ processingTasks.length }}</span>
              </div>
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
import { ref, computed, onMounted, watch } from 'vue'
import { useTaskStore } from '@/stores/task'
import MainLayout from '@/layouts/MainLayout.vue'
import StepIndicator from '@/components/common/StepIndicator.vue'
import FileUploader from '@/components/upload/FileUploader.vue'
import FormatSelector from '@/components/upload/FormatSelector.vue'
import TaskCard from '@/components/task/TaskCard.vue'
import ResultPreview from '@/components/result/ResultPreview.vue'
import ToastNotification from '@/components/common/ToastNotification.vue'
import {
  FileText,
  ArrowRight,
  ArrowLeft,
  Download,
  Copy,
  Rocket,
  Loader2,
  CheckCircle
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
const totalTasks = computed(() => tasks.value.length)
const hasCompletedTasks = computed(() => completedTasks.value.length > 0)
const hasProcessingTasks = computed(() => processingTasks.value.length > 0)

onMounted(async () => {
  await taskStore.fetchTasks()
})

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
    for (const file of selectedFiles.value) {
      await taskStore.createTask(file, outputFormat.value)
    }

    showToast(`已提交 ${selectedFiles.value.length} 个任务`, 'success')
    uploaderRef.value?.clearFiles()
    selectedFiles.value = []

    for (const task of processingTasks.value) {
      pollTaskStatus(task.task_id)
    }
  } catch (error: any) {
    showToast(error.message || '解析失败', 'error')
  }
}

function pollTaskStatus(taskId: string) {
  const poll = async () => {
    const task = tasks.value.find(t => t.task_id === taskId)
    if (!task) return

    if (task.status === 'completed' || task.status === 'failed') {
      if (task.status === 'completed') {
        await taskStore.fetchTaskResult(taskId)
        showToast('解析完成', 'success')
      }
      return
    }

    await taskStore.updateTaskStatus(taskId, task.status, (task.progress || 0) + 10)
    setTimeout(poll, 2000)
  }

  poll()
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

function handleCopy() {
  showToast('已复制到剪贴板', 'success')
}

function handleDownload(_format: string) {
  if (currentTask.value) {
    downloadTask(currentTask.value)
  }
}

async function downloadAllResults() {
  for (const task of completedTasks.value) {
    await downloadTask(task)
  }
}

async function copyAllResults() {
  for (const task of completedTasks.value) {
    await copyTask(task)
  }
}

function showToast(message: string, type: 'success' | 'error' | 'info' | 'warning' = 'success') {
  toast.value = { visible: true, message, type }
  setTimeout(() => {
    toast.value.visible = false
  }, 3000)
}
</script>
