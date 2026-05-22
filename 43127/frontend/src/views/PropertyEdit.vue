<template>
  <div class="property-edit">
    <div class="page-header">
      <h2 class="page-title">{{ isEdit ? '编辑房源' : '发布房源' }}</h2>
      <el-button @click="goBack">返回</el-button>
    </div>

    <div class="card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="房源标题" prop="title">
              <el-input v-model="form.title" placeholder="请输入房源标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="小区名称" prop="community">
              <el-input v-model="form.community" placeholder="请输入小区名称" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="详细地址" prop="address">
              <el-input v-model="form.address" placeholder="请输入详细地址" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所在区域" prop="region">
              <el-select v-model="form.region" placeholder="请选择区域" style="width: 100%">
                <el-option label="浦东" value="浦东" />
                <el-option label="徐汇" value="徐汇" />
                <el-option label="静安" value="静安" />
                <el-option label="长宁" value="长宁" />
                <el-option label="黄浦" value="黄浦" />
                <el-option label="普陀" value="普陀" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="面积(㎡)" prop="area">
              <el-input-number v-model="form.area" :min="0" :precision="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="户型" prop="layout">
              <el-select v-model="form.layout" placeholder="请选择户型" style="width: 100%">
                <el-option label="一室" value="一室" />
                <el-option label="两室" value="两室" />
                <el-option label="三室" value="三室" />
                <el-option label="四室" value="四室" />
                <el-option label="五室及以上" value="五室及以上" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="楼层" prop="floor">
              <el-input v-model="form.floor" placeholder="如: 6/18" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="月租金(元)" prop="rent">
              <el-input-number v-model="form.rent" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="押金(元)" prop="deposit">
              <el-input-number v-model="form.deposit" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="付款方式" prop="paymentType">
              <el-select v-model="form.paymentType" placeholder="请选择" style="width: 100%">
                <el-option label="押一付一" value="押一付一" />
                <el-option label="押一付三" value="押一付三" />
                <el-option label="押一付六" value="押一付六" />
                <el-option label="年付" value="年付" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="楼栋" prop="building">
              <el-input v-model="form.building" placeholder="如: A栋" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="房间号" prop="roomNo">
              <el-input v-model="form.roomNo" placeholder="如: 1801" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="房源描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请输入房源描述"
          />
        </el-form-item>

        <el-form-item label="房源图片">
          <div class="image-upload">
            <div
              v-for="(url, index) in form.imageUrls"
              :key="index"
              class="image-item"
            >
              <img :src="url" alt="" />
              <button class="delete-btn" @click="removeImage(index)">×</button>
            </div>
            <el-upload
              :show-file-list="false"
              :before-upload="beforeUpload"
              :http-request="handleUpload"
              accept="image/*"
            >
              <div class="el-upload--picture-card">
                <el-icon><Plus /></el-icon>
              </div>
            </el-upload>
          </div>
        </el-form-item>

        <el-form-item label="配套设施">
          <el-checkbox-group v-model="form.facilityIds">
            <el-checkbox
              v-for="facility in facilities"
              :key="facility.id"
              :label="facility.id"
            >
              {{ facility.icon }} {{ facility.name }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            保存
          </el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadProps } from 'element-plus'
import type { Facility } from '@/types'
import { getProperty, createProperty, updateProperty, uploadImage, getFacilities } from '@/api/property'

const route = useRoute()
const router = useRouter()

const isEdit = ref(!!route.params.id)
const loading = ref(false)
const formRef = ref<FormInstance>()
const facilities = ref<Facility[]>([])

const form = reactive({
  title: '',
  community: '',
  address: '',
  region: '',
  area: 0,
  layout: '',
  floor: '',
  rent: 0,
  deposit: 0,
  paymentType: '',
  building: '',
  roomNo: '',
  description: '',
  imageUrls: [] as string[],
  facilityIds: [] as number[]
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入房源标题', trigger: 'blur' }],
  community: [{ required: true, message: '请输入小区名称', trigger: 'blur' }],
  rent: [{ required: true, message: '请输入月租金', trigger: 'blur' }],
  layout: [{ required: true, message: '请选择户型', trigger: 'change' }]
}

onMounted(async () => {
  await loadFacilities()
  if (isEdit.value) {
    await loadProperty()
  }
})

async function loadFacilities() {
  try {
    const res = await getFacilities()
    facilities.value = res.data
  } catch (error) {
    console.error('Failed to load facilities:', error)
  }
}

async function loadProperty() {
  try {
    const res = await getProperty(Number(route.params.id))
    const data = res.data
    form.title = data.title
    form.community = data.community
    form.address = data.address
    form.region = data.region
    form.area = data.area
    form.layout = data.layout
    form.floor = data.floor
    form.rent = data.rent
    form.deposit = data.deposit
    form.paymentType = data.paymentType
    form.building = data.building
    form.roomNo = data.roomNo
    form.description = data.description
    form.imageUrls = data.images?.map((img: any) => img.url) || []
    form.facilityIds = data.facilities?.map((f: any) => f.id) || []
  } catch (error) {
    console.error('Failed to load property:', error)
  }
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过5MB')
    return false
  }
  return true
}

async function handleUpload(options: any) {
  try {
    const res = await uploadImage(options.file)
    form.imageUrls.push(res.data.url)
  } catch (error) {
    console.error('Upload failed:', error)
  }
}

function removeImage(index: number) {
  form.imageUrls.splice(index, 1)
}

async function handleSubmit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        if (isEdit.value) {
          await updateProperty(Number(route.params.id), form)
          ElMessage.success('更新成功')
        } else {
          await createProperty(form)
          ElMessage.success('创建成功')
        }
        router.push('/properties')
      } catch (error) {
        console.error('Submit failed:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.property-edit {
  padding: 0;
}

.image-upload {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.image-upload .image-item {
  width: 100px;
  height: 100px;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
  border: 1px solid #dcdfe6;
}

.image-upload .image-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-upload .image-item .delete-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 20px;
  height: 20px;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}
</style>
