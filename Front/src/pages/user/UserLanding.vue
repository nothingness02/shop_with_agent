<template>
  <section class="jd-top">
    <div class="top-left">中国大陆 · 江苏</div>
    <div class="top-right">
      <span>欢迎来到悦购</span>
      <template v-if="isAuthed">
        <RouterLink to="/user/home">我的账号</RouterLink>
      </template>
      <template v-else>
        <RouterLink to="/user/login">请登录</RouterLink>
        <RouterLink to="/user/register">免费注册</RouterLink>
      </template>
      <RouterLink to="/chat">客服</RouterLink>
    </div>
  </section>

  <section class="jd-header">
    <div class="logo">悦购</div>
    <div class="search-block">
      <div class="search-bar">
        <input type="text" placeholder="搜索商品 / 店铺" aria-label="搜索" />
        <button type="button">搜索</button>
      </div>
      <div class="hotwords">
        <span v-for="word in hotWords" :key="word">{{ word }}</span>
      </div>
    </div>
    <div class="promo-box">
      <div class="promo-title">新人礼包</div>
      <p>下单立减 + 专属补贴</p>
      <RouterLink class="promo-cta" to="/user/register">立即领取</RouterLink>
    </div>
  </section>

  <section class="jd-nav">
    <span v-for="item in navItems" :key="item">{{ item }}</span>
  </section>

  <section class="jd-main">
    <aside class="category">
      <h3>全部分类</h3>
      <ul>
        <li v-for="item in categories" :key="item.label">
          <span class="icon">{{ item.icon }}</span>
          <span>{{ item.label }}</span>
        </li>
      </ul>
    </aside>

    <div class="hero">
      <div class="hero-banner">
        <div>
          <p class="eyebrow">今日热卖</p>
          <h2>超级品类日 · 低至 5 折</h2>
          <p class="lead">精选爆款、品质直降、限时福利</p>
          <div class="actions">
            <template v-if="isAuthed">
              <RouterLink to="/user/home" class="primary">进入用户中心</RouterLink>
              <RouterLink to="/chat" class="ghost">联系客服</RouterLink>
            </template>
            <template v-else>
              <RouterLink to="/user/register" class="primary">立即注册</RouterLink>
              <RouterLink to="/user/login" class="ghost">登录查看</RouterLink>
            </template>
          </div>
        </div>
        <div class="hero-card">
          <h4>今日推荐</h4>
          <ul>
            <li>爆款 1 元购</li>
            <li>百亿补贴专区</li>
            <li>限时秒杀</li>
          </ul>
        </div>
      </div>
      <div class="hero-grid">
        <div class="hero-mini" v-for="item in miniPromos" :key="item.title">
          <span class="badge">{{ item.badge }}</span>
          <h4>{{ item.title }}</h4>
          <p>{{ item.desc }}</p>
        </div>
      </div>
    </div>

    <aside class="user-card">
      <div class="avatar">😊</div>
      <div class="user-info">
        <p class="muted">{{ isAuthed ? `已登录 · ID ${userId}` : "请登录享受更多权益" }}</p>
        <div class="user-actions">
          <RouterLink v-if="!isAuthed" to="/user/login">登录</RouterLink>
          <RouterLink v-if="!isAuthed" to="/user/register">注册</RouterLink>
          <RouterLink v-if="isAuthed" to="/user/home">订单中心</RouterLink>
        </div>
      </div>
      <div class="user-stats">
        <div>
          <strong>{{ shopCount }}</strong>
          <span>店铺</span>
        </div>
        <div>
          <strong>{{ productCount }}</strong>
          <span>商品</span>
        </div>
      </div>
    </aside>
  </section>

  <section class="jd-floor">
    <div class="floor-title">
      <h2>精选店铺</h2>
      <p v-if="shopsLoading">加载中...</p>
      <p v-else>来自后端实时数据</p>
    </div>
    <div class="shop-grid" v-if="shops.length">
      <RouterLink
        v-for="(shop, index) in shops"
        :key="getShopKey(shop, index)"
        :to="getShopId(shop) ? `/shops/${getShopId(shop)}` : '/user/home'"
        class="shop-card"
      >
        <div class="shop-avatar">🏬</div>
        <div>
          <h3>{{ shop.Name || shop.name || "未命名店铺" }}</h3>
          <p>{{ shop.Description || shop.description || "暂无描述" }}</p>
        </div>
      </RouterLink>
    </div>
    <p v-else class="empty">暂无店铺</p>
  </section>

  <section class="jd-floor">
    <div class="floor-title">
      <h2>精选商品</h2>
      <p v-if="productsLoading">加载中...</p>
      <p v-else>来自后端实时数据</p>
    </div>
    <div class="product-grid" v-if="products.length">
      <RouterLink
        v-for="(item, index) in products"
        :key="getProductKey(item, index)"
        :to="getProductId(item) ? `/products/${getProductId(item)}` : '/user/home'"
        class="product-card"
      >
        <div class="thumb"></div>
        <div>
          <h3>{{ item.Name || item.name || "未命名商品" }}</h3>
          <p>{{ item.Description || item.description || "" }}</p>
          <span class="price">{{ formatPrice(item.Price || item.price) }}</span>
        </div>
      </RouterLink>
    </div>
    <p v-else class="empty">暂无商品</p>
  </section>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import { listProductsByShop, listShops } from "../../services/api.js";
import { getAuth } from "../../services/storage.js";

const isAuthed = ref(false);
const userId = ref("");
const shops = ref([]);
const products = ref([]);
const shopsLoading = ref(false);
const productsLoading = ref(false);

const categories = [
  { label: "手机数码", icon: "📱" },
  { label: "家用电器", icon: "📺" },
  { label: "电脑办公", icon: "💻" },
  { label: "服饰美妆", icon: "👗" },
  { label: "食品生鲜", icon: "🥑" },
  { label: "家居生活", icon: "🛋️" },
  { label: "运动户外", icon: "🎽" },
  { label: "母婴玩具", icon: "🧸" },
  { label: "图书文娱", icon: "📚" },
  { label: "汽车用品", icon: "🚗" }
];

const navItems = [
  "秒杀",
  "超市",
  "家电",
  "服饰",
  "生鲜",
  "家具",
  "美妆",
  "数码",
  "运动",
  "母婴"
];

const hotWords = ["空调", "羽绒服", "手机", "笔记本", "零食"];

const miniPromos = [
  { badge: "闪购", title: "限时秒杀", desc: "爆款直降" },
  { badge: "新品", title: "新品首发", desc: "抢先体验" },
  { badge: "补贴", title: "百亿补贴", desc: "价格更低" }
];

const syncAuth = () => {
  const auth = getAuth();
  isAuthed.value = Boolean(auth && auth.access_token);
  userId.value = auth?.user_id || "";
};

const formatPrice = (price) => {
  if (price === undefined || price === null || price === "") {
    return "";
  }
  return `¥${price}`;
};

const getShopKey = (shop, index) => shop.ID || shop.id || shop.Name || shop.name || index;
const getShopId = (shop) => shop.ID || shop.id || null;
const getProductKey = (item, index) => item.ID || item.id || item.Name || item.name || index;
const getProductId = (item) => item.ID || item.id || null;

const shopCount = ref(0);
const productCount = ref(0);

const loadShops = async () => {
  shopsLoading.value = true;
  try {
    const data = await listShops();
    shops.value = Array.isArray(data) ? data : [];
    shopCount.value = shops.value.length;
  } catch {
    shops.value = [];
    shopCount.value = 0;
  } finally {
    shopsLoading.value = false;
  }
};

const loadProducts = async () => {
  productsLoading.value = true;
  const collected = [];
  try {
    const source = shops.value.slice(0, 3);
    for (const shop of source) {
      const id = getShopId(shop);
      if (!id) {
        continue;
      }
      const data = await listProductsByShop(id);
      if (Array.isArray(data)) {
        collected.push(...data.slice(0, 4));
      }
    }
  } catch {
    // ignore
  } finally {
    products.value = collected;
    productCount.value = collected.length;
    productsLoading.value = false;
  }
};

onMounted(async () => {
  syncAuth();
  window.addEventListener("storage", syncAuth);
  window.addEventListener("auth-changed", syncAuth);
  await loadShops();
  await loadProducts();
});

onBeforeUnmount(() => {
  window.removeEventListener("storage", syncAuth);
  window.removeEventListener("auth-changed", syncAuth);
});
</script>

<style scoped>
.jd-top {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  color: var(--text-muted);
}

.jd-top a {
  margin-left: 12px;
  color: inherit;
}

.jd-header {
  display: grid;
  grid-template-columns: 160px 1fr 200px;
  gap: 20px;
  align-items: center;
  padding: 18px 0 10px;
}

.logo {
  font-family: var(--font-display);
  font-size: 2.2rem;
  color: var(--accent-strong);
}

.search-block {
  display: grid;
  gap: 8px;
}

.search-bar {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 8px;
}

.search-bar input {
  padding: 12px 16px;
  border-radius: 999px;
  border: 1px solid var(--surface-border);
  background: #fff;
}

.search-bar button {
  padding: 12px 18px;
  border-radius: 999px;
  border: none;
  background: var(--accent-strong);
  color: var(--button-text);
  font-weight: 600;
}

.hotwords {
  display: flex;
  gap: 12px;
  font-size: 0.8rem;
  color: var(--text-muted);
}

.promo-box {
  padding: 14px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  text-align: center;
}

.promo-title {
  font-weight: 600;
}

.promo-cta {
  display: inline-block;
  margin-top: 10px;
  padding: 6px 12px;
  border-radius: 999px;
  background: var(--accent-strong);
  color: var(--button-text);
}

.jd-nav {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  padding: 10px 0 18px;
  color: var(--text-primary);
}

.jd-main {
  display: grid;
  grid-template-columns: 220px 1fr 240px;
  gap: 18px;
  align-items: start;
}

.category {
  padding: 16px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
}

.category ul {
  margin: 12px 0 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: 8px;
}

.category li {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
}

.hero {
  display: grid;
  gap: 16px;
}

.hero-banner {
  display: grid;
  grid-template-columns: 1.4fr 1fr;
  gap: 16px;
  padding: 20px;
  border-radius: 20px;
  background: linear-gradient(120deg, #fff4f3, #ffffff);
  border: 1px solid var(--surface-border);
}

.hero-card {
  padding: 16px;
  border-radius: 16px;
  background: #fff;
  border: 1px solid var(--surface-border);
}

.hero-card ul {
  margin: 10px 0 0;
  padding-left: 16px;
  color: var(--text-muted);
}

.hero-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(140px, 1fr));
  gap: 12px;
}

.hero-mini {
  padding: 12px;
  border-radius: 14px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
}

.badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  background: var(--accent-soft);
  color: var(--accent-strong);
  font-size: 0.75rem;
}

.eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.18em;
  font-size: 0.7rem;
  color: var(--text-muted);
}

.lead {
  margin: 8px 0 16px;
  color: var(--text-muted);
}

.actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.primary,
.ghost {
  padding: 10px 18px;
  border-radius: 999px;
  font-weight: 600;
}

.primary {
  background: var(--accent-strong);
  color: var(--button-text);
}

.ghost {
  border: 1px solid var(--surface-border);
  color: var(--text-primary);
}

.user-card {
  padding: 16px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  display: grid;
  gap: 12px;
  text-align: center;
}

.user-card .avatar {
  font-size: 2rem;
}

.user-actions {
  display: flex;
  justify-content: center;
  gap: 10px;
}

.user-stats {
  display: flex;
  justify-content: space-around;
}

.jd-floor {
  margin-top: 24px;
  display: grid;
  gap: 14px;
}

.floor-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.shop-grid,
.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 14px;
}

.shop-card,
.product-card {
  display: grid;
  grid-template-columns: 60px 1fr;
  gap: 12px;
  padding: 12px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--surface-border);
  color: inherit;
}

.shop-avatar,
.thumb {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  background: rgba(225, 37, 27, 0.12);
  display: grid;
  place-items: center;
}

.product-card h3,
.shop-card h3 {
  margin: 0 0 4px;
  font-size: 1rem;
}

.product-card p,
.shop-card p {
  margin: 0 0 6px;
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

.muted {
  color: var(--text-muted);
}

@media (max-width: 1100px) {
  .jd-header {
    grid-template-columns: 1fr;
  }

  .jd-main {
    grid-template-columns: 1fr;
  }

  .hero-banner {
    grid-template-columns: 1fr;
  }
}
</style>
