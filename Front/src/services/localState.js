const PREFIX = "shop_front_";

export function loadLocal(key, fallback) {
  const raw = localStorage.getItem(`${PREFIX}${key}`);
  if (!raw) {
    return fallback;
  }
  try {
    return JSON.parse(raw);
  } catch {
    return fallback;
  }
}

export function saveLocal(key, value) {
  localStorage.setItem(`${PREFIX}${key}`, JSON.stringify(value));
}
