<template>
  <div
    :class="[
      'inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium',
      statusClasses[status],
    ]"
  >
    <component :is="statusIcons[status]" class="w-3 h-3" />
    {{ statusText[status] }}
  </div>
</template>

<script setup lang="ts">
import { Clock, Loader2, CheckCircle2, XCircle } from 'lucide-vue-next'

type TaskStatus = 'pending' | 'processing' | 'completed' | 'failed'

defineProps<{
  status: TaskStatus
}>()

const statusClasses: Record<TaskStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-700',
  processing: 'bg-blue-100 text-blue-700',
  completed: 'bg-green-100 text-green-700',
  failed: 'bg-red-100 text-red-700',
}

const statusText: Record<TaskStatus, string> = {
  pending: '等待中',
  processing: '处理中',
  completed: '已完成',
  failed: '失败',
}

const statusIcons: Record<TaskStatus, any> = {
  pending: Clock,
  processing: Loader2,
  completed: CheckCircle2,
  failed: XCircle,
}
</script>
