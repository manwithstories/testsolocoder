<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>顾客管理</span>
          <div class="header-actions">
            <el-input
              v-model="keyword"
              placeholder="搜索顾客名称/手机号"
              clearable
              style="width: 240px"
              :prefix-icon="Search"
              @keyup.enter="fetchList"
            />
            <el-button type="primary" :icon="Search" @click="fetchList">搜索</el-button>
          </div>
        </div>
      </template>

      <el-table :data="customers" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar :size="40" :src="row.user?.avatar">
              {{ row.name?.[0] || row.user?.nickname?.[0] || 'U' }}
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="姓名" />
        <el-table-column label="手机号">
          <template #default="{ row }">
            {{ row.user?.phone }}
          </template>
        </el-table-column>
        <el-table-column prop="gender" label="性别" width="80" />
        <el-table-column prop="age" label="年龄" width="80" />
        <el-table-column prop="skin_type" label="皮肤类型" />
        <el-table-column label="会员等级">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.member_level)">{{ getLevelText(row.member_level) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="points" label="积分" width="100" />
        <el-table-column prop="total_spent" label="消费金额" width="120">
          <template #default="{ row }">
            ¥{{ row.total_spent?.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="visit_count" label="到店次数" width="100" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/customers/${row.id}`)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCustomers } from '@/api/auth'
import { Search } from '@element-plus/icons-vue'
import type { Customer } from '@/types'

const loading = ref(false)
const keyword = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const customers = ref<Customer[]>([])

const getLevelType = (level: number) => {
  const types = ['', 'info', '', 'warning', 'success', 'danger']
  return types[level] || 'info'
}

const getLevelText = (level: number) => {
  const texts = ['', '普通会员', '', '银卡会员', '金卡会员', '钻石会员']
  return texts[level] || '普通会员'
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getCustomers({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value
    })
    customers.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
