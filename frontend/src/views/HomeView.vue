<template>
  <div class="min-h-[calc(100vh-200px)]">
    <div class="container mx-auto px-4 py-8 max-w-4xl">
      <div class="text-center mb-8">
        <h2 class="text-3xl font-bold text-gray-900 mb-2">PDF文档解析</h2>
        <p class="text-gray-600">快速将PDF文档转换为Markdown或纯文本格式</p>
      </div>

      <div class="grid md:grid-cols-3 gap-8">
        <div class="md:col-span-2 space-y-6">
          <div class="bg-white rounded-xl border shadow-sm p-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">上传PDF文件</h3>

            <FileUploader
              @file-selected="handleFileSelected"
              @error="handleError"
            />

            <div v-if="selectedFile" class="mt-6 space-y-4">
              <FormatSelector v-model="outputFormat" />

              <button
                type="button"
                :disabled="isUploading"
                class="w-full bg-blue-600 text-white px-6 py-3 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                @click="handleUpload"
              >
                <Loader2 v-if="isUploading" class="w-5 h-5 animate-spin" />
                <Upload v-else class="w-5 h-5" />
                {{ isUploading ? '上传中...' : '开始解析' }}
              </button>
            </div>

            <div v-if="error" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p class="text-sm text-red-600">{{ error }}</p>
            </div>

            <div v-if="successMessage" class="mt-4 p-4 bg-green-50 border border-green-200 rounded-lg">
              <p class="text-sm text-green-600">{{ successMessage }}</p>
            </div>
          </div>

          <div v-if="recentTasks.length > 0" class="space-y-4">
            <h3 class="text-lg font-medium text-gray-900">最近任务</h3>
            <TaskList
              :tasks="recentTasks"
              :loading="isLoadingTasks"
              :total="totalTasks"
              :page="page"
              :page-size="pageSize"
              @page-change="handlePageChange"
            />
          </div>
        </div>

        <div class="space-y-6">
          <div class="bg-white rounded-xl border shadow-sm p-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">功能特点</h3>
            <ul class="space-y-3">
              <li class="flex items-start gap-3">
                <div class="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                  <Check class="w-4 h-4 text-blue-600" />
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-700">智能解析</p>
                  <p class="text-xs text-gray-500">保留文档结构和格式</p>
                </div>
              </li>
              <li class="flex items-start gap-3">
                <div class="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                  <Check class="w-4 h-4 text-blue-600" />
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-700">表格识别</p>
                  <p class="text-xs text-gray-500">自动识别并转换表格</p>
                </div>
              </li>
              <li class="flex items-start gap-3">
                <div class="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                  <Check class="w-4 h-4 text-blue-600" />
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-700">公式转换</p>
                  <p class="text-xs text-gray-500">支持数学公式识别</p>
                </div>
              </li>
              <li class="flex items-start gap-3">
                <div class="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                  <Check class="w-4 h-4 text-blue-600" />
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-700">多格式输出</p>
                  <p class="text-xs text-gray-500">支持Markdown和纯文本</p>
                </div>
              </li>
            </ul>
          </div>

          <div class="bg-white rounded-xl border shadow-sm p-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">使用说明</h3>
            <ol class="space-y-3 text-sm text-gray-600">
              <li class="flex gap-3">
                <span class="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center flex-shrink-0 text-xs font-medium">1</span>
                <span>上传PDF文件</span>
              </li>
              <li class="flex gap-3">
                <span class="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center flex-shrink-0 text-xs font-medium">2</span>
                <span>选择输出格式</span>
              </li>
              <li class="flex gap-3">
                <span class="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center flex-shrink-0 text-xs font-medium">3</span>
                <span>等待解析完成</span>
              </li>
              <li class="flex gap-3">
                <span class="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center flex-shrink-0 text-xs font-medium">4</span>
                <span>查看或下载结果</span>
              </li>
            </ol>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Upload, Loader2, Check } from 'lucide-vue-next'
import { useTaskStore } from '@/stores/task'
import FileUploader from '@/components/upload/FileUploader.vue'
import FormatSelector from '@/components/upload/FormatSelector.vue'
import TaskList from '@/components/task/TaskList.vue'

const taskStore = useTaskStore()

const selectedFile = ref<File | null>(null)
const outputFormat = ref<'markdown' | 'txt'>('markdown')
const isUploading = ref(false)
const isLoadingTasks = ref(false)
const error = ref('')
const successMessage = ref('')

const recentTasks = ref(taskStore.tasks)
const totalTasks = ref(taskStore.total)
const page = ref(1)
const pageSize = ref(5)

async function handleFileSelected(file: File) {
  selectedFile.value = file
  error.value = ''
  successMessage.value = ''
}

function handleError(message: string) {
  error.value = message
}

async function handleUpload() {
  if (!selectedFile.value) return

  isUploading.value = true
  error.value = ''
  successMessage.value = ''

  try {
    const result = await taskStore.createTask(selectedFile.value, outputFormat.value)
    successMessage.value = `任务创建成功！任务ID: ${result?.task_id}`
    selectedFile.value = null

    await loadTasks()
  } catch (e: any) {
    error.value = e.message || '上传失败，请重试'
  } finally {
    isUploading.value = false
  }
}

async function loadTasks() {
  isLoadingTasks.value = true
  await taskStore.fetchTasks({ page: page.value, page_size: pageSize.value })
  recentTasks.value = taskStore.tasks
  totalTasks.value = taskStore.total
  isLoadingTasks.value = false
}

function handlePageChange(newPage: number) {
  page.value = newPage
  loadTasks()
}

onMounted(() => {
  loadTasks()
})
</script>
