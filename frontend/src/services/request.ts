import axios from 'axios';

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api/v1',
  timeout: 10000,
});

// Request interceptor
request.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor
request.interceptors.response.use(
  response => {
    if (response.data.code !== 0) {
      return Promise.reject(new Error(response.data.message));
    }
    return response.data;
  },
  error => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export default request;
