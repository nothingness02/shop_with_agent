<template>
  <section class="chat-page">
    <div>
      <p class="eyebrow">共享通道</p>
      <h2>统一聊天空间</h2>
      <p class="lead">
        面向用户与商家的共享界面。后端接口开放前，将使用同一套 API 客户端与本地消息缓存。
      </p>
    </div>
    <ChatWidget
      title="店铺聊天"
      subtitle="用户与商家在同一空间沟通。"
      storage-key="shared-chat"
    />
  </section>
</template>

<script setup>
import ChatWidget from "../../components/ChatWidget.vue";
</script>

<style scoped>
.chat-page {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 24px;
  align-items: start;
}

.eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.2em;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.lead {
  color: var(--text-muted);
  line-height: 1.6;
}
</style>