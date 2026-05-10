import { ref, computed } from 'vue'
import axios from 'axios'

export interface UploadFile {
  id: string
  file: File
  progress: number
  status: 'pending' | 'uploading' | 'completed' | 'error'
  error?: string
}

export interface UploadOptions {
  maxSize?: number
  allowedTypes?: string[]
  onProgress?: (fileId: string, progress: number) => void
  onComplete?: (fileId: string, response: any) => void
  onError?: (fileId: string, error: string) => void
}

const DEFAULT_MAX_SIZE = 100 * 1024 * 1024
const DEFAULT_ALLOWED_TYPES = ['.pdf', '.doc', '.docx', '.ppt', '.pptx']

export function useFileUpload(options: UploadOptions = {}) {
  const files = ref<UploadFile[]>([])
  const isUploading = ref(false)
  const maxSize = options.maxSize || DEFAULT_MAX_SIZE
  const allowedTypes = options.allowedTypes || DEFAULT_ALLOWED_TYPES

  const hasFiles = computed(() => files.value.length > 0)
  const uploadProgress = computed(() => {
    if (files.value.length === 0) return 0
    const total = files.value.reduce((sum, f) => sum + f.progress, 0)
    return Math.round(total / files.value.length)
  })

  function generateId(): string {
    return `file_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  }

  function validateFile(file: File): string | null {
    const ext = '.' + file.name.split('.').pop()?.toLowerCase()

    if (!allowedTypes.includes(ext)) {
      return `不支持 ${ext} 格式，请上传 PDF、DOC、DOCX、PPT 或 PPTX 文件`
    }

    if (file.size > maxSize) {
      return `${file.name} 超过 ${Math.round(maxSize / 1024 / 1024)}MB 限制`
    }

    return null
  }

  function addFiles(newFiles: File[]): { valid: File[]; errors: string[] } {
    const valid: File[] = []
    const errors: string[] = []

    for (const file of newFiles) {
      const error = validateFile(file)
      if (error) {
        errors.push(error)
        continue
      }

      const existsIndex = files.value.findIndex(f => f.file.name === file.name)
      if (existsIndex === -1) {
        files.value.push({
          id: generateId(),
          file,
          progress: 0,
          status: 'pending'
        })
      }
      valid.push(file)
    }

    return { valid, errors }
  }

  function removeFile(fileId: string) {
    const index = files.value.findIndex(f => f.id === fileId)
    if (index > -1) {
      files.value.splice(index, 1)
    }
  }

  function clearFiles() {
    files.value = []
  }

  async function uploadFile(
    fileId: string,
    endpoint: string,
    additionalData?: Record<string, any>
  ): Promise<any> {
    const fileRecord = files.value.find(f => f.id === fileId)
    if (!fileRecord) throw new Error('File not found')

    fileRecord.status = 'uploading'

    const formData = new FormData()
    formData.append('file', fileRecord.file)

    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        formData.append(key, value)
      })
    }

    try {
      const response = await axios.post(endpoint, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        },
        onUploadProgress: (e) => {
          const progress = Math.round((e.loaded * 100) / (e.total || 1))
          fileRecord.progress = progress
          options.onProgress?.(fileId, progress)
        }
      })

      fileRecord.status = 'completed'
      fileRecord.progress = 100
      options.onComplete?.(fileId, response.data)

      return response.data
    } catch (error: any) {
      fileRecord.status = 'error'
      const errorMsg = error.response?.data?.error || '上传失败'
      fileRecord.error = errorMsg
      options.onError?.(fileId, errorMsg)
      throw error
    }
  }

  async function uploadAll(endpoint: string, additionalData?: Record<string, any>) {
    if (isUploading.value) return

    isUploading.value = true

    try {
      const uploadPromises = files.value.map(f => uploadFile(f.id, endpoint, additionalData))
      await Promise.all(uploadPromises)
    } finally {
      isUploading.value = false
    }
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

  return {
    files,
    isUploading,
    hasFiles,
    uploadProgress,
    addFiles,
    removeFile,
    clearFiles,
    uploadFile,
    uploadAll,
    getFileColor,
    validateFile
  }
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
