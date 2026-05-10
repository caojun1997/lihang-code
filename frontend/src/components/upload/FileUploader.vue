<template>
  <div
    class="relative"
    @dragenter.prevent="handleDragEnter"
    @dragleave.prevent="handleDragLeave"
    @dragover.prevent
    @drop.prevent="handleDrop"
    @paste="handlePaste"
  >
    <div
      :class="[
        'border-2 border-dashed rounded-2xl p-16 text-center transition-all duration-300 cursor-pointer relative overflow-hidden',
        isDragOver
          ? 'border-primary bg-primary/5 scale-[1.02]'
          : 'border-border hover:border-primary/50 hover:bg-muted/50',
        isUploading ? 'pointer-events-none' : '',
      ]"
      :style="isDragOver ? { borderStyle: 'dashed' } : {}"
      @click="triggerFileInput"
    >
      <input
        ref="fileInput"
        type="file"
        accept=".pdf,.doc,.docx,.ppt,.pptx"
        multiple
        class="hidden"
        @change="handleFileSelect"
      />

      <div class="flex flex-col items-center gap-6 animate-fade-in">
        <div
          :class="[
            'w-20 h-20 rounded-2xl flex items-center justify-center transition-all duration-300',
            isDragOver ? 'bg-primary scale-110 shadow-lg shadow-primary/30' : 'bg-muted',
          ]"
        >
          <Upload
            :class="[
              'w-10 h-10 transition-all duration-300',
              isDragOver ? 'text-white' : 'text-muted-foreground',
            ]"
          />
        </div>

        <div class="space-y-2">
          <p class="text-xl font-semibold text-foreground">
            {{ isDragOver ? '释放文件开始上传' : '拖拽文件到此处' }}
          </p>
          <p class="text-sm text-muted-foreground">
            或<span class="text-primary font-medium hover:underline">点击选择文件</span>
          </p>
        </div>

        <div class="flex items-center gap-4 text-xs text-muted-foreground/70">
          <span class="flex items-center gap-1.5 px-3 py-1.5 bg-muted/50 rounded-full">
            <FileText class="w-3.5 h-3.5" />
            PDF
          </span>
          <span class="flex items-center gap-1.5 px-3 py-1.5 bg-muted/50 rounded-full">
            <FileText class="w-3.5 h-3.5" />
            DOCX
          </span>
          <span class="flex items-center gap-1.5 px-3 py-1.5 bg-muted/50 rounded-full">
            <FileText class="w-3.5 h-3.5" />
            PPTX
          </span>
          <span class="px-3 py-1.5 bg-primary/10 text-primary rounded-full font-medium">
            最大 50MB
          </span>
        </div>

        <p v-if="isDragOver" class="text-sm text-primary font-medium animate-pulse">
          松开鼠标即可上传
        </p>
      </div>

      <div
        v-if="isDragOver"
        class="absolute inset-0 bg-gradient-to-br from-primary/10 to-transparent pointer-events-none"
      />
    </div>

    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0 translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-2"
    >
      <div v-if="selectedFiles.length > 0" class="mt-4 space-y-3">
        <div
          v-for="(file, index) in selectedFiles"
          :key="index"
          class="bg-card border rounded-xl p-4 flex items-center gap-4 hover:shadow-md transition-shadow"
        >
          <div
            :class="[
              'w-12 h-12 rounded-xl flex items-center justify-center',
              getFileColor(file.name),
            ]"
          >
            <FileText class="w-6 h-6 text-white" />
          </div>

          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-foreground truncate">{{ file.name }}</p>
            <p class="text-xs text-muted-foreground">{{ formatFileSize(file.size) }}</p>
          </div>

          <div v-if="uploadProgress[index] !== undefined" class="w-24">
            <div class="h-2 bg-muted rounded-full overflow-hidden">
              <div
                class="h-full bg-primary rounded-full transition-all duration-300"
                :style="{ width: `${uploadProgress[index]}%` }"
              />
            </div>
            <p class="text-xs text-muted-foreground text-center mt-1">{{ uploadProgress[index] }}%</p>
          </div>

          <button
            type="button"
            class="p-2 text-muted-foreground hover:text-destructive transition-colors"
            @click.stop="removeFile(index)"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
      </div>
    </Transition>

    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
    >
      <div v-if="error" class="mt-4 p-4 bg-destructive/10 border border-destructive/20 rounded-xl">
        <div class="flex items-start gap-3">
          <AlertCircle class="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
          <div>
            <p class="text-sm font-medium text-destructive">上传失败</p>
            <p class="text-xs text-destructive/80 mt-1">{{ error }}</p>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Upload, FileText, X, AlertCircle } from 'lucide-vue-next'
import { formatFileSize } from '@/lib/utils'

const emit = defineEmits<{
  (e: 'files-selected', files: File[]): void
  (e: 'error', message: string): void
  (e: 'upload-progress', index: number, progress: number): void
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFiles = ref<File[]>([])
const uploadProgress = ref<Record<number, number>>({})
const isDragOver = ref(false)
const isUploading = ref(false)
const error = ref('')

const MAX_FILE_SIZE = 50 * 1024 * 1024
const ALLOWED_EXTENSIONS = ['.pdf', '.doc', '.docx', '.ppt', '.pptx']

function triggerFileInput() {
  if (!isUploading.value) {
    fileInput.value?.click()
  }
}

function handleDragEnter(_e: DragEvent) {
  isDragOver.value = true
}

function handleDragLeave(e: DragEvent) {
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  if (
    e.clientX <= rect.left ||
    e.clientX >= rect.right ||
    e.clientY <= rect.top ||
    e.clientY >= rect.bottom
  ) {
    isDragOver.value = false
  }
}

function handleDrop(e: DragEvent) {
  isDragOver.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    handleFiles(Array.from(files))
  }
}

function handleFileSelect(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files && input.files.length > 0) {
    handleFiles(Array.from(input.files))
  }
}

function handlePaste(e: ClipboardEvent) {
  const items = e.clipboardData?.items
  if (!items) return

  const files: File[] = []
  for (const item of items) {
    if (item.type.startsWith('application/pdf') || item.type.includes('document')) {
      const file = item.getAsFile()
      if (file) files.push(file)
    }
  }

  if (files.length > 0) {
    handleFiles(files)
  }
}

function handleFiles(files: File[]) {
  error.value = ''

  for (const file of files) {
    const ext = '.' + file.name.split('.').pop()?.toLowerCase()

    if (!ALLOWED_EXTENSIONS.includes(ext)) {
      error.value = `不支持 ${ext} 格式，请上传 PDF、DOC、DOCX、PPT 或 PPTX 文件`
      emit('error', error.value)
      return
    }

    if (file.size > MAX_FILE_SIZE) {
      error.value = `${file.name} 超过 50MB 限制`
      emit('error', error.value)
      return
    }

    const existsIndex = selectedFiles.value.findIndex(f => f.name === file.name)
    if (existsIndex === -1) {
      selectedFiles.value.push(file)
    }
  }

  emit('files-selected', selectedFiles.value)
}

function removeFile(index: number) {
  selectedFiles.value.splice(index, 1)
  delete uploadProgress.value[index]
  emit('files-selected', selectedFiles.value)
}

function getFileColor(filename: string): string {
  const ext = filename.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'pdf':
      return 'bg-red-500'
    case 'doc':
    case 'docx':
      return 'bg-blue-600'
    case 'ppt':
    case 'pptx':
      return 'bg-orange-500'
    default:
      return 'bg-gray-500'
  }
}

function setUploading(value: boolean) {
  isUploading.value = value
}

function updateProgress(index: number, progress: number) {
  uploadProgress.value[index] = progress
  emit('upload-progress', index, progress)
}

function clearFiles() {
  selectedFiles.value = []
  uploadProgress.value = {}
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

defineExpose({
  setUploading,
  updateProgress,
  clearFiles,
})
</script>
