const tokenKey = 'nudgepay_token';

export function getToken(): string | null {
  if (typeof window === 'undefined') {
    return null;
  }
  return window.localStorage.getItem(tokenKey);
}

export function setToken(token: string) {
  if (typeof window === 'undefined') {
    return;
  }
  window.localStorage.setItem(tokenKey, token);
}

export function clearToken() {
  if (typeof window === 'undefined') {
    return;
  }
  window.localStorage.removeItem(tokenKey);
}
