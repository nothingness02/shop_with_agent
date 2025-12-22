const AUTH_KEY = "shop_front_auth";

export function getAuth() {
  const raw = localStorage.getItem(AUTH_KEY);
  if (!raw) {
    return null;
  }
  try {
    const cleaned = raw.replace(/^\uFEFF/, "").trim();
    return JSON.parse(cleaned);
  } catch {
    const tokenMatch = raw.match(/\"access_token\"\\s*:\\s*\"([^\"]+)\"/);
    if (tokenMatch) {
      return { access_token: tokenMatch[1] };
    }
    return null;
  }
}

export function setAuth(payload) {
  localStorage.setItem(AUTH_KEY, JSON.stringify(payload));
  window.dispatchEvent(new Event("auth-changed"));
}

export function clearAuth() {
  localStorage.removeItem(AUTH_KEY);
  window.dispatchEvent(new Event("auth-changed"));
}
