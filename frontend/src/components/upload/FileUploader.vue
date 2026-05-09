<template>
  <div
    class="relative"
    @dragenter.prevent="handleDragEnter"
    @dragleave.prevent="handleDragLeave"
    @dragover.prevent
    @drop.prevent="handleDrop"
  >
    <div
      :class="[
        'border-2 border-dashed rounded-xl p-12 text-center transition-all duration-300 cursor-pointer',
        isDragOver
          ? 'border-blue-500 bg-blue-50'
          : 'border-gray-300 hover:border-blue-400 hover:bg-gray-50',
        isUploading ? 'pointer-events-none opacity-60' : '',
      ]"
      @click="triggerFileInput"
    >
      <input
        ref="fileInput"
        type="file"
        accept=".pdf"
        class="hidden"
        @change="handleFileSelect"
      />

      <div class="flex flex-col items-center gap-4">
        <div
          :class="[
            'w-16 h-16 rounded-full flex items-center justify-center transition-colors',
            isDragOver ? 'bg-blue-500' : 'bg-gray-100',
          ]"
        >
          <Upload
            :class="[
              'w-8 h-8 transition-colors',
              isDragOver ? 'text-white' : 'text-gray-400',
            ]"
          />
        </div>

        <div>
          <p class="text-lg font-medium text-gray-700">
            {{ isDragOver ? '释放文件开始上传' : '拖拽PDF文件到此处' }}
          </p>
          <p class="text-sm text-gray-500 mt-1">或点击选择文件</p>
        </div>

        <div class="flex items-center gap-2 text-xs text-gray-400">
          <FileText class="w-4 h-4" />
          <span>支持 PDF 文件，最大 10MB</span>
        </div>
      </div>
    </div>

    <div v-if="selectedFile" class="mt-4">
      <div class="bg-white border rounded-lg p-4 flex items-center gap-4">
        <div class="w-10 h-10 bg-red-100 rounded-lg flex items-center justify-center">
          <FileText class="w-5 h-5 text-red-500" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-gray-900 truncate">{{ selectedFile.name }}</p>
          <p class="text-xs text-gray-500">{{ formatFileSize(selectedFile.size) }}</p>
        </div>
        <button
          type="button"
          class="p-2 text-gray-400 hover:text-red-500 transition-colors"
          @click.stop="clearFile"
        >
          <X class="w-5 h-5" />
        </button>
      </div>
    </div>

    <div v-if="error" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
      <p class="text-sm text-red-600">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Upload, FileText, X } from 'lucide-vue-next'
import { formatFileSize } from '@/lib/utils'

const emit = defineEmits<{
  (e: 'file-selected', file: File): void
  (e: 'error', message: string): void
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const isDragOver = ref(false)
const error = ref('')
const isUploading = ref(false)

function triggerFileInput() {
  if (!isUploading.value) {
    fileInput.value?.click()
  }
}

function handleDragEnter(e: DragEvent) {
  isDragOver.value = true
}

function handleDragLeave(e: DragEvent) {
  isDragOver.value = false
}

function handleDrop(e: DragEvent) {
  isDragOver.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    handleFile(files[0])
  }
}

function handleFileSelect(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files && input.files.length > 0) {
    handleFile(input.files[0])
  }
}

function handleFile(file: File) {
  error.value = ''

  if (!file.name.toLowerCase().endsWith('.pdf')) {
    error.value = '只支持 PDF 格式的文件'
    emit('error', error.value)
    return
  }

  if (file.size > 10 * 1024 * 1024) {
    error.value = '文件大小不能超过 10MB'
    emit('error', error.value)
    return
  }

  selectedFile.value = file
  emit('file-selected', file)
}

function clearFile() {
  selectedFile.value = null
  error.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

function setUploading(value: boolean) {
  isUploading.value = value
}
</script>
