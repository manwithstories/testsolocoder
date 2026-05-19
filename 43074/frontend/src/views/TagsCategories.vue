<template>
  <div class="tags-categories-page">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>标签管理</span>
              <el-button type="primary" size="small" @click="showAddTag">
                <el-icon><Plus /></el-icon>
                新建标签
              </el-button>
            </div>
          </template>
          <div class="tags-list">
            <div
              v-for="tag in tags"
              :key="tag.id"
              class="tag-item"
            >
              <div class="tag-info">
                <span
                  class="tag-dot"
                  :style="{ backgroundColor: tag.color }"
                ></span>
                <span class="tag-name">{{ tag.name }}</span>
              </div>
              <div class="tag-actions">
                <el-button type="primary" link size="small" @click="editTag(tag)">编辑</el-button>
                <el-button type="danger" link size="small" @click="deleteTag(tag)">删除</el-button>
              </div>
            </div>
            <div v-if="!tags.length" class="empty">
              <el-empty description="暂无标签" :image-size="80" />
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>分类管理</span>
              <el-button type="primary" size="small" @click="showAddCategory">
                <el-icon><Plus /></el-icon>
                新建分类
              </el-button>
            </div>
          </template>
          <el-tree
            :data="categories"
            :props="{ label: 'name', children: 'children' }"
            node-key="id"
            default-expand-all
          >
            <template #default="{ node, data }">
              <div class="tree-node">
                <span>{{ data.name }}</span>
                <div class="node-actions">
                  <el-button type="primary" link size="small" @click.stop="addSubCategory(data)">
                    添加子分类
                  </el-button>
                  <el-button type="danger" link size="small" @click.stop="deleteCategory(data)">
                    删除
                  </el-button>
                </div>
              </div>
            </template>
          </el-tree>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="tagDialogVisible" :title="editingTag ? '编辑标签' : '新建标签'" width="400px">
      <el-form :model="tagForm" label-width="60px">
        <el-form-item label="名称">
          <el-input v-model="tagForm.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="颜色">
          <el-color-picker v-model="tagForm.color" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tagDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTag">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="categoryDialogVisible" :title="editingCategory ? '编辑分类' : '新建分类'" width="400px">
      <el-form :model="categoryForm" label-width="60px">
        <el-form-item label="名称">
          <el-input v-model="categoryForm.name" placeholder="请输入分类名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCategory">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTags, createTag, updateTag, deleteTag as apiDeleteTag, getCategories, createCategory, deleteCategory as apiDeleteCategory } from '@/api/common'
import type { Tag, Category } from '@/types'

const tags = ref<Tag[]>([])
const categories = ref<Category[]>([])

const tagDialogVisible = ref(false)
const editingTag = ref<Tag | null>(null)
const tagForm = reactive({ name: '', color: '#409eff' })

const categoryDialogVisible = ref(false)
const editingCategory = ref<Category | null>(null)
const parentCategory = ref<Category | null>(null)
const categoryForm = reactive({ name: '' })

const loadTags = async () => {
  try {
    tags.value = await getTags()
  } catch (e) {}
}

const loadCategories = async () => {
  try {
    categories.value = await getCategories(false)
  } catch (e) {}
}

const showAddTag = () => {
  editingTag.value = null
  tagForm.name = ''
  tagForm.color = '#409eff'
  tagDialogVisible.value = true
}

const editTag = (tag: Tag) => {
  editingTag.value = tag
  tagForm.name = tag.name
  tagForm.color = tag.color
  tagDialogVisible.value = true
}

const saveTag = async () => {
  if (!tagForm.name.trim()) {
    ElMessage.warning('请输入标签名称')
    return
  }
  try {
    if (editingTag.value) {
      await updateTag(editingTag.value.id, { name: tagForm.name, color: tagForm.color })
      ElMessage.success('更新成功')
    } else {
      await createTag({ name: tagForm.name, color: tagForm.color })
      ElMessage.success('创建成功')
    }
    tagDialogVisible.value = false
    loadTags()
  } catch (e) {}
}

const deleteTag = async (tag: Tag) => {
  try {
    await ElMessageBox.confirm(`确定删除标签「${tag.name}」吗？`, '确认', {
      type: 'warning'
    })
    await apiDeleteTag(tag.id)
    ElMessage.success('删除成功')
    loadTags()
  } catch (e) {}
}

const showAddCategory = () => {
  editingCategory.value = null
  parentCategory.value = null
  categoryForm.name = ''
  categoryDialogVisible.value = true
}

const addSubCategory = (category: Category) => {
  editingCategory.value = null
  parentCategory.value = category
  categoryForm.name = ''
  categoryDialogVisible.value = true
}

const saveCategory = async () => {
  if (!categoryForm.name.trim()) {
    ElMessage.warning('请输入分类名称')
    return
  }
  try {
    await createCategory({
      name: categoryForm.name,
      parent_id: parentCategory.value?.id
    })
    ElMessage.success('创建成功')
    categoryDialogVisible.value = false
    loadCategories()
  } catch (e) {}
}

const deleteCategory = async (category: Category) => {
  try {
    await ElMessageBox.confirm(`确定删除分类「${category.name}」吗？`, '确认', {
      type: 'warning'
    })
    await apiDeleteCategory(category.id)
    ElMessage.success('删除成功')
    loadCategories()
  } catch (e) {}
}

onMounted(() => {
  loadTags()
  loadCategories()
})
</script>

<style scoped lang="scss">
.tags-categories-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .tags-list {
    min-height: 400px;

    .tag-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 16px;
      border: 1px solid #ebeef5;
      border-radius: 6px;
      margin-bottom: 8px;

      .tag-info {
        display: flex;
        align-items: center;
        gap: 8px;

        .tag-dot {
          width: 12px;
          height: 12px;
          border-radius: 50%;
        }

        .tag-name {
          font-weight: 500;
        }
      }
    }
  }

  .tree-node {
    flex: 1;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .node-actions {
      opacity: 0;
      transition: opacity 0.2s;
    }

    &:hover .node-actions {
      opacity: 1;
    }
  }
}
</style>
