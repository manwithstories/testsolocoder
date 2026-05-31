<template>
  <AppLayout>
    <div class="page">
      <h2>接单处理</h2>
      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="订单号" width="80" />
        <el-table-column label="服务">
          <template #default="{ row }">{{ row.booking?.service?.name }}</template>
        </el-table-column>
        <el-table-column prop="total_amount" label="金额" width="100" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag>{{ { created: '已创建', reported: '待验收', confirmed: '已确认', paid: '已结算' }[row.status] || row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" type="primary"
              v-if="row.status==='created' || row.status==='reported'"
              @click="report(row)">提交报告</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="reportVisible" title="提交服务报告">
      <el-form :model="reportForm" label-width="80px">
        <el-form-item label="描述"><el-input v-model="reportForm.report_text" type="textarea" :rows="4" /></el-form-item>
        <el-form-item label="图片"><el-input v-model="reportForm.report_images" placeholder="逗号分隔" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reportVisible=false">取消</el-button>
        <el-button type="primary" @click="submitReport">提交</el-button>
      </template>
    </el-dialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import AppLayout from '../../components/AppLayout.vue'
import { listOrders, submitReport } from '../../api/order'

const list = ref<any[]>([])
const reportVisible = ref(false)
const currentId = ref<number | null>(null)
const reportForm = reactive({ report_text: '', report_images: '' })

async function load() {
  const res = await listOrders({})
  list.value = (res.data as any).data || []
}

function report(row: any) {
  currentId.value = row.id
  reportForm.report_text = ''
  reportForm.report_images = ''
  reportVisible.value = true
}
async function submitReport() {
  if (!reportForm.report_text) { ElMessage.warning('请填写描述'); return }
  await submitReport(currentId.value!, reportForm)
  ElMessage.success('报告已提交')
  reportVisible.value = false
  load()
}

onMounted(load)
</script>
