<template>
  <section class="dashboard">
    <div>
      <h2>商家指挥中心</h2>
      <p class="lead">管理商品、处理订单、运营店铺。</p>
      <div class="stats">
        <div class="stat">
          <span>我的店铺</span>
          <strong>{{ shopId || "-" }}</strong>
        </div>
        <div class="stat">
          <span>商品数量</span>
          <strong>{{ products.length }}</strong>
        </div>
        <div class="stat">
          <span>订单数量</span>
          <strong>{{ orders.length }}</strong>
        </div>
      </div>
    </div>
    <ChatWidget
      title="商家聊天"
      subtitle="与买家对话的统一通道。"
      storage-key="merchant-chat"
    />
  </section>

  <section class="grid">
    <div class="panel">
      <div class="panel-header">
        <h3>店铺信息</h3>
      </div>
      <form class="form" @submit.prevent="submitShop">
        <label>
          店铺名称
          <input v-model="shopForm.name" type="text" required />
        </label>
        <label>
          店铺描述
          <input v-model="shopForm.description" type="text" />
        </label>
        <button type="submit" class="primary">创建/更新店铺</button>
        <p v-if="shopStatus" class="status">{{ shopStatus }}</p>
        <p v-if="shopError" class="error">{{ shopError }}</p>
      </form>
      <div class="inline-form">
        <label>
          当前店铺ID
          <input v-model="shopId" type="number" min="1" />
        </label>
        <button type="button" class="ghost" @click="saveShopId">保存</button>
      </div>
    </div>

    <div class="panel">
      <div class="panel-header">
        <h3>新增商品</h3>
      </div>
      <form class="form" @submit.prevent="submitProduct">
        <label>
          商品名称
          <input v-model="productForm.name" type="text" required />
        </label>
        <label>
          商品描述
          <input v-model="productForm.description" type="text" />
        </label>
        <div class="form-row">
          <label>
            价格
            <input v-model="productForm.price" type="number" step="0.01" min="0" required />
          </label>
          <label>
            库存
            <input v-model="productForm.stock" type="number" min="0" required />
          </label>
        </div>
        <label>
          图片URL
          <input v-model="productForm.img" type="text" />
        </label>
        <button type="submit" class="primary">发布商品</button>
        <p v-if="productStatus" class="status">{{ productStatus }}</p>
        <p v-if="productError" class="error">{{ productError }}</p>
      </form>
    </div>

    <div class="panel">
      <div class="panel-header">
        <h3>商品管理</h3>
        <button type="button" class="ghost" @click="loadProducts">刷新</button>
      </div>
      <ul v-if="products.length" class="list">
        <li v-for="(item, index) in products" :key="getProductKey(item, index)">
          <div>
            <strong>{{ item.Name || item.name || "未命名商品" }}</strong>
            <span class="muted">¥{{ item.Price || item.price }} · 库存 {{ item.Stock || item.stock }}</span>
          </div>
        </li>
      </ul>
      <p v-else class="empty">暂无商品</p>
    </div>

    <div class="panel">
      <div class="panel-header">
        <h3>订单管理</h3>
        <button type="button" class="ghost" @click="loadOrders">刷新</button>
      </div>
      <ul v-if="orders.length" class="list">
        <li v-for="(order, index) in orders" :key="getOrderKey(order, index)" class="order-item">
          <div>
            <strong>订单 {{ order.OrderID || order.order_id || order.ID || order.id }}</strong>
            <span class="muted">状态：{{ formatStatus(order.Status || order.status) }}</span>
          </div>
          <div class="order-actions">
            <select v-model="orderStates[getOrderKey(order, index)]">
              <option v-for="status in statusOptions" :key="status.value" :value="status.value">
                {{ status.label }}
              </option>
            </select>
            <button type="button" class="primary" @click="submitStatus(order, index)">
              更新
            </button>
          </div>
        </li>
      </ul>
      <p v-else class="empty">暂无订单</p>
      <p v-if="orderStatus" class="status">{{ orderStatus }}</p>
      <p v-if="orderError" class="error">{{ orderError }}</p>
    </div>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import ChatWidget from "../../components/ChatWidget.vue";
import {
  createProduct,
  createShop,
  listOrders,
  listProductsByShop,
  updateOrderStatus
} from "../../services/api.js";
import { getAuth } from "../../services/storage.js";

const shopId = ref("");
const shopForm = reactive({ name: "", description: "" });
const productForm = reactive({ name: "", description: "", price: "", stock: 0, img: "" });
const products = ref([]);
const orders = ref([]);
const orderStates = reactive({});

const shopStatus = ref("");
const shopError = ref("");
const productStatus = ref("");
const productError = ref("");
const orderStatus = ref("");
const orderError = ref("");

const statusOptions = [
  { value: "pending", label: "待支付" },
  { value: "paid", label: "已支付" },
  { value: "shipped", label: "已发货" },
  { value: "delivered", label: "已送达" },
  { value: "completed", label: "已完成" },
  { value: "cancelled", label: "已取消" },
  { value: "refunded", label: "已退款" }
];

const formatStatus = (status) => {
  const found = statusOptions.find((item) => item.value === (status || "").toLowerCase());
  return found ? found.label : "未知状态";
};

const getProductKey = (item, index) => item.ID || item.id || item.Name || item.name || index;
const getOrderKey = (order, index) => order.ID || order.id || order.OrderID || order.order_id || index;

const saveShopId = () => {
  if (shopId.value) {
    localStorage.setItem("merchant_shop_id", String(shopId.value));
  }
};

const submitShop = async () => {
  shopStatus.value = "";
  shopError.value = "";
  const auth = getAuth();
  if (!auth?.user_id) {
    shopError.value = "请先登录商家账号";
    return;
  }
  try {
    const shop = await createShop({
      name: shopForm.name,
      description: shopForm.description,
      owner_id: auth.user_id
    });
    shopId.value = shop?.ID || shop?.id || shopId.value;
    saveShopId();
    shopStatus.value = "店铺已保存";
    loadProducts();
  } catch (err) {
    shopError.value = err.message;
  }
};

const submitProduct = async () => {
  productStatus.value = "";
  productError.value = "";
  if (!shopId.value) {
    productError.value = "请先设置店铺ID";
    return;
  }
  try {
    await createProduct(shopId.value, {
      name: productForm.name,
      description: productForm.description,
      price: Number(productForm.price),
      stock: Number(productForm.stock),
      product_img: productForm.img
    });
    productStatus.value = "商品已发布";
    productForm.name = "";
    productForm.description = "";
    productForm.price = "";
    productForm.stock = 0;
    productForm.img = "";
    loadProducts();
  } catch (err) {
    productError.value = err.message;
  }
};

const loadProducts = async () => {
  if (!shopId.value) {
    products.value = [];
    return;
  }
  try {
    const data = await listProductsByShop(shopId.value);
    products.value = Array.isArray(data) ? data : [];
  } catch {
    products.value = [];
  }
};

const loadOrders = async () => {
  try {
    const data = await listOrders();
    orders.value = Array.isArray(data) ? data : [];
    orders.value.forEach((order, index) => {
      const key = getOrderKey(order, index);
      if (!orderStates[key]) {
        orderStates[key] = (order.Status || order.status || "pending").toLowerCase();
      }
    });
  } catch (err) {
    orderError.value = err.message;
    orders.value = [];
  }
};

const submitStatus = async (order, index) => {
  orderStatus.value = "";
  orderError.value = "";
  const id = order.ID || order.id;
  if (!id) {
    orderError.value = "缺少订单ID";
    return;
  }
  const key = getOrderKey(order, index);
  const status = orderStates[key];
  try {
    await updateOrderStatus(id, status);
    orderStatus.value = "订单状态已更新";
    loadOrders();
  } catch (err) {
    orderError.value = err.message;
  }
};

onMounted(() => {
  shopId.value = localStorage.getItem("merchant_shop_id") || "";
  loadProducts();
  loadOrders();
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
  border-radius: 12px;
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
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 18px;
}

.panel {
  padding: 20px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  display: grid;
  gap: 12px;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.form {
  display: grid;
  gap: 12px;
}

.form label {
  display: grid;
  gap: 6px;
  font-size: 0.9rem;
  color: var(--text-muted);
}

.form input {
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 12px;
}

.inline-form {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 10px;
  align-items: end;
}

.inline-form input {
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
}

.list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 10px;
}

.order-item {
  display: grid;
  gap: 10px;
}

.order-actions {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 8px;
  align-items: center;
}

select {
  padding: 8px 10px;
  border-radius: 10px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
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
  color: #f25c78;
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