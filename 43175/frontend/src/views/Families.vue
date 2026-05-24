<template>
  <div class="page-container">
    <div class="page-header">
      <h2>家庭管理</h2>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
        创建家庭
      </el-button>
    </div>

    <div class="card-grid">
      <div v-for="family in families" :key="family.id" class="device-card">
        <div class="device-header">
          <span class="device-name">{{ family.name }}</span>
          <el-tag v-if="family.ownerId === currentUserId" type="warning" size="small">管理员</el-tag>
        </div>
        <div class="device-info">
          <p v-if="family.description">{{ family.description }}</p>
          <p>成员数: {{ family.members?.length || 0 }}</p>
        </div>
        <div class="device-actions">
          <el-button size="small" type="primary" @click="viewMembers(family)">成员</el-button>
          <el-button size="small" type="success" @click="inviteMember(family)" v-if="family.ownerId === currentUserId">邀请</el-button>
          <el-button size="small" type="danger" @click="handleDeleteFamily(family)" v-if="family.ownerId === currentUserId">删除</el-button>
        </div>
      </div>
    </div>

    <el-empty v-if="families.length === 0 && !loading" description="暂无家庭，创建一个开始管理智能家居" />

    <el-dialog v-model="showCreateDialog" title="创建家庭" width="450px">
      <el-form :model="createForm" ref="createFormRef" label-width="80px">
        <el-form-item label="家庭名称">
          <el-input v-model="createForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="createForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveFamily">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showMembersDialog" :title="`${currentFamily?.name || ''} - 成员管理`" width="600px">
      <div v-if="currentFamily">
        <el-table :data="currentFamily.members || []" style="width: 100%">
          <el-table-column label="用户">
            <template #default="{ row }">
              <div class="flex-center">
                <el-avatar :size="32" style="margin-right: 8px;">
                  {{ row.user?.username?.charAt(0) || 'U' }}
                </el-avatar>
                <span>{{ row.user?.username }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="user.email" label="邮箱" />
          <el-table-column prop="role" label="角色" width="120">
            <template #default="{ row }">
              <el-tag :type="getRoleType(row.role)">{{ getRoleLabel(row.role) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" v-if="currentFamily.ownerId === currentUserId">
            <template #default="{ row }">
              <el-select v-model="row.role" size="small" style="width: 100px;" @change="updateRole(row)">
                <el-option label="管理员" value="admin" />
                <el-option label="普通用户" value="member" />
                <el-option label="访客" value="guest" />
              </el-select>
              <el-button type="danger" link size="small" @click="removeMember(row)" v-if="row.userId !== currentUserId">移除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>

    <el-dialog v-model="showInviteDialog" title="邀请成员" width="450px">
      <el-form :model="inviteForm" ref="inviteFormRef" label-width="80px">
        <el-form-item label="邮箱">
          <el-input v-model="inviteForm.email" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="inviteForm.role" style="width: 100%">
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="member" />
            <el-option label="访客" value="guest" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showInviteDialog = false">取消</el-button>
        <el-button type="primary" :loading="inviting" @click="sendInvitation">发送邀请</el-button>
      </template>
    </el-dialog>

    <div v-if="invitations.length > 0" class="mt-20">
      <h3>待处理的邀请</h3>
      <el-table :data="invitations" style="width: 100%">
        <el-table-column label="家庭">
          <template #default="{ row }">{{ row.family?.name }}</template>
        </el-table-column>
        <el-table-column prop="email" label="邀请邮箱" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)">{{ getRoleLabel(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expiresAt" label="有效期至" width="180" />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button type="success" size="small" @click="acceptInvitation(row)">接受</el-button>
            <el-button type="danger" size="small" @click="rejectInvitation(row)">拒绝</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  listFamilies, createFamily, deleteFamily, inviteMember as apiInvite,
  removeMember as apiRemove, updateMemberRole as apiUpdateRole,
  listInvitations, acceptInvitation as apiAccept, rejectInvitation as apiReject,
  type Family
} from '@/api/family'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const currentUserId = computed(() => userStore.user?.id || 0)

const loading = ref(false)
const saving = ref(false)
const inviting = ref(false)
const families = ref<Family[]>([])
const invitations = ref<any[]>([])
const showCreateDialog = ref(false)
const showMembersDialog = ref(false)
const showInviteDialog = ref(false)
const currentFamily = ref<Family | null>(null)
const createFormRef = ref<FormInstance>()
const inviteFormRef = ref<FormInstance>()

const createForm = reactive({ name: '', description: '' })
const inviteForm = reactive({ email: '', role: 'member' })

onMounted(async () => {
  await loadFamilies()
  await loadInvitations()
})

async function loadFamilies() {
  loading.value = true
  try {
    families.value = await listFamilies()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadInvitations() {
  try {
    invitations.value = await listInvitations()
  } catch (e) {
    console.error(e)
  }
}

async function saveFamily() {
  saving.value = true
  try {
    await createFamily(createForm)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    Object.assign(createForm, { name: '', description: '' })
    loadFamilies()
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}

async function handleDeleteFamily(family: Family) {
  try {
    await ElMessageBox.confirm(`确定要删除家庭"${family.name}"吗？此操作不可撤销。`, '提示', { type: 'warning' })
    await deleteFamily(family.id)
    ElMessage.success('删除成功')
    loadFamilies()
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}

function viewMembers(family: Family) {
  currentFamily.value = family
  showMembersDialog.value = true
}

async function updateRole(member: any) {
  if (!currentFamily.value) return
  try {
    await apiUpdateRole(currentFamily.value.id, member.id, { role: member.role })
    ElMessage.success('角色已更新')
  } catch (e) {
    console.error(e)
  }
}

async function removeMember(member: any) {
  if (!currentFamily.value) return
  try {
    await ElMessageBox.confirm('确定要移除该成员吗？', '提示', { type: 'warning' })
    await apiRemove(currentFamily.value.id, member.id)
    ElMessage.success('移除成功')
    const idx = currentFamily.value.members?.findIndex(m => m.id === member.id)
    if (idx !== -1) currentFamily.value.members?.splice(idx, 1)
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}

function inviteMember(family: Family) {
  currentFamily.value = family
  Object.assign(inviteForm, { email: '', role: 'member' })
  showInviteDialog.value = true
}

async function sendInvitation() {
  if (!currentFamily.value) return
  inviting.value = true
  try {
    await apiInvite(currentFamily.value.id, inviteForm)
    ElMessage.success('邀请已发送')
    showInviteDialog.value = false
  } catch (e) {
    console.error(e)
  } finally {
    inviting.value = false
  }
}

async function acceptInvitation(inv: any) {
  try {
    await apiAccept(inv.id)
    ElMessage.success('已接受邀请')
    loadInvitations()
    loadFamilies()
  } catch (e) {
    console.error(e)
  }
}

async function rejectInvitation(inv: any) {
  try {
    await apiReject(inv.id)
    ElMessage.success('已拒绝邀请')
    loadInvitations()
  } catch (e) {
    console.error(e)
  }
}

function getRoleLabel(role: string) {
  const map: Record<string, string> = { admin: '管理员', member: '普通用户', guest: '访客' }
  return map[role] || role
}

function getRoleType(role: string) {
  const map: Record<string, string> = { admin: 'warning', member: '', guest: 'info' }
  return map[role] || ''
}
</script>
