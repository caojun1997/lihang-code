<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 translate-y-4"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-4"
    >
      <div
        v-if="visible"
        class="fixed bottom-6 right-6 z-50 max-w-sm"
      >
        <div
          :class="[
            'flex items-start gap-3 p-4 rounded-xl shadow-lg border backdrop-blur-sm',
            toastClasses
          ]"
        >
          <div class="flex-shrink-0 mt-0.5">
            <component :is="iconComponent" class="w-5 h-5" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium">{{ message }}</p>
          </div>
          <button
            class="flex-shrink-0 p-1 rounded hover:bg-black/10 transition-colors"
            @click="emit('close')"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CheckCircle, XCircle, AlertCircle, Info, X } from 'lucide-vue-next'

interface Props {
  visible: boolean
  message: string
  type?: 'success' | 'error' | 'info' | 'warning'
}

const props = withDefaults(defineProps<Props>(), {
  type: 'info'
})

const emit = defineEmits<{
  (e: 'close'): void
}>()

const toastClasses = computed(() => {
  const base = 'bg-background/95'
  switch (props.type) {
    case 'success':
      return `${base} border-green-500/20 text-green-600`
    case 'error':
      return `${base} border-red-500/20 text-red-600`
    case 'warning':
      return `${base} border-amber-500/20 text-amber-600`
    default:
      return `${base} border-blue-500/20 text-blue-600`
  }
})

const iconComponent = computed(() => {
  switch (props.type) {
    case 'success':
      return CheckCircle
    case 'error':
      return XCircle
    case 'warning':
      return AlertCircle
    default:
      return Info
  }
})
</script>
