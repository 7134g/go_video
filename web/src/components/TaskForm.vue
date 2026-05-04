<template>
  <el-dialog :model-value="modelValue" :title="task ? '编辑任务' : '新建任务'" width="500" @close="handleClose">
    <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="URL" prop="url">
        <el-input v-model="form.url" />
      </el-form-item>
      <el-form-item label="请求头" prop="header">
        <el-input v-model="form.header" type="textarea" :rows="2" />
      </el-form-item>
      <el-form-item label="类型" prop="type">
        <el-select v-model="form.type" style="width: 100%">
          <el-option label="MP4" value="mp4" />
          <el-option label="M3U8" value="m3u8" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { taskApi, type Task } from '../api/task'

const props = defineProps<{ modelValue: boolean; task?: Task }>()
const emit = defineEmits<{ 'update:modelValue': [v: boolean]; success: []; close: [] }>()

const formRef = ref<FormInstance>()
const form = ref({ name: '', url: '', header: '', type: 'mp4' })

const rules: FormRules = {
  name: [{ required: true, message: '请输入名称' }],
  url: [{ required: true, message: '请输入URL' }],
  type: [{ required: true, message: '请选择类型' }],
}

watch(() => props.task, (t) => {
  if (t) {
    form.value = { name: t.name, url: t.url, header: t.header, type: t.type }
  } else {
    form.value = { name: '', url: '', header: '', type: 'mp4' }
  }
}, { immediate: true })

function handleClose() {
  emit('update:modelValue', false)
  emit('close')
  formRef.value?.resetFields()
}

async function handleSubmit() {
  await formRef.value?.validate()
  if (props.task) {
    await taskApi.update({ id: props.task.id, ...form.value })
    ElMessage.success('更新成功')
  } else {
    await taskApi.create(form.value)
    ElMessage.success('创建成功')
  }
  emit('success')
  handleClose()
}
</script>
