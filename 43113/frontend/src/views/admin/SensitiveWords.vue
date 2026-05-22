<template>
  <div class="sensitive-words-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>敏感词管理</span>
          <el-button type="primary" @click="showAddDialog = true">
            <el-icon><Plus /></el-icon>
            添加敏感词
          </el-button>
        </div>
      </template>

      <el-table :data="words" style="width: 100%">
        <el-table-column label="ID" prop="id" width="80" />
        <el-table-column label="敏感词" prop="word" />
        <el-table-column label="分类" prop="category" width="120" />
        <el-table-column label="替换为" prop="replaceTo" width="150" />
        <el-table-column label="级别" prop="level" width="80" />
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="danger" size="small" @click="deleteWord(row.id)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        @current-change="fetchWords"
      />
    </el-card>

    <el-dialog v-model="showAddDialog" title="添加敏感词" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="敏感词">
          <el-input v-model="form.word" placeholder="请输入敏感词" />
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="form.category" placeholder="请输入分类" />
        </el-form-item>
        <el-form-item label="替换为">
          <el-input v-model="form.replaceTo" placeholder="请输入替换文本" />
        </el-form-item>
        <el-form-item label="级别">
          <el-input-number v-model="form.level" :min="1" :max="10" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="addWord">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { auditApi } from '@/api'
import type { SensitiveWord } from '@/types'
import dayjs from 'dayjs'

const words = ref<SensitiveWord[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const showAddDialog = ref(false)

const form = reactive({
  word: '',
  category: '',
  replaceTo: '',
  level: 1
})

const fetchWords = async () => {
  try {
    const res = await auditApi.getSensitiveWords({
      page: page.value,
      pageSize: pageSize.value
    })
    words.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  }
}

const addWord = async () => {
  if (!form.word) return
  try {
    await auditApi.createSensitiveWord(form.word, form.category, form.replaceTo, form.level)
    showAddDialog.value = false
    form.word = ''
    form.category = ''
    form.replaceTo = ''
    form.level = 1
    fetchWords()
  } catch (e) {
    console.error(e)
  }
}

const deleteWord = async (id: number) => {
  try {
    await auditApi.deleteSensitiveWord(id)
    fetchWords()
  } catch (e) {
    console.error(e)
  }
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchWords()
})
</script>

<style scoped lang="scss">
.sensitive-words-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
