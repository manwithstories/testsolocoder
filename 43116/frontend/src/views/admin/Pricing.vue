<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">定价规则</h1>
      <el-button type="primary" :icon="Plus" @click="showAddDialog = true">添加规则</el-button>
    </div>

    <div class="search-bar">
      <el-select v-model="filters.ruleType" placeholder="规则类型" clearable style="width: 140px">
        <el-option label="周末" value="weekend" />
        <el-option label="节假日" value="holiday" />
        <el-option label="特殊" value="special" />
      </el-select>
      <el-button type="primary" @click="loadRules">搜索</el-button>
    </div>

    <el-table :data="rules" v-loading="loading" style="width: 100%">
      <el-table-column prop="name" label="规则名称" min-width="150" />
      <el-table-column prop="rule_type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag size="small">{{ getRuleTypeText(row.rule_type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="时间范围" width="240">
        <template #default="{ row }">
          <span v-if="row.start_date && row.end_date">
            {{ formatDate(row.start_date) }} - {{ formatDate(row.end_date) }}
          </span>
          <span v-else>{{ row.weekdays || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="倍率" width="100">
        <template #default="{ row }">
          <el-tag :type="row.multiplier > 1 ? 'danger' : 'success'" size="small">
            {{ row.multiplier }}x
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="priority" label="优先级" width="80" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-switch
            :model-value="row.is_active"
            @change="() => toggleActive(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="editRule(row)">编辑</el-button>
          <el-button type="danger" link size="small" @click="deleteRule(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="loadRules"
      />
    </div>

    <el-dialog v-model="showAddDialog" :title="editingRule ? '编辑规则' : '添加规则'" width="500px">
      <el-form :model="ruleForm" label-width="100px">
        <el-form-item label="规则名称">
          <el-input v-model="ruleForm.name" />
        </el-form-item>
        <el-form-item label="规则类型">
          <el-select v-model="ruleForm.rule_type">
            <el-option label="周末" value="weekend" />
            <el-option label="节假日" value="holiday" />
            <el-option label="特殊" value="special" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="ruleForm.rule_type === 'holiday'" label="开始时间">
          <el-date-picker v-model="ruleForm.start_date" type="date" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="ruleForm.rule_type === 'holiday'" label="结束时间">
          <el-date-picker v-model="ruleForm.end_date" type="date" style="width: 100%" />
        </el-form-item>
        <el-form-item label="倍率">
          <el-input-number v-model="ruleForm.multiplier" :min="0.1" :step="0.1" :precision="1" />
        </el-form-item>
        <el-form-item label="优先级">
          <el-input-number v-model="ruleForm.priority" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveRule">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { pricingApi } from '@/api'
import type { PricingRule } from '@/types'

const rules = ref<PricingRule[]>([])
const loading = ref(false)
const showAddDialog = ref(false)
const editingRule = ref<PricingRule | null>(null)

const filters = reactive({
  ruleType: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const ruleForm = reactive({
  name: '',
  rule_type: 'weekend',
  start_date: undefined as Date | undefined,
  end_date: undefined as Date | undefined,
  multiplier: 1.0,
  priority: 0
})

onMounted(() => {
  loadRules()
})

const loadRules = async () => {
  loading.value = true
  try {
    const res = await pricingApi.getRules({
      page: pagination.page,
      page_size: pagination.pageSize,
      type: filters.ruleType
    })
    rules.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const editRule = (row: PricingRule) => {
  editingRule.value = row
  Object.assign(ruleForm, {
    name: row.name,
    rule_type: row.rule_type,
    start_date: row.start_date ? new Date(row.start_date) : undefined,
    end_date: row.end_date ? new Date(row.end_date) : undefined,
    multiplier: row.multiplier,
    priority: row.priority
  })
  showAddDialog.value = true
}

const saveRule = async () => {
  try {
    const data = {
      ...ruleForm,
      start_date: ruleForm.start_date?.toISOString(),
      end_date: ruleForm.end_date?.toISOString()
    }
    if (editingRule.value) {
      await pricingApi.updateRule(editingRule.value.id, data)
      ElMessage.success('更新成功')
    } else {
      await pricingApi.createRule(data)
      ElMessage.success('添加成功')
    }
    showAddDialog.value = false
    editingRule.value = null
    loadRules()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const deleteRule = async (row: PricingRule) => {
  try {
    await ElMessageBox.confirm(`确定要删除规则 ${row.name} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await pricingApi.deleteRule(row.id)
    ElMessage.success('删除成功')
    loadRules()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '删除失败')
    }
  }
}

const toggleActive = async (row: PricingRule) => {
  try {
    await pricingApi.toggleRuleActive(row.id)
    ElMessage.success('操作成功')
    loadRules()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const getRuleTypeText = (type: string) => {
  const map: Record<string, string> = {
    weekend: '周末',
    holiday: '节假日',
    special: '特殊'
  }
  return map[type] || type
}
</script>
