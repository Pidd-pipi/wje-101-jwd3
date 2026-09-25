<template>
  <div class="page">
    <h1>豆种库</h1>
    <el-tabs v-model="activeTab" class="tabs" @tab-change="onTabChange">
      <el-tab-pane name="all">
        <template #label>
          <span><el-icon><Coffee /></el-icon> 全部豆子</span>
        </template>
      </el-tab-pane>
      <el-tab-pane name="favorites">
        <template #label>
          <span><el-icon><Star /></el-icon> 我的收藏</span>
        </template>
      </el-tab-pane>
    </el-tabs>

    <SearchFilter @search="onSearch" @reset="onReset">
      <template #filters>
        <el-form-item label="产地">
          <el-select v-model="origin" clearable placeholder="全部产地" style="width: 150px" @change="load">
            <el-option v-for="o in ORIGINS" :key="o" :label="o" :value="o" />
          </el-select>
        </el-form-item>
        <el-form-item label="处理法">
          <el-select v-model="process" clearable placeholder="全部处理法" style="width: 150px" @change="load">
            <el-option v-for="(label, value) in ProcessMethodMap" :key="value" :label="label" :value="value" />
          </el-select>
        </el-form-item>
      </template>
    </SearchFilter>

    <el-row :gutter="16">
      <el-col v-for="b in displayedBeans" :key="b.id" :xs="24" :sm="12" :md="8">
        <el-card class="bean-card" shadow="hover">
          <div class="card-head">
            <h3>
              {{ b.name }}
              <el-tag size="small" type="warning">{{ ProcessMethodMap[b.process_method] }}</el-tag>
            </h3>
            <el-button
              :type="b.is_favorited ? 'warning' : 'default'"
              size="small"
              :icon="b.is_favorited ? StarFilled : Star"
              :loading="pendingIds.has(b.id)"
              @click="toggleFavorite(b)"
            >
              {{ b.favorite_count }}
            </el-button>
          </div>
          <div class="meta">{{ b.origin || '-' }}</div>
          <FlavorTags :tags="b.flavor_tags" />
          <p class="desc">{{ b.description }}</p>
          <div class="card-foot">
            <el-button v-if="isAdmin" size="small" type="danger" plain @click="removeBean(b.id)">删除</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <EmptyState v-if="!displayedBeans.length" :description="activeTab === 'favorites' ? '还没有收藏任何豆种，去全部豆子里看看吧' : '暂无豆种'" />

    <el-button v-if="isAdmin" type="primary" style="margin-top: 16px" @click="showAdd = true">新增豆种</el-button>
    <el-dialog v-model="showAdd" title="新增豆种" width="480px">
      <el-form label-width="80px">
        <el-form-item label="名称"><el-input v-model="addForm.name" /></el-form-item>
        <el-form-item label="产地"><el-input v-model="addForm.origin" /></el-form-item>
        <el-form-item label="处理法">
          <el-select v-model="addForm.process_method" style="width: 200px">
            <el-option v-for="(label, value) in ProcessMethodMap" :key="value" :label="label" :value="value" />
          </el-select>
        </el-form-item>
        <el-form-item label="风味标签"><el-input v-model="addForm.flavor_tags" placeholder='如 ["坚果","焦糖"]' /></el-form-item>
        <el-form-item label="描述"><el-input v-model="addForm.description" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAdd = false">取消</el-button>
        <el-button type="primary" @click="addBean">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Star, StarFilled, Coffee } from '@element-plus/icons-vue'
import SearchFilter from '@/components/common/SearchFilter.vue'
import FlavorTags from '@/components/common/FlavorTags.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { useBeanStore } from '@/stores/useBeanStore'
import { useAuth } from '@/hooks/useAuth'
import { createBean, deleteBean } from '@/api/bean'
import { setPendingFavoriteBean } from '@/utils/storage'
import { ProcessMethodMap, type BeanItem, type ProcessMethod } from '@/constants/bean'
import { useRoute, useRouter } from 'vue-router'

const store = useBeanStore()
const { isAdmin, isLoggedIn } = useAuth()
const route = useRoute()
const router = useRouter()
const activeTab = ref(route.query.tab === 'favorites' ? 'favorites' : 'all')
const origin = ref('')
const process = ref('')
const keyword = ref('')
const showAdd = ref(false)
const pendingIds = ref<Set<number>>(new Set())
const addForm = reactive({ name: '', origin: '', process_method: 'washed', flavor_tags: '[]', description: '' })

const ORIGINS = ['埃塞俄比亚', '哥伦比亚', '哥斯达黎加', '印度尼西亚']

const displayedBeans = computed<BeanItem[]>(() =>
  activeTab.value === 'favorites' ? store.favorites : store.beans,
)

onMounted(() => {
  if (activeTab.value === 'favorites' && !isLoggedIn.value) {
    ElMessage.info('登录后查看我的收藏')
    router.replace({ path: '/login', query: { redirect: route.fullPath } })
    return
  }
  load()
})

async function load() {
  const params = { page: 1, page_size: 100, origin: origin.value, process: process.value, keyword: keyword.value }
  if (activeTab.value === 'favorites') {
    if (!isLoggedIn.value) return
    await store.loadFavorites(params)
  } else {
    await store.load(params)
  }
}

function onTabChange(name: string | number) {
  activeTab.value = String(name)
  router.replace({ query: activeTab.value === 'favorites' ? { tab: 'favorites' } : {} })
  if (activeTab.value === 'favorites' && !isLoggedIn.value) {
    ElMessage.info('登录后查看我的收藏')
    router.push({ path: '/login', query: { redirect: route.fullPath } })
    return
  }
  load()
}

function onSearch(kw: string) {
  keyword.value = kw
  load()
}
function onReset() {
  origin.value = ''
  process.value = ''
  keyword.value = ''
  load()
}

async function toggleFavorite(b: BeanItem) {
  if (!isLoggedIn.value) {
    // Remember the bean so it can be favorited right after coming back.
    setPendingFavoriteBean(b.id)
    ElMessage.info('请先登录，登录后将自动为你收藏该豆种')
    router.push({ path: '/login', query: { redirect: route.fullPath } })
    return
  }
  if (pendingIds.value.has(b.id)) return
  pendingIds.value.add(b.id)
  try {
    if (b.is_favorited) {
      await store.removeFavorite(b.id)
      ElMessage.success('已取消收藏')
    } else {
      await store.addFavorite(b.id)
      ElMessage.success('已收藏')
    }
  } finally {
    pendingIds.value.delete(b.id)
  }
}

async function addBean() {
  if (!addForm.name) {
    ElMessage.warning('请填写名称')
    return
  }
  await createBean({ ...addForm, process_method: addForm.process_method as ProcessMethod })
  ElMessage.success('豆种已添加')
  showAdd.value = false
  await load()
}
async function removeBean(id: number) {
  await deleteBean(id)
  ElMessage.success('已下架，相关收藏已清理')
  await load()
}
</script>

<style scoped>
.page { max-width: 1200px; margin: 0 auto; }
.tabs { margin-bottom: 8px; }
.bean-card { margin-bottom: 16px; }
.card-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 8px; }
.card-head h3 { margin: 0; font-size: 16px; display: flex; align-items: center; gap: 6px; flex-wrap: wrap; }
.card-foot { margin-top: 8px; }
.meta { color: #999; font-size: 12px; margin: 6px 0; }
.desc { color: #666; margin-top: 8px; }
</style>
