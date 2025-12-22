import { getAuth } from "./storage.js";

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8080";

async function request(path, options = {}) {
  const auth = getAuth();

  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {})
  };
  if (auth && auth.access_token) {
    headers.Authorization = `Bearer ${auth.access_token}`;
  }
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    const message = data.error || "请求失败";
    throw new Error(message);
  }
  return data;
}

export async function registerUser(payload) {
  return request("/api/v0/users", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function login(payload) {
  return request("/api/v0/auth/login", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function refreshToken(payload) {
  return request("/api/v0/auth/refresh", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function logout() {
  return request("/api/v0/auth/logout", {
    method: "POST"
  });
}

export async function listShops() {
  return request("/api/v2/shops", { method: "GET" });
}

export async function getShop(id) {
  return request(`/api/v2/shops/${id}`, { method: "GET" });
}

export async function listProductsByShop(id) {
  return request(`/api/v2/shops/${id}/products`, { method: "GET" });
}

export async function getProduct(id) {
  return request(`/api/v2/products/${id}`, { method: "GET" });
}

export async function listOrders() {
  return request("/api/v1/orders", { method: "GET" });
}

export async function createOrder(payload) {
  return request("/api/v1/orders", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function listCart() {
  return request("/api/v1/cart", { method: "GET" });
}

export async function addCartItem(payload) {
  return request("/api/v1/cart/items", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function updateCartItem(id, payload) {
  return request(`/api/v1/cart/items/${id}`, {
    method: "PATCH",
    body: JSON.stringify(payload)
  });
}

export async function deleteCartItem(id) {
  return request(`/api/v1/cart/items/${id}`, { method: "DELETE" });
}

export async function clearCart() {
  return request("/api/v1/cart", { method: "DELETE" });
}
export async function updateOrderStatus(id, status) {
  return request(`/api/v1/orders/${id}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status })
  });
}

export async function createShop(payload) {
  return request("/api/v2/shops", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function createProduct(shopId, payload) {
  return request(`/api/v2/shops/${shopId}/products`, {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

// TODO: Wire chat endpoints once they are exposed by the backend.
