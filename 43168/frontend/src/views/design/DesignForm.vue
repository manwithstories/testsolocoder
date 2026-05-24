<template>
  <div class="design-form">
    <el-page-header @back="router.back()">
      <template #content>
        <span>{{ isEdit ? '编辑方案' : '新建方案' }}</span>
      </template>
    </el-page-header>

    <el-card shadow="never" style="margin-top: 16px">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        style="max-width: 800px"
      >
        <el-form-item label="方案名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入方案名称" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item label="业主" prop="ownerId">
          <el-input v-model="form.ownerName" placeholder="请输入业主姓名" />
        </el-form-item>
        <el-form-item label="房型" prop="houseType">
          <el-select v-model="form.houseType" placeholder="请选择房型" clearable style="width: 100%">
            <el-option label="一居室" value="一居室" />
            <el-option label="两居室" value="两居室" />
            <el-option label="三居室" value="三居室" />
            <el-option label="四居室" value="四居室" />
            <el-option label="复式" value="复式" />
            <el-option label="别墅" value="别墅" />
          </el-select>
        </el-form-item>
        <el-form-item label="面积(㎡)" prop="area">
          <el-input-number v-model="form.area" :min="0" :precision="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="预算(元)" prop="budget">
          <el-input-number v-model="form.budget" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="封面图" prop="coverImage">
          <el-upload
            class="cover-uploader"
            :show-file-list="false"
            :auto-upload="false"
            :on-change="handleCoverChange"
            accept="image/*"
          >
            <el-image v-if="form.coverImage" :src="form.coverImage" fit="cover" class="cover-image" />
            <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item label="方案图片">
          <el-upload
            list-type="picture-card"
            :auto-upload="false"
            :file-list="imageFileList"
            :on-change="handleImageChange"
            :on-remove="handleImageRemove"
            accept="image/*"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="4" placeholder="请输入方案描述" maxlength="500" show-word-limit />
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
import type { FormInstance, FormRules, UploadFile } from 'element-plus'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getDesign, createDesign, updateDesign, type DesignFormData } from '@/api/design'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const imageFileList = ref<UploadFile[]>([])

const isEdit = computed(() => !!route.params.id)

const form = reactive<DesignFormData & { ownerName?: string }>({
  name: '',
  description: '',
  ownerName: '',
  houseType: '',
  area: 0,
  budget: 0,
  coverImage: '',
  images: [],
  status: 0
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入方案名称', trigger: 'blur' }],
  description: [{ max: 500, message: '描述不超过 500 字', trigger: 'blur' }]
}

function handleCoverChange(file: UploadFile) {
  const reader = new FileReader()
  reader.onload = (e) => {
    form.coverImage = e.target?.result as string
  }
  if (file.raw) reader.readAsDataURL(file.raw)
}

function handleImageChange(file: UploadFile) {
  const reader = new FileReader()
  reader.onload = (e) => {
    const url = e.target?.result as string
    if (!form.images) form.images = []
    form.images.push(url)
    imageFileList.value.push({
      name: file.name,
      url,
      uid: file.uid
    })
  }
  if (file.raw) reader.readAsDataURL(file.raw)
}

function handleImageRemove(file: UploadFile) {
  const idx = imageFileList.value.findIndex((f) => f.uid === file.uid)
  if (idx > -1) {
    imageFileList.value.splice(idx, 1)
    form.images?.splice(idx, 1)
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (isEdit.value) {
        await updateDesign(route.params.id as string, form)
        ElMessage.success('更新成功')
      } else {
        await createDesign(form)
        ElMessage.success('创建成功')
      }
      router.back()
    } finally {
      submitting.value = false
    }
  })
}

async function fetchDetail() {
  if (!isEdit.value) return
  const data = await getDesign(route.params.id as string)
  Object.assign(form, data)
  if (data.images?.length) {
    imageFileList.value = data.images.map((img, idx) => ({
      name: img.name || `image-${idx}`,
      url: img.url,
      uid: img.id
    }))
  }
}

onMounted(fetchDetail)
</script>

<style lang="scss" scoped>
.design-form {
  :deep(.cover-uploader) {
    :deep(.el-upload) {
      border: 1px dashed #d9d9d9;
      border-radius: 6px;
      cursor: pointer;
      position: relative;
      overflow: hidden;
      width: 120px;
      height: 120px;
      display: flex;
      align-items: center;
      justify-content: center;
      &:hover {
        border-color: #409eff;
      }
    }
  }
  .cover-image {
    width: 120px;
    height: 120px;
  }
  .cover-uploader-icon {
    font-size: 28px;
    color: #8c939d;
  }
}
</style>
