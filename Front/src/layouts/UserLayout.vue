<template>
  <div class="page-shell user-shell">
    <header class="site-header">
      <div class="brand">悦购</div>
      <nav>
        <RouterLink to="/user">首页</RouterLink>
        <RouterLink to="/cart">购物车</RouterLink>
        <RouterLink to="/chat">客服聊天</RouterLink>
        <RouterLink to="/merchant" class="merchant-link">商家入口</RouterLink>
        <template v-if="isAuthed">
          <RouterLink to="/user/home" class="cta">我的</RouterLink>
          <button type="button" class="ghost" @click="handleLogout">退出</button>
        </template>
        <template v-else>
          <RouterLink to="/user/login">登录</RouterLink>
          <RouterLink to="/user/register" class="cta">立即注册</RouterLink>
        </template>
      </nav>
    </header>
    <main>
      <slot />
    </main>
    <footer>
      <div>统一 API ・ 用户体验</div>
      <RouterLink to="/merchant">商家控制台</RouterLink>
    </footer>
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref } from "vue";
import { useRouter, RouterLink } from "vue-router";
import { clearAuth, getAuth } from "../services/storage.js";
import { logout } from "../services/api.js";

const router = useRouter();
const isAuthed = ref(false);

const syncAuth = () => {
  const auth = getAuth();
  isAuthed.value = Boolean(auth && auth.access_token);
};

const handleLogout = async () => {
  try {
    await logout();
  } catch {
    // ignore logout errors, still clear local state
  }
  clearAuth();
  router.push("/user");
};

onMounted(() => {
  document.documentElement.setAttribute("data-theme", "user");
  syncAuth();
  window.addEventListener("storage", syncAuth);
  window.addEventListener("auth-changed", syncAuth);
});

onBeforeUnmount(() => {
  document.documentElement.setAttribute("data-theme", "neutral");
  window.removeEventListener("storage", syncAuth);
  window.removeEventListener("auth-changed", syncAuth);
});
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  padding: 28px clamp(20px, 5vw, 64px) 40px;
  display: grid;
  gap: 40px;
}

.site-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.brand {
  font-family: var(--font-display);
  font-size: 1.4rem;
  letter-spacing: 0.04em;
}

nav {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

nav a {
  color: var(--text-primary);
  font-weight: 500;
}

nav a.cta {
  padding: 10px 16px;
  border-radius: 999px;
  background: var(--accent-strong);
  color: var(--button-text);
}

nav .merchant-link {
  font-size: 0.85rem;
  color: var(--text-muted);
}

nav button.ghost {
  padding: 10px 16px;
  border-radius: 999px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
  cursor: pointer;
}

main {
  display: grid;
  gap: 32px;
}

footer {
  display: flex;
  justify-content: space-between;
  color: var(--text-muted);
  font-size: 0.9rem;
}

@media (max-width: 720px) {
  footer {
    flex-direction: column;
    gap: 12px;
  }
}
</style>
