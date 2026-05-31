<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>翻译记忆库</span>
          <el-button type="primary" @click="showAddDialog = true">
            <el-icon><Plus /></el-icon>添加记忆
          </el-button>
        </div>
      </template>

      <div class="filter-bar">
        <el-input
          v-model="keyword"
          placeholder="搜索原文或译文"
          clearable
          style="width: 200px"
          @keyup.enter="loadData"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-button type="primary" @click="loadData">搜索</el-button>
      </div>

      <el-table :data="memories" stripe v-loading="loading">
        <el-table-column prop="source_text" label="原文" min-width="200">
          <template #default="{ row }">{{ row.source_text }}</template>
        </el-table-column>
        <el-table-column prop="translated_text" label="译文" min-width="200">
          <template #default="{ row }">{{ row.translated_text }}</template>
        </el-table-column>
        <el-table-column prop="source_lang" label="源语言" width="100" />
        <el-table-column prop="target_lang" label="目标语言" width="100" />
        <el-table-column prop="usage_count" label="使用次数" width="100" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @change="loadData"
        />
      </div>
    </el-card>

    <el-dialog v-model="showAddDialog" :title="editing ? '编辑记忆' : '添加记忆'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="源语言">
          <el-input v-model="form.source_lang" />
        </el-form-item>
        <el-form-item label="目标语言">
          <el-input v-model="form.target_lang" />
        </el-form-item>
        <el-form-item label="原文">
          <el-input v-model="form.source_text" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="译文">
          <el-input v-model="form.translated_text" type="textarea" :rows="3" />
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
import { listTranslationMemories, createTranslationMemory, updateTranslationMemory, deleteTranslationMemory } from '@/api/memory'

const memories = ref<any[]>([])
const loading = ref(false)
const keyword = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const showAddDialog = ref(false)
const editing = ref(false)
const form = reactive({
  id: 0,
  source_text: '',
  translated_text: '',
  source_lang: '',
  target_lang: ''
})

async function loadData() {
  loading.value = true
  try {
    const res = await listTranslationMemories({
      keyword: keyword.value,
      page: page.value,
      page_size: pageSize.value
    }) as any
    if (Array.isArray(res)) {
      memories.value = res
    } else {
      memories.value = res?.list || []
      total.value = res?.total || 0
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
      await updateTranslationMemory(form.id, {
        source_text: form.source_text,
        translated_text: form.translated_text
      })
      ElMessage.success('更新成功')
    } else {
      await createTranslationMemory(form)
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
    await ElMessageBox.confirm('确定删除该翻译记忆？', '提示', { type: 'warning' })
    await deleteTranslationMemory(row.id)
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

  .pagination {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
