<template>
  <div class="delivery-form">
    <el-page-header @back="router.back()">
      <template #content>
        <span>{{ isEdit ? '编辑预约' : '预约安装' }}</span>
      </template>
    </el-page-header>

    <el-card shadow="never" style="margin-top: 16px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width: 720px">
        <el-form-item label="订单号" prop="orderNo">
          <el-input v-model="form.orderNo" placeholder="请输入订单号" />
        </el-form-item>
        <el-form-item label="联系人" prop="contactName">
          <el-input v-model="form.contactName" placeholder="请输入联系人姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="contactPhone">
          <el-input v-model="form.contactPhone" placeholder="请输入联系电话" maxlength="11" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input
            v-model="form.address"
            type="textarea"
            :rows="2"
            placeholder="请输入详细地址"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="常用地址">
          <el-select
            v-model="selectedAddressId"
            placeholder="从常用地址选择"
            clearable
            filterable
            style="width: 100%"
            @change="handleSelectAddress"
          >
            <el-option
              v-for="addr in addressBook"
              :key="addr.id"
              :label="`${addr.name} - ${addr.address}`"
              :value="addr.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="配送日期" prop="deliveryDate">
          <el-date-picker
            v-model="form.deliveryDate"
            type="date"
            placeholder="请选择配送日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="时间段" prop="timeSlot">
          <el-select v-model="form.timeSlot" placeholder="请选择时间段" style="width: 100%">
            <el-option label="上午 09:00 - 12:00" value="上午 09:00 - 12:00" />
            <el-option label="下午 13:00 - 17:00" value="下午 13:00 - 17:00" />
            <el-option label="晚间 18:00 - 21:00" value="晚间 18:00 - 21:00" />
          </el-select>
        </el-form-item>
        <el-form-item label="安装师傅">
          <el-input v-model="form.installer" placeholder="请输入安装师傅姓名" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="请输入备注信息" maxlength="200" show-word-limit />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">提交</el-button>
          <el-button @click="router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import {
  getDelivery,
  createDelivery,
  updateDelivery,
  listAddressBook,
  type DeliveryFormData,
  type AddressBook
} from '@/api/delivery'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const selectedAddressId = ref<number | ''>('')
const addressBook = ref<AddressBook[]>([])

const isEdit = computed(() => !!route.params.id)

const form = reactive<DeliveryFormData & { orderNo?: string }>({
  orderNo: '',
  contactName: '',
  contactPhone: '',
  address: '',
  deliveryDate: '',
  timeSlot: '',
  installer: '',
  remark: '',
  status: 0
})

const rules: FormRules = {
  contactName: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  contactPhone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  address: [{ required: true, message: '请输入地址', trigger: 'blur' }],
  deliveryDate: [{ required: true, message: '请选择配送日期', trigger: 'change' }],
  timeSlot: [{ required: true, message: '请选择时间段', trigger: 'change' }]
}

function handleSelectAddress(id: number) {
  const addr = addressBook.value.find((a) => a.id === id)
  if (addr) {
    form.contactName = addr.contactName
    form.contactPhone = addr.contactPhone
    form.address = addr.address
  }
}

async function fetchAddressBook() {
  try {
    addressBook.value = await listAddressBook()
  } catch {
    addressBook.value = []
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (isEdit.value) {
        await updateDelivery(route.params.id as string, form)
        ElMessage.success('更新成功')
      } else {
        await createDelivery(form)
        ElMessage.success('预约成功')
      }
      router.back()
    } finally {
      submitting.value = false
    }
  })
}

async function fetchDetail() {
  if (!isEdit.value) return
  const data = await getDelivery(route.params.id as string)
  Object.assign(form, data)
}

onMounted(() => {
  fetchDetail()
  fetchAddressBook()
})
</script>

<style lang="scss" scoped>
.delivery-form {
}
</style>
