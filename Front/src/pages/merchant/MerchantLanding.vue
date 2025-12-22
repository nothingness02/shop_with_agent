<template>
  <section class="hero">
    <div>
      <p class="eyebrow">商家系统</p>
      <h1>实时掌控店铺，洞察经营数据。</h1>
      <p class="lead">
        面向商家的专属界面，整合数据、库存与运营聊天。角色由入口自动识别。
      </p>
      <div class="actions">
        <RouterLink to="/merchant/register" class="primary">创建工作台</RouterLink>
        <RouterLink to="/merchant/login" class="ghost">登录</RouterLink>
      </div>
    </div>
    <div class="hero-card">
      <div class="metric">
        <span>今日成交</span>
        <strong>$4,920</strong>
      </div>
      <div class="metric">
        <span>待处理订单</span>
        <strong>18</strong>
      </div>
      <div class="metric">
        <span>在线聊天</span>
        <strong>7 个会话</strong>
      </div>
    </div>
  </section>

  <section class="features">
    <FeatureCard
      title="库存聚焦"
      description="批量更新、SKU 检索与清晰审计。"
      icon="📊"
    />
    <FeatureCard
      title="订单指挥"
      description="随时跟踪订单并快速更新状态。"
      icon="🧭"
    />
    <FeatureCard
      title="商家沟通"
      description="无需离开控制台即可处理买家问题。"
      icon="⚡"
    />
  </section>
</template>

<script setup>
import { RouterLink } from "vue-router";
import FeatureCard from "../../components/FeatureCard.vue";
</script>

<style scoped>
.hero {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 28px;
  align-items: center;
}

.eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.3em;
  font-size: 0.7rem;
  color: var(--text-muted);
}

h1 {
  font-size: clamp(2.2rem, 3vw, 3.4rem);
  margin: 0 0 16px;
}

.lead {
  font-size: 1.05rem;
  color: var(--text-muted);
  line-height: 1.7;
}

.actions {
  display: flex;
  gap: 16px;
  margin-top: 20px;
  flex-wrap: wrap;
}

.primary,
.ghost {
  padding: 12px 20px;
  border-radius: 10px;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.08em;
}

.primary {
  background: var(--accent-strong);
  color: var(--button-text);
}

.ghost {
  border: 1px solid var(--surface-border);
  color: var(--text-primary);
}

.hero-card {
  padding: 26px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  display: grid;
  gap: 16px;
  box-shadow: var(--shadow-strong);
}

.metric {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  font-size: 0.95rem;
  color: var(--text-muted);
}

.metric strong {
  font-size: 1.4rem;
  color: var(--text-primary);
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 18px;
}
</style>