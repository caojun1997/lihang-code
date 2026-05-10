<template>
  <div class="space-y-3">
    <label class="text-sm font-medium text-foreground">输出格式</label>
    <div class="grid grid-cols-3 gap-3">
      <button
        v-for="option in formatOptions"
        :key="option.value"
        type="button"
        :class="[
          'relative flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all duration-200',
          modelValue === option.value
            ? 'border-primary bg-primary/5 shadow-sm'
            : 'border-border hover:border-primary/30 hover:bg-muted/50',
        ]"
        @click="emit('update:modelValue', option.value)"
      >
        <component :is="option.icon" class="w-6 h-6" :class="modelValue === option.value ? 'text-primary' : 'text-muted-foreground'" />
        <span class="text-sm font-medium" :class="modelValue === option.value ? 'text-primary' : 'text-foreground'">
          {{ option.label }}
        </span>
        <span class="text-xs text-muted-foreground">{{ option.description }}</span>

        <div
          v-if="modelValue === option.value"
          class="absolute -top-2 -right-2 w-5 h-5 bg-primary rounded-full flex items-center justify-center"
        >
          <Check class="w-3 h-3 text-primary-foreground" />
        </div>
      </button>
    </div>

    <div class="flex flex-wrap gap-2 mt-4">
      <span class="text-xs text-muted-foreground">高级选项：</span>
      <label
        v-for="opt in advancedOptions"
        :key="opt.key"
        class="flex items-center gap-1.5 text-xs cursor-pointer"
      >
        <input
          type="checkbox"
          v-model="advancedState[opt.key]"
          class="rounded border-input text-primary focus:ring-primary"
        />
        <span class="text-muted-foreground hover:text-foreground transition-colors">{{ opt.label }}</span>
      </label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { FileText, Database, Check } from 'lucide-vue-next'

const props = defineProps<{
  modelValue: 'markdown' | 'txt' | 'json'
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: 'markdown' | 'txt' | 'json'): void
}>()

const formatOptions = [
  {
    value: 'markdown' as const,
    label: 'Markdown',
    description: '结构化文本',
    icon: FileText,
  },
  {
    value: 'txt' as const,
    label: '纯文本',
    description: '无格式',
    icon: FileText,
  },
  {
    value: 'json' as const,
    label: 'JSON',
    description: '结构化数据',
    icon: Database,
  },
]

const advancedOptions = [
  { key: 'tables' as const, label: '保留表格结构' },
  { key: 'formulas' as const, label: '保留公式' },
  { key: 'images' as const, label: '提取图片' },
]

const advancedState = reactive({
  tables: true,
  formulas: true,
  images: true,
})

void props
</script>
