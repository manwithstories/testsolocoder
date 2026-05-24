<template>
  <el-config-provider :locale="locale">
    <router-view />
  </el-config-provider>
</template>

<script setup lang="ts">
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { computed } from 'vue'

const locale = computed(() => zhCn)
</script>

<style>
#app {
  min-height: 100vh;
  background-color: #f5f7fa;
}
</style>
