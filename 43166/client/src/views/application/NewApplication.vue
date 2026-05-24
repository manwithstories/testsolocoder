<template>
  <div class="new-application">
    <div class="page-header flex-between">
      <h2 class="page-title">新建公司注册申请</h2>
      <el-button @click="goBack">返回</el-button>
    </div>

    <div class="card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        @submit.prevent
      >
        <div class="form-section">
          <div class="form-section-title">基本信息</div>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="公司名称" prop="companyName">
                <el-input v-model="form.companyName" placeholder="请输入公司名称" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="公司类型" prop="companyType">
                <el-select v-model="form.companyType" placeholder="请选择公司类型" style="width: 100%">
                  <el-option label="有限责任公司" value="llc" />
                  <el-option label="股份有限公司" value="joint_stock" />
                  <el-option label="个人独资" value="sole" />
                  <el-option label="合伙企业" value="partnership" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="注册资本" prop="registeredCapital">
                <el-input-number
                  v-model="form.registeredCapital"
                  :min="1"
                  :step="10000"
                  style="width: 100%"
                  placeholder="请输入注册资本"
                />
                <span style="margin-left: 8px; color: #909399">元</span>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="注册地址" prop="registeredAddress">
                <el-input v-model="form.registeredAddress" placeholder="请输入注册地址" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <div class="form-section">
          <div class="form-section-title">经营范围</div>
          <el-form-item prop="businessScope">
            <el-input
              v-model="form.businessScope"
              type="textarea"
              :rows="3"
              placeholder="请输入经营范围"
            />
          </el-form-item>
        </div>

        <div class="form-section">
          <div class="form-section-title">股东信息</div>
          <el-form-item prop="shareholderInfo">
            <el-input
              v-model="form.shareholderInfo"
              type="textarea"
              :rows="3"
              placeholder="请输入股东信息，格式：姓名、持股比例、身份证号"
            />
          </el-form-item>
        </div>

        <div class="form-section">
          <div class="form-section-title">材料上传</div>
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="身份证正面">
                <el-upload
                  class="uploader"
                  :show-file-list="false"
                  :before-upload="(file: any) => handleBeforeUpload(file, 'idCardFront')"
                  :http-request="(options: any) => handleUpload(options, 'idCardFront')"
                  accept="image/*,.pdf"
                >
                  <img v-if="form.idCardFront" :src="getImageUrl(form.idCardFront)" class="preview-image" />
                  <div v-else class="uploader-placeholder">
                    <el-icon :size="28"><Plus /></el-icon>
                    <div class="uploader-text">点击上传</div>
                  </div>
                </el-upload>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="身份证反面">
                <el-upload
                  class="uploader"
                  :show-file-list="false"
                  :before-upload="(file: any) => handleBeforeUpload(file, 'idCardBack')"
                  :http-request="(options: any) => handleUpload(options, 'idCardBack')"
                  accept="image/*,.pdf"
                >
                  <img v-if="form.idCardBack" :src="getImageUrl(form.idCardBack)" class="preview-image" />
                  <div v-else class="uploader-placeholder">
                    <el-icon :size="28"><Plus /></el-icon>
                    <div class="uploader-text">点击上传</div>
                  </div>
                </el-upload>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="营业执照预审">
                <el-upload
                  class="uploader"
                  :show-file-list="false"
                  :before-upload="(file: any) => handleBeforeUpload(file, 'licensePreview')"
                  :http-request="(options: any) => handleUpload(options, 'licensePreview')"
                  accept="image/*,.pdf"
                >
                  <img v-if="form.licensePreview" :src="getImageUrl(form.licensePreview)" class="preview-image" />
                  <div v-else class="uploader-placeholder">
                    <el-icon :size="28"><Plus /></el-icon>
                    <div class="uploader-text">点击上传</div>
                  </div>
                </el-upload>
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="handleSubmit">
            保存草稿
          </el-button>
          <el-button size="large" @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadRequestOptions } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { applicationApi } from '@/api/application'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  companyName: '',
  companyType: 'llc',
  registeredCapital: 100000,
  businessScope: '',
  registeredAddress: '',
  shareholderInfo: '',
  idCardFront: '',
  idCardBack: '',
  licensePreview: ''
})

const rules: FormRules = {
  companyName: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
  companyType: [{ required: true, message: '请选择公司类型', trigger: 'change' }],
  registeredCapital: [{ required: true, message: '请输入注册资本', trigger: 'blur' }],
  registeredAddress: [{ required: true, message: '请输入注册地址', trigger: 'blur' }],
  businessScope: [{ required: true, message: '请输入经营范围', trigger: 'blur' }],
  shareholderInfo: [{ required: true, message: '请输入股东信息', trigger: 'blur' }]
}

const handleBeforeUpload = (file: File, field: string) => {
  const isImage = file.type.startsWith('image/') || file.type === 'application/pdf'
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isImage) {
    ElMessage.error('只能上传图片或PDF文件!')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过10MB!')
    return false
  }
  return true
}

const handleUpload = async (options: UploadRequestOptions, field: string) => {
  try {
    const res = await applicationApi.uploadMaterials(0, field, options.file as File)
    if (res) {
      ;(form as any)[field] = res.filePath
      ElMessage.success('上传成功')
    }
  } catch (error) {
    console.error('上传失败:', error)
    options.onError!(error as any)
  }
}

const getImageUrl = (path: string) => {
  if (!path) return ''
  return `/uploads/${path}`
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await applicationApi.create({
          companyName: form.companyName,
          companyType: form.companyType as any,
          registeredCapital: form.registeredCapital,
          businessScope: form.businessScope,
          registeredAddress: form.registeredAddress,
          shareholderInfo: form.shareholderInfo
        })

        ElMessage.success('创建成功')
        router.push(`/applications/${res?.id}`)
      } catch (error: any) {
        ElMessage.error(error.message || '创建失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const goBack = () => {
  router.back()
}
</script>

<style scoped>
.uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  width: 150px;
  height: 150px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  transition: border-color 0.3s;
}

.uploader:hover {
  border-color: #409eff;
}

.uploader-placeholder {
  text-align: center;
  color: #8c939d;
}

.uploader-text {
  font-size: 12px;
  margin-top: 4px;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
