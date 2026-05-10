<template>
  <div class="w-full">
    <div
      v-if="label || showPercentage"
      class="flex items-center justify-between mb-2 text-sm"
    >
      <span v-if="label" class="text-muted-foreground">{{ label }}</span>
      <span v-if="showPercentage" class="font-medium text-primary">{{ percentage }}%</span>
    </div>

    <div
      class="h-2 bg-muted rounded-full overflow-hidden"
      :class="sizeClass"
    >
      <div
        :class="[
          'h-full rounded-full transition-all duration-300 ease-out',
          progressClass
        ]"
        :style="{ width: `${percentage}%` }"
      />
    </div>

    <p v-if="message" class="mt-2 text-xs text-muted-foreground">
      {{ message }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  value: number
  max?: number
  size?: 'sm' | 'md' | 'lg'
  variant?: 'primary' | 'success' | 'warning' | 'error'
  label?: string
  message?: string
  showPercentage?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  max: 100,
  size: 'md',
  variant: 'primary',
  showPercentage: false
})

const percentage = computed(() => {
  const pct = Math.min(Math.max((props.value / props.max) * 100, 0), 100)
  return Math.round(pct)
})

const sizeClass = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'h-1'
    case 'lg':
      return 'h-3'
    default:
      return 'h-2'
  }
})

const progressClass = computed(() => {
  switch (props.variant) {
    case 'success':
      return 'bg-emerald-500'
    case 'warning':
      return 'bg-amber-500'
    case 'error':
      return 'bg-red-500'
    default:
      return 'bg-primary'
  }
})
</script>
