<template>
  <div class="space-y-4">
    <div class="space-y-2">
      <label class="text-sm font-medium text-foreground">输出格式</label>
      <div class="grid grid-cols-3 gap-3">
        <button
          v-for="option in formatOptions"
          :key="option.value"
          :class="[
            'relative p-4 rounded-xl border-2 transition-all duration-200 text-center',
            'hover:border-primary/50 hover:bg-muted/30',
            modelValue === option.value
              ? 'border-primary bg-primary/5 shadow-sm'
              : 'border-border bg-card'
          ]"
          @click="emit('update:modelValue', option.value)"
        >
          <component :is="option.icon" class="w-6 h-6 mx-auto mb-2" />
          <p class="text-sm font-medium text-foreground">{{ option.label }}</p>
          <p class="text-xs text-muted-foreground mt-0.5">{{ option.description }}</p>
          <div
            v-if="modelValue === option.value"
            class="absolute top-2 right-2 w-5 h-5 bg-primary rounded-full flex items-center justify-center"
          >
            <Check class="w-3 h-3 text-white" />
          </div>
        </button>
      </div>
    </div>

    <div class="space-y-3">
      <label class="text-sm font-medium text-foreground">解析选项</label>
      <div class="grid grid-cols-2 gap-3">
        <label
          v-for="option in optionItems"
          :key="option.key"
          class="flex items-center gap-3 p-3 rounded-lg border border-border bg-card cursor-pointer hover:bg-muted/30 transition-colors"
        >
          <input
            type="checkbox"
            :checked="options[option.key]"
            class="w-4 h-4 rounded border-input text-primary focus:ring-primary"
            @change="toggleOption(option.key)"
          />
          <div class="flex-1">
            <p class="text-sm font-medium text-foreground">{{ option.label }}</p>
            <p class="text-xs text-muted-foreground">{{ option.description }}</p>
          </div>
        </label>
      </div>
    </div>

    <div class="space-y-2">
      <label class="text-sm font-medium text-foreground">语言设置</label>
      <select
        v-model="options.language"
        class="w-full h-10 px-3 rounded-lg border border-border bg-background text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
      >
        <option value="auto">自动检测</option>
        <option value="zh">中文</option>
        <option value="en">英文</option>
        <option value="ja">日文</option>
        <option value="ko">韩文</option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import { FileText, Code, Check } from 'lucide-vue-next'

const props = defineProps<{
  modelValue: 'markdown' | 'txt' | 'json'
  options: {
    preserve_tables: boolean
    extract_images: boolean
    latex_formulas: boolean
    preserve_headings: boolean
    language: string
  }
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: 'markdown' | 'txt' | 'json'): void
  (e: 'update:options', value: typeof props.options): void
}>()

const formatOptions = [
  {
    value: 'markdown' as const,
    label: 'Markdown',
    description: '保留格式结构',
    icon: FileText
  },
  {
    value: 'txt' as const,
    label: '纯文本',
    description: '仅提取文本',
    icon: Code
  },
  {
    value: 'json' as const,
    label: 'JSON',
    description: '含元数据',
    icon: FileText
  }
]

const optionItems = [
  {
    key: 'preserve_tables' as const,
    label: '保留表格结构',
    description: '表格转为Markdown格式'
  },
  {
    key: 'extract_images' as const,
    label: '提取图片',
    description: '保留图片引用'
  },
  {
    key: 'latex_formulas' as const,
    label: '识别公式',
    description: '转为LaTeX格式'
  },
  {
    key: 'preserve_headings' as const,
    label: '保留目录层级',
    description: '识别标题结构'
  }
]

function toggleOption(key: keyof typeof props.options) {
  if (key === 'language') return
  const newOptions = { ...props.options, [key]: !props.options[key] }
  emit('update:options', newOptions)
}
</script>
