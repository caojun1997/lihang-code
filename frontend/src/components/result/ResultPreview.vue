<template>
  <div class="h-full flex flex-col bg-card rounded-xl border border-border overflow-hidden">
    <div class="flex items-center justify-between px-4 py-3 border-b border-border bg-muted/30">
      <div class="flex items-center gap-3">
        <FileText class="w-5 h-5 text-primary" />
        <div>
          <h3 class="font-semibold text-foreground text-sm">{{ fileName || '预览' }}</h3>
          <p class="text-xs text-muted-foreground">{{ wordCount }} 字 · {{ format?.toUpperCase() }}</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <div class="flex bg-muted rounded-lg p-0.5">
          <button
            v-for="fmt in formats"
            :key="fmt"
            :class="[
              'px-2.5 py-1 text-xs font-medium rounded-md transition-all',
              selectedFormat === fmt
                ? 'bg-primary text-primary-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            ]"
            @click="selectedFormat = fmt"
          >
            {{ fmt.toUpperCase() }}
          </button>
        </div>

        <div class="w-px h-5 bg-border mx-1" />

        <button
          class="p-1.5 text-muted-foreground hover:text-primary transition-colors"
          title="复制全部"
          @click="copyContent"
        >
          <Copy class="w-4 h-4" />
        </button>
        <button
          class="p-1.5 text-muted-foreground hover:text-primary transition-colors"
          title="下载"
          @click="downloadContent"
        >
          <Download class="w-4 h-4" />
        </button>
      </div>
    </div>

    <div class="flex-1 overflow-auto scrollbar-thin">
      <div class="max-w-3xl mx-auto p-6">
        <template v-if="loading">
          <div class="space-y-4">
            <div class="h-6 bg-muted rounded animate-shimmer w-3/4" />
            <div class="h-4 bg-muted rounded animate-shimmer w-full" />
            <div class="h-4 bg-muted rounded animate-shimmer w-5/6" />
            <div class="h-4 bg-muted rounded animate-shimmer w-full" />
            <div class="h-6 bg-muted rounded animate-shimmer w-2/3 mt-6" />
            <div class="h-4 bg-muted rounded animate-shimmer w-full" />
            <div class="h-4 bg-muted rounded animate-shimmer w-4/5" />
          </div>
        </template>

        <template v-else-if="content">
          <div
            class="prose-custom"
            v-html="renderedContent"
          />
        </template>

        <template v-else>
          <div class="text-center py-16">
            <div class="w-16 h-16 bg-muted rounded-2xl flex items-center justify-center mx-auto mb-4">
              <FileText class="w-8 h-8 text-muted-foreground/50" />
            </div>
            <p class="text-muted-foreground">暂无预览内容</p>
            <p class="text-sm text-muted-foreground/70 mt-1">
              上传文件并解析后即可预览结果
            </p>
          </div>
        </template>
      </div>
    </div>

    <div class="px-4 py-3 border-t border-border bg-muted/30">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4 text-xs text-muted-foreground">
          <button
            :class="[
              'flex items-center gap-1.5 transition-colors',
              { 'text-primary': showSearch }
            ]"
            @click="showSearch = !showSearch"
          >
            <Search class="w-3.5 h-3.5" />
            搜索
          </button>
          <span>{{ wordCount }} 字</span>
          <span v-if="createdAt">{{ formatTime(createdAt) }}</span>
        </div>

        <div v-if="searchQuery" class="text-xs text-primary">
          找到 {{ searchCount }} 处
        </div>
      </div>

      <Transition
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="opacity-0 -translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 -translate-y-2"
      >
        <div v-if="showSearch" class="mt-3 flex items-center gap-2">
          <div class="relative flex-1">
            <Search class="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-muted-foreground" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索内容..."
              class="w-full pl-8 pr-3 py-2 text-sm border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary/50"
            />
          </div>
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { FileText, Copy, Download, Search } from 'lucide-vue-next'
import { marked } from 'marked'

const props = defineProps<{
  content?: string
  fileName?: string
  format?: string
  createdAt?: string
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'copy'): void
  (e: 'download', format: string): void
}>()

const formats = ['md', 'txt', 'json'] as const
const selectedFormat = ref<'md' | 'txt' | 'json'>('md')
const showSearch = ref(false)
const searchQuery = ref('')

const wordCount = computed(() => {
  if (!props.content) return 0
  return props.content.replace(/\s/g, '').length
})

const renderedContent = computed(() => {
  if (!props.content) return ''
  try {
    if (selectedFormat.value === 'md') {
      return marked.parse(props.content) as string
    }
    return props.content
  } catch {
    return props.content
  }
})

const searchCount = computed(() => {
  if (!searchQuery.value || !props.content) return 0
  const regex = new RegExp(escapeRegExp(searchQuery.value), 'gi')
  return (props.content.match(regex) || []).length
})

function escapeRegExp(string: string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

function copyContent() {
  if (props.content) {
    navigator.clipboard.writeText(props.content)
    emit('copy')
  }
}

function downloadContent() {
  emit('download', selectedFormat.value)
}

function formatTime(dateString: string): string {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>
