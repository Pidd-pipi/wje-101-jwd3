<template>
  <div class="login-page">
    <el-card class="card">
      <h2>{{ mode === 'login' ? '登录' : '注册' }}</h2>
      <el-form :model="form" label-width="70px">
        <el-form-item label="用户名"><el-input v-model="form.username" /></el-form-item>
        <el-form-item v-if="mode === 'register'" label="邮箱"><el-input v-model="form.email" /></el-form-item>
        <el-form-item label="密码"><el-input v-model="form.password" type="password" show-password /></el-form-item>
        <el-form-item v-if="mode === 'register'" label="简介"><el-input v-model="form.bio" type="textarea" :rows="2" /></el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" class="submit" @click="submit">
            {{ mode === 'login' ? '登录' : '注册并登录' }}
          </el-button>
        </el-form-item>
      </el-form>
      <el-button link type="primary" @click="toggle">
        {{ mode === 'login' ? '没有账号？去注册' : '已有账号？去登录' }}
      </el-button>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/useUserStore'
import { useBeanStore } from '@/stores/useBeanStore'
import { takePendingFavoriteBean } from '@/utils/storage'

const store = useUserStore()
const beanStore = useBeanStore()
const route = useRoute()
const router = useRouter()
const mode = ref<'login' | 'register'>('login')
const loading = ref(false)
const form = reactive({ username: '', email: '', password: '', bio: '' })

function toggle() {
  mode.value = mode.value === 'login' ? 'register' : 'login'
}

async function submit() {
  if (!form.username || !form.password) {
    ElMessage.warning('请填写用户名和密码')
    return
  }
  loading.value = true
  try {
    if (mode.value === 'login') {
      await store.login(form.username, form.password)
      ElMessage.success('登录成功')
    } else {
      await store.register({ username: form.username, email: form.email, password: form.password, bio: form.bio })
      ElMessage.success('注册成功')
    }
    // Coming back from a favorite click: keep the bean the user wanted.
    const pendingBeanId = takePendingFavoriteBean()
    if (pendingBeanId) {
      try {
        await beanStore.addFavorite(pendingBeanId)
        ElMessage.success('已为你收藏刚才选中的豆种')
      } catch {
        // Error toast already shown; the user lands back on the bean library.
      }
    }
    router.push((route.query.redirect as string) || '/')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page { display: flex; justify-content: center; padding-top: 60px; }
.card { width: 420px; }
.submit { width: 100%; }
</style>
