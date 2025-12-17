<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

const API_BASE = 'http://localhost:8080';
const VIEWS = {
  HOME: 'home',
  DETAIL: 'product_detail',
  CART: 'cart',
  LOGIN: 'login',
  ORDERS: 'orders',
  MERCHANT: 'merchant',
} as const;
type View = typeof VIEWS[keyof typeof VIEWS];

const ui = reactive({
  view: VIEWS.HOME as View,
  loading: false,
  submitting: false,
  search: '',
  toast: { show: false, message: '' },
});

const session = reactive({
  user: loadUserFromStorage(),
  token: localStorage.getItem('token') || '',
});

const shops = ref<any[]>([]);
const products = ref<any[]>([]);
const activeProduct = ref<any | null>(null);
const cart = ref<{ product: any; quantity: number }[]>([]);
const orders = ref<any[]>([]);
const merchant = reactive({
  shops: [] as any[],
  products: [] as any[],
  selectedShop: null as any,
  loading: false,
  creating: false,
  shopForm: { name: '', description: '' },
  productForm: { name: '', description: '', price: '', stock: 0, product_img: '' },
});

const isAuthed = computed(() => !!session.token);
const isMerchant = computed(() => session.user?.role === 5 || session.user?.role === 10);
const cartTotalCount = computed(() => cart.value.reduce((acc, item) => acc + item.quantity, 0));
const cartTotalPrice = computed(() =>
  cart.value.reduce((acc, item) => acc + Number(item.product.Price || 0) * item.quantity, 0),
);
const merchantStock = computed(() =>
  merchant.products.reduce((acc, item) => acc + Number(item.Stock || 0), 0),
);
const merchantCatalogValue = computed(() =>
  merchant.products
    .reduce((acc, item) => acc + Number(item.Price || 0) * Number(item.Stock || 0), 0)
    .toFixed(2),
);

onMounted(() => {
  if (isAuthed.value) {
    hydrateHome();
  } else {
    changeView(VIEWS.LOGIN);
  }
});

function loadUserFromStorage() {
  try {
    const stored = localStorage.getItem('user');
    return stored ? JSON.parse(stored) : null;
  } catch (err) {
    console.warn('Failed to parse user from storage', err);
    return null;
  }
}

function changeView(view: View) {
  if (view === VIEWS.MERCHANT && !guardMerchant()) return;
  ui.view = view;
  window.scrollTo({ top: 0, behavior: 'smooth' });
  if (view === VIEWS.ORDERS && isAuthed.value) {
    fetchOrders();
  }
  if (view === VIEWS.MERCHANT) {
    loadMerchantShops();
  }
}

function showToast(message: string) {
  ui.toast.message = message;
  ui.toast.show = true;
  setTimeout(() => (ui.toast.show = false), 2200);
}

function guardAuth(message = '请先登录后再继续') {
  if (!isAuthed.value) {
    changeView(VIEWS.LOGIN);
    showToast(message);
    return false;
  }
  return true;
}

function guardMerchant(message = 'Merchant access only') {
  if (!guardAuth(message)) return false;
  if (!isMerchant.value) {
    showToast(message);
    return false;
  }
  return true;
}

async function apiCall(endpoint: string, method = 'GET', body?: any) {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (session.token) {
    headers.Authorization = `Bearer ${session.token}`;
  }

  try {
    const res = await fetch(`${API_BASE}${endpoint}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : null,
    });

    if (res.status === 401) {
      logout();
      showToast('登录失效，请重新登录');
      return null;
    }

    const data = await res.json();
    if (!res.ok) throw new Error(data.error || '请求失败');
    return data;
  } catch (err: any) {
    console.error(err);
    showToast(err.message || '网络异常');
    return null;
  }
}

// --------- Data loaders ---------
async function hydrateHome() {
  ui.loading = true;
  await loadShops();
  await loadInitialFeed();
  ui.loading = false;
}

async function loadShops() {
  if (!guardAuth()) return;
  const data = await apiCall('/api/v2/shops');
  if (data) shops.value = data;
}

async function loadInitialFeed() {
  if (!guardAuth() || shops.value.length === 0) return;
  const shopId = shops.value[0].ID;
  const data = await apiCall(`/api/v2/shops/${shopId}/products`);
  if (data) products.value = data;
}

async function performSearch() {
  if (!ui.search) return;
  if (!guardAuth()) return;
  ui.loading = true;
  const data = await apiCall(`/api/v2/search/products?q=${encodeURIComponent(ui.search)}`);
  if (data) {
    products.value = data;
    changeView(VIEWS.HOME);
  }
  ui.loading = false;
}

async function selectShop(shop: any) {
  if (!guardAuth()) return;
  ui.loading = true;
  const data = await apiCall(`/api/v2/shops/${shop.ID}/products`);
  if (data) products.value = data;
  ui.loading = false;
  showToast(`已切换至 ${shop.Name}`);
}

function showProductDetail(product: any) {
  activeProduct.value = product;
  changeView(VIEWS.DETAIL);
}

function handleImageError(e: Event) {
  const target = e.target as HTMLImageElement;
  target.src = 'https://via.placeholder.com/360?text=JD+Shop';
}

// --------- Merchant dashboard ---------
async function loadMerchantShops() {
  if (!guardMerchant()) return;
  merchant.loading = true;
  const data = await apiCall('/api/v2/shops');
  if (data) {
    const owned = data.filter((shop: any) => shop.OwnerID === session.user?.ID);
    merchant.shops = owned;

    if (owned.length > 0) {
      const alreadySelected = owned.find((s: any) => s.ID === merchant.selectedShop?.ID);
      merchant.selectedShop = alreadySelected || owned[0];
      await loadMerchantProducts(merchant.selectedShop.ID);
    } else {
      merchant.selectedShop = null;
      merchant.products = [];
    }
  }
  merchant.loading = false;
}

async function loadMerchantProducts(shopId: number) {
  if (!guardMerchant()) return;
  merchant.loading = true;
  const data = await apiCall(`/api/v2/shops/${shopId}/products`);
  if (data) merchant.products = data;
  merchant.loading = false;
}

function selectMerchantShop(shop: any) {
  merchant.selectedShop = shop;
  loadMerchantProducts(shop.ID);
}

async function createShop() {
  if (!guardMerchant()) return;
  if (!merchant.shopForm.name.trim()) {
    showToast('请输入店铺名称');
    return;
  }
  merchant.creating = true;
  const payload = {
    name: merchant.shopForm.name,
    description: merchant.shopForm.description,
    owner_id: session.user?.ID,
  };
  const res = await apiCall('/api/v2/shops', 'POST', payload);
  if (res) {
    merchant.shops.unshift(res);
    merchant.selectedShop = res;
    merchant.shopForm = { name: '', description: '' };
    await loadMerchantProducts(res.ID);
    showToast('店铺已创建');
  }
  merchant.creating = false;
}

async function createProduct() {
  if (!guardMerchant()) return;
  if (!merchant.selectedShop) {
    showToast('请先选择店铺');
    return;
  }
  if (!merchant.productForm.name || !merchant.productForm.price) {
    showToast('商品名称与价格为必填');
    return;
  }
  merchant.creating = true;
  const payload = {
    name: merchant.productForm.name,
    description: merchant.productForm.description,
    price: Number(merchant.productForm.price),
    stock: Number(merchant.productForm.stock || 0),
    product_img: merchant.productForm.product_img,
  };
  const res = await apiCall(`/api/v2/shops/${merchant.selectedShop.ID}/products`, 'POST', payload);
  if (res) {
    merchant.products.unshift(res);
    merchant.productForm = { name: '', description: '', price: '', stock: 0, product_img: '' };
    showToast('商品已上架');
  }
  merchant.creating = false;
}

// --------- Cart ---------
function addToCart(product: any) {
  if (!guardAuth()) return;
  const existing = cart.value.find((i) => i.product.ID === product.ID);
  if (existing) {
    existing.quantity += 1;
  } else {
    cart.value.push({ product, quantity: 1 });
  }
  showToast('已加入购物车');
}

function buyNow(product: any) {
  addToCart(product);
  changeView(VIEWS.CART);
}

function updateQuantity(item: { product: any; quantity: number }, delta: number) {
  item.quantity += delta;
  if (item.quantity <= 0) {
    cart.value = cart.value.filter((i) => i !== item);
  }
}

function clearCart() {
  cart.value = [];
}

// --------- Orders ---------
async function submitOrder() {
  if (cart.value.length === 0 || !guardAuth()) return;
  ui.submitting = true;

  const payload = {
    user_id: session.user?.ID,
    items: cart.value.map((item) => ({
      product_id: item.product.ID,
      product_name: item.product.Name,
      product_img: item.product.ProductImg,
      price: item.product.Price,
      quantity: item.quantity,
    })),
  };

  const res = await apiCall('/api/v1/orders', 'POST', payload);
  if (res) {
    showToast('下单成功');
    cart.value = [];
    changeView(VIEWS.ORDERS);
  }
  ui.submitting = false;
}

async function fetchOrders() {
  if (!guardAuth()) return;
  ui.loading = true;
  const data = await apiCall('/api/v1/orders');
  if (data) orders.value = data;
  ui.loading = false;
}

function calculateOrderTotal(order: any) {
  if (order.TotalAmount) return order.TotalAmount;
  if (order.Items) {
    return order.Items.reduce(
      (acc: number, item: any) => acc + Number(item.Price || 0) * item.Quantity,
      0,
    ).toFixed(2);
  }
  return '0.00';
}

// --------- Auth ---------
const isRegister = ref(false);
const authForm = reactive({ username: '', password: '', email: '', role: 1 });

async function handleAuth() {
  if (isRegister.value) {
    const res = await apiCall('/api/v0/users', 'POST', authForm);
    if (res) {
      showToast('注册成功，请登录');
      isRegister.value = false;
    }
    return;
  }

  const res = await apiCall('/api/v0/auth/login', 'POST', {
    username: authForm.username,
    password: authForm.password,
  });

  if (res && res.access_token) {
    session.token = res.access_token;
    const userData = { ID: res.user_id, username: authForm.username, role: res.role };
    session.user = userData;
    localStorage.setItem('token', session.token);
    localStorage.setItem('user', JSON.stringify(userData));
    showToast('登录成功');
    changeView(VIEWS.HOME);
    hydrateHome();
  }
}

function logout() {
  session.user = null;
  session.token = '';
  cart.value = [];
  localStorage.removeItem('user');
  localStorage.removeItem('token');
  changeView(VIEWS.LOGIN);
}
</script>

<template>
  <div class="relative min-h-screen bg-gradient-to-br from-[#fef2f2] via-white to-[#fff4f1]">
    <div class="pointer-events-none absolute inset-0 opacity-50">
      <div class="absolute -left-20 top-10 h-64 w-64 rounded-full bg-[#ffe1e7] blur-3xl"></div>
      <div class="absolute -right-10 bottom-10 h-72 w-72 rounded-full bg-[#ffe2cc] blur-3xl"></div>
    </div>

    <!-- HEADER -->
    <header class="sticky top-0 z-40 border-b border-white/60 bg-white/70 backdrop-blur">
      <div class="mx-auto flex max-w-6xl items-center gap-3 px-4 py-3 md:px-8">
        <button
          class="flex items-center gap-2 rounded-full bg-[#f10215] px-3 py-1 text-white shadow-lg shadow-red-200"
          @click="changeView(VIEWS.HOME)"
        >
          <span class="text-lg font-extrabold tracking-tight">JD</span>
          <span class="hidden text-xs font-semibold uppercase md:block">Smart Shop</span>
        </button>

        <div class="flex flex-1 items-center gap-2 rounded-full bg-white px-4 py-2 shadow-sm ring-1 ring-gray-100">
          <i class="fas fa-search text-gray-400"></i>
          <input
            v-model="ui.search"
            @keyup.enter="performSearch"
            type="text"
            placeholder="搜索商品 / 店铺 / 关键字"
            class="w-full bg-transparent text-sm text-gray-700 outline-none"
          />
          <button
            v-if="ui.search"
            @click="performSearch"
            class="rounded-full bg-[#f10215] px-3 py-1 text-xs font-semibold text-white transition hover:brightness-95"
          >
            搜索
          </button>
        </div>

        <div class="flex items-center gap-3 text-sm">
          <button
            v-if="!session.user"
            class="rounded-full border border-transparent bg-black text-white px-3 py-1.5 transition hover:border-black hover:bg-white hover:text-black"
            @click="changeView(VIEWS.LOGIN)"
          >
            登录
          </button>
          <div v-else class="relative group">
            <button class="flex items-center gap-2 rounded-full bg-white px-3 py-1.5 shadow-sm ring-1 ring-gray-100">
              <i class="fas fa-user text-gray-500"></i>
              <span class="hidden md:block text-gray-700">{{ session.user.username }}</span>
            </button>
            <div
              class="absolute right-0 mt-2 hidden w-40 rounded-xl bg-white p-2 text-sm text-gray-700 shadow-lg ring-1 ring-gray-100 md:group-hover:block"
            >
              <div
                v-if="isMerchant"
                class="cursor-pointer rounded-lg px-3 py-2 hover:bg-gray-50"
                @click="changeView(VIEWS.MERCHANT)"
              >
                商家中心
              </div>
              <div class="cursor-pointer rounded-lg px-3 py-2 hover:bg-gray-50" @click="changeView(VIEWS.ORDERS)">
                我的订单
              </div>
              <div class="cursor-pointer rounded-lg px-3 py-2 text-[#f10215] hover:bg-gray-50" @click="logout">
                退出登录
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-6xl px-3 pb-24 pt-4 md:px-8 md:pb-12">
      <!-- HOME -->
      <section v-if="ui.view === VIEWS.HOME" class="space-y-6">
        <div class="relative overflow-hidden rounded-2xl bg-gradient-to-r from-[#f10215] to-[#ff7a3d] p-6 text-white shadow-xl">
          <div class="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(255,255,255,0.18)_0,_transparent_45%)]"></div>
          <div class="relative flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <p class="text-sm uppercase tracking-[0.25em] opacity-80">AI Shop</p>
              <h1 class="mt-2 text-3xl font-extrabold leading-tight md:text-4xl">智能优选，惊喜直达</h1>
              <p class="mt-2 text-sm md:text-base md:opacity-90">一键发现爆款，极速送达，逛店更轻松</p>
              <div class="mt-4 flex gap-3">
                <button
                  class="rounded-full bg-white px-4 py-2 text-sm font-semibold text-[#f10215] shadow-lg shadow-red-200 transition hover:brightness-95"
                  @click="ui.search ? performSearch() : hydrateHome()"
                >
                  刷新推荐
                </button>
                <button
                  class="rounded-full border border-white/60 px-4 py-2 text-sm font-semibold transition hover:bg-white/10"
                  @click="changeView(VIEWS.CART)"
                >
                  查看购物车
                </button>
              </div>
            </div>
            <div class="flex gap-4 md:gap-6">
              <div class="rounded-2xl bg-white/15 px-5 py-4 text-center backdrop-blur">
                <div class="text-2xl font-bold">{{ shops.length || '—' }}</div>
                <p class="text-xs opacity-80">精选店铺</p>
              </div>
              <div class="rounded-2xl bg-white/15 px-5 py-4 text-center backdrop-blur">
                <div class="text-2xl font-bold">{{ products.length || '—' }}</div>
                <p class="text-xs opacity-80">热销商品</p>
              </div>
            </div>
          </div>
        </div>

        <div v-if="shops.length" class="space-y-2 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="h-6 w-0.5 rounded-full bg-[#f10215]"></span>
              <h3 class="text-lg font-bold text-gray-800">推荐店铺</h3>
            </div>
            <span class="text-xs text-gray-500">可快速切换热卖商品</span>
          </div>
          <div class="no-scrollbar flex gap-3 overflow-x-auto pt-1">
            <button
              v-for="shop in shops"
              :key="shop.ID"
              class="group relative min-w-[140px] flex-shrink-0 rounded-xl bg-gradient-to-br from-white to-[#fff7f5] p-3 text-left shadow-sm ring-1 ring-gray-100 transition hover:-translate-y-0.5 hover:shadow-md"
              @click="selectShop(shop)"
            >
              <div class="flex items-center gap-3">
                <div
                  class="flex h-10 w-10 items-center justify-center rounded-full bg-white text-lg text-[#f10215] ring-1 ring-red-100 shadow-sm"
                >
                  <i class="fas fa-store"></i>
                </div>
                <div class="flex-1">
                  <p class="truncate text-sm font-semibold text-gray-800">{{ shop.Name }}</p>
                  <p class="text-[11px] text-gray-500">{{ shop.Description || '优质好店' }}</p>
                </div>
              </div>
              <div class="absolute bottom-1 right-2 text-[10px] text-[#f10215] opacity-80">点此切换</div>
            </button>
          </div>
        </div>

        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="h-6 w-0.5 rounded-full bg-[#f10215]"></span>
              <h3 class="text-lg font-bold text-gray-800">猜你喜欢</h3>
            </div>
            <span class="text-xs text-gray-500">实时匹配热度与销量</span>
          </div>

          <div v-if="ui.loading" class="flex justify-center py-12">
            <div class="loader" />
          </div>

          <div v-else class="grid grid-cols-2 gap-3 md:grid-cols-3 lg:grid-cols-4">
            <article
              v-for="product in products"
              :key="product.ID"
              class="group flex cursor-pointer flex-col overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-gray-100 transition hover:-translate-y-1 hover:shadow-lg"
              @click="showProductDetail(product)"
            >
              <div class="relative aspect-square bg-gradient-to-br from-gray-50 to-gray-100">
                <img
                  :src="product.ProductImg || 'https://via.placeholder.com/400x400?text=No+Image'"
                  class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                  @error="handleImageError"
                />
                <span
                  class="absolute left-2 top-2 rounded-full bg-white/90 px-2 py-1 text-[10px] font-semibold text-[#f10215] shadow-sm"
                  >自营 · 极速达</span
                >
              </div>
              <div class="flex flex-1 flex-col gap-2 p-3">
                <h4 class="line-clamp-2 text-sm font-semibold text-gray-900">{{ product.Name }}</h4>
                <p class="line-clamp-1 text-[11px] text-gray-500">{{ product.Description || '品质保障' }}</p>
                <div class="mt-auto flex items-center justify-between pt-1">
                  <span class="text-lg font-bold text-[#f10215]">¥{{ product.Price }}</span>
                  <button
                    class="rounded-full bg-gray-100 p-2 text-[#f10215] transition hover:bg-[#f10215] hover:text-white"
                    @click.stop="addToCart(product)"
                  >
                    <i class="fas fa-cart-plus"></i>
                  </button>
                </div>
              </div>
            </article>
          </div>
        </div>
      </section>

      <!-- PRODUCT DETAIL -->
      <section
        v-if="ui.view === VIEWS.DETAIL && activeProduct"
        class="mt-4 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100 md:mt-8 md:flex md:gap-8 md:p-8"
      >
        <div class="relative w-full md:w-1/2">
          <button
            class="absolute left-2 top-2 flex h-10 w-10 items-center justify-center rounded-full bg-black/40 text-white md:hidden"
            @click="changeView(VIEWS.HOME)"
          >
            <i class="fas fa-arrow-left"></i>
          </button>
          <div class="overflow-hidden rounded-2xl bg-gray-50">
            <img
              :src="activeProduct.ProductImg || 'https://via.placeholder.com/600?text=Product'"
              class="h-full w-full object-contain"
            />
          </div>
        </div>

        <div class="mt-4 flex flex-1 flex-col gap-4 md:mt-0">
          <div>
            <div class="flex items-start justify-between">
              <p class="text-3xl font-bold text-[#f10215]">¥{{ activeProduct.Price }}</p>
              <span class="text-xs text-gray-500">库存：{{ activeProduct.Stock ?? '—' }}</span>
            </div>
            <h2 class="mt-2 text-xl font-bold text-gray-900">{{ activeProduct.Name }}</h2>
            <p class="mt-2 rounded-xl bg-gray-50 p-3 text-sm text-gray-600">
              {{ activeProduct.Description || '暂无详细描述' }}
            </p>
          </div>

          <div class="space-y-2 text-xs text-gray-600">
            <div class="flex items-center gap-2">
              <span class="font-semibold text-gray-500">发货</span>
              <span class="text-gray-900">京东物流</span>
              <span class="text-gray-400">|</span>
              <span>次日达</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="font-semibold text-gray-500">服务</span>
              <span>7天无理由 · 正品保障</span>
            </div>
          </div>

          <div class="sticky bottom-4 mt-auto flex flex-col gap-3 md:static md:flex-row">
            <button
              class="flex-1 rounded-full bg-amber-400 px-4 py-3 text-center text-white shadow-md transition hover:brightness-95"
              @click="addToCart(activeProduct)"
            >
              加入购物车
            </button>
            <button
              class="flex-1 rounded-full bg-[#f10215] px-4 py-3 text-center text-white shadow-md transition hover:brightness-95"
              @click="buyNow(activeProduct)"
            >
              立即购买
            </button>
          </div>
        </div>
      </section>

      <!-- CART -->
      <section v-if="ui.view === VIEWS.CART" class="mt-4 space-y-4 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100 md:mt-8">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold text-gray-900">购物车 ({{ cartTotalCount }})</h2>
          <button v-if="cart.length" class="text-xs text-gray-500 hover:text-gray-700" @click="clearCart">清空</button>
        </div>

        <div v-if="!cart.length" class="flex flex-col items-center gap-3 py-20 text-gray-400">
          <i class="fas fa-shopping-cart text-4xl"></i>
          <p>购物车空空如也</p>
          <button
            class="rounded-full border border-gray-200 px-4 py-2 text-sm text-gray-700 hover:border-gray-300"
            @click="changeView(VIEWS.HOME)"
          >
            去逛逛
          </button>
        </div>

        <div v-else class="space-y-3">
          <article
            v-for="(item, index) in cart"
            :key="index"
            class="flex gap-4 rounded-xl bg-gradient-to-r from-white to-[#fff8f5] p-3 shadow-sm ring-1 ring-gray-100"
          >
            <div class="flex h-6 w-6 items-center justify-center text-[#f10215]">
              <i class="fas fa-check-circle"></i>
            </div>
            <div class="h-20 w-20 overflow-hidden rounded-xl bg-gray-50 ring-1 ring-gray-100">
              <img :src="item.product.ProductImg" class="h-full w-full object-cover" @error="handleImageError" />
            </div>
            <div class="flex flex-1 flex-col justify-between">
              <h4 class="line-clamp-2 text-sm font-semibold text-gray-900">{{ item.product.Name }}</h4>
              <div class="flex items-end justify-between">
                <span class="text-lg font-bold text-[#f10215]">¥{{ item.product.Price }}</span>
                <div class="flex items-center rounded-full border border-gray-200">
                  <button class="px-3 py-1 text-gray-500" @click="updateQuantity(item, -1)">-</button>
                  <span class="w-8 text-center text-sm">{{ item.quantity }}</span>
                  <button class="px-3 py-1 text-gray-500" @click="updateQuantity(item, 1)">+</button>
                </div>
              </div>
            </div>
          </article>
        </div>

        <div
          v-if="cart.length"
          class="sticky bottom-4 flex items-center justify-between rounded-full bg-white px-4 py-3 shadow-lg ring-1 ring-gray-100 md:static md:rounded-xl md:px-6"
        >
          <div>
            <span class="text-sm text-gray-500">合计</span>
            <span class="ml-2 text-xl font-bold text-[#f10215]">¥{{ cartTotalPrice.toFixed(2) }}</span>
          </div>
          <button
            class="rounded-full bg-[#f10215] px-6 py-2 text-sm font-semibold text-white transition hover:brightness-95 disabled:opacity-50"
            :disabled="ui.submitting"
            @click="submitOrder"
          >
            {{ ui.submitting ? '提交中...' : '去结算' }}
          </button>
        </div>
      </section>

      <!-- MERCHANT -->
      <section v-if="ui.view === VIEWS.MERCHANT" class="mt-6 space-y-6">
        <div class="relative overflow-hidden rounded-2xl bg-gradient-to-r from-[#0f172a] via-[#c80b1a] to-[#ff7a3d] p-6 text-white shadow-xl">
          <div class="absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,rgba(255,255,255,0.16),transparent_45%)]"></div>
          <div class="relative flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <div>
              <p class="text-xs uppercase tracking-[0.3em] text-white/80">Merchant Desk</p>
              <h2 class="mt-2 text-3xl font-extrabold leading-tight">经营面板，轻松管理店铺</h2>
              <p class="mt-2 max-w-2xl text-sm text-white/80">店铺、商品、库存都在同一处完成，保持与消费者端一致的轻量体验。</p>
            </div>
            <div class="flex items-center gap-3">
              <div class="rounded-xl bg-white/15 px-3 py-2 text-sm">店铺 {{ merchant.shops.length }}</div>
              <div
                v-if="merchant.selectedShop"
                class="rounded-xl bg-white px-4 py-2 text-sm font-semibold text-[#c80b1a] shadow-md shadow-red-200"
              >
                {{ merchant.selectedShop.Name }}
              </div>
            </div>
          </div>
        </div>

        <div class="grid gap-4 md:grid-cols-3">
          <div class="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100">
            <div class="text-xs text-gray-500">店铺数量</div>
            <div class="mt-2 flex items-end justify-between">
              <span class="text-2xl font-bold text-gray-900">{{ merchant.shops.length }}</span>
              <span class="rounded-full bg-[#f10215]/10 px-3 py-1 text-xs text-[#f10215]">Owner #{{ session.user?.ID }}</span>
            </div>
          </div>
          <div class="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100">
            <div class="text-xs text-gray-500">商品数</div>
            <div class="mt-2 flex items-end justify-between">
              <span class="text-2xl font-bold text-gray-900">{{ merchant.products.length }}</span>
              <span class="rounded-full bg-[#0ea5e9]/10 px-3 py-1 text-xs text-[#0ea5e9]">在售</span>
            </div>
          </div>
          <div class="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100">
            <div class="text-xs text-gray-500">库存价值 (估算)</div>
            <div class="mt-2 flex items-end justify-between">
              <span class="text-2xl font-bold text-gray-900">￥{{ merchantCatalogValue }}</span>
              <span class="rounded-full bg-[#16a34a]/10 px-3 py-1 text-xs text-[#15803d]">库存 {{ merchantStock }}</span>
            </div>
          </div>
        </div>

        <div class="grid gap-6 lg:grid-cols-3">
          <div class="space-y-4 lg:col-span-2">
            <div class="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100">
              <div class="flex flex-wrap items-center justify-between gap-3">
                <div class="flex items-center gap-2">
                  <span class="h-6 w-0.5 rounded-full bg-[#f10215]"></span>
                  <div>
                    <h3 class="text-lg font-bold text-gray-800">我的店铺</h3>
                    <p class="text-xs text-gray-500">选择店铺查看商品，刷新同步最新数据</p>
                  </div>
                </div>
                <button
                  class="rounded-full border border-white/0 bg-[#f10215] px-4 py-2 text-xs font-semibold text-white shadow-sm transition hover:brightness-95 disabled:opacity-60"
                  :disabled="merchant.loading"
                  @click="loadMerchantShops"
                >
                  {{ merchant.loading ? '同步中...' : '刷新店铺' }}
                </button>
              </div>

              <div v-if="merchant.shops.length" class="mt-3 grid grid-cols-1 gap-3 md:grid-cols-2">
                <button
                  v-for="shop in merchant.shops"
                  :key="shop.ID"
                  class="group flex items-start gap-3 rounded-2xl bg-gradient-to-br from-white to-[#fff7f5] p-4 text-left shadow-sm ring-1 ring-gray-100 transition hover:-translate-y-0.5 hover:shadow-md"
                  :class="{ 'ring-2 ring-[#f10215] shadow-lg shadow-red-100': merchant.selectedShop?.ID === shop.ID }"
                  @click="selectMerchantShop(shop)"
                >
                  <div class="flex h-12 w-12 items-center justify-center rounded-full bg-white text-lg text-[#f10215] ring-1 ring-red-100 shadow-sm">
                    <i class="fas fa-store"></i>
                  </div>
                  <div class="flex-1">
                    <div class="flex items-center justify-between gap-2">
                      <p class="truncate text-base font-semibold text-gray-900">{{ shop.Name }}</p>
                      <span class="rounded-full bg-gray-100 px-2 py-1 text-[10px] text-gray-600">ID #{{ shop.ID }}</span>
                    </div>
                    <p class="mt-1 line-clamp-2 text-xs text-gray-500">{{ shop.Description || '暂无描述' }}</p>
                  </div>
                </button>
              </div>
              <div v-else class="mt-4 rounded-xl bg-gray-50 p-6 text-center text-sm text-gray-500">
                还没有店铺，先创建一个吧。
              </div>
            </div>

            <div class="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100">
              <div class="flex flex-wrap items-center justify-between gap-3">
                <div>
                  <h3 class="text-lg font-bold text-gray-800">商品清单</h3>
                  <p class="text-xs text-gray-500">
                    {{ merchant.selectedShop ? `正在查看 ${merchant.selectedShop.Name}` : '请选择店铺后查看商品' }}
                  </p>
                </div>
                <span class="rounded-full bg-gray-100 px-3 py-1 text-xs text-gray-600">{{ merchant.products.length }} 件</span>
              </div>

              <div v-if="merchant.loading" class="flex justify-center py-10">
                <div class="loader" />
              </div>

              <div
                v-else-if="merchant.selectedShop && merchant.products.length"
                class="mt-3 divide-y divide-gray-100 rounded-xl bg-white ring-1 ring-gray-50"
              >
                <article v-for="product in merchant.products" :key="product.ID" class="flex gap-4 p-3">
                  <div class="h-16 w-16 overflow-hidden rounded-xl bg-gray-50 ring-1 ring-gray-100">
                    <img
                      :src="product.ProductImg || 'https://via.placeholder.com/120?text=Product'"
                      class="h-full w-full object-cover"
                      @error="handleImageError"
                    />
                  </div>
                  <div class="flex flex-1 flex-col gap-1">
                    <div class="flex items-start justify-between gap-2">
                      <p class="text-base font-semibold text-gray-900">{{ product.Name }}</p>
                      <span class="rounded-full bg-[#f10215]/10 px-3 py-1 text-xs font-semibold text-[#f10215]">￥{{ product.Price }}</span>
                    </div>
                    <p class="text-xs text-gray-500 line-clamp-2">{{ product.Description || '暂无描述' }}</p>
                    <div class="flex items-center gap-4 text-[11px] text-gray-500">
                      <span>库存 {{ product.Stock ?? 0 }}</span>
                      <span>ID #{{ product.ID }}</span>
                      <span class="rounded-full bg-gray-100 px-2 py-0.5 text-[10px] text-gray-600">Shop #{{ product.ShopID }}</span>
                    </div>
                  </div>
                </article>
              </div>
              <div v-else class="mt-4 rounded-xl bg-gray-50 p-6 text-center text-sm text-gray-500">
                {{ merchant.selectedShop ? '暂无商品，上架一件试试。' : '先选择店铺或新建店铺后再上架商品。' }}
              </div>
            </div>
          </div>

          <div class="space-y-4">
            <div class="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100">
              <div class="flex items-center justify-between">
                <h3 class="text-base font-semibold text-gray-900">新建店铺</h3>
                <span class="text-xs text-gray-400">1 分钟开店</span>
              </div>
              <form class="mt-3 space-y-3" @submit.prevent="createShop">
                <div>
                  <label class="text-xs text-gray-500">店铺名称</label>
                  <input
                    v-model="merchant.shopForm.name"
                    type="text"
                    class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
                    placeholder="输入店铺名"
                  />
                </div>
                <div>
                  <label class="text-xs text-gray-500">描述</label>
                  <textarea
                    v-model="merchant.shopForm.description"
                    rows="2"
                    class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
                    placeholder="主营品类 / 服务亮点"
                  ></textarea>
                </div>
                <button
                  type="submit"
                  class="flex w-full items-center justify-center rounded-full bg-[#f10215] px-4 py-2 text-sm font-semibold text-white transition hover:brightness-95 disabled:opacity-60"
                  :disabled="merchant.creating"
                >
                  {{ merchant.creating ? '创建中...' : '创建店铺' }}
                </button>
              </form>
            </div>

            <div class="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100">
              <div class="flex items-center justify-between">
                <h3 class="text-base font-semibold text-gray-900">快速上架</h3>
                <span class="text-xs text-gray-400">绑定当前店铺</span>
              </div>
              <form class="mt-3 space-y-3" @submit.prevent="createProduct">
                <div class="rounded-xl bg-gray-50 px-3 py-2 text-xs text-gray-600">
                  当前店铺：{{ merchant.selectedShop ? merchant.selectedShop.Name : '请选择店铺' }}
                </div>
                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <label class="text-xs text-gray-500">商品名称</label>
                    <input
                      v-model="merchant.productForm.name"
                      type="text"
                      class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
                      placeholder="新品标题"
                    />
                  </div>
                  <div>
                    <label class="text-xs text-gray-500">价格</label>
                    <input
                      v-model="merchant.productForm.price"
                      type="number"
                      step="0.01"
                      class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
                      placeholder="如 199.00"
                    />
                  </div>
                  <div>
                    <label class="text-xs text-gray-500">库存</label>
                    <input
                      v-model="merchant.productForm.stock"
                      type="number"
                      class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
                      placeholder="100"
                    />
                  </div>
                  <div>
                    <label class="text-xs text-gray-500">封面图 URL</label>
                    <input
                      v-model="merchant.productForm.product_img"
                      type="url"
                      class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
                      placeholder="https://..."
                    />
                  </div>
                </div>
                <div>
                  <label class="text-xs text-gray-500">卖点 / 描述</label>
                  <textarea
                    v-model="merchant.productForm.description"
                    rows="2"
                    class="mt-1 w-full rounded-xl border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
                    placeholder="一句话吸引买家"
                  ></textarea>
                </div>
                <button
                  type="submit"
                  class="flex w-full items-center justify-center rounded-full bg-[#111827] px-4 py-2 text-sm font-semibold text-white transition hover:brightness-110 disabled:opacity-60"
                  :disabled="merchant.creating || !merchant.selectedShop"
                >
                  {{ merchant.creating ? '上架中...' : '发布到当前店铺' }}
                </button>
              </form>
            </div>
          </div>
        </div>
      </section>

      <!-- LOGIN / REGISTER -->
      <section
        v-if="ui.view === VIEWS.LOGIN"
        class="mt-6 flex min-h-[70vh] items-center justify-center md:mt-12"
      >
        <div class="w-full max-w-md rounded-2xl bg-white p-6 shadow-lg ring-1 ring-gray-100 md:p-8">
          <h2 class="text-center text-2xl font-bold text-gray-900">
            {{ isRegister ? '注册新账号' : '欢迎登录' }}
          </h2>
          <p class="mt-2 text-center text-sm text-gray-500">
            访问受保护的店铺、搜索与下单接口需先登录
          </p>

          <div class="mt-6 space-y-4">
            <div>
              <label class="block text-sm text-gray-600">用户名</label>
              <input
                v-model="authForm.username"
                type="text"
                class="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm outline-none ring-0 focus:border-[#f10215]"
              />
            </div>
            <div v-if="isRegister">
              <label class="block text-sm text-gray-600">邮箱</label>
              <input
                v-model="authForm.email"
                type="email"
                class="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
              />
            </div>
            <div>
              <label class="block text-sm text-gray-600">密码</label>
              <input
                v-model="authForm.password"
                type="password"
                class="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
              />
            </div>
            <div v-if="isRegister">
              <label class="block text-sm text-gray-600">角色 (1: 用户, 10: 管理员)</label>
              <select
                v-model="authForm.role"
                class="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm outline-none focus:border-[#f10215]"
              >
                <option :value="1">普通用户</option>
                <option :value="10">管理员</option>
              </select>
            </div>

            <button
              class="w-full rounded-full bg-[#f10215] px-4 py-2.5 text-sm font-semibold text-white transition hover:brightness-95"
              @click="handleAuth"
            >
              {{ isRegister ? '立即注册' : '登录' }}
            </button>

            <button
              class="w-full text-center text-sm text-[#f10215] transition hover:text-[#cf0512]"
              @click="isRegister = !isRegister"
            >
              {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
            </button>
          </div>
        </div>
      </section>

      <!-- ORDERS -->
      <section v-if="ui.view === VIEWS.ORDERS" class="mt-6 space-y-4">
        <div class="flex items-center gap-2">
          <span class="h-6 w-0.5 rounded-full bg-[#f10215]"></span>
          <h2 class="text-xl font-bold text-gray-900">我的订单</h2>
        </div>

        <div v-if="ui.loading" class="flex justify-center py-10">
          <div class="loader" />
        </div>

        <div v-else-if="orders.length === 0" class="rounded-2xl bg-white p-10 text-center text-gray-400 shadow-sm ring-1 ring-gray-100">
          暂无订单
        </div>

        <div v-else class="space-y-4">
          <article
            v-for="order in orders"
            :key="order.ID"
            class="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100"
          >
            <div class="flex items-center justify-between border-b border-gray-100 pb-2">
              <div class="text-xs text-gray-500">订单号 {{ order.ID }}</div>
              <span class="rounded-full bg-gray-100 px-3 py-1 text-xs text-gray-700">{{ order.Status }}</span>
            </div>
            <div class="divide-y divide-gray-100">
              <div v-for="item in order.Items" :key="item.ID" class="flex items-center gap-3 py-3">
                <img
                  :src="item.ProductImg || 'https://via.placeholder.com/60'"
                  class="h-12 w-12 rounded-lg bg-gray-50 object-cover"
                  @error="handleImageError"
                />
                <div class="flex-1">
                  <p class="line-clamp-1 text-sm font-semibold text-gray-900">{{ item.ProductName }}</p>
                  <div class="mt-1 flex justify-between text-xs text-gray-500">
                    <span>x{{ item.Quantity }}</span>
                    <span>¥{{ item.Price }}</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="mt-3 flex items-center justify-end gap-2 text-gray-700">
              <span class="text-xs text-gray-500">实付</span>
              <span class="text-lg font-bold">¥{{ calculateOrderTotal(order) }}</span>
            </div>
          </article>
        </div>
      </section>
    </main>

    <!-- MOBILE NAV -->
    <nav
      class="fixed bottom-0 left-0 right-0 z-40 flex items-center justify-around border-t border-gray-100 bg-white/95 py-2 text-xs text-gray-600 shadow-lg backdrop-blur md:hidden"
    >
      <button
        class="flex flex-col items-center"
        :class="{ 'text-[#f10215]': ui.view === VIEWS.HOME }"
        @click="changeView(VIEWS.HOME)"
      >
        <i class="fas fa-home text-xl"></i>
        <span>首页</span>
      </button>
      <div class="flex flex-col items-center text-gray-400">
        <i class="fas fa-compass text-xl"></i>
        <span>分类</span>
      </div>
      <button
        class="relative flex flex-col items-center"
        :class="{ 'text-[#f10215]': ui.view === VIEWS.CART }"
        @click="changeView(VIEWS.CART)"
      >
        <i class="fas fa-shopping-cart text-xl"></i>
        <span>购物车</span>
        <span
          v-if="cartTotalCount > 0"
          class="absolute -right-2 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-[#f10215] text-[10px] text-white"
          >{{ cartTotalCount }}</span
        >
      </button>
      <button
        class="flex flex-col items-center"
        :class="{ 'text-[#f10215]': ui.view === VIEWS.ORDERS || ui.view === VIEWS.LOGIN }"
        @click="session.user ? changeView(VIEWS.ORDERS) : changeView(VIEWS.LOGIN)"
      >
        <i class="fas fa-user text-xl"></i>
        <span>我的</span>
      </button>
    </nav>

    <!-- TOAST -->
    <transition name="fade">
      <div
        v-if="ui.toast.show"
        class="fixed left-1/2 top-10 -translate-x-1/2 rounded-full bg-black/80 px-4 py-2 text-sm text-white shadow-lg"
      >
        {{ ui.toast.message }}
      </div>
    </transition>
  </div>
</template>

<style scoped>
:global(:root) {
  --jd-red: #f10215;
  --jd-dark-red: #d9000b;
}

.loader {
  border: 4px solid #f3f3f3;
  border-radius: 50%;
  border-top: 4px solid var(--jd-red);
  width: 36px;
  height: 36px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
