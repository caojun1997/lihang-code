<template>
  <div
    :class="[
      'inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium transition-colors',
      badgeClasses
    ]"
  >
    <component :is="iconComponent" v-if="showIcon" class="w-3.5 h-3.5" />
    <span>{{ label }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  Clock,
  Upload,
  Loader2,
  CheckCircle,
  AlertTriangle,
  XCircle,
  Ban
} from 'lucide-vue-next'
import type { TaskStatus } from '@/api/types'

interface Props {
  status: TaskStatus
  showIcon?: boolean
  size?: 'sm' | 'md'
}

const props = withDefaults(defineProps<Props>(), {
  showIcon: true,
  size: 'md'
})

const statusConfig: Record<TaskStatus, { label: string; class: string }> = {
  pending: {
    label: '等待中',
    class: 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'
  },
  uploading: {
    label: '上传中',
    class: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
  },
  processing: {
    label: '解析中',
    class: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
  },
  completed: {
    label: '已完成',
    class: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400'
  },
  partial_failed: {
    label: '部分失败',
    class: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
  },
  failed: {
    label: '解析失败',
    class: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
  },
  cancelled: {
    label: '已取消',
    class: 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'
  }
}

const badgeConfig = computed(() => statusConfig[props.status] || statusConfig.pending)
const label = computed(() => badgeConfig.value.label)
const badgeClasses = computed(() => badgeConfig.value.class)

const iconComponent = computed(() => {
  const icons: Record<TaskStatus, any> = {
    pending: Clock,
    uploading: Upload,
    processing: Loader2,
    completed: CheckCircle,
    partial_failed: AlertTriangle,
    failed: XCircle,
    cancelled: Ban
  }
  return icons[props.status] || Clock
})
</script>
