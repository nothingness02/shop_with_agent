<template>
  <section class="form-wrap">
    <div>
      <p class="eyebrow">商家入驻</p>
      <h2>创建商家账号</h2>
      <p class="lead">角色由入口自动识别，无需手动选择。</p>
    </div>
    <form @submit.prevent="submit">
      <label>
        用户名
        <input v-model="form.username" type="text" required />
      </label>
      <label>
        邮箱
        <input v-model="form.email" type="email" required />
      </label>
      <label>
        密码
        <input v-model="form.password" type="password" minlength="6" required />
      </label>
      <button type="submit">创建商家</button>
      <p v-if="status" class="status">{{ status }}</p>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </section>
</template>

<script setup>
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { registerUser } from "../../services/api.js";

const router = useRouter();
const form = reactive({
  username: "",
  email: "",
  password: ""
});
const status = ref("");
const error = ref("");

const submit = async () => {
  error.value = "";
  status.value = "正在创建商家...";
  try {
    await registerUser({ ...form, role: 5 });
    status.value = "商家已创建，正在跳转到登录...";
    setTimeout(() => router.push("/merchant/login"), 800);
  } catch (err) {
    status.value = "";
    error.value = err.message;
  }
};
</script>

<style scoped>
.form-wrap {
  display: grid;
  gap: 18px;
  max-width: 480px;
}

.eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.3em;
  font-size: 0.7rem;
  color: var(--text-muted);
}

h2 {
  margin: 6px 0 8px;
}

form {
  display: grid;
  gap: 14px;
  padding: 24px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  box-shadow: var(--shadow-strong);
}

label {
  display: grid;
  gap: 6px;
  font-size: 0.9rem;
  color: var(--text-muted);
}

input {
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
}

button {
  padding: 12px 16px;
  border-radius: 10px;
  border: none;
  background: var(--accent-strong);
  color: var(--button-text);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-size: 0.75rem;
}

.status {
  color: var(--accent-strong);
  font-size: 0.9rem;
}

.error {
  color: #f25c78;
  font-size: 0.9rem;
}
</style>