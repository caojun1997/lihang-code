# PDF解析系统 — 全新前端交互设计方案

## 一、设计理念与愿景

本方案基于对市场上主流PDF解析工具的深度调研，融合Smallpdf的简洁体验、iLovePDF的功能完整性、UPDF的AI智能特性，以及NotebookLM的对话式交互优势，打造一款**「极简操作 · 智能体验 · 专业输出」**的下一代PDF解析工具。

**核心设计原则：**
- **三步完成**：拖拽上传 → 智能解析 → 多元输出
- **所见即所得**：实时预览解析过程，即时反馈处理状态
- **专业而不复杂**：面向普通用户，隐藏技术细节
- **AI赋能**：智能识别文档结构，自动标注重点内容

---

## 二、设计风格与视觉语言

### 2.1 整体视觉风格

**风格定位**：极简现代、专业高效、友好亲和

**设计参考**：
- Smallpdf的干净利落
- Notion的卡片式布局
- Linear的精致动效
- Vercel的深邃配色

### 2.2 色彩系统

```
主色调 (Primary)     : #6366F1  (Indigo 500) - 代表智能与专业
辅助色 (Secondary)   : #8B5CF6  (Violet 500) - 代表创新与活力
成功色 (Success)     : #10B981  (Emerald 500) - 解析成功
警告色 (Warning)     : #F59E0B  (Amber 500)   - 部分成功
错误色 (Error)       : #EF4444  (Red 500)     - 解析失败
信息色 (Info)        : #3B82F6  (Blue 500)    - 处理中

背景色 (Background) :
  - 主背景        : #FAFAFA (浅色) / #0A0A0B (深色)
  - 卡片背景      : #FFFFFF (浅色) / #18181B (深色)
  - 次级背景      : #F4F4F5 (浅色) / #27272A (深色)

文字色 (Text) :
  - 主文字        : #18181B (浅色) / #FAFAFA (深色)
  - 次要文字      : #71717A (浅色) / #A1A1AA (深色)
  - 占位文字      : #A1A1AA (浅色) / #71717A (深色)

边框色 (Border)      : #E4E4E7 (浅色) / #27272A (深色)
```

### 2.3 字体系统

```
主字体    : Inter (英文) / "PingFang SC", "Microsoft YaHei" (中文)
代码字体  : "JetBrains Mono", "Fira Code", monospace

字号层级  :
  - 标题 H1    : 32px / 2rem     , font-weight: 700
  - 标题 H2    : 24px / 1.5rem   , font-weight: 600
  - 标题 H3    : 18px / 1.125rem , font-weight: 600
  - 正文       : 14px / 0.875rem , font-weight: 400
  - 辅助文字   : 12px / 0.75rem  , font-weight: 400
  - 标签/徽章  : 11px / 0.6875rem, font-weight: 500
```

### 2.4 间距与圆角

```
间距系统 (基于 4px 网格) :
  - xs  : 4px   (1)
  - sm  : 8px   (2)
  - md  : 16px  (4)
  - lg  : 24px  (6)
  - xl  : 32px  (8)
  - 2xl : 48px  (12)
  - 3xl : 64px  (16)

圆角系统 :
  - 小元素 (按钮、输入框) : 8px
  - 卡片、面板           : 12px
  - 大容器、模态框       : 16px
  - 头像、图标容器       : 9999px (全圆角)
```

### 2.5 阴影系统

```
阴影层级  :
  - 阴影-sm  : 0 1px 2px rgba(0,0,0,0.05)
  - 阴影-md  : 0 4px 6px -1px rgba(0,0,0,0.1)
  - 阴影-lg  : 0 10px 15px -3px rgba(0,0,0,0.1)
  - 阴影-xl  : 0 20px 25px -5px rgba(0,0,0,0.1)
  - 阴影-glow: 0 0 20px rgba(99,102,241,0.3) (主色调光晕)
```

---

## 三、页面布局与结构

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│  顶部导航栏 (64px 固定高度)                                    │
│  ┌─────────────────────────────────────────────────────────┐│
│  │ Logo + 产品名    搜索框        主题切换 │ 用户头像/登录  ││
│  └─────────────────────────────────────────────────────────┘│
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  主内容区域                                                  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                                                       │  │
│  │  左侧主操作区 (flex-1)                                │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │                                                 │ │  │
│  │  │     [步骤指示器 - 3步流程]                       │ │  │
│  │  │                                                 │ │  │
│  │  │     ┌─────────────────────────────────────┐    │ │  │
│  │  │     │                                     │    │ │  │
│  │  │     │     上传区域 / 步骤内容              │    │ │  │
│  │  │     │                                     │    │ │  │
│  │  │     └─────────────────────────────────────┘    │ │  │
│  │  │                                                 │ │  │
│  │  │     [底部操作按钮]                             │ │  │
│  │  │                                                 │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  右侧预览面板 (400px 可折叠)                          │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │  实时预览 / 结果展示                            │ │  │
│  │  │                                                 │ │  │
│  │  │                                                 │ │  │
│  │  │                                                 │ │  │
│  │  │                                                 │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  底部状态栏 (可选，根据页面状态显示)                           │
│  ┌─────────────────────────────────────────────────────────┐│
│  │  处理进度 │ 统计信息 │ 快捷操作                           ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

### 3.2 响应式断点

```
桌面端   : >= 1280px  - 完整三栏布局
笔记本端 : 1024-1279px - 两侧布局，右侧面板可折叠
平板端   : 768-1023px  - 单栏布局，Tab切换
移动端   : < 768px     - 卡片式布局，底部浮动操作栏
```

---

## 四、核心交互流程

### 4.1 三步完成解析流程

**第一步：上传文件**

```
┌─────────────────────────────────────────────────────────────┐
│  ① 上传                                        ② 配置  ③ 下载 │
│  ● ───── ○ ───── ○                                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                                                         │ │
│  │              ┌───────────────────────────┐              │ │
│  │              │                           │              │ │
│  │              │    📄                      │              │ │
│  │              │    拖拽文件到此处          │              │ │
│  │              │    或点击选择文件          │              │ │
│  │              │                           │              │ │
│  │              │    支持 PDF、DOC、PPT      │              │ │
│  │              │    最大 100MB              │              │ │
│  │              │                           │              │ │
│  │              └───────────────────────────┘              │ │
│  │                                                         │ │
│  │         [  从电脑选择 ]    [ 使用网址 ]                  │ │
│  │                                                         │ │
│  │  ───────────────────────────────────────────────────── │ │
│  │                                                         │ │
│  │  已选择文件：                                            │ │
│  │  ┌─────────────────────────────────────────────────┐   │ │
│  │  │ 📄 技术文档_V2.0.pdf          12.3 MB      ✕    │   │ │
│  │  │ 📄 财务报表_Q1.pptx           8.7 MB       ✕    │   │ │
│  │  └─────────────────────────────────────────────────┘   │ │
│  │                                                         │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
│                              [ 下一步：配置解析选项 → ]       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**第二步：配置选项**

```
┌─────────────────────────────────────────────────────────────┐
│  ① 上传                                        ② 配置  ③ 下载 │
│      ○ ───── ● ───── ○                                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                                                         │ │
│  │  输出格式                                               │ │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐                     │ │
│  │  │ Markdown│ │ 纯文本  │ │   JSON  │                     │ │
│  │  │   ✓    │ │        │ │        │                     │ │
│  │  └─────────┘ └─────────┘ └─────────┘                     │ │
│  │                                                         │ │
│  │  ───────────────────────────────────────────────────── │ │
│  │                                                         │ │
│  │  解析选项                                               │ │
│  │  ┌─────────────────────────────────────────────────┐   │ │
│  │  │ ☑ 保留表格结构      ☑ 提取图片                   │   │ │
│  │  │ ☑ 识别公式(Latex)  ☑ 保留目录层级              │   │ │
│  │  │ ☐ OCR识别(扫描件)  ☐ 双栏布局优化              │   │ │
│  │  └─────────────────────────────────────────────────┘   │ │
│  │                                                         │ │
│  │  语言设置                                               │ │
│  │  [ 自动检测 ▼ ]                                         │ │
│  │                                                         │ │
│  │  ───────────────────────────────────────────────────── │ │
│  │                                                         │ │
│  │  高级设置                                               │ │
│  │  [ 展开 ▼ ]                                             │ │
│  │                                                         │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
│  [ ← 上一步 ]                    [ 开始解析 → ]              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**第三步：解析与下载**

```
┌─────────────────────────────────────────────────────────────┐
│  ① 上传                                        ② 配置  ③ 下载 │
│      ○ ───── ○ ───── ●                                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌───────────────────────────────────┬─────────────────────┐│
│  │                                   │                     ││
│  │  正在解析...                       │   预览面板           ││
│  │                                   │                     ││
│  │  ████████████░░░░░░░  68%        │   [实时渲染结果]     ││
│  │                                   │                     ││
│  │  正在处理：第 15 / 22 页          │                     ││
│  │  当前：识别表格结构...             │                     ││
│  │                                   │                     ││
│  │  ┌─────────────────────────────┐ │                     ││
│  │  │ ⏳ 等待中 - 财务报表_Q1.pptx │ │                     ││
│  │  └─────────────────────────────┘ │                     ││
│  │                                   │                     ││
│  │  ┌─────────────────────────────┐ │                     ││
│  │  │ 🔄 处理中 - 技术文档_V2...  │ │                     ││
│  │  └─────────────────────────────┘ │                     ││
│  │                                   │                     ││
│  │  ┌─────────────────────────────┐ │                     ││
│  │  │ ✅ 完成 - 论文摘要.pdf      │ │                     ││
│  │  │    [查看] [下载] [复制]     │ │                     ││
│  │  └─────────────────────────────┘ │                     ││
│  │                                   │                     ││
│  └───────────────────────────────────┴─────────────────────┘│
│                                                             │
│  [ ← 重新解析 ]     [📥 下载全部]  [📋 复制全部]              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 交互状态与反馈

**拖拽交互状态**

| 状态 | 视觉表现 | 动效 |
|------|---------|------|
| 默认 | 虚线边框，灰色背景 | 无 |
| 悬停 | 边框加深，背景微变 | scale(1.01), 200ms |
| 拖入 | 主色边框，蓝色光晕背景 | scale(1.02), pulse动画 |
| 释放 | 边框闪烁，成功图标 | bounce, 300ms |

**上传进度状态**

```
状态流程：
  [选择文件] → [上传中 0-100%] → [上传完成 ✓] → [开始解析]
                                              ↓
                                        [解析中 0-100%]
                                              ↓
                                    [部分成功 ⚠] / [完成 ✓] / [失败 ✗]
```

---

## 五、组件设计规范

### 5.1 上传组件 (FileUploader)

**功能规格**：
- 拖拽上传：支持拖拽多个文件
- 点击选择：支持文件选择器
- URL导入：支持粘贴URL下载
- 粘贴上传：支持Ctrl+V粘贴
- 文件预览：显示缩略图和元数据
- 进度显示：实时上传/解析进度
- 错误提示：格式错误、文件过大等

**视觉规格**：
- 容器：圆角12px，虚线边框
- 拖拽区：最小高度240px
- 图标：64px容器，居中
- 文件列表：最大高度320px，滚动

### 5.2 格式选择器 (FormatSelector)

**功能规格**：
- 单选格式：Markdown / TXT / JSON
- 可视化切换：图标+文字
- 悬停预览：显示格式示例

**视觉规格**：
- 选项卡片：120px × 80px
- 选中态：主色边框+背景色
- 图标：24px，文字14px

### 5.3 任务卡片 (TaskCard)

**功能规格**：
- 显示：文件名、状态、大小、时间
- 操作：查看、下载、复制、删除
- 状态：7种状态的不同样式
- 进度：处理中显示进度条

**视觉规格**：
- 卡片：圆角12px，阴影-sm
- 间距：内部16px，间距12px
- 状态徽章：圆角6px，12px字号

### 5.4 结果预览 (ResultPreview)

**功能规格**：
- 实时预览：边解析边显示
- 格式切换：MD/TXT/JSON
- 搜索高亮：关键词搜索
- 复制下载：快捷操作
- 目录导航：Markdown目录

**视觉规格**：
- 工具栏：固定顶部，高度56px
- 内容区：prose样式，最大宽度720px
- 底部栏：固定底部，高度64px

---

## 六、状态设计

### 6.1 任务状态系统

| 状态值 | 显示文案 | 颜色 | 图标 | 说明 |
|--------|---------|------|------|------|
| `pending` | 等待中 | gray | Clock | 队列等待 |
| `uploading` | 上传中 | blue | Upload | 文件上传 |
| `processing` | 解析中 | blue | Loader | 内容解析 |
| `completed` | 已完成 | green | CheckCircle | 成功完成 |
| `partial_failed` | 部分失败 | amber | AlertTriangle | 部分成功 |
| `failed` | 解析失败 | red | XCircle | 完全失败 |
| `cancelled` | 已取消 | gray | Ban | 用户取消 |

### 6.2 状态流转图

```
                    ┌──────────┐
                    │ pending  │
                    └────┬─────┘
                         │
                    文件上传
                         │
                    ┌────▼─────┐
              ┌─────│ uploading│─────┐
              │     └──────────┘     │
           取消                        成功
              │                        │
              ▼                        ▼
        ┌──────────┐           ┌──────────────┐
        │ cancelled│           │  processing  │
        └──────────┘           └───────┬──────┘
                                      │
                           ┌──────────┼──────────┐
                           │          │          │
                        成功        部分失败     失败
                           │          │          │
                           ▼          ▼          ▼
                     ┌──────────┐ ┌─────────┐ ┌────────┐
                     │ completed│ │partial_ │ │ failed │
                     └──────────┘ │ failed  │ └────────┘
                                  └─────────┘
```

### 6.3 错误处理规范

| 错误类型 | HTTP状态码 | 错误码 | 用户提示 | 解决方案 |
|---------|-----------|--------|---------|---------|
| 文件过大 | 413 | `FILE_TOO_LARGE` | 文件超过100MB限制 | 提示分割文件 |
| 格式不支持 | 400 | `UNSUPPORTED_FORMAT` | 不支持此文件格式 | 显示支持格式 |
| 解析失败 | 500 | `PARSE_FAILED` | 解析失败，请重试 | 显示重试按钮 |
| 认证失败 | 401 | `UNAUTHORIZED` | 请先登录 | 跳转登录页 |
| 限流 | 429 | `RATE_LIMITED` | 操作过于频繁 | 显示等待时间 |
| 服务器错误 | 500 | `INTERNAL_ERROR` | 服务器异常 | 显示联系支持 |

---

## 七、动画与过渡

### 7.1 页面过渡

```css
/* 页面进入 */
.page-enter-active {
  animation: fadeSlideIn 0.3s ease-out;
}
@keyframes fadeSlideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 页面离开 */
.page-leave-active {
  animation: fadeSlideOut 0.2s ease-in;
}
@keyframes fadeSlideOut {
  from {
    opacity: 1;
    transform: translateY(0);
  }
  to {
    opacity: 0;
    transform: translateY(-10px);
  }
}
```

### 7.2 交互反馈动画

| 交互场景 | 动画效果 | 时长 | 缓动函数 |
|---------|---------|------|---------|
| 按钮点击 | scale(0.98) → scale(1) | 150ms | ease-out |
| 拖拽进入 | border-color变化 + 背景渐变 | 200ms | ease |
| 文件添加 | slideIn + fadeIn | 300ms | ease-out |
| 删除文件 | slideOut + fadeOut | 200ms | ease-in |
| 成功反馈 | 绿色闪烁 + 轻微弹跳 | 400ms | spring |
| 复制成功 | 按钮变色 + 提示文字淡入 | 200ms | ease |

### 7.3 加载状态动画

```css
/* 骨架屏脉冲 */
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

.skeleton {
  background: linear-gradient(
    90deg,
    var(--muted) 25%,
    var(--muted-foreground/10) 50%,
    var(--muted) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

/* 进度条动画 */
.progress-bar {
  transition: width 0.3s ease-out;
}

/* 旋转加载 */
.spin {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
```

---

## 八、API接口设计

### 8.1 任务相关接口

| 接口路径 | 方法 | 描述 | 请求体 | 响应 |
|---------|------|------|-------|------|
| `/api/v1/tasks` | POST | 创建解析任务 | `{file, format, options}` | `{task_id, status}` |
| `/api/v1/tasks` | GET | 获取任务列表 | `?status=&page=&size=` | `{tasks[], total}` |
| `/api/v1/tasks/:id` | GET | 获取任务详情 | - | `{task}` |
| `/api/v1/tasks/:id/status` | GET | 获取任务状态 | - | `{status, progress}` |
| `/api/v1/tasks/:id/result` | GET | 获取解析结果 | - | `{content, metadata}` |
| `/api/v1/tasks/:id/download` | GET | 下载结果文件 | `?format=md` | Blob |
| `/api/v1/tasks/:id` | DELETE | 删除任务 | - | `{success}` |

### 8.2 请求/响应示例

**创建任务请求**：
```json
{
  "file_name": "技术文档.pdf",
  "file_size": 12582912,
  "format": "markdown",
  "options": {
    "preserve_tables": true,
    "extract_images": true,
    "latex_formulas": true,
    "preserve_headings": true,
    "language": "auto"
  }
}
```

**任务状态响应**：
```json
{
  "task_id": "tsk_abc123",
  "file_name": "技术文档.pdf",
  "status": "processing",
  "progress": 68,
  "current_page": 15,
  "total_pages": 22,
  "message": "正在识别表格结构...",
  "created_at": "2026-05-10T10:30:00Z",
  "updated_at": "2026-05-10T10:32:45Z"
}
```

**解析结果响应**：
```json
{
  "task_id": "tsk_abc123",
  "content": "# 技术文档\n\n## 第一章 概述\n\n...",
  "format": "markdown",
  "metadata": {
    "title": "技术文档",
    "author": "张三",
    "pages": 22,
    "words": 15234,
    "tables": 5,
    "images": 12,
    "formulas": 8
  },
  "created_at": "2026-05-10T10:30:00Z"
}
```

---

## 九、技术实现方案

### 9.1 技术栈

| 层级 | 技术选型 | 说明 |
|------|---------|------|
| 框架 | Vue 3.4+ + TypeScript 5 | 组合式API (Composition API) |
| 构建 | Vite 5 | 快速开发体验 |
| UI库 | shadcn/ui + Tailwind CSS | 基于Radix的组件库 |
| 状态 | Pinia | Vue官方推荐状态管理 |
| 路由 | Vue Router 4 | SPA路由管理 |
| HTTP | Axios + Fetch API | 请求拦截与错误处理 |
| 拖拽 | @vueuse/core | 拖拽、粘贴等工具函数 |
| 渲染 | marked + DOMPurify | Markdown渲染与安全过滤 |
| 动画 | Vue Transition + CSS | 页面过渡与微交互 |
| 主题 | CSS Variables + dark mode | 亮色/暗色主题支持 |

### 9.2 目录结构

```
frontend/src/
├── api/
│   ├── index.ts          # Axios实例配置
│   ├── task.ts           # 任务相关API
│   └── types.ts          # TypeScript类型定义
├── assets/
│   └── styles/
│       ├── main.css      # 全局样式
│       └── themes.css    # 主题变量
├── components/
│   ├── common/
│   │   ├── Toast.vue        # 轻提示组件
│   │   ├── Modal.vue        # 模态框
│   │   ├── Button.vue       # 按钮组件
│   │   └── Loading.vue      # 加载组件
│   ├── upload/
│   │   ├── FileUploader.vue     # 文件上传组件
│   │   ├── FileDropZone.vue     # 拖拽区域
│   │   ├── FileList.vue         # 文件列表
│   │   └── FormatSelector.vue   # 格式选择器
│   ├── task/
│   │   ├── TaskList.vue      # 任务列表
│   │   ├── TaskCard.vue      # 任务卡片
│   │   ├── TaskProgress.vue  # 任务进度
│   │   └── StatusBadge.vue   # 状态徽章
│   └── preview/
│       ├── ResultPreview.vue     # 结果预览
│       ├── MarkdownView.vue      # Markdown渲染
│       └── SearchHighlight.vue    # 搜索高亮
├── composables/
│   ├── useFileUpload.ts      # 文件上传逻辑
│   ├── useTaskPolling.ts     # 任务轮询
│   └── useClipboard.ts       # 剪贴板操作
├── layouts/
│   └── MainLayout.vue        # 主布局组件
├── router/
│   └── index.ts              # 路由配置
├── stores/
│   ├── task.ts               # 任务状态管理
│   ├── upload.ts             # 上传状态管理
│   └── user.ts               # 用户状态管理
├── utils/
│   ├── format.ts             # 格式化工具
│   └── validate.ts           # 验证工具
├── views/
│   ├── HomeView.vue          # 首页/主工作区
│   ├── ParserView.vue        # 解析页面
│   └── HistoryView.vue       # 历史记录
├── App.vue                   # 根组件
└── main.ts                   # 应用入口
```

### 9.3 关键实现细节

**1. 文件上传流程**：
```typescript
// 使用 Composable 封装上传逻辑
export function useFileUpload() {
  const files = ref<File[]>([])
  const uploading = ref(false)
  const progress = ref<Record<string, number>>({})

  async function upload(file: File, onProgress: (p: number) => void) {
    const formData = new FormData()
    formData.append('file', file)

    const response = await axios.post('/api/v1/tasks', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (e) => {
        const percent = Math.round((e.loaded * 100) / (e.total || 1))
        onProgress(percent)
      }
    })

    return response.data
  }

  return { files, uploading, progress, upload }
}
```

**2. 任务状态轮询**：
```typescript
// Composable 封装轮询逻辑
export function useTaskPolling(taskId: Ref<string>) {
  const status = ref<TaskStatus>('pending')
  const progress = ref(0)
  let intervalId: number | null = null

  async function poll() {
    try {
      const response = await taskApi.getStatus(taskId.value)
      status.value = response.status
      progress.value = response.progress

      if (status.value === 'completed' || status.value === 'failed') {
        stopPolling()
      }
    } catch (error) {
      console.error('轮询失败:', error)
    }
  }

  function startPolling(interval = 2000) {
    poll()
    intervalId = window.setInterval(poll, interval)
  }

  function stopPolling() {
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }
  }

  onUnmounted(stopPolling)

  return { status, progress, startPolling, stopPolling }
}
```

**3. Markdown安全渲染**：
```typescript
import { marked } from 'marked'
import DOMPurify from 'dompurify'

marked.setOptions({
  breaks: true,
  gfm: true
})

export function renderMarkdown(content: string): string {
  const raw = marked.parse(content) as string
  return DOMPurify.sanitize(raw, {
    ALLOWED_TAGS: ['h1', 'h2', 'h3', 'h4', 'p', 'ul', 'ol', 'li',
                    'blockquote', 'pre', 'code', 'table', 'thead', 'tbody',
                    'tr', 'th', 'td', 'strong', 'em', 'a', 'img', 'br'],
    ALLOWED_ATTR: ['href', 'src', 'alt', 'class']
  })
}
```

---

## 十、Accessibility (无障碍访问)

### 10.1 键盘导航

| 元素 | Tab键 | 方向键 | Enter键 | Escape键 |
|------|-------|--------|---------|----------|
| 文件上传区 | focus | - | 打开选择器 | 取消 |
| 格式选项 | focus | 切换选项 | 选择 | 取消 |
| 任务列表 | focus | 上下移动 | 打开详情 | 返回 |
| 模态框 | 聚焦首元素 | 导航 | 确认 | 关闭 |

### 10.2 ARIA标签

```html
<!-- 上传区域 -->
<div role="button"
     tabindex="0"
     aria-label="拖拽文件到此处或点击选择"
     aria-describedby="upload-hint">
</div>

<!-- 进度条 -->
<div role="progressbar"
     aria-valuenow="68"
     aria-valuemin="0"
     aria-valuemax="100"
     aria-label="文件解析进度">
</div>

<!-- 状态徽章 -->
<span role="status"
      aria-live="polite">
  解析完成
</span>
```

---

## 十一、性能优化

### 11.1 加载策略

- **路由懒加载**：使用 `defineAsyncComponent` 异步加载页面组件
- **组件懒加载**：非关键组件延迟加载
- **图片优化**：使用 `loading="lazy"` 和 WebP 格式
- **代码分割**：每个页面独立 chunk

### 11.2 渲染优化

- **虚拟列表**：文件列表和任务列表使用虚拟滚动
- **防抖节流**：搜索输入防抖300ms
- **计算缓存**：使用 `computed` 缓存计算结果
- **DOM复用**：使用 `v-memo` 减少不必要的重渲染

### 11.3 缓存策略

- **本地存储**：用户偏好、主题设置
- **会话存储**：上传临时状态
- **HTTP缓存**：GET 请求缓存控制

---

## 十二、测试计划

| 测试类型 | 覆盖范围 | 测试工具 |
|---------|---------|---------|
| 单元测试 | 工具函数、Composable | Vitest |
| 组件测试 | 组件渲染、交互 | Vue Test Utils |
| E2E测试 | 完整用户流程 | Playwright |
| 视觉测试 | UI一致性 | Percy / Chromatic |
| 性能测试 | 首屏加载、交互响应 | Lighthouse |
| 兼容性测试 | 浏览器、设备 | BrowserStack |

---

*文档版本：v2.0*
*最后更新：2026-05-10*
*作者：PDF解析助手开发团队*
