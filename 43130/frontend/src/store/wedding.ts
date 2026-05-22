import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { Wedding } from '@/types'
import { weddingApi } from '@/api/wedding'

export const useWeddingStore = defineStore('wedding', () => {
  const currentWedding = ref<Wedding | null>(null)
  const weddings = ref<Wedding[]>([])
  const loading = ref(false)

  const currentWeddingId = computed(() => currentWedding.value?.id || null)

  async function fetchWeddings(params?: any) {
    loading.value = true
    try {
      const res = await weddingApi.getList(params)
      weddings.value = res.data.list
      if (!currentWedding.value && weddings.value.length > 0) {
        currentWedding.value = weddings.value[0]
      }
    } finally {
      loading.value = false
    }
  }

  function setCurrentWedding(wedding: Wedding) {
    currentWedding.value = wedding
  }

  function clearCurrentWedding() {
    currentWedding.value = null
  }

  return {
    currentWedding,
    weddings,
    loading,
    currentWeddingId,
    fetchWeddings,
    setCurrentWedding,
    clearCurrentWedding
  }
})
