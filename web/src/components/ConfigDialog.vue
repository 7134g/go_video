<template>
  <el-dialog v-model="visible" title="系统配置" width="500px" @close="$emit('close')">
    <el-form :model="form" label-width="140px" v-loading="loading">
      <el-form-item label="并发任务数">
        <el-input-number v-model="form.max_concurrent_tasks" :min="1" :max="20" />
      </el-form-item>
      <el-form-item label="分片并发数">
        <el-input-number v-model="form.max_segment_workers" :min="1" :max="50" />
      </el-form-item>
      <el-form-item label="下载目录">
        <el-input v-model="form.download_dir" />
      </el-form-item>
      <el-form-item label="最大连续错误数">
        <el-input-number v-model="form.max_consecutive_errors" :min="1" :max="100" />
      </el-form-item>
      <el-form-item label="默认请求头">
        <div style="width: 100%">
          <div v-for="(_, key) in form.default_headers" :key="key" style="display: flex; gap: 8px; margin-bottom: 8px;">
            <el-input :model-value="key" placeholder="Header名称" style="width: 40%" @input="(v: string) => renameHeader(key, v)" />
            <el-input v-model="form.default_headers[key]" placeholder="Header值" style="flex: 1" />
            <el-button type="danger" :icon="Delete" circle size="small" @click="removeHeader(key)" />
          </div>
          <el-button type="primary" size="small" @click="addHeader">添加请求头</el-button>
        </div>
      </el-form-item>
      <el-form-item label="启用拦截器">
        <el-switch v-model="form.interceptor_enabled" />
      </el-form-item>
      <el-form-item label="拦截代理地址">
        <el-input v-model="form.agent_address" placeholder="127.0.0.1:8888" />
      </el-form-item>
      <el-form-item label="HTTP代理地址">
        <el-input v-model="form.http_proxy_address" placeholder="127.0.0.1:7890" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'
import { configApi, type Config } from '../api/config'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits(['update:modelValue', 'close'])

const visible = ref(props.modelValue)
const loading = ref(false)
const saving = ref(false)
const form = ref<Config>({
  max_concurrent_tasks: 3,
  max_segment_workers: 5,
  download_dir: './downloads',
  max_consecutive_errors: 10,
  default_headers: {},
  interceptor_enabled: false,
  agent_address: '127.0.0.1:8888',
  http_proxy_address: ''
})

watch(() => props.modelValue, async (val) => {
  visible.value = val
  if (val) await loadConfig()
})

watch(visible, (val) => emit('update:modelValue', val))

function addHeader() {
  const key = `Header-${Object.keys(form.value.default_headers).length + 1}`
  form.value.default_headers[key] = ''
}

function removeHeader(key: string) {
  delete form.value.default_headers[key]
}

function renameHeader(oldKey: string, newKey: string) {
  if (oldKey === newKey) return
  const val = form.value.default_headers[oldKey]
  delete form.value.default_headers[oldKey]
  if (val !== undefined) {
    form.value.default_headers[newKey] = val
  }
}

async function loadConfig() {
  loading.value = true
  try {
    const { data } = await configApi.get()
    form.value = data
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await configApi.update(form.value)
    ElMessage.success('保存成功')
    visible.value = false
  } finally {
    saving.value = false
  }
}
</script>
