<template>
  <div class="my-items">
    <div class="header">
      <h3>我的拍卖品</h3>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        发布拍卖品
      </el-button>
    </div>

    <el-table :data="items" v-loading="loading">
      <el-table-column prop="title" label="标题" min-width="200">
        <template #default="{ row }">
          <div class="item-title" @click="$router.push(`/items/${row.id}`)">
            <img :src="getMainImage(row)" :alt="row.title" />
            <span>{{ row.title }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="current_price" label="当前价" width="120">
        <template #default="{ row }">¥{{ row.current_price.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="bid_count" label="出价次数" width="100" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="editItem(row)">编辑</el-button>
          <el-button
            v-if="row.status === 0"
            link
            type="success"
            size="small"
            @click="onlineItem(row)"
          >上架</el-button>
          <el-button
            v-if="row.status === 1"
            link
            type="warning"
            size="small"
            @click="offlineItem(row)"
          >下架</el-button>
          <el-button
            v-if="row.status === 0"
            link
            type="danger"
            size="small"
            @click="deleteItem(row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchItems"
    />

    <el-dialog v-model="showCreateDialog" :title="editingItem ? '编辑拍卖品' : '发布拍卖品'" width="600px">
      <el-form :model="itemForm" label-width="100px">
        <el-form-item label="标题">
          <el-input v-model="itemForm.title" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="itemForm.category_id" style="width: 100%">
            <el-option label="艺术品" :value="1" />
            <el-option label="收藏品" :value="2" />
            <el-option label="奢侈品" :value="3" />
            <el-option label="其他" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="起拍价">
          <el-input-number v-model="itemForm.start_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="保留价">
          <el-input-number v-model="itemForm.reserve_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="itemForm.description" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="所在地">
          <el-input v-model="itemForm.location" />
        </el-form-item>
        <el-form-item label="成色">
          <el-select v-model="itemForm.condition" style="width: 100%">
            <el-option label="全新" value="全新" />
            <el-option label="99新" value="99新" />
            <el-option label="95新" value="95新" />
            <el-option label="9成新" value="9成新" />
            <el-option label="8成新" value="8成新" />
            <el-option label="有瑕疵" value="有瑕疵" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="saveItem">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import type { AuctionItem } from '@/types'
import { itemApi } from '@/api'

const items = ref<AuctionItem[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const showCreateDialog = ref(false)
const editingItem = ref<AuctionItem | null>(null)

const itemForm = reactive({
  title: '',
  description: '',
  category_id: 1,
  start_price: 0,
  reserve_price: 0,
  location: '',
  condition: '全新',
})

const fetchItems = async () => {
  loading.value = true
  try {
    const res = await itemApi.getMyItems({ page: page.value, page_size: pageSize.value })
    items.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const getMainImage = (item: AuctionItem) => {
  if (item.images && item.images.length > 0) {
    return item.images[0].url
  }
  return 'https://via.placeholder.com/50x50?text=No'
}

const statusType = (status: number) => {
  const map: Record<number, string> = { 0: 'info', 1: 'success', 2: 'danger', 3: 'primary' }
  return map[status] || 'info'
}

const statusText = (status: number) => {
  const map: Record<number, string> = { 0: '草稿', 1: '拍卖中', 2: '已下架', 3: '已售出' }
  return map[status] || ''
}

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

const editItem = (item: AuctionItem) => {
  editingItem.value = item
  Object.assign(itemForm, {
    title: item.title,
    description: item.description,
    category_id: item.category_id,
    start_price: item.start_price,
    reserve_price: item.reserve_price,
    location: item.location,
    condition: item.condition,
  })
  showCreateDialog.value = true
}

const saveItem = async () => {
  try {
    if (editingItem.value) {
      await itemApi.update(editingItem.value.id, itemForm)
      ElMessage.success('更新成功')
    } else {
      await itemApi.create(itemForm)
      ElMessage.success('创建成功')
    }
    showCreateDialog.value = false
    fetchItems()
  } catch (e) {}
}

const onlineItem = async (item: AuctionItem) => {
  try {
    await itemApi.online(item.id)
    ElMessage.success('上架成功')
    fetchItems()
  } catch (e) {}
}

const offlineItem = async (item: AuctionItem) => {
  try {
    await itemApi.offline(item.id)
    ElMessage.success('下架成功')
    fetchItems()
  } catch (e) {}
}

const deleteItem = async (item: AuctionItem) => {
  try {
    await ElMessageBox.confirm('确定要删除这个拍卖品吗？', '提示', { type: 'warning' })
    await itemApi.delete(item.id)
    ElMessage.success('删除成功')
    fetchItems()
  } catch (e) {}
}

onMounted(() => {
  fetchItems()
})
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h3 {
  margin: 0;
}

.item-title {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.item-title img {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 4px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
