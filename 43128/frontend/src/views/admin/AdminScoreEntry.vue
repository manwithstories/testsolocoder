<template>
  <div class="page">
    <div class="page-title">成绩录入</div>

    <div class="card">
      <h3>手动录入</h3>
      <el-form :model="form" label-width="100px">
        <el-form-item label="赛事ID">
          <el-input-number v-model="form.event_id" :min="1" />
        </el-form-item>
        <el-form-item label="项目ID">
          <el-input-number v-model="form.event_item_id" :min="1" />
        </el-form-item>
        <el-form-item label="成绩列表">
          <div v-for="(s, idx) in form.scores" :key="idx" style="display:flex;gap:8px;margin-bottom:8px">
            <el-input v-model="s.user_id" placeholder="用户ID" style="width:120px" />
            <el-input-number v-model="s.score" :min="0" :precision="2" placeholder="成绩" style="width:140px" />
            <el-input v-model="s.time_used" placeholder="用时(可选)" style="width:140px" />
            <el-input v-model="s.remarks" placeholder="备注(可选)" style="width:200px" />
            <el-button link type="danger" @click="form.scores.splice(idx, 1)">删除</el-button>
          </div>
          <el-button plain @click="form.scores.push({ user_id: '', score: 0, time_used: '', remarks: '' })">+ 添加成绩</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit">提交录入</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="card">
      <h3>批量导入 Excel</h3>
      <p class="muted">Excel 格式：用户ID | 成绩 | 用时 | 备注（首行为表头）</p>
      <el-button link type="primary" :href="scoreApi.template">下载模板</el-button>
      <el-upload
        :action="`${location.origin}/api/v1/admin/scores/import`"
        :headers="uploadHeaders"
        :data="{ event_id: form.event_id, event_item_id: form.event_item_id }"
        :show-file-list="false"
        :on-success="onImportSuccess"
        accept=".xlsx,.xls"
      >
        <el-button>选择 Excel 文件上传</el-button>
      </el-upload>
    </div>

    <div class="card">
      <h3>查看项目成绩</h3>
      <el-input v-model="queryItemId" placeholder="输入项目ID" style="width:160px" />
      <el-button style="margin-left:8px" @click="fetchItemScores">查询</el-button>
      <el-table v-if="itemScores.length" :data="itemScores" stripe style="margin-top:12px">
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="score" label="成绩" />
        <el-table-column prop="rank" label="排名" width="80" />
        <el-table-column prop="points" label="积分" width="80" />
        <el-table-column prop="time_used" label="用时" />
        <el-table-column prop="remarks" label="备注" />
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { scoreApi } from '@/api'

const form = reactive({
  event_id: 1,
  event_item_id: 1,
  scores: [{ user_id: '', score: 0, time_used: '', remarks: '' }],
})

const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token') || ''}`,
}))

const queryItemId = ref<number | ''>('')
const itemScores = ref<any[]>([])

async function submit() {
  const scores = form.scores.filter((s: any) => s.user_id && s.score > 0)
  if (scores.length === 0) {
    ElMessage.warning('请至少填写一条有效成绩')
    return
  }
  await scoreApi.entry({
    event_id: form.event_id,
    event_item_id: form.event_item_id,
    scores: scores.map((s: any) => ({ ...s, user_id: Number(s.user_id) })),
  })
  ElMessage.success('成绩录入成功')
  form.scores = [{ user_id: '', score: 0, time_used: '', remarks: '' }]
}

function onImportSuccess(resp: any) {
  if (resp.code === 0) {
    ElMessage.success(`导入成功，共 ${resp.data?.imported || 0} 条`)
  } else {
    ElMessage.error(resp.message || '导入失败')
  }
}

async function fetchItemScores() {
  if (!queryItemId.value) return
  const res = await scoreApi.listByItem(Number(queryItemId.value))
  itemScores.value = (res.data as any[]) || []
}
</script>
