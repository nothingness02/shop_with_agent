<template>
  <div class="page-shell merchant-shell">
    <header class="site-header">
      <div class="brand">铸造台</div>
      <nav>
        <RouterLink to="/merchant">总览</RouterLink>
        <RouterLink to="/merchant/login">登录</RouterLink>
        <RouterLink to="/merchant/register" class="cta">创建工作台</RouterLink>
        <RouterLink to="/chat">聊天</RouterLink>
      </nav>
    </header>
    <main>
      <slot />
    </main>
    <footer>
      <div>统一 API ・ 商家运营</div>
      <RouterLink to="/user">用户空间</RouterLink>
    </footer>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount } from "vue";
import { RouterLink } from "vue-router";

onMounted(() => {
  document.documentElement.setAttribute("data-theme", "merchant");
});

onBeforeUnmount(() => {
  document.documentElement.setAttribute("data-theme", "neutral");
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
  font-size: 1.3rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
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
  border-radius: 10px;
  background: var(--accent-strong);
  color: var(--button-text);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-size: 0.75rem;
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