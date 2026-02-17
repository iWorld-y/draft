import axios from 'axios';
import type { AxiosRequestConfig } from 'axios';
import { clearAuth, getAuthToken, refreshAccessToken, setAuthToken } from './auth';

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api/v1',
  timeout: 10000,
  withCredentials: true,
});

request.interceptors.request.use((config) => {
  const token = getAuthToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

request.interceptors.response.use(
  (response) => response.data,
  async (error) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };

    if (error?.response?.status === 401 && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true;
      try {
        const newToken = await refreshAccessToken();
        setAuthToken(newToken);
        originalRequest.headers = {
          ...(originalRequest.headers || {}),
          Authorization: `Bearer ${newToken}`,
        };
        return request(originalRequest);
      } catch (refreshError) {
        clearAuth();
        return Promise.reject(refreshError);
      }
    }

    const serverMessage = error?.response?.data?.message;
    if (serverMessage) {
      return Promise.reject(new Error(serverMessage));
    }

    return Promise.reject(error);
  },
);

export default request;
