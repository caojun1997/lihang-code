<template>
  <span :class="['inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium', config.bgColor, config.textColor]">
    <component :is="config.icon" class="w-3.5 h-3.5" />
    {{ config.label }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { 
  Clock, 
  Upload, 
  Loader2, 
  CheckCircle2, 
  AlertTriangle, 
  XCircle, 
  Ban
} from 'lucide-vue-next'
import type { TaskStatus } from '@/api/types'

const props = defineProps<{
  status: TaskStatus
}>()

const config = computed(() => {
  const configs: Record<TaskStatus, {
    icon: any
    label: string
    bgColor: string
    textColor: string
  }> = {
    pending: {
      icon: Clock,
      label: '等待中',
      bgColor: 'bg-gray-100 dark:bg-gray-800',
      textColor: 'text-gray-600 dark:text-gray-300',
    },
    uploading: {
      icon: Upload,
      label: '上传中',
      bgColor: 'bg-blue-100 dark:bg-blue-900',
      textColor: 'text-blue-600 dark:text-blue-300',
    },
    processing: {
      icon: Loader2,
      label: '处理中',
      bgColor: 'bg-blue-100 dark:bg-blue-900',
      textColor: 'text-blue-600 dark:text-blue-300',
    },
    completed: {
      icon: CheckCircle2,
      label: '已完成',
      bgColor: 'bg-green-100 dark:bg-green-900',
      textColor: 'text-green-600 dark:text-green-300',
    },
    partial_failed: {
      icon: AlertTriangle,
      label: '部分失败',
      bgColor: 'bg-orange-100 dark:bg-orange-900',
      textColor: 'text-orange-600 dark:text-orange-300',
    },
    failed: {
      icon: XCircle,
      label: '失败',
      bgColor: 'bg-red-100 dark:bg-red-900',
      textColor: 'text-red-600 dark:text-red-300',
    },
    cancelled: {
      icon: Ban,
      label: '已取消',
      bgColor: 'bg-gray-100 dark:bg-gray-800',
      textColor: 'text-gray-500 dark:text-gray-400',
    },
  }
  return configs[props.status] || configs.pending
})
</script>
