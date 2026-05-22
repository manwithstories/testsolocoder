<template>
  <div class="page">
    <div class="page-title">{{ isEdit ? '编辑赛事' : '创建赛事' }}</div>
    <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
      <el-form-item prop="name" label="赛事名称"><el-input v-model="form.name" /></el-form-item>
      <el-form-item prop="description" label="简介"><el-input type="textarea" v-model="form.description" :rows="3" /></el-form-item>
      <el-form-item prop="location" label="地点"><el-input v-model="form.location" /></el-form-item>
      <el-form-item prop="start_date" label="开始时间"><el-date-picker v-model="form.start_date" type="datetime" value-format="YYYY-MM-DDTHH:mm:ssZ" /></el-form-item>
      <el-form-item prop="end_date" label="结束时间"><el-date-picker v-model="form.end_date" type="datetime" value-format="YYYY-MM-DDTHH:mm:ssZ" /></el-form-item>
      <el-form-item prop="registration_deadline" label="报名截止"><el-date-picker v-model="form.registration_deadline" type="datetime" value-format="YYYY-MM-DDTHH:mm:ssZ" /></el-form-item>
      <el-form-item prop="organizer" label="主办方"><el-input v-model="form.organizer" /></el-form-item>
      <el-form-item label="赛事项目">
        <div v-for="(it, idx) in form.items" :key="idx" class="card">
          <el-row :gutter="12">
            <el-col :span="6"><el-input v-model="it.name" placeholder="项目名称" /></el-col>
            <el-col :span="6"><el-input v-model="it.category" placeholder="类别" /></el-col>
            <el-col :span="4"><el-input v-model="it.gender" placeholder="性别" /></el-col>
            <el-col :span="3"><el-input-number v-model="it.min_age" :min="0" placeholder="最小年龄" /></el-col>
            <el-col :span="3"><el-input-number v-model="it.max_age" :min="0" placeholder="最大年龄" /></el-col>
            <el-col :span="6"><el-input-number v-model="it.quota" :min="1" placeholder="名额" /></el-col>
            <el-col :span="6"><el-input-number v-model="it.waitlist_quota" :min="0" placeholder="候补名额" /></el-col>
            <el-col :span="6"><el-input-number v-model="it.fee" :min="0" :precision="2" placeholder="费用" /></el-col>
            <el-col :span="24"><el-input v-model="it.requirements" placeholder="报名要求" /></el-col>
            <el-col :span="24"><el-button link type="danger" @click="removeItem(idx)">删除</el-button></el-col>
          </el-row>
        </div>
        <el-button type="primary" plain @click="addItem">添加项目</el-button>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submit">{{ isEdit ? '保存' : '创建' }}</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { eventApi } from '@/api'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)

function defaultItem() {
  return { name: '', category: '', gender: '', min_age: 0, max_age: 0, quota: 50, waitlist_quota: 20, fee: 0, requirements: '' }
}

const form = reactive<any>({
  name: '', description: '', location: '', start_date: '', end_date: '', registration_deadline: '',
  organizer: '', items: [defaultItem()],
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入赛事名称', trigger: 'blur' }],
  start_date: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  end_date: [{ required: true, message: '请选择结束时间', trigger: 'change' }],
  registration_deadline: [{ required: true, message: '请选择报名截止时间', trigger: 'change' }],
}

const formRef = ref<FormInstance>()

function addItem() { form.items.push(defaultItem()) }
function removeItem(idx: number) { form.items.splice(idx, 1) }

async function submit() {
  await formRef.value?.validate()
  if (isEdit.value) {
    await eventApi.update(Number(route.params.id), form)
  } else {
    await eventApi.create(form)
  }
  ElMessage.success('操作成功')
  router.push('/admin/events')
}

onMounted(async () => {
  if (isEdit.value) {
    const res = await eventApi.get(Number(route.params.id))
    const data = res.data as any
    Object.assign(form, data)
    if (data.items) {
      form.items = data.items.map((i: any) => ({ ...i }))
    }
  }
})
</script>
