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
        'border-2 border-dashed rounded-xl p-12 text-center transition-all duration-200 cursor-pointer',
        'hover:border-primary/50 hover:bg-muted/30',
        props.uploading
          ? 'border-primary bg-primary/5 scale-[1.01]'
          : 'border-border',
        props.uploading ? 'pointer-events-none opacity-75' : ''
      ]"
      @click="!props.uploading && triggerFileInput()"
    >
      <input
        ref="fileInput"
        type="file"
        accept=".pdf,.doc,.docx,.ppt,.pptx"
        multiple
        class="hidden"
        @change="handleFileSelect"
      />

      <div class="flex flex-col items-center gap-5">
        <div
          :class="[
            'w-16 h-16 rounded-2xl flex items-center justify-center transition-all duration-200',
            isDragOver ? 'bg-primary scale-110 shadow-lg shadow-primary/25' : 'bg-muted'
          ]"
        >
          <Upload
            :class="[
              'w-8 h-8 transition-colors',
              isDragOver ? 'text-white' : 'text-muted-foreground'
            ]"
          />
        </div>

        <div class="space-y-2">
          <p class="text-lg font-semibold text-foreground">
            {{ isDragOver ? '释放文件开始上传' : '拖拽文件到此处' }}
          </p>
          <p class="text-sm text-muted-foreground">
            或<span class="text-primary font-medium hover:underline">点击选择文件</span>
          </p>
        </div>

        <div class="flex flex-wrap items-center justify-center gap-2 text-xs text-muted-foreground">
          <span class="px-2.5 py-1 bg-muted/60 rounded-full flex items-center gap-1.5">
            <FileText class="w-3.5 h-3.5" />
            PDF
          </span>
          <span class="px-2.5 py-1 bg-muted/60 rounded-full flex items-center gap-1.5">
            <FileText class="w-3.5 h-3.5" />
            DOCX
          </span>
          <span class="px-2.5 py-1 bg-muted/60 rounded-full flex items-center gap-1.5">
            <FileText class="w-3.5 h-3.5" />
            PPTX
          </span>
          <span class="px-2.5 py-1 bg-primary/10 text-primary rounded-full font-medium">
            最大 {{ maxSizeMB }}MB
          </span>
        </div>
      </div>
    </div>

    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0 translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-2"
    >
      <div v-if="localFiles.length > 0" class="mt-4 space-y-2">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-foreground">已选择 {{ localFiles.length }} 个文件</span>
          <button
            class="text-xs text-muted-foreground hover:text-foreground transition-colors"
            @click="clearAll"
          >
            清空全部
          </button>
        </div>

        <div
          v-for="(file, index) in localFiles"
          :key="index"
          class="bg-card border rounded-lg p-3 flex items-center gap-3 hover:shadow-sm transition-shadow"
        >
          <div
            :class="[
              'w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0',
              getFileColor(file.name)
            ]"
          >
            <FileText class="w-4 h-4 text-white" />
          </div>

          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-foreground truncate">{{ file.name }}</p>
            <p class="text-xs text-muted-foreground">{{ formatFileSize(file.size) }}</p>
          </div>

          <div v-if="uploadProgress[index] !== undefined" class="w-20">
            <div class="h-1.5 bg-muted rounded-full overflow-hidden">
              <div
                class="h-full bg-primary rounded-full transition-all duration-200"
                :style="{ width: `${uploadProgress[index]}%` }"
              />
            </div>
          </div>

          <button
            type="button"
            class="p-1.5 text-muted-foreground hover:text-destructive transition-colors"
            @click="removeFile(index)"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>
    </Transition>

    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
    >
      <div v-if="error" class="mt-4 p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
        <div class="flex items-start gap-2">
          <AlertCircle class="w-4 h-4 text-destructive flex-shrink-0 mt-0.5" />
          <p class="text-sm text-destructive">{{ error }}</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Upload, FileText, X, AlertCircle } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  modelValue?: File[]
  maxSizeMB?: number
  uploading?: boolean
}>(), {
  maxSizeMB: 100,
  uploading: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', files: File[]): void
  (e: 'files-selected', files: File[]): void
  (e: 'error', message: string): void
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const localFiles = ref<File[]>([])
const uploadProgress = ref<Record<number, number>>({})
const isDragOver = ref(false)
const error = ref('')

const ALLOWED_EXTENSIONS = ['.pdf', '.doc', '.docx', '.ppt', '.pptx']

watch(() => props.modelValue, (newFiles) => {
  if (newFiles) {
    localFiles.value = [...newFiles]
  }
}, { immediate: true })

function triggerFileInput() {
  fileInput.value?.click()
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
  const validFiles: File[] = []

  for (const file of files) {
    const ext = '.' + file.name.split('.').pop()?.toLowerCase()

    if (!ALLOWED_EXTENSIONS.includes(ext)) {
      error.value = `不支持 ${ext} 格式，请上传 PDF、DOC、DOCX、PPT 或 PPTX 文件`
      emit('error', error.value)
      continue
    }

    if (file.size > props.maxSizeMB * 1024 * 1024) {
      error.value = `${file.name} 超过 ${props.maxSizeMB}MB 限制`
      emit('error', error.value)
      continue
    }

    const existsIndex = localFiles.value.findIndex(f => f.name === file.name)
    if (existsIndex === -1) {
      localFiles.value.push(file)
      validFiles.push(file)
    }
  }

  if (validFiles.length > 0) {
    emit('update:modelValue', localFiles.value)
    emit('files-selected', localFiles.value)
  }
}

function removeFile(index: number) {
  localFiles.value.splice(index, 1)
  delete uploadProgress.value[index]
  emit('update:modelValue', localFiles.value)
  emit('files-selected', localFiles.value)
}

function clearAll() {
  localFiles.value = []
  uploadProgress.value = {}
  if (fileInput.value) {
    fileInput.value.value = ''
  }
  emit('update:modelValue', [])
  emit('files-selected', [])
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

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

defineExpose({
  updateProgress: (index: number, progress: number) => {
    uploadProgress.value[index] = progress
  },
  clearFiles: clearAll
})
</script>
