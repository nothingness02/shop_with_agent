<template>
  <section class="not-found">
    <h2>页面未找到</h2>
    <p>返回用户端或商家端。</p>
    <div class="actions">
      <RouterLink to="/user" class="primary">用户端</RouterLink>
      <RouterLink to="/merchant" class="ghost">商家端</RouterLink>
    </div>
  </section>
</template>

<script setup>
import { RouterLink } from "vue-router";
</script>

<style scoped>
.not-found {
  display: grid;
  gap: 12px;
  max-width: 420px;
}

.actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.primary,
.ghost {
  padding: 10px 16px;
  border-radius: 12px;
  font-weight: 600;
}

.primary {
  background: var(--accent-strong);
  color: var(--button-text);
}

.ghost {
  border: 1px solid var(--surface-border);
  color: var(--text-primary);
}
</style>