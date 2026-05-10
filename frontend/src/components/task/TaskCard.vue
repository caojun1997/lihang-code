<template>
  <div
    :class="[
      'group relative p-4 rounded-xl border transition-all duration-200 cursor-pointer',
      'hover:shadow-md hover:border-primary/30',
      isActive
        ? 'border-primary bg-primary/5 shadow-sm'
        : 'border-border bg-card hover:bg-muted/30'
    ]"
    @click="emit('click')"
  >
    <div class="flex items-start gap-3">
      <div
        :class="[
          'w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0',
          fileColorClass
        ]"
      >
        <FileText class="w-5 h-5 text-white" />
      </div>

      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 mb-1">
          <p class="text-sm font-medium text-foreground truncate">
            {{ task.file_name }}
          </p>
          <StatusBadge :status="task.status" size="sm" />
        </div>

        <div class="flex items-center gap-3 text-xs text-muted-foreground">
          <span>{{ formatSize(task.file_size) }}</span>
          <span>·</span>
          <span>{{ formatTime(task.created_at) }}</span>
        </div>

        <div v-if="task.status === 'processing' || task.status === 'uploading'" class="mt-2">
          <ProgressBar
            :value="task.progress || 0"
            size="sm"
            :show-percentage="true"
          />
        </div>
      </div>
    </div>

    <div
      v-if="showActions"
      class="absolute right-2 top-2 opacity-0 group-hover:opacity-100 transition-opacity"
    >
      <div class="flex items-center gap-1 p-1 bg-background/95 backdrop-blur-sm rounded-lg shadow-sm border">
        <button
          v-if="task.status === 'completed'"
          class="p-1.5 rounded hover:bg-muted transition-colors"
          title="查看"
          @click.stop="emit('view')"
        >
          <Eye class="w-4 h-4 text-muted-foreground" />
        </button>
        <button
          v-if="task.status === 'completed'"
          class="p-1.5 rounded hover:bg-muted transition-colors"
          title="下载"
          @click.stop="emit('download')"
        >
          <Download class="w-4 h-4 text-muted-foreground" />
        </button>
        <button
          v-if="task.status === 'completed'"
          class="p-1.5 rounded hover:bg-muted transition-colors"
          title="复制"
          @click.stop="emit('copy')"
        >
          <Copy class="w-4 h-4 text-muted-foreground" />
        </button>
        <button
          class="p-1.5 rounded hover:bg-muted transition-colors"
          title="删除"
          @click.stop="emit('delete')"
        >
          <Trash2 class="w-4 h-4 text-muted-foreground hover:text-destructive" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { FileText, Eye, Download, Copy, Trash2 } from 'lucide-vue-next'
import StatusBadge from './StatusBadge.vue'
import ProgressBar from '@/components/common/ProgressBar.vue'
import type { Task } from '@/api/types'

interface Props {
  task: Task
  isActive?: boolean
  showActions?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isActive: false,
  showActions: true
})

const emit = defineEmits<{
  (e: 'click'): void
  (e: 'view'): void
  (e: 'download'): void
  (e: 'copy'): void
  (e: 'delete'): void
}>()

const fileColorClass = computed(() => {
  const ext = props.task.file_name.split('.').pop()?.toLowerCase()
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
})

function formatSize(bytes?: number): string {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function formatTime(dateString?: string): string {
  if (!dateString) return ''
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`

  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}
</script>
