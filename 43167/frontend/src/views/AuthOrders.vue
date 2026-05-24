<template>
  <el-card>
    <template #header>鉴定申请</template>
    <el-button type="primary" @click="showCreate = true">提交鉴定申请</el-button>
    <el-table :data="list" style="margin-top: 16px">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="status" label="状态" />
      <el-table-column prop="note" label="备注" />
      <el-table-column label="报告">
        <template #default="{ row }">
          <el-link v-if="row.report" type="primary" :href="row.report.pdf_path" target="_blank">下载PDF</el-link>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="showCreate" title="提交鉴定">
      <el-form :model="form" label-width="100px">
        <el-form-item label="手表ID"><el-input-number v-model="form.watch_id" :min="0" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.note" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" @click="submit">提交</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const list = ref<any[]>([])
const showCreate = ref(false)
const form = reactive({ watch_id: 0, note: '' })

async function loadList() {
  const res: any = await request.get('/auth-orders')
  list.value = res.list || []
}

onMounted(loadList)

async function submit() {
  await request.post('/auth-orders', form)
  ElMessage.success('已提交')
  showCreate.value = false
  loadList()
}
</script>
