<template>
  <div class="page-container">
    <div class="page-header">
      <h2>设备详情</h2>
      <el-button :icon="ArrowLeft" @click="$router.back()">返回</el-button>
    </div>

    <el-row :gutter="20" v-if="device">
      <el-col :span="12">
        <div class="chart-container">
          <h3 class="mb-20">设备信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="设备名称">{{ device.name }}</el-descriptions-item>
            <el-descriptions-item label="设备类型">{{ getDeviceTypeLabel(device.deviceType) }}</el-descriptions-item>
            <el-descriptions-item label="厂商">{{ device.vendor || '-' }}</el-descriptions-item>
            <el-descriptions-item label="安装位置">{{ device.location || '-' }}</el-descriptions-item>
            <el-descriptions-item label="功率">{{ device.power }} W</el-descriptions-item>
            <el-descriptions-item label="通信协议">{{ getProtocolLabel(device.protocol) }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(device.status)">{{ getStatusLabel(device.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="最后在线">
              {{ device.lastOnlineTime || '-' }}
            </el-descriptions-item>
          </el-descriptions>

          <div class="mt-20 flex-between">
            <el-button type="primary" :icon="Edit" @click="showEdit = true">编辑设备</el-button>
            <el-button :type="device.status === 'on' || device.status === 'online' ? 'warning' : 'success'"
              @click="toggleDeviceStatus">
              {{ device.status === 'on' || device.status === 'online' ? '关闭设备' : '开启设备' }}
            </el-button>
          </div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="chart-container">
          <div class="flex-between mb-20">
            <h3>能耗趋势</h3>
            <el-radio-group v-model="energyPeriod" size="small" @change="loadEnergy">
              <el-radio-button value="hour">近1小时</el-radio-button>
              <el-radio-button value="day">近1天</el-radio-button>
              <el-radio-button value="week">近7天</el-radio-button>
              <el-radio-button value="month">近30天</el-radio-button>
            </el-radio-group>
          </div>
          <div ref="energyChartRef" style="height: 300px;"></div>
          <div class="stat-card mt-20">
            <div class="flex-between">
              <div>
                <div class="stat-value">{{ totalEnergy.toFixed(2) }} kWh</div>
                <div class="stat-label">总能耗</div>
              </div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-dialog v-model="showEdit" title="编辑设备" width="500px">
      <el-form :model="editForm" ref="editFormRef" label-width="80px">
        <el-form-item label="设备名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="厂商">
          <el-input v-model="editForm.vendor" />
        </el-form-item>
        <el-form-item label="安装位置">
          <el-input v-model="editForm.location" />
        </el-form-item>
        <el-form-item label="功率(W)">
          <el-input-number v-model="editForm.power" :min="0" :max="10000" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEdit = false">取消</el-button>
        <el-button type="primary" @click="saveEdit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, type FormInstance } from 'element-plus'
import { ArrowLeft, Edit } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import {
  getDevice, updateDevice, updateDeviceStatus, getDeviceEnergy,
  deviceTypeOptions, protocolOptions, type Device
} from '@/api/device'

const route = useRoute()
const deviceId = Number(route.params.id)
const device = ref<Device | null>(null)
const showEdit = ref(false)
const editFormRef = ref<FormInstance>()
const energyPeriod = ref('day')
const totalEnergy = ref(0)
const energyChartRef = ref<HTMLElement>()
let energyChart: echarts.ECharts | null = null

const editForm = reactive({
  name: '',
  vendor: '',
  location: '',
  power: 0
})

onMounted(async () => {
  await loadDevice()
  energyChart = echarts.init(energyChartRef.value!)
  loadEnergy()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  energyChart?.dispose()
})

function handleResize() {
  energyChart?.resize()
}

async function loadDevice() {
  try {
    const res = await getDevice(deviceId)
    device.value = res
    Object.assign(editForm, {
      name: res.name,
      vendor: res.vendor,
      location: res.location,
      power: res.power
    })
  } catch (e) {
    console.error(e)
  }
}

async function loadEnergy() {
  try {
    const res = await getDeviceEnergy(deviceId, energyPeriod.value)
    totalEnergy.value = res.totalEnergy || 0
    const data = res.data || []
    energyChart?.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: 50, right: 20, top: 20, bottom: 40 },
      xAxis: {
        type: 'category',
        data: data.map((d: any) => d.timestamp?.substring(11, 16) || d.date)
      },
      yAxis: { type: 'value', name: 'W' },
      series: [{
        data: data.map((d: any) => d.power),
        type: 'line',
        smooth: true,
        areaStyle: { opacity: 0.3 },
        lineStyle: { color: '#67c23a', width: 2 },
        itemStyle: { color: '#67c23a' }
      }]
    })
  } catch (e) {
    console.error(e)
  }
}

async function toggleDeviceStatus() {
  if (!device.value) return
  const newStatus = device.value.status === 'on' || device.value.status === 'online' ? 'off' : 'on'
  try {
    await updateDeviceStatus(deviceId, newStatus)
    ElMessage.success('状态已更新')
    loadDevice()
  } catch (e) {
    console.error(e)
  }
}

async function saveEdit() {
  try {
    await updateDevice(deviceId, editForm)
    ElMessage.success('更新成功')
    showEdit.value = false
    loadDevice()
  } catch (e) {
    console.error(e)
  }
}

function getDeviceTypeLabel(type: string) {
  return deviceTypeOptions.find(o => o.value === type)?.label || type
}

function getProtocolLabel(proto: string) {
  return protocolOptions.find(o => o.value === proto)?.label || proto
}

function getStatusLabel(status: string) {
  const map: Record<string, string> = { online: '在线', offline: '离线', on: '开启', off: '关闭' }
  return map[status] || status
}

function getStatusType(status: string) {
  const map: Record<string, string> = { online: 'success', offline: 'info', on: 'success', off: 'warning' }
  return map[status] || 'info'
}
</script>
