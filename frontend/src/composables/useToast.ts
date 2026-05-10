import { ref } from 'vue'

export interface ToastOptions {
  message: string
  type?: 'success' | 'error' | 'info' | 'warning'
  duration?: number
}

export interface Toast extends ToastOptions {
  id: string
}

const toasts = ref<Toast[]>([])

let toastId = 0

export function useToast() {
  function show(options: ToastOptions): string {
    const id = `toast_${++toastId}`
    const toast: Toast = {
      id,
      type: 'info',
      duration: 3000,
      ...options
    }

    toasts.value.push(toast)

    if (toast.duration && toast.duration > 0) {
      setTimeout(() => {
        dismiss(id)
      }, toast.duration)
    }

    return id
  }

  function dismiss(id: string) {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }

  function success(message: string, duration?: number) {
    return show({ message, type: 'success', duration })
  }

  function error(message: string, duration?: number) {
    return show({ message, type: 'error', duration: duration || 5000 })
  }

  function info(message: string, duration?: number) {
    return show({ message, type: 'info', duration })
  }

  function warning(message: string, duration?: number) {
    return show({ message, type: 'warning', duration: duration || 4000 })
  }

  function clear() {
    toasts.value = []
  }

  return {
    toasts,
    show,
    dismiss,
    success,
    error,
    info,
    warning,
    clear
  }
}

export function useClipboard() {
  const copied = ref(false)
  const copyTimeout = ref<number | null>(null)

  async function copy(text: string): Promise<boolean> {
    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(text)
      } else {
        const textArea = document.createElement('textarea')
        textArea.value = text
        textArea.style.position = 'fixed'
        textArea.style.left = '-999999px'
        textArea.style.top = '-999999px'
        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()
        document.execCommand('copy')
        textArea.remove()
      }

      copied.value = true

      if (copyTimeout.value) {
        clearTimeout(copyTimeout.value)
      }

      copyTimeout.value = window.setTimeout(() => {
        copied.value = false
      }, 2000)

      return true
    } catch (err) {
      console.error('复制失败:', err)
      return false
    }
  }

  function resetCopied() {
    copied.value = false
    if (copyTimeout.value) {
      clearTimeout(copyTimeout.value)
      copyTimeout.value = null
    }
  }

  return {
    copied,
    copy,
    resetCopied
  }
}
