<template>
  <div
    :class="[
      'bg-card border rounded-2xl p-5 transition-all duration-300 cursor-pointer group',
      isActive
        ? 'border-primary shadow-lg shadow-primary/10 ring-2 ring-primary/20'
        : 'border-border hover:border-primary/30 hover:shadow-md',
    ]"
    @click="emit('click')"
  >
    <div class="flex items-start gap-4">
      <div
        :class="[
          'w-14 h-14 rounded-xl flex items-center justify-center flex-shrink-0 transition-all duration-300',
          statusConfig.bgColor,
          isActive ? 'scale-105' : 'group-hover:scale-105',
        ]"
      >
        <component :is="statusConfig.icon" class="w-7 h-7 text-white" />
      </div>

      <div class="flex-1 min-w-0 space-y-2">
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <h3 class="font-semibold text-foreground truncate group-hover:text-primary transition-colors">
              {{ task.file_name }}
            </h3>
            <p class="text-sm text-muted-foreground">
              {{ formatFileSize(task.file_size) }}
            </p>
          </div>
          <StatusBadge :status="task.status" class="flex-shrink-0" />
        </div>

        <div v-if="isProcessing" class="space-y-2">
          <div class="flex items-center justify-between text-xs text-muted-foreground">
            <span>{{ statusConfig.label }}</span>
            <span>{{ task.progress || 0 }}%</span>
          </div>
          <div class="h-2 bg-muted rounded-full overflow-hidden">
            <div
              :class="['h-full rounded-full transition-all duration-500', statusConfig.progressColor]"
              :style="{ width: `${task.progress || 0}%` }"
            />
          </div>
          <p v-if="task.current_page && task.total_pages" class="text-xs text-muted-foreground">
            正在处理第 {{ task.current_page }} / {{ task.total_pages }} 页
          </p>
        </div>

        <div class="flex items-center justify-between pt-2 border-t border-border/50">
          <span class="text-xs text-muted-foreground">
            {{ formatTime(task.created_at) }}
          </span>
          <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              v-if="task.status === 'completed'"
              class="p-1.5 text-muted-foreground hover:text-primary transition-colors"
              title="查看结果"
              @click.stop="emit('view')"
            >
              <Eye class="w-4 h-4" />
            </button>
            <button
              v-if="task.status === 'completed'"
              class="p-1.5 text-muted-foreground hover:text-primary transition-colors"
              title="下载"
              @click.stop="emit('download')"
            >
              <Download class="w-4 h-4" />
            </button>
            <button
              v-if="task.status === 'completed'"
              class="p-1.5 text-muted-foreground hover:text-primary transition-colors"
              title="复制"
              @click.stop="emit('copy')"
            >
              <Copy class="w-4 h-4" />
            </button>
            <button
              class="p-1.5 text-muted-foreground hover:text-destructive transition-colors"
              title="删除"
              @click.stop="emit('delete')"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
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
  Ban,
  Eye,
  Download,
  Copy,
  Trash2
} from 'lucide-vue-next'
import type { Task, TaskStatus } from '@/api/types'
import StatusBadge from './StatusBadge.vue'
import { formatFileSize } from '@/lib/utils'

const props = defineProps<{
  task: Task
  isActive?: boolean
}>()

const emit = defineEmits<{
  (e: 'click'): void
  (e: 'view'): void
  (e: 'download'): void
  (e: 'copy'): void
  (e: 'delete'): void
}>()

const statusConfig = computed(() => {
  const configs: Record<TaskStatus, {
    icon: any
    label: string
    bgColor: string
    progressColor: string
  }> = {
    pending: {
      icon: Clock,
      label: '等待中',
      bgColor: 'bg-gray-500',
      progressColor: 'bg-gray-500',
    },
    uploading: {
      icon: Upload,
      label: '上传中',
      bgColor: 'bg-blue-500',
      progressColor: 'bg-blue-500',
    },
    processing: {
      icon: Loader2,
      label: '解析中',
      bgColor: 'bg-blue-500 animate-pulse',
      progressColor: 'bg-gradient-to-r from-blue-500 to-primary',
    },
    completed: {
      icon: CheckCircle2,
      label: '已完成',
      bgColor: 'bg-green-500',
      progressColor: 'bg-green-500',
    },
    partial_failed: {
      icon: AlertTriangle,
      label: '部分失败',
      bgColor: 'bg-orange-500',
      progressColor: 'bg-orange-500',
    },
    failed: {
      icon: XCircle,
      label: '失败',
      bgColor: 'bg-red-500',
      progressColor: 'bg-red-500',
    },
    cancelled: {
      icon: Ban,
      label: '已取消',
      bgColor: 'bg-gray-400',
      progressColor: 'bg-gray-400',
    },
  }
  return configs[props.task.status] || configs.pending
})

const isProcessing = computed(() => 
  ['pending', 'uploading', 'processing'].includes(props.task.status)
)

function formatTime(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  return date.toLocaleDateString('zh-CN')
}
</script>
