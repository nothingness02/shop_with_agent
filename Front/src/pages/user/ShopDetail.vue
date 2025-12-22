<template>
  <section class="detail">
    <div class="header">
      <div>
        <h2>{{ shopName }}</h2>
        <p class="lead">{{ shopDesc }}</p>
      </div>
      <RouterLink to="/user/home" class="ghost">返回用户中心</RouterLink>
    </div>

    <div class="panel">
      <div class="panel-header">
        <h3>店铺商品</h3>
        <span v-if="loading">加载中...</span>
      </div>
      <div v-if="products.length" class="product-grid">
        <RouterLink
          v-for="(item, index) in products"
          :key="getProductKey(item, index)"
          :to="getProductId(item) ? `/products/${getProductId(item)}` : '/user/home'"
          class="product-card"
        >
          <div class="thumb"></div>
          <div>
            <h4>{{ item.Name || item.name || "未命名商品" }}</h4>
            <p>{{ item.Description || item.description || "" }}</p>
            <span class="price">{{ formatPrice(item.Price || item.price) }}</span>
          </div>
        </RouterLink>
      </div>
      <p v-else class="empty">暂无商品</p>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, RouterLink } from "vue-router";
import { getShop, listProductsByShop } from "../../services/api.js";

const route = useRoute();
const shop = ref(null);
const products = ref([]);
const loading = ref(false);

const shopName = computed(() => shop.value?.Name || shop.value?.name || "店铺详情");
const shopDesc = computed(() => shop.value?.Description || shop.value?.description || "暂无简介");

const formatPrice = (price) => {
  if (price === undefined || price === null || price === "") {
    return "";
  }
  return `¥${price}`;
};

const getProductKey = (item, index) => item.ID || item.id || item.Name || item.name || index;
const getProductId = (item) => item.ID || item.id || null;

const load = async (id) => {
  if (!id) {
    return;
  }
  loading.value = true;
  try {
    shop.value = await getShop(id);
  } catch {
    shop.value = null;
  }
  try {
    const data = await listProductsByShop(id);
    products.value = Array.isArray(data) ? data : [];
  } catch {
    products.value = [];
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
  gap: 20px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.lead {
  color: var(--text-muted);
  margin: 8px 0 0;
}

.ghost {
  padding: 10px 16px;
  border-radius: 999px;
  border: 1px solid var(--surface-border);
  color: var(--text-primary);
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

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.product-card {
  display: grid;
  grid-template-columns: 52px 1fr;
  gap: 12px;
  padding: 12px;
  border-radius: 16px;
  border: 1px solid var(--surface-border);
  background: rgba(225, 37, 27, 0.04);
  color: inherit;
}

.thumb {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  background: rgba(225, 37, 27, 0.12);
}

.product-card h4 {
  margin: 0 0 4px;
  font-size: 1rem;
}

.product-card p {
  margin: 0 0 8px;
  color: var(--text-muted);
  font-size: 0.85rem;
}

.price {
  color: var(--accent-strong);
  font-weight: 600;
}

.empty {
  margin: 0;
  color: var(--text-muted);
}
</style>