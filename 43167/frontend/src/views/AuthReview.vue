<template>
  <el-card>
    <template #header>鉴定审核</template>
    <el-table :data="list">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="status" label="状态" />
      <el-table-column prop="note" label="备注" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" size="small" type="primary" @click="assign(row)">接单</el-button>
          <el-button v-if="row.status === 'assigned'" size="small" type="success" @click="openReport(row)">出具报告</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showReport" title="出具鉴定报告">
      <el-form :model="report" label-width="100px">
        <el-form-item label="结论">
          <el-select v-model="report.conclusion">
            <el-option label="正品" value="authentic" />
            <el-option label="赝品" value="fake" />
            <el-option label="存疑" value="uncertain" />
          </el-select>
        </el-form-item>
        <el-form-item label="是否正品"><el-switch v-model="report.authentic" /></el-form-item>
        <el-form-item label="估值"><el-input-number v-model="report.estimated_value" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="详情"><el-input v-model="report.details" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showReport = false">取消</el-button>
        <el-button type="primary" @click="submitReport">提交</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const list = ref<any[]>([])
const showReport = ref(false)
const current = ref<number | null>(null)
const report = reactive({ conclusion: 'authentic', authentic: true, estimated_value: 0, details: '' })

async function loadList() {
  const res: any = await request.get('/auth-orders')
  list.value = res.list || []
}

onMounted(loadList)

async function assign(row: any) {
  await request.post(`/auth-orders/${row.id}/assign`)
  ElMessage.success('已接单')
  loadList()
}
function openReport(row: any) {
  current.value = row.id
  showReport.value = true
}
async function submitReport() {
  await request.post(`/auth-orders/${current.value}/report`, report)
  ElMessage.success('报告已提交')
  showReport.value = false
  loadList()
}
</script>
