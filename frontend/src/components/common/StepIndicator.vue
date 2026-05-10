<template>
  <div class="flex items-center justify-center gap-2">
    <div
      v-for="(step, index) in steps"
      :key="index"
      class="flex items-center"
    >
      <div class="flex items-center gap-3">
        <div
          :class="[
            'w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-all duration-300',
            getStepClass(index)
          ]"
        >
          <component
            :is="getStepIcon(index)"
            v-if="getStepIcon(index)"
            class="w-4 h-4"
          />
          <span v-else>{{ index + 1 }}</span>
        </div>
        <span
          :class="[
            'text-sm font-medium transition-colors hidden sm:block',
            index === currentStep
              ? 'text-foreground'
              : index < currentStep
              ? 'text-primary'
              : 'text-muted-foreground'
          ]"
        >
          {{ step }}
        </span>
      </div>

      <div
        v-if="index < steps.length - 1"
        :class="[
          'w-12 sm:w-20 h-0.5 mx-2 transition-colors duration-300',
          index < currentStep ? 'bg-primary' : 'bg-border'
        ]"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Check, Upload, Settings, Download } from 'lucide-vue-next'

interface Props {
  steps: string[]
  currentStep: number
}

const props = defineProps<Props>()

function getStepClass(index: number): string {
  if (index === props.currentStep) {
    return 'bg-primary text-primary-foreground ring-4 ring-primary/20'
  }
  if (index < props.currentStep) {
    return 'bg-primary text-primary-foreground'
  }
  return 'bg-muted text-muted-foreground'
}

function getStepIcon(index: number) {
  if (index < props.currentStep) {
    return Check
  }
  if (index === 0) return Upload
  if (index === 1) return Settings
  if (index === 2) return Download
  return null
}
</script>
