<template>
  <div class="page-container">
    <div class="page-header">
      <h2>场景联动</h2>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
        创建场景
      </el-button>
    </div>

    <div class="card-grid">
      <div v-for="scene in scenes" :key="scene.id" class="device-card">
        <div class="device-header">
          <div>
            <span class="scene-icon">{{ scene.icon }}</span>
            <span class="device-name" style="margin-left: 8px;">{{ scene.name }}</span>
          </div>
          <el-switch v-model="scene.isActive" @change="toggleScene(scene)" />
        </div>
        <div class="device-info">
          <p v-if="scene.description">{{ scene.description }}</p>
          <p>条件数: {{ scene.conditions?.length || 0 }}</p>
          <p>动作数: {{ scene.actions?.length || 0 }}</p>
        </div>
        <div class="device-actions">
          <el-button size="small" type="primary" :icon="VideoPlay" @click="execute(scene)">执行</el-button>
          <el-button size="small" @click="editScene(scene)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDeleteScene(scene)">删除</el-button>
        </div>
      </div>
    </div>

    <el-empty v-if="scenes.length === 0 && !loading" description="暂无场景" />

    <el-dialog v-model="showCreateDialog" :title="editingScene ? '编辑场景' : '创建场景'" width="600px">
      <el-form :model="sceneForm" ref="sceneFormRef" label-width="80px">
        <el-form-item label="场景名称">
          <el-input v-model="sceneForm.name" />
        </el-form-item>
        <el-form-item label="图标">
          <el-select v-model="sceneForm.icon" style="width: 100%">
            <el-option v-for="icon in sceneIcons" :key="icon.value" :label="icon.label" :value="icon.value">
              <span>{{ icon.value }} {{ icon.label }}</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="sceneForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="触件">
          <div v-for="(cond, idx) in sceneForm.conditions" :key="idx" class="condition-item">
            <el-select v-model="cond.type" style="width: 120px;">
              <el-option v-for="opt in conditionTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>
            <el-input v-if="cond.type === 'time'" v-model="cond.timeExpr" placeholder="HH:mm" style="width: 120px;" />
            <el-select v-else v-model="cond.deviceId" placeholder="选择设备" style="width: 150px;">
              <el-option v-for="d in devices" :key="d.id" :label="d.name" :value="d.id" />
            </el-select>
            <el-select v-if="cond.type !== 'time'" v-model="cond.operator" style="width: 80px;">
              <el-option label="=" value="eq" />
              <el-option label="!=" value="neq" />
            </el-select>
            <el-input v-if="cond.type !== 'time'" v-model="cond.value" placeholder="值" style="width: 100px;" />
            <el-button type="danger" link :icon="Delete" @click="removeCondition(idx)" />
          </div>
          <el-button type="primary" link size="small" @click="addCondition">+ 添加条件</el-button>
        </el-form-item>
        <el-form-item label="执行动作">
          <div v-for="(act, idx) in sceneForm.actions" :key="idx" class="condition-item">
            <el-select v-model="act.deviceId" placeholder="选择设备" style="width: 150px;">
              <el-option v-for="d in devices" :key="d.id" :label="d.name" :value="d.id" />
            </el-select>
            <el-select v-model="act.action" style="width: 120px;">
              <el-option v-for="opt in actionOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>
            <el-input v-model="act.value" placeholder="值(可选)" style="width: 100px;" />
            <el-button type="danger" link :icon="Delete" @click="removeAction(idx)" />
          </div>
          <el-button type="primary" link size="small" @click="addAction">+ 添加动作</el-button>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveScene">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, VideoPlay, Delete } from '@element-plus/icons-vue'
import {
  listScenes, createScene, updateScene, deleteScene, executeScene,
  sceneIcons, conditionTypeOptions, actionOptions
} from '@/api/scene'
import { listDevices } from '@/api/device'
import { useFamilyStore } from '@/stores/family'

const familyStore = useFamilyStore()
const loading = ref(false)
const saving = ref(false)
const scenes = ref<any[]>([])
const devices = ref<any[]>([])
const showCreateDialog = ref(false)
const editingScene = ref<any | null>(null)
const sceneFormRef = ref<FormInstance>()

const sceneForm = reactive({
  name: '',
  icon: '🏠',
  description: '',
  isActive: true,
  conditions: [{ type: 'time', timeExpr: '08:00' }] as any[],
  actions: [] as any[]
})

onMounted(async () => {
  await loadScenes()
  devices.value = await listDevices()
})

watch(() => familyStore.currentFamilyId, () => {
  loadScenes()
})

async function loadScenes() {
  loading.value = true
  try {
    scenes.value = await listScenes({ familyId: familyStore.currentFamilyId || undefined })
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function addCondition() {
  sceneForm.conditions.push({ type: 'device', operator: 'eq' })
}

function removeCondition(idx: number) {
  sceneForm.conditions.splice(idx, 1)
}

function addAction() {
  sceneForm.actions.push({ action: 'on' })
}

function removeAction(idx: number) {
  sceneForm.actions.splice(idx, 1)
}

function editScene(scene: any) {
  editingScene.value = scene
  Object.assign(sceneForm, {
    name: scene.name,
    icon: scene.icon,
    description: scene.description,
    isActive: scene.isActive,
    conditions: JSON.parse(JSON.stringify(scene.conditions || [])),
    actions: JSON.parse(JSON.stringify(scene.actions || []))
  })
  showCreateDialog.value = true
}

async function saveScene() {
  saving.value = true
  try {
    if (!familyStore.currentFamilyId) {
      ElMessage.warning('请先选择或创建家庭')
      return
    }
    const data = { ...sceneForm, familyId: familyStore.currentFamilyId }
    if (editingScene.value) {
      await updateScene(editingScene.value.id, data)
      ElMessage.success('更新成功')
    } else {
      await createScene(data)
      ElMessage.success('创建成功')
    }
    showCreateDialog.value = false
    resetForm()
    loadScenes()
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}

async function handleDeleteScene(scene: any) {
  try {
    await ElMessageBox.confirm(`确定要删除场景"${scene.name}"吗？`, '提示', { type: 'warning' })
    await deleteScene(scene.id)
    ElMessage.success('删除成功')
    loadScenes()
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}

async function toggleScene(scene: any) {
  try {
    await updateScene(scene.id, { isActive: scene.isActive })
    ElMessage.success(scene.isActive ? '已启用' : '已禁用')
  } catch (e) {
    console.error(e)
  }
}

async function execute(scene: any) {
  try {
    await executeScene(scene.id)
    ElMessage.success('场景执行成功')
  } catch (e) {
    console.error(e)
  }
}

function resetForm() {
  editingScene.value = null
  Object.assign(sceneForm, {
    name: '',
    icon: '🏠',
    description: '',
    isActive: true,
    conditions: [{ type: 'time', timeExpr: '08:00' }],
    actions: []
  })
}
</script>

<style lang="scss" scoped>
.scene-icon {
  font-size: 20px;
}

.condition-item {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}
</style>
