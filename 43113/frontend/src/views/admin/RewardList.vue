<template>
  <div class="reward-list-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>奖品管理</span>
          <el-button type="primary" @click="showAddDialog = true">
            <el-icon><Plus /></el-icon>
            添加奖品
          </el-button>
        </div>
      </template>

      <el-table :data="rewards" style="width: 100%">
        <el-table-column label="ID" prop="id" width="80" />
        <el-table-column label="图标" width="80">
          <template #default="{ row }">
            <span style="font-size: 24px">{{ row.image }}</span>
          </template>
        </el-table-column>
        <el-table-column label="名称" prop="name" />
        <el-table-column label="描述" prop="description" show-overflow-tooltip />
        <el-table-column label="积分消耗" prop="pointsCost" width="120" />
        <el-table-column label="库存" width="100">
          <template #default="{ row }">
            {{ row.stock === -1 ? '无限' : row.stock }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '上架' : '下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button
              :type="row.status === 'active' ? 'warning' : 'success'"
              size="small"
              @click="toggleStatus(row)"
            >
              {{ row.status === 'active' ? '下架' : '上架' }}
            </el-button>
            <el-button type="danger" size="small" @click="deleteReward(row.id)">
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
        @current-change="fetchRewards"
      />
    </el-card>

    <el-dialog v-model="showAddDialog" title="添加奖品" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="form.name" placeholder="请输入奖品名称" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.image" placeholder="请输入图标(emoji)" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述"
          />
        </el-form-item>
        <el-form-item label="积分消耗">
          <el-input-number v-model="form.pointsCost" :min="0" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="form.stock" :min="-1" />
          <div class="tip">-1 表示无限库存</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="addReward">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { rewardApi } from '@/api'
import type { Reward } from '@/types'

const rewards = ref<Reward[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const showAddDialog = ref(false)

const form = reactive({
  name: '',
  image: '',
  description: '',
  pointsCost: 0,
  stock: -1
})

const fetchRewards = async () => {
  try {
    const res = await rewardApi.getRewardList({
      page: page.value,
      pageSize: pageSize.value
    })
    rewards.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  }
}

const addReward = async () => {
  if (!form.name) return
  try {
    await rewardApi.createReward(form)
    showAddDialog.value = false
    form.name = ''
    form.image = ''
    form.description = ''
    form.pointsCost = 0
    form.stock = -1
    fetchRewards()
  } catch (e) {
    console.error(e)
  }
}

const toggleStatus = async (row: Reward) => {
  try {
    await rewardApi.updateReward(row.id, {
      status: row.status === 'active' ? 'deleted' : 'active'
    })
    fetchRewards()
  } catch (e) {
    console.error(e)
  }
}

const deleteReward = async (id: number) => {
  try {
    await rewardApi.deleteReward(id)
    fetchRewards()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  fetchRewards()
})
</script>

<style scoped lang="scss">
.reward-list-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .tip {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
