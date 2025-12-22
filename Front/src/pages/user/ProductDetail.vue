<template>
  <section class="detail">
    <div class="breadcrumb">
      <span>首页</span>
      <span>/</span>
      <span>商品详情</span>
      <span>/</span>
      <span>{{ productName }}</span>
    </div>

    <div class="top-bar">
      <div class="store">
        <div class="store-avatar">🏬</div>
        <div>
          <strong>{{ shopName }}</strong>
          <div class="muted">官方旗舰 · 快速发货</div>
        </div>
      </div>
      <div class="store-actions">
        <button type="button" class="ghost">关注店铺</button>
        <button type="button" class="ghost">联系客服</button>
        <RouterLink v-if="shopId" :to="`/shops/${shopId}`" class="ghost-link">
          进入店铺
        </RouterLink>
      </div>
    </div>

    <div class="main">
      <div class="gallery">
        <div class="thumb-list">
          <div class="thumb-mini" v-for="n in 4" :key="n"></div>
        </div>
        <div class="hero-img"></div>
      </div>

      <div class="purchase">
        <div class="promo-badge">礼遇季</div>
        <h2>{{ productName }}</h2>
        <p class="lead">{{ productDesc }}</p>
        <div class="price">{{ formatPrice(productPrice) }}</div>
        <div class="meta">
          <span>库存：{{ productStock }}</span>
          <span>发货：24小时内</span>
        </div>
        <div class="specs">
          <span>规格</span>
          <div class="spec-tags">
            <button type="button" class="tag">标准款</button>
            <button type="button" class="tag">高级款</button>
            <button type="button" class="tag">礼盒装</button>
          </div>
        </div>
        <div class="qty">
          <span>数量</span>
          <div class="qty-control">
            <button type="button" @click="decrease">-</button>
            <input v-model.number="quantity" type="number" min="1" />
            <button type="button" @click="increase">+</button>
          </div>
        </div>
        <div class="buy-actions">
          <button type="button" class="primary" :disabled="submitting" @click="submitOrder">
            {{ submitting ? "提交中..." : "立即下单" }}
          </button>
          <button type="button" class="secondary" :disabled="addingCart" @click="addToCart">
            {{ addingCart ? "加入中..." : "加入购物车" }}
          </button>
        </div>
        <p v-if="loading" class="muted">加载中...</p>
        <p v-if="orderStatus" class="status">{{ orderStatus }}</p>
        <p v-if="orderError" class="error">{{ orderError }}</p>
        <p v-if="cartStatus" class="status">{{ cartStatus }}</p>
      </div>
    </div>

    <div class="detail-tabs">
      <button type="button" class="tab active">商品详情</button>
      <button type="button" class="tab">规格参数</button>
      <button type="button" class="tab">售后保障</button>
    </div>

    <div class="detail-panel">
      <h3>商品卖点</h3>
      <ul>
        <li>精选材质，工艺细节可见</li>
        <li>适合多种场景，送礼自用皆可</li>
        <li>售后无忧，支持退换</li>
      </ul>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, RouterLink } from "vue-router";
import { addCartItem, createOrder, getProduct } from "../../services/api.js";
import { getAuth } from "../../services/storage.js";

const route = useRoute();
const product = ref(null);
const loading = ref(false);
const quantity = ref(1);
const submitting = ref(false);
const addingCart = ref(false);
const orderStatus = ref("");
const orderError = ref("");
const cartStatus = ref("");

const productName = computed(() => product.value?.Name || product.value?.name || "商品详情");
const productDesc = computed(() => product.value?.Description || product.value?.description || "暂无描述");
const productPrice = computed(() => product.value?.Price || product.value?.price || "");
const productStock = computed(() => product.value?.Stock || product.value?.stock || "-");
const shopId = computed(() => product.value?.ShopID || product.value?.shop_id || null);
const shopName = computed(() => product.value?.Shop?.Name || "官方店铺");

const formatPrice = (price) => {
  if (price === undefined || price === null || price === "") {
    return "";
  }
  return `¥${price}`;
};

const getProductPayload = () => {
  if (!product.value) {
    return null;
  }
  return {
    product_id: product.value.ID || product.value.id,
    product_name: product.value.Name || product.value.name || "",
    product_img: product.value.ProductImg || product.value.product_img || "",
    price: Number(product.value.Price || product.value.price || 0),
    quantity: Number(quantity.value || 1)
  };
};

const submitOrder = async () => {
  orderStatus.value = "";
  orderError.value = "";
  cartStatus.value = "";
  const auth = getAuth();
  const userId = auth?.user_id;
  if (!userId) {
    orderError.value = "请先登录后再下单";
    return;
  }
  const item = getProductPayload();
  if (!item || !item.product_id || !item.product_name) {
    orderError.value = "商品信息不完整";
    return;
  }
  submitting.value = true;
  try {
    await createOrder({
      user_id: userId,
      items: [item]
    });
    orderStatus.value = "订单已提交";
  } catch (err) {
    orderError.value = err.message;
  } finally {
    submitting.value = false;
  }
};

const addToCart = async () => {
  cartStatus.value = "";
  orderError.value = "";
  const auth = getAuth();
  if (!auth?.user_id) {
    orderError.value = "请先登录后再加入购物车";
    return;
  }
  const item = getProductPayload();
  if (!item || !item.product_id) {
    orderError.value = "商品信息不完整";
    return;
  }
  addingCart.value = true;
  try {
    await addCartItem({ product_id: item.product_id, quantity: item.quantity });
    cartStatus.value = "已加入购物车";
  } catch (err) {
    orderError.value = err.message;
  } finally {
    addingCart.value = false;
  }
};

const increase = () => {
  quantity.value = Number(quantity.value || 1) + 1;
};

const decrease = () => {
  const next = Number(quantity.value || 1) - 1;
  quantity.value = next < 1 ? 1 : next;
};

const load = async (id) => {
  if (!id) {
    return;
  }
  loading.value = true;
  try {
    product.value = await getProduct(id);
  } catch {
    product.value = null;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  load(route.params.id);
});

watch(
  () => route.params.id,
  (id) => {
    load(id);
  }
);
</script>

<style scoped>
.detail {
  display: grid;
  gap: 18px;
}

.breadcrumb {
  display: flex;
  gap: 8px;
  font-size: 0.85rem;
  color: var(--text-muted);
}

.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
}

.store {
  display: flex;
  gap: 12px;
  align-items: center;
}

.store-avatar {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  background: rgba(225, 37, 27, 0.12);
  display: grid;
  place-items: center;
}

.store-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.ghost,
.ghost-link {
  padding: 8px 14px;
  border-radius: 999px;
  border: 1px solid var(--surface-border);
  background: transparent;
  color: var(--text-primary);
}

.main {
  display: grid;
  grid-template-columns: 1.6fr 1fr;
  gap: 18px;
  align-items: start;
}

.gallery {
  display: grid;
  grid-template-columns: 80px 1fr;
  gap: 16px;
}

.thumb-list {
  display: grid;
  gap: 10px;
}

.thumb-mini {
  width: 72px;
  height: 72px;
  border-radius: 12px;
  background: rgba(225, 37, 27, 0.12);
  border: 1px solid var(--surface-border);
}

.hero-img {
  height: 420px;
  border-radius: 20px;
  background: radial-gradient(circle at 30% 20%, rgba(225, 37, 27, 0.15), transparent 60%),
    #1a1a1a;
}

.purchase {
  padding: 18px;
  border-radius: 18px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  display: grid;
  gap: 12px;
}

.promo-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(225, 37, 27, 0.12);
  color: var(--accent-strong);
  font-size: 0.75rem;
  width: fit-content;
}

.lead {
  color: var(--text-muted);
  margin: 0;
}

.price {
  font-size: 1.6rem;
  color: var(--accent-strong);
  font-weight: 700;
}

.meta {
  display: flex;
  gap: 16px;
  color: var(--text-muted);
  font-size: 0.85rem;
}

.specs {
  display: grid;
  gap: 12px;
}

.spec-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tag {
  padding: 6px 12px;
  border-radius: 999px;
  border: 1px solid var(--surface-border);
  background: transparent;
  cursor: pointer;
}

.qty {
  display: flex;
  align-items: center;
  gap: 12px;
}

.qty-control {
  display: grid;
  grid-template-columns: 32px 60px 32px;
  gap: 6px;
  align-items: center;
}

.qty-control input {
  text-align: center;
  padding: 8px;
  border-radius: 8px;
  border: 1px solid var(--surface-border);
  background: transparent;
}

.qty-control button {
  border: 1px solid var(--surface-border);
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
}

.buy-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.primary,
.secondary {
  padding: 12px 16px;
  border-radius: 12px;
  border: none;
  font-weight: 600;
  cursor: pointer;
}

.primary {
  background: var(--accent-strong);
  color: var(--button-text);
}

.secondary {
  background: rgba(225, 37, 27, 0.12);
  color: var(--accent-strong);
}

.detail-tabs {
  display: flex;
  gap: 12px;
  padding-top: 8px;
}

.tab {
  padding: 10px 16px;
  border-radius: 999px;
  border: 1px solid var(--surface-border);
  background: transparent;
  cursor: pointer;
}

.tab.active {
  border-color: var(--accent-strong);
  color: var(--accent-strong);
}

.detail-panel {
  padding: 18px;
  border-radius: 18px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
}

.detail-panel ul {
  margin: 10px 0 0;
  padding-left: 18px;
  color: var(--text-muted);
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
}

@media (max-width: 1000px) {
  .main {
    grid-template-columns: 1fr;
  }

  .gallery {
    grid-template-columns: 1fr;
  }

  .thumb-list {
    grid-template-columns: repeat(4, 1fr);
  }
}
</style>