<template>
  <el-dialog
    v-model="visible"
    title="退房结算"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div v-if="checkIn" class="checkout-content">
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <span>入住信息</span>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="入住号">
            {{ checkIn.checkInNo }}
          </el-descriptions-item>
          <el-descriptions-item label="房间号">
            {{ checkIn.room?.roomNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="客人姓名">
            {{ checkIn.guestName }}
          </el-descriptions-item>
          <el-descriptions-item label="联系电话">
            {{ checkIn.guestPhone }}
          </el-descriptions-item>
          <el-descriptions-item label="入住时间">
            {{ formatDateTime(checkIn.checkInTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="预计退房">
            {{ formatDateTime(checkIn.expectedCheckOutTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="入住天数">
            {{ stayDays }} 天
          </el-descriptions-item>
          <el-descriptions-item label="押金">
            ¥{{ checkIn.deposit }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card class="fee-card">
        <template #header>
          <div class="card-header">
            <span>费用明细</span>
            <el-button type="primary" link size="small" @click="showAddCharge = true">
              添加费用
            </el-button>
          </div>
        </template>
        <el-table :data="feeList" border>
          <el-table-column prop="name" label="费用项目" />
          <el-table-column prop="amount" label="金额" width="120">
            <template #default="{ row }">
              ¥{{ row.amount }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80">
            <template #default="{ row, $index }">
              <el-button
                v-if="row.editable"
                type="danger"
                link
                size="small"
                @click="removeCharge($index)"
              >删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="fee-summary">
          <div class="fee-item">
            <span>房费总计：</span>
            <span class="fee-value">¥{{ roomFee }}</span>
          </div>
          <div class="fee-item">
            <span>额外费用：</span>
            <span class="fee-value">¥{{ extraFee }}</span>
          </div>
          <div class="fee-item">
            <span>押金抵扣：</span>
            <span class="fee-value deduction">-¥{{ checkIn.deposit }}</span>
          </div>
          <div class="fee-item total">
            <span>应付金额：</span>
            <span class="fee-value total-value">¥{{ payableAmount }}</span>
          </div>
        </div>
      </el-card>

      <el-card class="payment-card">
        <template #header>
          <div class="card-header">
            <span>支付方式</span>
          </div>
        </template>
        <el-radio-group v-model="paymentForm.method">
          <el-radio :value="PaymentMethod.CASH">现金</el-radio>
          <el-radio :value="PaymentMethod.WECHAT">微信支付</el-radio>
          <el-radio :value="PaymentMethod.ALIPAY">支付宝</el-radio>
          <el-radio :value="PaymentMethod.CREDIT_CARD">信用卡</el-radio>
          <el-radio :value="PaymentMethod.DEBIT_CARD">借记卡</el-radio>
          <el-radio :value="PaymentMethod.TRANSFER">转账</el-radio>
        </el-radio-group>
        <el-form :model="paymentForm" label-width="80px" class="payment-form">
          <el-form-item label="交易号">
            <el-input v-model="paymentForm.transactionId" placeholder="请输入交易号" />
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="paymentForm.remark" type="textarea" :rows="2" placeholder="请输入备注" />
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <el-dialog v-model="showAddCharge" title="添加费用" width="400px">
      <el-form :model="addChargeForm" label-width="80px">
        <el-form-item label="费用项目">
          <el-select v-model="addChargeForm.name" placeholder="请选择费用项目">
            <el-option label="迷你吧" value="迷你吧" />
            <el-option label="洗衣服务" value="洗衣服务" />
            <el-option label="餐饮" value="餐饮" />
            <el-option label="电话" value="电话" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="金额">
          <el-input-number v-model="addChargeForm.amount" :min="0" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="addChargeForm.description" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddCharge = false">取消</el-button>
        <el-button type="primary" @click="addCharge">确认添加</el-button>
      </template>
    </el-dialog>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleCheckOut">
        确认退房
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { PaymentMethod, type CheckIn } from '@/types'
import { checkOut } from '@/api/checkin'

const props = defineProps<{
  modelValue: boolean
  checkIn: CheckIn | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const submitting = ref(false)
const showAddCharge = ref(false)

const paymentForm = reactive({
  method: PaymentMethod.CASH,
  transactionId: '',
  remark: ''
})

const addChargeForm = reactive({
  name: '',
  amount: 0,
  description: ''
})

const extraCharges = ref<Array<{ name: string; amount: number; description?: string; editable: boolean }>>([])

const stayDays = computed(() => {
  if (!props.checkIn) return 0
  const start = new Date(props.checkIn.checkInTime)
  const end = new Date()
  const diff = end.getTime() - start.getTime()
  return Math.max(1, Math.ceil(diff / (1000 * 60 * 60 * 24)))
})

const roomFee = computed(() => {
  if (!props.checkIn) return 0
  const pricePerDay = props.checkIn.room?.roomType?.price || 200
  return pricePerDay * stayDays.value
})

const extraFee = computed(() => {
  return extraCharges.value.reduce((sum, item) => sum + item.amount, 0)
})

const payableAmount = computed(() => {
  if (!props.checkIn) return 0
  return Math.max(0, roomFee.value + extraFee.value - props.checkIn.deposit)
})

const feeList = computed(() => {
  const list = [
    { name: '房费', amount: roomFee.value, editable: false }
  ]
  return [...list, ...extraCharges.value]
})

const formatDateTime = (date: string) => {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

const addCharge = () => {
  if (!addChargeForm.name || !addChargeForm.amount) {
    ElMessage.warning('请填写完整信息')
    return
  }
  extraCharges.value.push({
    name: addChargeForm.name,
    amount: addChargeForm.amount,
    description: addChargeForm.description,
    editable: true
  })
  addChargeForm.name = ''
  addChargeForm.amount = 0
  addChargeForm.description = ''
  showAddCharge.value = false
  ElMessage.success('添加成功')
}

const removeCharge = (index: number) => {
  extraCharges.value.splice(index, 1)
}

const handleCheckOut = async () => {
  if (!props.checkIn) return
  submitting.value = true
  try {
    await checkOut(props.checkIn.id, {
      actualCheckOutTime: new Date().toISOString().slice(0, 19).replace('T', ' '),
      extraCharges: extraFee.value
    })
    ElMessage.success('退房成功')
    emit('success')
    handleClose()
  } catch (e: any) {
    ElMessage.error(e.message || '退房失败')
  } finally {
    submitting.value = false
  }
}

const handleClose = () => {
  extraCharges.value = []
  paymentForm.method = PaymentMethod.CASH
  paymentForm.transactionId = ''
  paymentForm.remark = ''
  visible.value = false
}

watch(() => props.modelValue, (val) => {
  if (val) {
    extraCharges.value = []
  }
})
</script>

<style scoped>
.checkout-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.fee-summary {
  margin-top: 16px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.fee-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.fee-item:last-child {
  margin-bottom: 0;
}

.fee-value {
  font-weight: 500;
}

.fee-value.deduction {
  color: #67c23a;
}

.fee-item.total {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e4e7ed;
  font-size: 16px;
}

.total-value {
  color: #f56c6c;
  font-weight: 600;
  font-size: 18px;
}

.payment-form {
  margin-top: 16px;
}

:deep(.el-descriptions__label) {
  width: 100px;
}
</style>
