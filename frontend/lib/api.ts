import axios from 'axios';

// Create an Axios instance with default configuration
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Important: Send cookies (refresh_token) with requests
});

// Add a request interceptor to attach the JWT token
api.interceptors.request.use(
  (config) => {
    // We'll store the token in localStorage
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

interface FailedRequest {
  resolve: (token: string) => void;
  reject: (error: any) => void;
}

let isRefreshing = false;
let failedQueue: FailedRequest[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token!);
    }
  });

  failedQueue = [];
};

// Response interceptor for handling token expiration
api.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;

    // Check if error is 401 and we haven't tried to refresh yet
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // If already refreshing, queue this request
        return new Promise(function (resolve, reject) {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers['Authorization'] = 'Bearer ' + token;
            return api(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // Call refresh endpoint
        // Note: we don't use 'api' instance here to avoid infinite loops if this fails
        // But we DO need withCredentials
        const response = await axios.post(
          `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/auth/refresh`,
          {},
          { withCredentials: true }
        );

        const { token } = response.data;

        if (token) {
            localStorage.setItem('token', token);
            api.defaults.headers.common['Authorization'] = 'Bearer ' + token;

            // Dispatch event for React components to sync state
            if (typeof window !== 'undefined') {
              window.dispatchEvent(new CustomEvent('auth:token-refreshed', { detail: token }));
            }

            processQueue(null, token);
            isRefreshing = false;

            // Update the header and retry
            originalRequest.headers['Authorization'] = 'Bearer ' + token;
            return api(originalRequest);
        }
      } catch (err) {
        processQueue(err, null);
        isRefreshing = false;

        // Refresh failed (token expired or revoked)
        // Redirect to login
        if (typeof window !== 'undefined') {
            localStorage.removeItem('token');
            localStorage.removeItem('user'); // Clean up user data too if exists
            window.location.href = '/login';
        }
        return Promise.reject(err);
      }
    }

    return Promise.reject(error);
  }
);

export default api;
