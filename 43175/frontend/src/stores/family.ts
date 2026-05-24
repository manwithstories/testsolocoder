import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { listFamilies, type Family } from '@/api/family'

export const useFamilyStore = defineStore('family', () => {
  const families = ref<Family[]>([])
  const currentFamilyId = ref<number | null>(
    Number(localStorage.getItem('currentFamilyId')) || null
  )
  const loading = ref(false)

  const currentFamily = computed(() => {
    return families.value.find(f => f.id === currentFamilyId.value) || null
  })

  async function loadFamilies() {
    loading.value = true
    try {
      families.value = await listFamilies()
      if (families.value.length > 0) {
        if (!currentFamilyId.value || !families.value.find(f => f.id === currentFamilyId.value)) {
          currentFamilyId.value = families.value[0].id
          localStorage.setItem('currentFamilyId', String(currentFamilyId.value))
        }
      } else {
        currentFamilyId.value = null
        localStorage.removeItem('currentFamilyId')
      }
    } catch (e) {
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  function setCurrentFamily(familyId: number) {
    currentFamilyId.value = familyId
    localStorage.setItem('currentFamilyId', String(familyId))
  }

  return {
    families,
    currentFamilyId,
    currentFamily,
    loading,
    loadFamilies,
    setCurrentFamily
  }
})
