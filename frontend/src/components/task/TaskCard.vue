<template>
  <div class="bg-white rounded-xl border shadow-sm overflow-hidden">
    <div class="p-6">
      <div class="flex items-start justify-between mb-4">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center">
            <FileText class="w-6 h-6 text-red-500" />
          </div>
          <div>
            <h3 class="font-medium text-gray-900">{{ task.file_name }}</h3>
            <p class="text-sm text-gray-500">{{ formatDate(task.created_at) }}</p>
          </div>
        </div>
        <TaskStatus :status="task.status" />
      </div>

      <div class="flex items-center gap-6 text-sm text-gray-500 mb-4">
        <div class="flex items-center gap-1">
          <HardDrive class="w-4 h-4" />
          <span>{{ formatFileSize(task.file_size) }}</span>
        </div>
        <div class="flex items-center gap-1">
          <Tag class="w-4 h-4" />
          <span>{{ task.output_format === 'markdown' ? 'Markdown' : '纯文本' }}</span>
        </div>
      </div>

      <div v-if="task.status === 'processing'" class="mb-4">
        <div class="flex items-center justify-between text-sm mb-2">
          <span class="text-gray-500">处理进度</span>
          <span class="font-medium text-blue-600">{{ progress }}%</span>
        </div>
        <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
          <div
            class="h-full bg-blue-500 rounded-full transition-all duration-300"
            :style="{ width: `${progress}%` }"
          />
        </div>
      </div>

      <div v-if="task.error_message" class="mb-4 p-3 bg-red-50 rounded-lg">
        <p class="text-sm text-red-600">{{ task.error_message }}</p>
      </div>

      <div class="flex items-center gap-3">
        <button
          v-if="task.status === 'completed'"
          type="button"
          class="flex-1 bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors"
          @click="emit('view')"
        >
          查看结果
        </button>
        <button
          v-if="task.status === 'completed'"
          type="button"
          class="flex items-center gap-1 bg-gray-100 text-gray-700 px-4 py-2 rounded-lg text-sm font-medium hover:bg-gray-200 transition-colors"
          @click="emit('download')"
        >
          <Download class="w-4 h-4" />
          下载
        </button>
        <button
          type="button"
          class="flex items-center gap-1 text-red-600 px-4 py-2 rounded-lg text-sm font-medium hover:bg-red-50 transition-colors"
          @click="emit('delete')"
        >
          <Trash2 class="w-4 h-4" />
          删除
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { FileText, HardDrive, Tag, Download, Trash2 } from 'lucide-vue-next'
import { formatFileSize, formatDate } from '@/lib/utils'
import type { Task } from '@/api/types'
import TaskStatus from './TaskStatus.vue'

const props = defineProps<{
  task: Task
  progress?: number
}>()

const emit = defineEmits<{
  (e: 'view'): void
  (e: 'download'): void
  (e: 'delete'): void
}>()
</script>
