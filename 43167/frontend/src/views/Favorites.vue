<template>
  <el-card>
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span>收藏夹</span>
        <el-button type="primary" @click="showGroup = true">新建分组</el-button>
      </div>
    </template>
    <el-row :gutter="16">
      <el-col :span="6" v-for="g in groups" :key="g.id">
        <el-card>
          <h4>{{ g.name }} ({{ g.brand }})</h4>
          <el-button size="small" type="danger" @click="delGroup(g.id)">删除</el-button>
        </el-card>
      </el-col>
    </el-row>
    <el-divider />
    <h3>我的收藏</h3>
    <el-table :data="favorites">
      <el-table-column prop="watch_id" label="手表ID" />
      <el-table-column prop="group_id" label="分组ID" />
    </el-table>

    <el-dialog v-model="showGroup" title="新建分组">
      <el-form :model="group" label-width="80px">
        <el-form-item label="品牌"><el-input v-model="group.brand" /></el-form-item>
        <el-form-item label="名称"><el-input v-model="group.name" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showGroup = false">取消</el-button>
        <el-button type="primary" @click="submitGroup">提交</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const groups = ref<any[]>([])
const favorites = ref<any[]>([])
const showGroup = ref(false)
const group = reactive({ brand: '', name: '' })

onMounted(async () => {
  groups.value = await request.get('/favorites/groups') || []
  const res: any = await request.get('/favorites')
  favorites.value = res?.list || res || []
})
async function submitGroup() {
  await request.post('/favorites/groups', group)
  ElMessage.success('已创建')
  showGroup.value = false
}
async function delGroup(id: number) {
  await request.delete(`/favorites/groups/${id}`)
  ElMessage.success('已删除')
}
</script>
