<template>
  <AppLayout>
    <div class="page">
      <h2>工单管理</h2>
      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="编号" width="80" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type==='refund' ? 'danger' : 'warning'">{{ row.type==='refund' ? '退款' : '投诉' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag>{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button size="small" type="primary" v-if="row.status==='open' || row.status==='escalated'" @click="assign(row)">分配</el-button>
            <el-button size="small" type="success" v-if="row.status!=='resolved' && row.status!=='closed'" @click="resolve(row)">解决</el-button>
            <el-button size="small" v-if="row.status!=='closed'" @click="close(row)">关闭</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="assignVisible" title="分配工单">
      <el-form label-width="80px">
        <el-form-item label="客服ID"><el-input-number v-model="agentId" :min="1" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignVisible=false">取消</el-button>
        <el-button type="primary" @click="submitAssign">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="resolveVisible" title="处理工单">
      <el-form label-width="80px">
        <el-form-item label="处理结果"><el-input v-model="resolveResult" type="textarea" :rows="4" /></el-form-item>
        <el-form-item label="执行退款"><el-switch v-model="resolveRefund" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resolveVisible=false">取消</el-button>
        <el-button type="primary" @click="submitResolve">确定</el-button>
      </template>
    </el-dialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppLayout from '../../components/AppLayout.vue'
import { listTickets, assignTicket, resolveTicket, closeTicket } from '../../api/ticket'

const list = ref<any[]>([])
const assignVisible = ref(false)
const resolveVisible = ref(false)
const currentId = ref<number | null>(null)
const agentId = ref<number>(1)
const resolveResult = ref('')
const resolveRefund = ref(false)

async function load() {
  const res = await listTickets({})
  list.value = (res.data as any).data || []
}

function statusLabel(s: string) {
  return { open: '待处理', assigned: '已分配', pending: '处理中', resolved: '已解决', closed: '已关闭', escalated: '已升级' }[s] || s
}

function assign(row: any) { currentId.value = row.id; agentId.value = 1; assignVisible.value = true }
async function submitAssign() { await assignTicket(currentId.value!, agentId.value); ElMessage.success('已分配'); assignVisible.value = false; load() }
function resolve(row: any) { currentId.value = row.id; resolveResult.value = ''; resolveRefund.value = false; resolveVisible.value = true }
async function submitResolve() {
  await resolveTicket(currentId.value!, { result: resolveResult.value, refund: resolveRefund.value })
  ElMessage.success('已解决'); resolveVisible.value = false; load()
}
async function close(row: any) {
  await ElMessageBox.confirm('确认关闭?', '提示')
  await closeTicket(row.id); ElMessage.success('已关闭'); load()
}

onMounted(load)
</script>
