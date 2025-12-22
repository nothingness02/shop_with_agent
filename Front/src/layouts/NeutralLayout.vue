<template>
  <div class="page-shell neutral-shell">
    <header class="site-header">
      <div class="brand">商店系统</div>
      <nav>
        <RouterLink to="/user">用户端</RouterLink>
        <RouterLink to="/merchant">商家端</RouterLink>
      </nav>
    </header>
    <main>
      <slot />
    </main>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount } from "vue";
import { RouterLink } from "vue-router";
import { getAuth } from "../services/storage.js";

const getTheme = () => {
  const auth = getAuth();
  if (!auth || !auth.role) {
    return "neutral";
  }
  return auth.role === 5 ? "merchant" : "user";
};

onMounted(() => {
  document.documentElement.setAttribute("data-theme", getTheme());
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
  gap: 32px;
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
  letter-spacing: 0.03em;
}

nav {
  display: flex;
  gap: 16px;
}

nav a {
  color: var(--text-primary);
  font-weight: 500;
}
</style>