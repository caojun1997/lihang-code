<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="opacity-0 translate-y-4"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-4"
  >
    <div
      v-if="show"
      :class="[
        'fixed bottom-6 right-6 z-50 flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg backdrop-blur-sm',
        typeConfig.bgColor,
        typeConfig.borderColor,
        'border'
      ]"
    >
      <component :is="typeConfig.icon" class="w-5 h-5" :class="typeConfig.iconColor" />
      <span class="text-sm font-medium text-foreground">{{ message }}</span>
      <button
        class="p-1 hover:bg-black/10 rounded transition-colors"
        @click="emit('close')"
      >
        <X class="w-4 h-4 text-muted-foreground" />
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CheckCircle2, XCircle, Info, X } from 'lucide-vue-next'

const props = defineProps<{
  message: string
  type?: 'success' | 'error' | 'info'
  show?: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const typeConfig = computed(() => {
  const configs = {
    success: {
      icon: CheckCircle2,
      bgColor: 'bg-green-50/95 dark:bg-green-900/90',
      borderColor: 'border-green-200 dark:border-green-700',
      iconColor: 'text-green-600 dark:text-green-400',
    },
    error: {
      icon: XCircle,
      bgColor: 'bg-red-50/95 dark:bg-red-900/90',
      borderColor: 'border-red-200 dark:border-red-700',
      iconColor: 'text-red-600 dark:text-red-400',
    },
    info: {
      icon: Info,
      bgColor: 'bg-blue-50/95 dark:bg-blue-900/90',
      borderColor: 'border-blue-200 dark:border-blue-700',
      iconColor: 'text-blue-600 dark:text-blue-400',
    },
  }
  return configs[props.type || 'success']
})
</script>
