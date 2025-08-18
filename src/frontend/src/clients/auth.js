import { client } from "./initial";

// Auth API functions
const auth = {
  // Login user
  login: async (credentials) => {
    try {
      const response = await client.post('/login', credentials);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Login failed');
    }
  },

  // Register user
  register: async (userData) => {
    try {
      const response = await client.post('/register', userData);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Registration failed');
    }
  },

  // Get user
  getMe: async () => {
    try {
      const response = await client.get('/me', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to get user');
    }
  },

  // Get users

  getUsers: async () => {
    try {
      const response = await client.get('/users');
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || 'Failed to get users');
    }
  },

  // Logout user
  logout: async () => {
    try {
      await client.post('/logout');
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      // Clear auth data regardless of API response
      localStorage.removeItem('authToken');
      localStorage.removeItem('user');
    }
  },
};

export default auth;
