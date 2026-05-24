<template>
  <el-card>
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span>我的交易</span>
        <el-button v-if="auth.user?.role === 'seller'" type="primary" @click="showPublish = true">发布新交易</el-button>
      </div>
    </template>
    <el-table :data="list">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="watch_id" label="手表ID" />
      <el-table-column prop="start_price" label="起拍价" />
      <el-table-column prop="final_price" label="成交价" />
      <el-table-column prop="status" label="状态" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button v-if="row.status === 'open' || row.status === 'bidding'" size="small" @click="bid(row)">出价</el-button>
          <el-button v-if="row.status === 'bidding'" size="small" type="success" @click="accept(row)">接受出价</el-button>
          <el-button v-if="row.status === 'pending_deal'" size="small" type="success" @click="ship(row)">发货</el-button>
          <el-button v-if="row.status === 'shipped'" size="small" type="success" @click="complete(row)">确认收货</el-button>
          <el-button size="small" type="danger" @click="cancel(row)">取消</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showPublish" title="发布新交易" width="500px">
      <el-form :model="publishForm" label-width="100px">
        <el-form-item label="手表ID"><el-input-number v-model="publishForm.watch_id" :min="1" /></el-form-item>
        <el-form-item label="起拍价"><el-input-number v-model="publishForm.start_price" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="publishForm.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPublish = false">取消</el-button>
        <el-button type="primary" :loading="publishing" @click="submitPublish">发布</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import request from '@/utils/request'

const auth = useAuthStore()
const list = ref<any[]>([])
const showPublish = ref(false)
const publishing = ref(false)
const publishForm = reactive({ watch_id: 0, start_price: 0, remark: '' })

onMounted(loadList)

async function loadList() {
  const res: any = await request.get('/trades')
  list.value = res.list || []
}

async function submitPublish() {
  if (publishForm.watch_id <= 0) {
    ElMessage.warning('请选择手表')
    return
  }
  if (publishForm.start_price <= 0) {
    ElMessage.warning('请输入起拍价')
    return
  }
  publishing.value = true
  try {
    await request.post('/trades', publishForm)
    ElMessage.success('已发布')
    showPublish.value = false
    publishForm.watch_id = 0
    publishForm.start_price = 0
    publishForm.remark = ''
    loadList()
  } finally {
    publishing.value = false
  }
}

async function bid(row: any) {
  const { value } = await ElMessageBox.prompt('出价金额', '出价')
  await request.post(`/trades/${row.id}/bids`, { price: Number(value) })
  ElMessage.success('已出价')
  loadList()
}
async function accept(row: any) {
  const { value } = await ElMessageBox.prompt('接受的出价ID', '接受出价')
  await request.post(`/trades/${row.id}/accept`, { bid_id: Number(value) })
  ElMessage.success('已接受')
  loadList()
}
async function ship(row: any) {
  await request.post(`/trades/${row.id}/ship`)
  ElMessage.success('已发货')
  loadList()
}
async function complete(row: any) {
  await request.post(`/trades/${row.id}/complete`)
  ElMessage.success('交易完成')
  loadList()
}
async function cancel(row: any) {
  await ElMessageBox.confirm('确定取消该交易?', '提示')
  await request.post(`/trades/${row.id}/cancel`)
  ElMessage.success('已取消')
  loadList()
}
</script>
