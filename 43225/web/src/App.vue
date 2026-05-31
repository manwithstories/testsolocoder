<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

onMounted(async () => {
  if (userStore.token && !userStore.user) {
    try {
      await userStore.fetchProfile()
    } catch (error) {
      console.error('Failed to fetch profile:', error)
    }
  }
})
</script>
