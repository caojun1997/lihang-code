<template>
  <div class="space-y-2">
    <label class="text-sm font-medium text-gray-700">输出格式</label>
    <div class="grid grid-cols-2 gap-3">
      <button
        v-for="option in options"
        :key="option.value"
        type="button"
        :class="[
          'relative rounded-lg border p-4 transition-all hover:border-blue-300',
          modelValue === option.value
            ? 'border-blue-500 bg-blue-50 ring-2 ring-blue-200'
            : 'border-gray-200 bg-white',
        ]"
        @click="emit('update:modelValue', option.value)"
      >
        <component
          :is="option.icon"
          :class="[
            'w-5 h-5 mb-2',
            modelValue === option.value ? 'text-blue-500' : 'text-gray-400',
          ]"
        />
        <p
          :class="[
            'text-sm font-medium',
            modelValue === option.value ? 'text-blue-700' : 'text-gray-700',
          ]"
        >
          {{ option.label }}
        </p>
        <p class="text-xs text-gray-500 mt-1">{{ option.description }}</p>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { FileCode, FileText } from 'lucide-vue-next'

interface FormatOption {
  value: 'markdown' | 'txt'
  label: string
  description: string
  icon: any
}

defineProps<{
  modelValue: 'markdown' | 'txt'
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: 'markdown' | 'txt'): void
}>()

const options: FormatOption[] = [
  {
    value: 'markdown',
    label: 'Markdown',
    description: '保留文档结构和格式',
    icon: FileCode,
  },
  {
    value: 'txt',
    label: '纯文本',
    description: '去除格式，纯文本内容',
    icon: FileText,
  },
]
</script>
