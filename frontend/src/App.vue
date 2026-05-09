<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <header class="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
      <div class="container mx-auto px-4 py-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center">
            <FileText class="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 class="text-xl font-bold text-gray-900">PDF Parser</h1>
            <p class="text-xs text-gray-500">智能文档解析平台</p>
          </div>
        </div>
        <nav class="flex items-center gap-4">
          <router-link
            to="/"
            class="text-sm font-medium text-gray-600 hover:text-blue-600 transition-colors"
          >
            首页
          </router-link>
          <router-link
            to="/tasks"
            class="text-sm font-medium text-gray-600 hover:text-blue-600 transition-colors"
          >
            任务列表
          </router-link>
        </nav>
      </div>
    </header>

    <main>
      <router-view />
    </main>

    <footer class="border-t bg-white/50 mt-auto">
      <div class="container mx-auto px-4 py-6 text-center text-sm text-gray-500">
        PDF Parser - 快速将PDF文档转换为Markdown或纯文本格式
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { FileText } from 'lucide-vue-next'
</script>
