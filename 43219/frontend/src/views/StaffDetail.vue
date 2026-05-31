<template>
  <AppLayout>
    <div class="page">
      <el-descriptions v-if="staff" title="家政人员档案" :column="2" bordered>
        <el-descriptions-item label="姓名">{{ staff.real_name || staff.username }}</el-descriptions-item>
        <el-descriptions-item label="评分">
          <el-rate :model-value="staff.rating" disabled /> ({{ staff.rating }})
        </el-descriptions-item>
        <el-descriptions-item label="等级">Lv.{{ staff.level }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="staff.suspended ? 'danger' : 'success'">
            {{ staff.suspended ? '已暂停派单' : '正常' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="技能">{{ staff.skills || '-' }}</el-descriptions-item>
        <el-descriptions-item label="简介">{{ profile?.intro || '-' }}</el-descriptions-item>
      </el-descriptions>
      <div class="card" style="margin-top:16px">
        <h3>客户评价 ({{ reviews.length }})</h3>
        <div v-for="r in reviews" :key="r.id" class="review">
          <div class="row">
            <el-rate :model-value="r.rating" disabled />
            <span class="muted">{{ r.created_at }}</span>
          </div>
          <p>{{ r.content }}</p>
        </div>
        <el-empty v-if="!reviews.length" description="暂无评价" />
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '../components/AppLayout.vue'
import { getStaffDetail } from '../api/user'

const route = useRoute()
const staff = ref<any>(null)
const profile = ref<any>(null)
const reviews = ref<any[]>([])

onMounted(async () => {
  const res = await getStaffDetail(Number(route.params.id))
  const d = (res.data as any).data
  staff.value = d.user
  profile.value = d.profile
  reviews.value = d.reviews || []
})
</script>

<style scoped>
.review {
  border-bottom: 1px solid #eee;
  padding: 12px 0;
}
</style>
