<template>
  <section class="cart">
    <div class="header">
      <h2>购物车</h2>
      <div class="actions">
        <button type="button" class="ghost" @click="loadCart">刷新</button>
        <button type="button" class="ghost" @click="clearAll" :disabled="clearing">清空</button>
      </div>
    </div>

    <div class="panel">
      <div v-if="items.length" class="list">
        <div class="row" v-for="item in items" :key="item.ID || item.id">
          <div class="info">
            <strong>{{ item.ProductName || item.product_name }}</strong>
            <span class="muted">¥{{ item.Price || item.price }} · 商品ID {{ item.ProductID || item.product_id }}</span>
          </div>
          <div class="qty">
            <button type="button" @click="adjust(item, -1)">-</button>
            <span>{{ item.Quantity || item.quantity }}</span>
            <button type="button" @click="adjust(item, 1)">+</button>
          </div>
          <button type="button" class="ghost" @click="remove(item)">移除</button>
        </div>
      </div>
      <p v-else class="empty">购物车为空</p>
    </div>

    <div class="panel">
      <div class="checkout">
        <div>
          <strong>合计</strong>
          <span class="price">¥{{ totalPrice }}</span>
        </div>
        <button type="button" class="primary" :disabled="submitting" @click="checkout">
          {{ submitting ? "提交中..." : "下单" }}
        </button>
      </div>
      <p v-if="status" class="status">{{ status }}</p>
      <p v-if="error" class="error">{{ error }}</p>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { clearCart, createOrder, deleteCartItem, listCart, updateCartItem } from "../../services/api.js";
import { getAuth } from "../../services/storage.js";

const items = ref([]);
const status = ref("");
const error = ref("");
const submitting = ref(false);
const clearing = ref(false);

const totalPrice = computed(() => {
  return items.value.reduce((sum, item) => {
    const price = Number(item.Price || item.price || 0);
    const qty = Number(item.Quantity || item.quantity || 0);
    return sum + price * qty;
  }, 0).toFixed(2);
});

const loadCart = async () => {
  status.value = "";
  error.value = "";
  try {
    const data = await listCart();
    items.value = Array.isArray(data) ? data : [];
  } catch (err) {
    error.value = err.message;
    items.value = [];
  }
};

const adjust = async (item, delta) => {
  const id = item.ID || item.id;
  if (!id) {
    return;
  }
  const current = Number(item.Quantity || item.quantity || 0);
  const next = current + delta;
  if (next < 1) {
    return;
  }
  try {
    await updateCartItem(id, { quantity: next });
    item.Quantity = next;
    item.quantity = next;
  } catch (err) {
    error.value = err.message;
  }
};

const remove = async (item) => {
  const id = item.ID || item.id;
  if (!id) {
    return;
  }
  try {
    await deleteCartItem(id);
    items.value = items.value.filter((row) => (row.ID || row.id) !== id);
  } catch (err) {
    error.value = err.message;
  }
};

const clearAll = async () => {
  clearing.value = true;
  try {
    await clearCart();
    items.value = [];
  } catch (err) {
    error.value = err.message;
  } finally {
    clearing.value = false;
  }
};

const checkout = async () => {
  status.value = "";
  error.value = "";
  const auth = getAuth();
  if (!auth?.user_id) {
    error.value = "请先登录";
    return;
  }
  if (!items.value.length) {
    error.value = "购物车为空";
    return;
  }
  submitting.value = true;
  try {
    const payloadItems = items.value.map((item) => ({
      product_id: item.ProductID || item.product_id,
      product_name: item.ProductName || item.product_name,
      product_img: item.ProductImg || item.product_img,
      price: Number(item.Price || item.price || 0),
      quantity: Number(item.Quantity || item.quantity || 1)
    }));
    await createOrder({
      user_id: auth.user_id,
      items: payloadItems
    });
    await clearCart();
    items.value = [];
    status.value = "订单已提交";
  } catch (err) {
    error.value = err.message;
  } finally {
    submitting.value = false;
  }
};

onMounted(() => {
  loadCart();
});
</script>

<style scoped>
.cart {
  display: grid;
  gap: 16px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.actions {
  display: flex;
  gap: 10px;
}

.panel {
  padding: 18px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
}

.list {
  display: grid;
  gap: 12px;
}

.row {
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 12px;
  align-items: center;
}

.qty {
  display: grid;
  grid-template-columns: 32px 40px 32px;
  gap: 6px;
  align-items: center;
  justify-items: center;
}

.qty button {
  border: 1px solid var(--surface-border);
  background: transparent;
  border-radius: 8px;
  width: 32px;
  height: 32px;
  cursor: pointer;
}

.checkout {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price {
  color: var(--accent-strong);
  font-size: 1.2rem;
  font-weight: 700;
}

.primary {
  padding: 10px 16px;
  border-radius: 12px;
  border: none;
  background: var(--accent-strong);
  color: var(--button-text);
  font-weight: 600;
}

.ghost {
  padding: 8px 14px;
  border-radius: 999px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
  cursor: pointer;
}

.status {
  color: var(--accent-strong);
  font-size: 0.9rem;
}

.error {
  color: #d6455d;
  font-size: 0.9rem;
}

.muted {
  color: var(--text-muted);
  font-size: 0.85rem;
}

.empty {
  margin: 0;
  color: var(--text-muted);
}
</style>
