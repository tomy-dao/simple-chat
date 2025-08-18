import axios from 'axios';

// Create axios instance with base configuration
export const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_API_URL || 'http://localhost/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
client.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);