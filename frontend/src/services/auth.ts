import axios from 'axios';

export interface AuthUser {
  id: number;
  username: string;
}

interface AuthPayload {
  access_token: string;
  user: AuthUser;
}

const CURRENT_USER_KEY = 'draft.currentUser';
const TOKEN_KEY = 'token';

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api/v1';

const authClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  withCredentials: true,
});

const persistSession = (payload: AuthPayload): AuthUser => {
  localStorage.setItem(CURRENT_USER_KEY, JSON.stringify(payload.user));
  localStorage.setItem(TOKEN_KEY, payload.access_token);
  return payload.user;
};

const unwrapAuthPayload = (response: { data: { code: number; message?: string; data?: AuthPayload } }): AuthPayload => {
  if (response.data?.code !== 0 || !response.data?.data?.access_token) {
    throw new Error(response.data?.message || '认证失败');
  }
  return response.data.data;
};

export const getCurrentUser = (): AuthUser | null => {
  const raw = localStorage.getItem(CURRENT_USER_KEY);
  if (!raw) {
    return null;
  }

  try {
    const user = JSON.parse(raw);
    if (!user?.id || !user?.username) {
      return null;
    }
    return { id: user.id, username: user.username };
  } catch {
    return null;
  }
};

export const getAuthToken = (): string | null => localStorage.getItem(TOKEN_KEY);

export const setAuthToken = (token: string) => {
  localStorage.setItem(TOKEN_KEY, token);
};

export const clearAuth = () => {
  localStorage.removeItem(CURRENT_USER_KEY);
  localStorage.removeItem(TOKEN_KEY);
};

export const register = async (username: string, password: string): Promise<AuthUser> => {
  const response = await authClient.post('/auth/register', { username, password });
  return persistSession(unwrapAuthPayload(response));
};

export const login = async (username: string, password: string): Promise<AuthUser> => {
  const response = await authClient.post('/auth/login', { username, password });
  return persistSession(unwrapAuthPayload(response));
};

export const refreshAccessToken = async (): Promise<string> => {
  const response = await authClient.post('/auth/refresh', {});
  const payload = unwrapAuthPayload(response);
  localStorage.setItem(CURRENT_USER_KEY, JSON.stringify(payload.user));
  setAuthToken(payload.access_token);
  return payload.access_token;
};

export const logout = async () => {
  try {
    await authClient.post('/auth/logout', {});
  } finally {
    clearAuth();
  }
};

export const fetchCurrentUser = async (): Promise<AuthUser> => {
  const token = getAuthToken();
  if (!token) {
    throw new Error('未登录');
  }

  const response = await authClient.get('/auth/me', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (response.data?.code !== 0 || !response.data?.data?.id) {
    throw new Error(response.data?.message || '获取用户信息失败');
  }

  const user = response.data.data as AuthUser;
  localStorage.setItem(CURRENT_USER_KEY, JSON.stringify(user));
  return user;
};
