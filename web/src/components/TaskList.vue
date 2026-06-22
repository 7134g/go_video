<template>
  <div class="page-container">
    <div class="main-content">
      <div class="toolbar">
        <el-button type="primary" @click="loadTasks">刷新</el-button>
        <el-button type="primary" @click="showForm = true">新建任务</el-button>
        <el-button type="success" @click="handleStart">开始下载</el-button>
        <el-button type="warning" @click="handleStopAll">暂停全部</el-button>
      </div>

      <div>
        <el-select v-model="statusFilter" placeholder="筛选状态" clearable @change="loadTasks" style="width: 120px">
          <el-option label="全部" :value="null" />
          <el-option label="未开始" :value="0" />
          <el-option label="执行中" :value="1" />
          <el-option label="完成" :value="2" />
          <el-option label="失败" :value="3" />
          <el-option label="已暂停" :value="4" />
        </el-select>
        <el-button @click="showConfig = true">配置</el-button>
      </div>


      <el-table :data="tasks" v-loading="loading" empty-text="暂无任务">
        <el-table-column prop="name" label="名称" />
        <el-table-column label="URL" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="url-cell" @dblclick="copyURL(row.url)">{{ row.url }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" width="150">
          <template #default="{ row }">
            <el-progress v-if="row.status === 1" :percentage="progress[row.id] || 0" :stroke-width="10" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <div  v-if="row.status === 0">
              <el-button size="small" type="success" @click="handleStartOne(row.id)">启动</el-button>
              <el-button size="small" type="info" @click="handleUpdateTitle(row.id)">更新标题</el-button>
            </div>
            <el-button size="small" @click="handleEdit(row)" v-if="row.status === 0 || row.status === 2 || row.status === 3 || row.status === 4">编辑</el-button>
            <el-button size="small" type="success" @click="handleRedownload(row.id)" v-if="row.status === 2">重新下载</el-button>
            <el-button size="small" type="warning" @click="handlePause(row.id)" v-if="row.status === 1">暂停</el-button>
            <el-button size="small" type="primary" @click="handleRetry(row.id)" v-if="row.status === 3 || row.status === 4">重试</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.id)" v-if="row.status !== 1">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <TaskForm v-model="showForm" :task="editTask" @success="loadTasks" @close="editTask = undefined" />
      <ConfigDialog v-model="showConfig" />
    </div>

    <div class="right-panel">
<!--      <div class="progress-panel">-->
<!--        <div class="panel-header">-->
<!--          <span>任务进度</span>-->
<!--        </div>-->
<!--        <div class="progress-content">-->
<!--          <div v-for="t in taskProgressList" :key="t.id" class="progress-item">-->
<!--            <div class="progress-name">{{ t.name }}</div>-->
<!--            <el-progress :percentage="t.percent" :stroke-width="8" />-->
<!--            <div class="progress-segment">{{ t.segment_done }}/{{ t.segment_all }} 段</div>-->
<!--          </div>-->
<!--          <div v-if="!taskProgressList.length" class="ws-empty">暂无进行中的任务</div>-->
<!--        </div>-->
<!--      </div>-->

      <div class="ws-log">
        <div class="ws-header">
          <span>WebSocket 日志</span>
          <el-tag :type="wsConnected ? 'success' : 'danger'" size="small">{{ wsConnected ? '已连接' : '未连接' }}</el-tag>
          <el-button size="small" @click="wsLogs = []">清空</el-button>
        </div>
        <div class="ws-content" ref="logContainer">
          <div v-for="(log, i) in wsLogs" :key="i" class="ws-line">
            <span class="ws-time">{{ log.time }}</span>
            <pre class="ws-data">{{ log.data }}</pre>
          </div>
          <div v-if="!wsLogs.length" class="ws-empty">暂无数据</div>
        </div>
      </div>
    </div>


  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { taskApi, type Task } from '../api/task'
import TaskForm from './TaskForm.vue'
import ConfigDialog from './ConfigDialog.vue'

const tasks = ref<Task[]>([])
const loading = ref(false)
const showForm = ref(false)
const showConfig = ref(false)
const editTask = ref<Task>()
const statusFilter = ref<number>()
const progress = ref<Record<number, number>>({})
interface TaskProgress {
  id: number
  name: string
  type: string
  done: number
  total: number
  percent: number
}
const taskProgressList = ref<TaskProgress[]>([])
const wsConnected = ref(false)
const wsLogs = ref<{ time: string; data: string }[]>([])
const logContainer = ref<HTMLElement>()
let ws: WebSocket | null = null
let refreshTimer: number | null = null

const statusText = (s: number) => ['待执行', '执行中', '完成', '失败', '已暂停'][s]
const statusType = (s: number) => {
  const types = ['info', 'warning', 'success', 'danger', 'info'] as const
  return types[s]
}

async function loadTasks() {
  loading.value = true
  try {
    const { data } = await taskApi.list(statusFilter.value)
    tasks.value = data
  } finally {
    loading.value = false
  }
}

async function copyURL(url: string) {
  await navigator.clipboard.writeText(url)
  ElMessage.success('URL 已复制')
}

function handleEdit(task: Task) {
  editTask.value = task
  showForm.value = true
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm('确定删除该任务？', '提示')
  await taskApi.delete(id)
  ElMessage.success('删除成功')
  loadTasks()
}

async function handleStart() {
  const { data } = await taskApi.start()
  ElMessage.success(`已启动 ${data.started} 个任务`)
  loadTasks()
}

async function handleStartOne(id: number) {
  await taskApi.startOne(id)
  ElMessage.success('任务已启动')
  loadTasks()
}

async function handlePause(id: number) {
  await taskApi.pause(id)
  ElMessage.success('任务已暂停')
  loadTasks()
}

async function handleRedownload(id: number) {
  await ElMessageBox.confirm('确定重新下载？已下载的文件将被删除。', '提示')
  await taskApi.redownload(id)
  ElMessage.success('任务已重置，可重新开始下载')
  loadTasks()
}

async function handleRetry(id: number) {
  await taskApi.retry(id)
  ElMessage.success('任务已重新启动')
  loadTasks()
}

async function handleStopAll() {
  await ElMessageBox.confirm('确定暂停所有进行中的任务？', '提示')
  await taskApi.stopAll()
  ElMessage.success('已停止所有任务')
  loadTasks()
}

async function handleUpdateTitle(id: number) {
  const { data } = await taskApi.updateTitle(id)
  if (data.name) {
    ElMessage.success(`标题已更新: ${data.name}`)
  } else {
    ElMessage.warning('未在 WebTree 中找到匹配的标题')
  }
  loadTasks()
}

function connectWS() {
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws = new WebSocket(`${protocol}//${location.host}/api/tasks/progress`)

  ws.onopen = () => {
    wsConnected.value = true
    addLog('[连接成功]')
  }

  ws.onmessage = (e) => {
    const data = JSON.parse(e.data)
    if (!data) return
    addLog(e.data)
    if (Array.isArray(data)) {
      taskProgressList.value = data as TaskProgress[]
      for (const item of data) {
        progress.value[item.id] = item.percent
      }
    }
  }

ws.onclose = () => {
    wsConnected.value = false
    addLog('[连接断开，3秒后重连...]')
    setTimeout(connectWS, 3000)
  }

  ws.onerror = () => {
    addLog('[连接错误]')
  }
}

function addLog(data: string) {
  const time = new Date().toLocaleTimeString()
  const el = logContainer.value
  const atBottom = el
    ? el.scrollHeight - el.scrollTop - el.clientHeight < 50
    : false
  wsLogs.value.push({ time, data })
  if (wsLogs.value.length > 100) wsLogs.value.shift()
  if (atBottom) {
    nextTick(() => {
      if (logContainer.value) {
        logContainer.value.scrollTop = logContainer.value.scrollHeight
      }
    })
  }
}

onMounted(() => {
  loadTasks()
  connectWS()
  refreshTimer = window.setInterval(() => {
    if (wsConnected.value) loadTasks()
  }, 5000)
})

onUnmounted(() => {
  ws?.close()
  if (refreshTimer) clearInterval(refreshTimer)
})
</script>

<style scoped>
.page-container {
  display: flex;
  height: 100vh;
}
.main-content {
  width: 66.67%;
  max-width: 66.67%;
  padding: 20px;
  overflow-x: hidden;
  overflow-y: auto;
  box-sizing: border-box;
}
.url-cell { cursor: pointer; }
.toolbar { margin-bottom: 16px; }

.main-content :deep(.el-table) {
  width: 100% !important;
}

.right-panel {
  display: flex;
  flex-direction: column;
  border-left: 1px solid #dcdfe6;
  background: #fff;
  position: fixed;
  right: 0;
  top: 0;
  bottom: 0;
  width: 33.33%;
  z-index: 1000;
}
.progress-panel {
  display: flex;
  flex-direction: column;
  height: 50%;
  border-bottom: 1px solid #dcdfe6;
}
.panel-header {
  padding: 10px 15px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
  font-weight: 500;
}
.progress-content {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}
.progress-item {
  margin-bottom: 12px;
}
.progress-name {
  font-size: 13px;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.progress-segment {
  font-size: 11px;
  color: #909399;
  margin-top: 2px;
}
.ws-log {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}
.ws-header {
  padding: 10px 15px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}
.ws-content {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  font-family: monospace;
  font-size: 12px;
  background: #1e1e1e;
  color: #d4d4d4;
}
.ws-line {
  display: flex;
  gap: 10px;
  margin-bottom: 4px;
}
.ws-time {
  color: #6a9955;
  flex-shrink: 0;
}
.ws-data {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}
.ws-empty {
  color: #666;
  text-align: center;
  padding: 20px;
}
</style>
