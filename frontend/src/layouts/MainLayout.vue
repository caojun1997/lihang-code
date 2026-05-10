<template>
  <div class="min-h-screen flex flex-col bg-background">
    <header class="sticky top-0 z-50 h-16 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div class="container h-full flex items-center justify-between px-4 md:px-6">
        <div class="flex items-center gap-6">
          <router-link to="/" class="flex items-center gap-3">
            <div class="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-lg shadow-primary/30">
              <FileText class="w-6 h-6 text-white" />
            </div>
            <div class="hidden sm:block">
              <h1 class="text-lg font-bold text-foreground">PDF解析助手</h1>
              <p class="text-xs text-muted-foreground">智能识别 · 高效转换</p>
            </div>
          </router-link>

          <nav class="hidden md:flex items-center gap-1">
            <router-link
              to="/"
              class="px-3 py-2 text-sm rounded-lg transition-colors"
              :class="$route.path === '/' ? 'bg-muted text-foreground font-medium' : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'"
            >
              上传解析
            </router-link>
            <router-link
              to="/history"
              class="px-3 py-2 text-sm rounded-lg transition-colors"
              :class="$route.path === '/history' ? 'bg-muted text-foreground font-medium' : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'"
            >
              历史记录
            </router-link>
          </nav>
        </div>

        <div class="flex items-center gap-2">
          <slot name="header-actions" />
          <button
            class="p-2.5 rounded-lg hover:bg-muted transition-colors"
            :aria-label="isDark ? '切换到亮色模式' : '切换到暗色模式'"
            @click="toggleTheme"
          >
            <Sun v-if="isDark" class="w-5 h-5 text-muted-foreground" />
            <Moon v-else class="w-5 h-5 text-muted-foreground" />
          </button>
          <div class="w-9 h-9 bg-primary/10 rounded-full flex items-center justify-center">
            <User class="w-5 h-5 text-primary" />
          </div>
        </div>
      </div>
    </header>

    <main class="flex-1">
      <slot />
    </main>

    <footer class="border-t border-border mt-auto">
      <div class="container px-4 md:px-6 py-4">
        <div class="flex flex-col sm:flex-row items-center justify-between gap-2 text-sm text-muted-foreground">
          <p>© 2026 PDF解析助手. All rights reserved.</p>
          <div class="flex items-center gap-4">
            <a href="#" class="hover:text-foreground transition-colors">帮助</a>
            <a href="#" class="hover:text-foreground transition-colors">隐私</a>
            <a href="#" class="hover:text-foreground transition-colors">使用条款</a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { FileText, Sun, Moon, User } from 'lucide-vue-next'

const isDark = ref(false)

onMounted(() => {
  checkDarkMode()
})

function checkDarkMode() {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

function toggleTheme() {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}
</script>
