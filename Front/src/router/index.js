import { createRouter, createWebHistory } from "vue-router";
import UserLanding from "../pages/user/UserLanding.vue";
import UserRegister from "../pages/user/UserRegister.vue";
import UserLogin from "../pages/user/UserLogin.vue";
import UserHome from "../pages/user/UserHome.vue";
import ShopDetail from "../pages/user/ShopDetail.vue";
import ProductDetail from "../pages/user/ProductDetail.vue";
import CartPage from "../pages/user/Cart.vue";
import MerchantLanding from "../pages/merchant/MerchantLanding.vue";
import MerchantRegister from "../pages/merchant/MerchantRegister.vue";
import MerchantLogin from "../pages/merchant/MerchantLogin.vue";
import MerchantHome from "../pages/merchant/MerchantHome.vue";
import ChatRoom from "../pages/shared/ChatRoom.vue";
import NotFound from "../pages/shared/NotFound.vue";
import { getAuth } from "../services/storage.js";

const routes = [
  { path: "/", redirect: "/user" },
  {
    path: "/user",
    component: UserLanding,
    meta: { layout: "user", title: "用户端" }
  },
  {
    path: "/user/register",
    component: UserRegister,
    meta: { layout: "user", title: "用户注册" }
  },
  {
    path: "/user/login",
    component: UserLogin,
    meta: { layout: "user", title: "用户登录" }
  },
  {
    path: "/user/home",
    component: UserHome,
    meta: { layout: "user", title: "用户中心" }
  },
  {
    path: "/cart",
    component: CartPage,
    meta: { layout: "user", title: "购物车" }
  },
  {
    path: "/shops/:id",
    component: ShopDetail,
    meta: { layout: "user", title: "店铺详情" }
  },
  {
    path: "/products/:id",
    component: ProductDetail,
    meta: { layout: "user", title: "商品详情" }
  },
  {
    path: "/merchant",
    component: MerchantLanding,
    meta: { layout: "merchant", title: "商家端" }
  },
  {
    path: "/merchant/register",
    component: MerchantRegister,
    meta: { layout: "merchant", title: "商家注册" }
  },
  {
    path: "/merchant/login",
    component: MerchantLogin,
    meta: { layout: "merchant", title: "商家登录" }
  },
  {
    path: "/merchant/home",
    component: MerchantHome,
    meta: { layout: "merchant", title: "商家控制台" }
  },
  {
    path: "/chat",
    component: ChatRoom,
    meta: { layout: "neutral", title: "聊天" }
  },
  {
    path: "/:pathMatch(.*)*",
    component: NotFound,
    meta: { layout: "neutral", title: "页面未找到" }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.afterEach((to) => {
  if (to.meta && to.meta.title) {
    document.title = `商店系统 - ${to.meta.title}`;
  } else {
    document.title = "商店系统";
  }
});

router.beforeEach((to) => {
  const protectedPaths = ["/cart", "/user/home", "/merchant/home"];
  const requiresAuth = protectedPaths.some((path) => to.path.startsWith(path));
  if (!requiresAuth) {
    return true;
  }
  const auth = getAuth();
  if (auth && auth.access_token) {
    return true;
  }
  return "/user/login";
});

export default router;
