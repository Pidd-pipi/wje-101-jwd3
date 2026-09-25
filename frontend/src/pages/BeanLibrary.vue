<template>
  <div class="page">
    <h1>豆种库</h1>
    <el-tabs v-model="activeTab" class="view-tabs" @tab-change="onTabChange">
      <el-tab-pane label="全部豆种" name="all" />
      <el-tab-pane label="我的收藏" name="favorites" />
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
    <el-row v-loading="loading" :gutter="16">
      <el-col v-for="b in beans" :key="b.id" :xs="24" :sm="12" :md="8">
        <el-card class="bean-card" shadow="hover">
          <h3>{{ b.name }} <el-tag size="small" type="warning">{{ ProcessMethodMap[b.process_method] }}</el-tag></h3>
          <div class="meta">{{ b.origin || '-' }}</div>
          <FlavorTags :tags="b.flavor_tags" />
          <p class="desc">{{ b.description }}</p>
          <div class="card-footer">
            <el-button
              :type="b.is_favorite ? 'warning' : 'default'"
              :loading="store.isPending(b.id)"
              size="small"
              @click="onToggleFavorite(b)"
            >
              <el-icon class="star-icon"><StarFilled v-if="b.is_favorite" /><Star v-else /></el-icon>
              {{ b.is_favorite ? '已收藏' : '收藏' }}
              <span class="fav-count">{{ b.favorite_count ?? 0 }}</span>
            </el-button>
            <el-button v-if="isAdmin" size="small" type="danger" plain @click="removeBean(b.id)">删除</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <EmptyState v-if="!loading && !beans.length" :description="activeTab === 'favorites' ? '还没有收藏任何豆种，去豆种库看看吧' : '暂无豆种'" />
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
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Star, StarFilled } from '@element-plus/icons-vue'
import SearchFilter from '@/components/common/SearchFilter.vue'
import FlavorTags from '@/components/common/FlavorTags.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { useBeanStore, type BeanViewMode } from '@/stores/useBeanStore'
import { useAuth } from '@/hooks/useAuth'
import { favoriteBean, createBean, deleteBean } from '@/api/bean'
import { setPendingBeanFavorite, takePendingBeanFavorite } from '@/utils/storage'
import { ProcessMethodMap, type CoffeeBean, type ProcessMethod } from '@/constants/bean'

const store = useBeanStore()
const { isLoggedIn, isAdmin } = useAuth()
const route = useRoute()
const router = useRouter()
const beans = computed(() => store.beans)
const activeTab = ref<BeanViewMode>(route.query.tab === 'favorites' ? 'favorites' : 'all')
const loading = ref(false)
const origin = ref('')
const process = ref('')
const keyword = ref('')
const showAdd = ref(false)
const addForm = reactive({ name: '', origin: '', process_method: 'washed', flavor_tags: '[]', description: '' })

const ORIGINS = ['埃塞俄比亚', '哥伦比亚', '哥斯达黎加', '印度尼西亚']

onMounted(async () => {
  // 未登录点收藏跳到登录页，回来后自动补上刚才想收藏的豆子。
  const pendingBeanId = takePendingBeanFavorite()
  if (isLoggedIn.value && pendingBeanId !== null) {
    try {
      await favoriteBean(pendingBeanId)
      ElMessage.success('已收藏')
    } catch {
      // 豆子可能已被管理员下架，错误提示已由请求拦截器给出
    }
  }
  if (activeTab.value === 'favorites' && !isLoggedIn.value) {
    activeTab.value = 'all'
    router.push({ path: '/login', query: { redirect: route.fullPath } })
    return
  }
  await load()
})

async function load() {
  loading.value = true
  try {
    await store.load({
      mode: activeTab.value,
      page: 1,
      page_size: 20,
      origin: origin.value,
      process: process.value,
      keyword: keyword.value,
    })
  } finally {
    loading.value = false
  }
}

function onTabChange(tab: string | number) {
  const next = tab === 'favorites' ? 'favorites' : 'all'
  if (next === 'favorites' && !isLoggedIn.value) {
    ElMessage.info('登录后即可查看我的收藏')
    router.push({ path: '/login', query: { redirect: route.fullPath } })
    // 保持原 tab 高亮，登录回来仍能看到收藏列表
    activeTab.value = 'all'
    return
  }
  activeTab.value = next
  router.replace({ query: { ...route.query, tab: next === 'all' ? undefined : next } })
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

async function onToggleFavorite(b: CoffeeBean) {
  if (!isLoggedIn.value) {
    setPendingBeanFavorite(b.id)
    ElMessage.info('请先登录，登录后将自动收藏该豆种')
    router.push({ path: '/login', query: { redirect: route.fullPath } })
    return
  }
  try {
    const favorited = await store.toggleFavorite(b.id)
    ElMessage.success(favorited ? '已收藏' : '已取消收藏')
  } catch {
    // 错误提示已由请求拦截器统一给出
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
  ElMessage.success('已删除')
  await load()
}
</script>

<style scoped>
.page { max-width: 1200px; margin: 0 auto; }
.view-tabs { margin-bottom: 8px; }
.bean-card { margin-bottom: 16px; }
.meta { color: #999; font-size: 12px; margin: 6px 0; }
.desc { color: #666; margin-top: 8px; }
.card-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 8px; }
.star-icon { margin-right: 4px; vertical-align: -2px; }
.fav-count { margin-left: 2px; }
</style>
