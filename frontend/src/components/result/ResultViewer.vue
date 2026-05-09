<template>
  <div class="bg-white rounded-xl border shadow-sm">
    <div class="p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">解析结果</h3>

      <div v-if="loading" class="flex items-center justify-center py-12">
        <Loader2 class="w-8 h-8 animate-spin text-blue-500" />
      </div>

      <div v-else-if="content" class="space-y-4">
        <div class="flex items-center justify-between text-sm text-gray-500">
          <span>字数统计: {{ wordCount }}</span>
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="flex items-center gap-1 text-blue-600 hover:text-blue-700"
              @click="copyContent"
            >
              <Copy class="w-4 h-4" />
              {{ copied ? '已复制' : '复制' }}
            </button>
          </div>
        </div>

        <div class="bg-gray-50 rounded-lg p-4 max-h-96 overflow-y-auto">
          <pre class="text-sm text-gray-700 whitespace-pre-wrap font-mono">{{ content }}</pre>
        </div>
      </div>

      <div v-else class="text-center py-12">
        <p class="text-gray-500">暂无解析结果</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Loader2, Copy } from 'lucide-vue-next'
import { taskApi } from '@/api'

const props = defineProps<{
  taskId: string
}>()

const content = ref('')
const wordCount = ref(0)
const loading = ref(true)
const copied = ref(false)

async function fetchResult() {
  loading.value = true
  try {
    const response = await taskApi.getTaskResult(props.taskId)
    content.value = response.data?.content || ''
    wordCount.value = response.data?.word_count || 0
  } catch (e) {
    console.error('Failed to fetch result:', e)
  } finally {
    loading.value = false
  }
}

async function copyContent() {
  try {
    await navigator.clipboard.writeText(content.value)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (e) {
    console.error('Failed to copy:', e)
  }
}

onMounted(() => {
  fetchResult()
})
</script>
