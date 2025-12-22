<template>
  <div class="feature-card">
    <div class="feature-icon">{{ icon }}</div>
    <div>
      <h3>{{ title }}</h3>
      <p>{{ description }}</p>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: String,
  description: String,
  icon: String
});
</script>

<style scoped>
.feature-card {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 16px;
  align-items: start;
  padding: 18px 20px;
  border-radius: 18px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  box-shadow: var(--shadow-soft);
}

.feature-icon {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  background: var(--accent-soft);
  color: var(--accent-strong);
  font-size: 20px;
}

h3 {
  margin: 0 0 6px;
  font-size: 1.05rem;
}

p {
  margin: 0;
  color: var(--text-muted);
  line-height: 1.5;
}
</style>
