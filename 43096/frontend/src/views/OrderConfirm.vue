<template>
  <Layout>
    <div class="order-confirm">
      <div class="form-card">
        <h2>订单确认</h2>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="100px"
        >
          <el-form-item label="真实姓名" prop="realName">
            <el-input v-model="form.realName" placeholder="请输入真实姓名" />
          </el-form-item>

          <el-form-item label="身份证号" prop="idCard">
            <el-input v-model="form.idCard" placeholder="请输入身份证号" maxlength="18" />
          </el-form-item>

          <el-form-item label="手机号" prop="phone">
            <el-input v-model="form.phone" placeholder="请输入手机号" maxlength="11" />
          </el-form-item>

          <el-form-item label="邮箱" prop="email">
            <el-input v-model="form.email" placeholder="请输入邮箱（用于接收电子票）" />
          </el-form-item>

          <el-form-item label="备注">
            <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="选填" />
          </el-form-item>
        </el-form>
      </div>

      <div class="order-summary">
        <h3>订单摘要</h3>
        <div class="summary-item">
          <span>已选座位：</span>
          <span>{{ seatLabels.join('、') }}</span>
        </div>
        <div class="summary-item">
          <span>座位数量：</span>
          <span>{{ seatIds.length }} 张</span>
        </div>
        <div class="summary-item total">
          <span>应付金额：</span>
          <span class="price">¥{{ totalPrice.toFixed(2) }}</span>
        </div>

        <div class="payment-method">
          <h4>支付方式</h4>
          <div class="payment-options">
            <label class="payment-option">
              <el-radio v-model="payType" :label="1">
                <img src="https://img.icons8.com/color/48/alipay.png" alt="支付宝" class="pay-icon" />
                支付宝
              </el-radio>
            </label>
            <label class="payment-option">
              <el-radio v-model="payType" :label="2">
                <img src="https://img.icons8.com/color/48/weixing.png" alt="微信支付" class="pay-icon" />
                微信支付
              </el-radio>
            </label>
          </div>
        </div>

        <el-button
          type="primary"
          size="large"
          class="pay-btn"
          @click="handleSubmit"
          :loading="submitting"
        >
          立即支付 ¥{{ totalPrice.toFixed(2) }}
        </el-button>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { orderApi, seatApi } from '@/api'
import { useSeatStore } from '@/store'
import Layout from '@/components/Layout.vue'

const route = useRoute()
const router = useRouter()
const seatStore = useSeatStore()

const formRef = ref<FormInstance>()
const submitting = ref(false)
const payType = ref(1)

const sessionId = ref(0)
const seatIds = ref<number[]>([])
const totalPrice = ref(0)
const seatLabels = ref<string[]>([])

const form = reactive({
  realName: '',
  idCard: '',
  phone: '',
  email: '',
  remark: ''
})

const rules: FormRules = {
  realName: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' },
    { pattern: /^[\u4e00-\u9fa5a-zA-Z·]{2,20}$/, message: '请输入正确的姓名', trigger: 'blur' }
  ],
  idCard: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, message: '请输入正确的身份证号', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱', trigger: 'blur' }
  ]
}

async function handleSubmit() {
  await formRef.value?.validate()

  submitting.value = true
  try {
    const order = await orderApi.create({
      session_id: sessionId.value,
      seat_ids: seatIds.value,
      real_name: form.realName,
      id_card: form.idCard,
      phone: form.phone,
      email: form.email,
      remark: form.remark
    })

    const payResult = await orderApi.pay({
      order_no: order.order_no,
      pay_type: payType.value
    })

    seatStore.clearSelected()
    seatStore.unlockSeats(seatIds.value)

    ElMessage.success('支付成功！')
    router.push(`/order/${order.order_no}`)
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  const query = route.query
  if (query) {
    sessionId.value = Number(query.sessionId) || 0
    seatIds.value = (query.seatIds as string)?.split(',').map(Number) || []
    totalPrice.value = Number(query.totalPrice) || 0
    seatLabels.value = []
  }
})
</script>

<style lang="scss" scoped>
.order-confirm {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
  display: flex;
  gap: 24px;
}

.form-card {
  flex: 1;
  background: white;
  padding: 30px;
  border-radius: 12px;

  h2 {
    margin: 0 0 24px 0;
    font-size: 20px;
  }
}

.order-summary {
  width: 320px;
  flex-shrink: 0;
  background: white;
  padding: 24px;
  border-radius: 12px;
  height: fit-content;
  position: sticky;
  top: 20px;

  h3 {
    margin: 0 0 16px 0;
    font-size: 18px;
    padding-bottom: 12px;
    border-bottom: 1px solid #eee;
  }

  .summary-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 12px;
    font-size: 14px;
    color: #666;

    &.total {
      margin-top: 16px;
      padding-top: 16px;
      border-top: 1px solid #eee;
      font-size: 16px;
      font-weight: 600;

      .price {
        color: #f56c6c;
        font-size: 24px;
      }
    }
  }

  .payment-method {
    margin-top: 20px;

    h4 {
      margin: 0 0 12px 0;
      font-size: 14px;
    }
  }

  .payment-options {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .payment-option {
    padding: 12px;
    border: 1px solid #e4e7ed;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;

    .pay-icon {
      width: 24px;
      height: 24px;
      margin-right: 8px;
    }
  }

  .pay-btn {
    width: 100%;
    height: 48px;
    font-size: 16px;
    margin-top: 20px;
  }
}
</style>
