<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="task"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
        @click.self="emit('close')"
      >
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 scale-95"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
        >
          <div class="w-full max-w-3xl max-h-[90vh] bg-card rounded-2xl shadow-xl border border-border overflow-hidden flex flex-col">
            <div class="flex items-center justify-between px-6 py-4 border-b border-border">
              <div class="flex items-center gap-3">
                <div
                  :class="[
                    'w-10 h-10 rounded-lg flex items-center justify-center',
                    fileColorClass
                  ]"
                >
                  <FileText class="w-5 h-5 text-white" />
                </div>
                <div>
                  <h2 class="font-semibold text-foreground">{{ task.file_name }}</h2>
                  <div class="flex items-center gap-2 text-xs text-muted-foreground">
                    <span>{{ formatFileSize(task.file_size) }}</span>
                    <span>·</span>
                    <StatusBadge :status="task.status" size="sm" />
                  </div>
                </div>
              </div>
              <button
                class="p-2 rounded-lg hover:bg-muted transition-colors"
                @click="emit('close')"
              >
                <X class="w-5 h-5 text-muted-foreground" />
              </button>
            </div>

            <div class="flex-1 overflow-auto p-6">
              <div v-if="loading" class="space-y-3">
                <div class="h-4 bg-muted rounded animate-shimmer w-3/4" />
                <div class="h-4 bg-muted rounded animate-shimmer w-full" />
                <div class="h-4 bg-muted rounded animate-shimmer w-5/6" />
              </div>

              <div v-else-if="result" class="space-y-4">
                <div class="flex items-center justify-between">
                  <span class="text-sm text-muted-foreground">解析结果</span>
                  <div class="flex items-center gap-2">
                    <button
                      class="px-3 py-1.5 text-sm bg-muted rounded-lg hover:bg-muted/80 transition-colors flex items-center gap-1.5"
                      @click="copyContent"
                    >
                      <Copy class="w-3.5 h-3.5" />
                      复制
                    </button>
                    <button
                      class="px-3 py-1.5 text-sm bg-primary text-primary-foreground rounded-lg hover:bg-primary/90 transition-colors flex items-center gap-1.5"
                      @click="emit('download')"
                    >
                      <Download class="w-3.5 h-3.5" />
                      下载
                    </button>
                  </div>
                </div>

                <div class="bg-muted/30 rounded-lg p-4 max-h-96 overflow-auto scrollbar-thin">
                  <pre class="text-sm text-foreground whitespace-pre-wrap break-words font-mono">{{ result.content?.slice(0, 5000) }}{{ (result.content?.length || 0) > 5000 ? '\n\n... (内容过长，请下载查看完整内容)' : '' }}</pre>
                </div>

                <div class="flex items-center gap-4 text-xs text-muted-foreground">
                  <span>字数: {{ result.content?.length || 0 }}</span>
                  <span>格式: {{ task.output_format?.toUpperCase() }}</span>
                  <span>创建: {{ formatTime(task.created_at) }}</span>
                </div>
              </div>

              <div v-else class="text-center py-8">
                <div class="w-12 h-12 bg-muted rounded-xl flex items-center justify-center mx-auto mb-3">
                  <FileText class="w-6 h-6 text-muted-foreground/50" />
                </div>
                <p class="text-muted-foreground">暂无解析结果</p>
              </div>
            </div>

            <div class="px-6 py-4 border-t border-border bg-muted/30">
              <div class="flex items-center justify-between">
                <div class="text-xs text-muted-foreground">
                  任务ID: {{ task.task_id }}
                </div>
                <button
                  class="px-4 py-2 bg-muted rounded-lg text-sm hover:bg-muted/80 transition-colors"
                  @click="emit('close')"
                >
                  关闭
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { FileText, X, Copy, Download } from 'lucide-vue-next'
import StatusBadge from './StatusBadge.vue'
import type { Task, TaskResult } from '@/api/types'

interface Props {
  task: Task | null
  result: TaskResult | null
  loading?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'download'): void
}>()

const fileColorClass = computed(() => {
  if (!props.task) return 'bg-gray-500'
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

function formatFileSize(bytes?: number): string {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function formatTime(dateString?: string): string {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function copyContent() {
  if (props.result?.content) {
    await navigator.clipboard.writeText(props.result.content)
  }
}
</script>
