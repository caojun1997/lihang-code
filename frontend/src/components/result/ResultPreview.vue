<template>
  <div class="h-full flex flex-col bg-card rounded-2xl border border-border overflow-hidden">
    <div class="flex items-center justify-between p-4 border-b border-border bg-muted/30">
      <div class="flex items-center gap-3">
        <FileText class="w-5 h-5 text-primary" />
        <div>
          <h3 class="font-semibold text-foreground">{{ fileName }}</h3>
          <p class="text-xs text-muted-foreground">{{ wordCount }} 字 · {{ format }}</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <div class="flex bg-muted rounded-lg p-1">
          <button
            v-for="fmt in formats"
            :key="fmt"
            :class="[
              'px-3 py-1.5 text-xs font-medium rounded-md transition-all',
              selectedFormat === fmt
                ? 'bg-primary text-primary-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            ]"
            @click="selectedFormat = fmt"
          >
            {{ fmt.toUpperCase() }}
          </button>
        </div>

        <div class="w-px h-6 bg-border mx-1" />

        <button
          class="p-2 text-muted-foreground hover:text-primary transition-colors"
          title="复制全部"
          @click="copyAll"
        >
          <Copy class="w-5 h-5" />
        </button>
        <button
          class="p-2 text-muted-foreground hover:text-primary transition-colors"
          title="下载"
          @click="downloadResult"
        >
          <Download class="w-5 h-5" />
        </button>
      </div>
    </div>

    <div class="flex-1 overflow-auto p-6">
      <div class="max-w-4xl mx-auto">
        <div
          v-if="content"
          class="prose prose-slate dark:prose-invert max-w-none"
          v-html="renderedContent"
        />
        <div v-else class="text-center text-muted-foreground py-12">
          <FileText class="w-12 h-12 mx-auto mb-4 opacity-50" />
          <p>暂无预览内容</p>
        </div>
      </div>
    </div>

    <div class="p-4 border-t border-border bg-muted/30">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4 text-sm text-muted-foreground">
          <button
            :class="[
              'flex items-center gap-1.5 hover:text-foreground transition-colors',
              { 'text-primary': searchQuery !== '' }
            ]"
            @click="showSearch = !showSearch"
          >
            <Search class="w-4 h-4" />
            搜索
          </button>
          <span>{{ wordCount }} 字</span>
          <span>{{ formatTime(createdAt || '') }}</span>
        </div>

        <div class="flex items-center gap-2">
          <button
            class="px-4 py-2 bg-primary text-primary-foreground rounded-lg hover:bg-primary/90 transition-colors flex items-center gap-2"
            @click="copyAll"
          >
            <Copy class="w-4 h-4" />
            复制全部
          </button>
        </div>
      </div>

      <div v-if="showSearch" class="mt-3 flex items-center gap-2">
        <div class="relative flex-1">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索内容..."
            class="w-full pl-10 pr-4 py-2 border border-border rounded-lg bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
          />
        </div>
        <span v-if="searchQuery" class="text-sm text-muted-foreground">
          找到 {{ searchResults }} 处
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { FileText, Copy, Download, Search } from 'lucide-vue-next'
import { marked } from 'marked'

const props = defineProps<{
  content: string
  fileName?: string
  format?: string
  createdAt?: string
}>()

const emit = defineEmits<{
  (e: 'copy'): void
  (e: 'download', format: string): void
}>()

const selectedFormat = ref('md')
const formats = ['md', 'txt', 'json']
const showSearch = ref(false)
const searchQuery = ref('')

const wordCount = computed(() => props.content?.length || 0)

const renderedContent = computed(() => {
  if (!props.content) return ''
  try {
    return marked.parse(props.content) as string
  } catch {
    return props.content
  }
})

const searchResults = computed(() => {
  if (!searchQuery.value) return 0
  const regex = new RegExp(searchQuery.value, 'gi')
  return (props.content.match(regex) || []).length
})

function copyAll() {
  navigator.clipboard.writeText(props.content)
  emit('copy')
}

function downloadResult() {
  emit('download', selectedFormat.value)
}

function formatTime(dateString: string): string {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}
</script>

<style>
.prose pre {
  @apply bg-muted rounded-lg p-4 overflow-x-auto;
}

.prose code {
  @apply text-sm;
}

.prose pre code {
  @apply bg-transparent;
}

.prose h1 {
  @apply text-2xl font-bold mb-4 mt-6;
}

.prose h2 {
  @apply text-xl font-semibold mb-3 mt-5;
}

.prose h3 {
  @apply text-lg font-medium mb-2 mt-4;
}

.prose p {
  @apply mb-4 leading-relaxed;
}

.prose ul, .prose ol {
  @apply mb-4 pl-6;
}

.prose li {
  @apply mb-1;
}

.prose table {
  @apply w-full border-collapse mb-4;
}

.prose th, .prose td {
  @apply border border-border px-4 py-2;
}

.prose th {
  @apply bg-muted font-semibold;
}
</style>
