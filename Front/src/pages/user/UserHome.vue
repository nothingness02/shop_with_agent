<template>
  <section class="dashboard">
    <div>
      <h2>用户中心</h2>
      <p class="lead">追踪订单、发现店铺、管理聊天。</p>
      <div class="stats">
        <div class="stat">
          <span>进行中的订单</span>
          <strong>{{ activeOrderCount }}</strong>
        </div>
        <div class="stat">
          <span>可见店铺</span>
          <strong>{{ shopCount }}</strong>
        </div>
        <div class="stat">
          <span>客服工单</span>
          <strong>0</strong>
        </div>
      </div>
    </div>
    <ChatWidget
      title="用户聊天"
      subtitle="一条线程联系商家与客服。"
      storage-key="user-chat"
    />
  </section>

  <section class="grid">
    <div class="panel">
      <div class="panel-header">
        <h3>订单进度</h3>
        <span v-if="ordersLoading">加载中...</span>
      </div>
      <ul v-if="orders.length">
        <li v-for="(order, index) in orders" :key="getOrderKey(order, index)">
          <span>{{ getOrderLabel(order) }}</span>
          <span class="muted">{{ formatStatus(order.Status || order.status) }}</span>
        </li>
      </ul>
      <p v-else class="empty">暂无订单</p>
    </div>
    <div class="panel">
      <div class="panel-header">
        <h3>推荐店铺</h3>
        <span v-if="shopsLoading">加载中...</span>
      </div>
      <ul v-if="shops.length">
        <li v-for="(shop, index) in shops" :key="getShopKey(shop, index)">
          <RouterLink
            v-if="getShopId(shop)"
            :to="`/shops/${getShopId(shop)}`"
            class="shop-link"
          >
            <span>{{ shop.Name || shop.name || "未命名店铺" }}</span>
            <span class="muted">{{ shop.Description || shop.description || "" }}</span>
          </RouterLink>
          <div v-else class="shop-link">
            <span>{{ shop.Name || shop.name || "未命名店铺" }}</span>
            <span class="muted">{{ shop.Description || shop.description || "" }}</span>
          </div>
        </li>
      </ul>
      <p v-else class="empty">暂无店铺</p>
    </div>
    <div class="panel">
      <div class="panel-header">
        <h3>快速下单</h3>
      </div>
      <form class="order-form" @submit.prevent="submitOrder">
        <div class="form-row">
          <label>
            商品ID
            <input v-model="draftItem.productId" type="number" min="1" required />
          </label>
          <label>
            商品名称
            <input v-model="draftItem.name" type="text" required />
          </label>
        </div>
        <div class="form-row">
          <label>
            价格
            <input v-model="draftItem.price" type="number" step="0.01" min="0" required />
          </label>
          <label>
            数量
            <input v-model="draftItem.quantity" type="number" min="1" required />
          </label>
        </div>
        <label>
          商品图片URL
          <input v-model="draftItem.img" type="text" placeholder="可选" />
        </label>
        <button type="button" class="ghost" @click="addItem">加入清单</button>

        <ul v-if="orderItems.length" class="order-items">
          <li v-for="(item, index) in orderItems" :key="index">
            <div>
              <strong>{{ item.product_name }}</strong>
              <span class="muted">ID {{ item.product_id }} · ¥{{ item.price }} × {{ item.quantity }}</span>
            </div>
            <button type="button" class="ghost" @click="removeItem(index)">移除</button>
          </li>
        </ul>
        <p v-else class="empty">暂无商品</p>

        <button type="submit" class="primary" :disabled="orderSubmitting">
          {{ orderSubmitting ? "提交中..." : "提交订单" }}
        </button>
        <p v-if="orderStatus" class="status">{{ orderStatus }}</p>
        <p v-if="orderError" class="error">{{ orderError }}</p>
      </form>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink } from "vue-router";
import ChatWidget from "../../components/ChatWidget.vue";
import { createOrder, listOrders, listShops } from "../../services/api.js";
import { getAuth } from "../../services/storage.js";

const orders = ref([]);
const shops = ref([]);
const ordersLoading = ref(false);
const shopsLoading = ref(false);
const orderSubmitting = ref(false);
const orderStatus = ref("");
const orderError = ref("");
const orderItems = ref([]);
const draftItem = reactive({
  productId: "",
  name: "",
  price: "",
  quantity: 1,
  img: ""
});

const inProgressStatuses = new Set(["pending", "paid", "shipped"]);

const activeOrderCount = computed(() =>
  orders.value.filter((order) => inProgressStatuses.has((order.Status || order.status || "").toLowerCase()))
    .length
);

const shopCount = computed(() => shops.value.length);

const formatStatus = (status) => {
  const value = (status || "").toLowerCase();
  const map = {
    pending: "待支付",
    paid: "已支付",
    shipped: "已发货",
    delivered: "已送达",
    completed: "已完成",
    cancelled: "已取消",
    refunded: "已退款"
  };
  return map[value] || "未知状态";
};

const getOrderKey = (order, index) => order.ID || order.id || order.OrderID || order.order_id || index;

const getOrderLabel = (order) => {
  const orderId = order.OrderID || order.order_id || order.ID || order.id || "-";
  const amount = order.ActualAmount || order.actual_amount || order.TotalAmount || order.total_amount || "";
  const amountText = amount !== "" ? `¥${amount}` : "";
  return `订单 ${orderId} ${amountText}`.trim();
};

const getShopKey = (shop, index) => shop.ID || shop.id || shop.Name || shop.name || index;
const getShopId = (shop) => shop.ID || shop.id || null;

const loadOrders = async () => {
  ordersLoading.value = true;
  try {
    const data = await listOrders();
    orders.value = Array.isArray(data) ? data : [];
  } catch {
    orders.value = [];
  } finally {
    ordersLoading.value = false;
  }
};

const loadShops = async () => {
  shopsLoading.value = true;
  try {
    const data = await listShops();
    shops.value = Array.isArray(data) ? data : [];
  } catch {
    shops.value = [];
  } finally {
    shopsLoading.value = false;
  }
};

const addItem = () => {
  const productId = Number(draftItem.productId);
  const price = Number(draftItem.price);
  const quantity = Number(draftItem.quantity);
  if (!productId || !draftItem.name || !price || !quantity) {
    orderError.value = "请完整填写商品信息";
    return;
  }
  orderItems.value.push({
    product_id: productId,
    product_name: draftItem.name.trim(),
    product_img: draftItem.img.trim(),
    price,
    quantity
  });
  orderError.value = "";
  draftItem.productId = "";
  draftItem.name = "";
  draftItem.price = "";
  draftItem.quantity = 1;
  draftItem.img = "";
};

const removeItem = (index) => {
  orderItems.value.splice(index, 1);
};

const submitOrder = async () => {
  orderStatus.value = "";
  orderError.value = "";
  const auth = getAuth();
  const userId = auth?.user_id;
  if (!userId) {
    orderError.value = "请先登录后再下单";
    return;
  }
  if (!orderItems.value.length) {
    orderError.value = "请先加入商品";
    return;
  }
  orderSubmitting.value = true;
  try {
    await createOrder({
      user_id: userId,
      items: orderItems.value
    });
    orderStatus.value = "订单已提交";
    orderItems.value = [];
    loadOrders();
  } catch (err) {
    orderError.value = err.message;
  } finally {
    orderSubmitting.value = false;
  }
};

onMounted(() => {
  loadOrders();
  loadShops();
});
</script>

<style scoped>
.dashboard {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 24px;
  align-items: start;
}

.lead {
  color: var(--text-muted);
  line-height: 1.6;
}

.stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 12px;
  margin-top: 18px;
}

.stat {
  padding: 16px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
}

.stat strong {
  display: block;
  font-size: 1.6rem;
  margin-top: 8px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 18px;
}

.panel {
  padding: 20px;
  border-radius: 20px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.panel ul {
  margin: 0;
  padding-left: 0;
  list-style: none;
  color: var(--text-primary);
  display: grid;
  gap: 10px;
}

.panel li {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.order-form {
  display: grid;
  gap: 12px;
}

.order-form label {
  display: grid;
  gap: 6px;
  font-size: 0.9rem;
  color: var(--text-muted);
}

.order-form input {
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 12px;
}

.order-items {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 10px;
}

.order-items li {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.order-form .primary {
  padding: 10px 16px;
  border-radius: 12px;
  border: none;
  background: var(--accent-strong);
  color: var(--button-text);
  font-weight: 600;
}

.order-form .ghost {
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

.shop-link {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  width: 100%;
  color: inherit;
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
