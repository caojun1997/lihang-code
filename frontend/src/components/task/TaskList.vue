<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <div class="w-1 h-6 bg-primary rounded-full" />
        <h2 class="text-lg font-semibold text-foreground">任务列表</h2>
      </div>
      <span class="text-sm text-muted-foreground">{{ tasks.length }} 个任务</span>
    </div>

    <div class="flex flex-wrap gap-2">
      <button
        v-for="status in statusOptions"
        :key="status.value"
        :class="[
          'px-3 py-1.5 text-sm rounded-lg transition-colors',
          selectedStatus === status.value
            ? 'bg-primary text-primary-foreground'
            : 'bg-muted text-muted-foreground hover:bg-muted/80'
        ]"
        @click="selectedStatus = status.value"
      >
        {{ status.label }}
        <span v-if="status.count > 0" class="ml-1">({{ status.count }})</span>
      </button>
    </div>

    <div class="space-y-3">
      <TransitionGroup
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="opacity-0 -translate-x-4"
        enter-to-class="opacity-100 translate-x-0"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="opacity-100 translate-x-0"
        leave-to-class="opacity-0 translate-x-4"
      >
        <TaskCard
          v-for="task in filteredTasks"
          :key="task.task_id"
          :task="task"
          :is-active="currentTask?.task_id === task.task_id"
          @click="emit('select', task)"
          @view="emit('view', task)"
          @download="emit('download', task)"
          @copy="emit('copy', task)"
          @delete="emit('delete', task)"
        />
      </TransitionGroup>

      <div v-if="filteredTasks.length === 0" class="text-center py-12">
        <div class="w-16 h-16 bg-muted rounded-2xl flex items-center justify-center mx-auto mb-4">
          <FileText class="w-8 h-8 text-muted-foreground" />
        </div>
        <p class="text-sm text-muted-foreground">暂无{{ selectedStatus === 'all' ? '' : selectedStatus }}任务</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { FileText } from 'lucide-vue-next'
import TaskCard from './TaskCard.vue'
import type { Task, TaskStatus } from '@/api/types'

const props = defineProps<{
  tasks: Task[]
  currentTask?: Task | null
}>()

const emit = defineEmits<{
  (e: 'select', task: Task): void
  (e: 'view', task: Task): void
  (e: 'download', task: Task): void
  (e: 'copy', task: Task): void
  (e: 'delete', task: Task): void
}>()

const selectedStatus = ref<TaskStatus | 'all'>('all')

const statusOptions = computed(() => {
  const statuses: { value: TaskStatus | 'all'; label: string; count: number }[] = [
    { value: 'all', label: '全部', count: props.tasks.length },
    { value: 'processing', label: '处理中', count: props.tasks.filter(t => t.status === 'processing').length },
    { value: 'completed', label: '已完成', count: props.tasks.filter(t => t.status === 'completed').length },
    { value: 'failed', label: '失败', count: props.tasks.filter(t => t.status === 'failed').length },
  ]
  return statuses
})

const filteredTasks = computed(() => {
  if (selectedStatus.value === 'all') {
    return props.tasks
  }
  return props.tasks.filter(t => t.status === selectedStatus.value)
})
</script>
