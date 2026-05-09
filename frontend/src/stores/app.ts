import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const isUploading = ref(false)
  const uploadProgress = ref(0)

  function setUploading(value: boolean) {
    isUploading.value = value
    if (!value) {
      uploadProgress.value = 0
    }
  }

  function setUploadProgress(value: number) {
    uploadProgress.value = value
  }

  return {
    isUploading,
    uploadProgress,
    setUploading,
    setUploadProgress,
  }
})
