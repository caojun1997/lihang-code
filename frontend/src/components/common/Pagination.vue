<template>
  <div class="flex items-center justify-between">
    <div class="text-sm text-muted-foreground">
      共 {{ total }} 条记录
    </div>
    
    <div class="flex items-center gap-1">
      <button
        :disabled="currentPage === 1"
        class="p-2 rounded-lg hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        @click="emit('change', currentPage - 1)"
      >
        <ChevronLeft class="w-4 h-4" />
      </button>

      <template v-for="page in visiblePages" :key="page">
        <button
          v-if="page === '...'"
          class="px-3 py-1.5 text-muted-foreground"
          disabled
        >
          ...
        </button>
        <button
          v-else
          :class="[
            'px-3 py-1.5 rounded-lg text-sm font-medium transition-colors',
            page === currentPage
              ? 'bg-primary text-primary-foreground'
              : 'hover:bg-muted text-muted-foreground hover:text-foreground'
          ]"
          @click="emit('change', page as number)"
        >
          {{ page }}
        </button>
      </template>

      <button
        :disabled="currentPage === totalPages"
        class="p-2 rounded-lg hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        @click="emit('change', currentPage + 1)"
      >
        <ChevronRight class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

interface Props {
  currentPage: number
  totalPages: number
  total: number
  pageSize: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'change', page: number): void
}>()

const visiblePages = computed(() => {
  const pages: (number | string)[] = []
  const total = props.totalPages
  const current = props.currentPage

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    pages.push(1)

    if (current > 3) {
      pages.push('...')
    }

    const start = Math.max(2, current - 1)
    const end = Math.min(total - 1, current + 1)

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    if (current < total - 2) {
      pages.push('...')
    }

    pages.push(total)
  }

  return pages
})
</script>
