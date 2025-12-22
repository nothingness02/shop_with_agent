<template>
  <section class="form-wrap">
    <div>
      <p class="eyebrow">欢迎回来</p>
      <h2>登录你的账号</h2>
      <p class="lead">登录后会自动识别角色。</p>
    </div>
    <form @submit.prevent="submit">
      <label>
        用户名
        <input v-model="form.username" type="text" required />
      </label>
      <label>
        密码
        <input v-model="form.password" type="password" required />
      </label>
      <button type="submit">登录</button>
      <p v-if="status" class="status">{{ status }}</p>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </section>
</template>

<script setup>
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { login } from "../../services/api.js";
import { setAuth } from "../../services/storage.js";

const router = useRouter();
const form = reactive({
  username: "",
  password: ""
});
const status = ref("");
const error = ref("");

const submit = async () => {
  error.value = "";
  status.value = "正在登录...";
  try {
    const data = await login(form);
    // console.log(data)
    setAuth(data);
    status.value = "登录成功。";
    setTimeout(() => router.push("/user/home"), 400);
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
  letter-spacing: 0.2em;
  font-size: 0.75rem;
  color: var(--text-muted);
}

h2 {
  margin: 6px 0 8px;
}

form {
  display: grid;
  gap: 14px;
  padding: 24px;
  border-radius: 20px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  box-shadow: var(--shadow-soft);
}

label {
  display: grid;
  gap: 6px;
  font-size: 0.9rem;
  color: var(--text-muted);
}

input {
  padding: 12px 14px;
  border-radius: 12px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
}

button {
  padding: 12px 16px;
  border-radius: 12px;
  border: none;
  background: var(--accent-strong);
  color: var(--button-text);
  font-weight: 600;
}

.status {
  color: var(--accent-strong);
  font-size: 0.9rem;
}

.error {
  color: #d6455d;
  font-size: 0.9rem;
}
</style>