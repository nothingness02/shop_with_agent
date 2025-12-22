<template>
  <section class="chat-widget">
    <header>
      <div>
        <h3>{{ title }}</h3>
        <p>{{ subtitle }}</p>
      </div>
      <span class="status">在线</span>
    </header>

    <div class="messages">
      <div
        v-for="(message, index) in messages"
        :key="`${message.id}-${index}`"
        :class="['message', message.from]"
      >
        <span class="bubble">{{ message.text }}</span>
        <span class="time">{{ message.time }}</span>
      </div>
    </div>

    <form class="composer" @submit.prevent="sendMessage">
      <input
        v-model="draft"
        type="text"
        placeholder="输入消息"
        aria-label="聊天消息"
      />
      <button type="submit">发送</button>
    </form>
  </section>
</template>

<script setup>
import { onBeforeUnmount, ref } from "vue";
import { loadLocal, saveLocal } from "../services/localState.js";

const props = defineProps({
  title: { type: String, default: "聊天" },
  subtitle: { type: String, default: "用户与商家统一沟通通道" },
  storageKey: { type: String, default: "chat" }
});

const demoMessages = [
  { id: 1, from: "other", text: "需要订单方面的帮助吗？", time: "09:30" },
  { id: 2, from: "self", text: "物流更新了，谢谢！", time: "09:31" }
];

const messages = ref(loadLocal(props.storageKey, demoMessages));
const draft = ref("");
let timer = null;

const sendMessage = () => {
  const value = draft.value.trim();
  if (!value) {
    return;
  }
  messages.value = [
    ...messages.value,
    {
      id: Date.now(),
      from: "self",
      text: value,
      time: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
    }
  ];
  draft.value = "";
  saveLocal(props.storageKey, messages.value);

  timer = setTimeout(() => {
    messages.value = [
      ...messages.value,
      {
        id: Date.now() + 1,
        from: "other",
        text: "收到，我们会尽快回复。",
        time: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
      }
    ];
    saveLocal(props.storageKey, messages.value);
  }, 700);
};

onBeforeUnmount(() => {
  if (timer) {
    clearTimeout(timer);
  }
});
</script>

<style scoped>
.chat-widget {
  display: grid;
  gap: 16px;
  padding: 20px;
  border-radius: 22px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  box-shadow: var(--shadow-soft);
}

header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

header h3 {
  margin: 0 0 4px;
}

header p {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.status {
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 0.75rem;
  background: var(--accent-soft);
  color: var(--accent-strong);
}

.messages {
  display: grid;
  gap: 10px;
  max-height: 220px;
  overflow-y: auto;
  padding-right: 6px;
}

.message {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.message.self {
  align-items: flex-end;
}

.message.other {
  align-items: flex-start;
}

.bubble {
  padding: 10px 14px;
  border-radius: 16px;
  background: var(--bubble-bg);
  color: var(--text-primary);
  max-width: 80%;
  line-height: 1.4;
}

.message.self .bubble {
  background: var(--accent-strong);
  color: var(--button-text);
}

.time {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.composer {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 10px;
}

.composer input {
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
}

.composer button {
  padding: 12px 16px;
  border-radius: 14px;
  border: none;
  background: var(--accent-strong);
  color: var(--button-text);
  font-weight: 600;
  cursor: pointer;
}
</style>