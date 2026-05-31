<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>术语库</span>
          <el-button type="primary" @click="showAddDialog = true">
            <el-icon><Plus /></el-icon>添加术语
          </el-button>
        </div>
      </template>

      <div class="filter-bar">
        <el-input v-model="keyword" placeholder="搜索术语" clearable style="width: 200px" @keyup.enter="loadData">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="domain" placeholder="领域" clearable style="width: 140px">
          <el-option label="通用" value="general" />
          <el-option label="技术" value="technology" />
          <el-option label="法律" value="legal" />
          <el-option label="医学" value="medical" />
          <el-option label="金融" value="finance" />
        </el-select>
        <el-button type="primary" @click="loadData">搜索</el-button>
      </div>

      <el-table :data="terms" stripe v-loading="loading">
        <el-table-column prop="source_term" label="源术语" width="200" />
        <el-table-column prop="target_term" label="目标术语" width="200" />
        <el-table-column prop="source_lang" label="源语言" width="100" />
        <el-table-column prop="target_lang" label="目标语言" width="100" />
        <el-table-column prop="domain" label="领域" width="100" />
        <el-table-column prop="part_of_speech" label="词性" width="100" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showAddDialog" :title="editing ? '编辑术语' : '添加术语'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="源术语">
              <el-input v-model="form.source_term" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="目标术语">
              <el-input v-model="form.target_term" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="源语言">
              <el-input v-model="form.source_lang" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="目标语言">
              <el-input v-model="form.target_lang" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="领域">
              <el-select v-model="form.domain" placeholder="选择领域" style="width: 100%">
                <el-option label="通用" value="general" />
                <el-option label="技术" value="technology" />
                <el-option label="法律" value="legal" />
                <el-option label="医学" value="medical" />
                <el-option label="金融" value="finance" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="词性">
              <el-input v-model="form.part_of_speech" placeholder="如: n., v., adj." />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="定义">
          <el-input v-model="form.definition" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listGlossaryTerms, createGlossaryTerm, updateGlossaryTerm, deleteGlossaryTerm } from '@/api/memory'

const terms = ref<any[]>([])
const loading = ref(false)
const keyword = ref('')
const domain = ref('')

const showAddDialog = ref(false)
const editing = ref(false)
const form = reactive({
  id: 0,
  source_term: '',
  target_term: '',
  source_lang: '',
  target_lang: '',
  domain: '',
  definition: '',
  part_of_speech: ''
})

async function loadData() {
  loading.value = true
  try {
    const res = await listGlossaryTerms({
      keyword: keyword.value,
      domain: domain.value
    }) as any
    if (Array.isArray(res)) {
      terms.value = res
    } else {
      terms.value = res?.list || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handleEdit(row: any) {
  editing.value = true
  Object.assign(form, row)
  showAddDialog.value = true
}

async function handleSave() {
  try {
    if (editing.value) {
      const { id, ...data } = form
      await updateGlossaryTerm(id, data)
      ElMessage.success('更新成功')
    } else {
      await createGlossaryTerm(form)
      ElMessage.success('添加成功')
    }
    showAddDialog.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  }
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm('确定删除该术语？', '提示', { type: 'warning' })
    await deleteGlossaryTerm(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (_) {}
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
  }
}
</style>
