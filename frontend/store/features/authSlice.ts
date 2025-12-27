import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { AuthState, User } from '@/types';

const initialState: AuthState = {
  user: null,
  token: typeof window !== 'undefined' ? localStorage.getItem('token') : null,
  isAuthenticated: false,
  isLoading: true, // Start true to check for existing token
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setCredentials: (
      state,
      action: PayloadAction<{ user: User; token: string }>
    ) => {
      state.user = action.payload.user;
      state.token = action.payload.token;
      state.isAuthenticated = true;
      state.isLoading = false;
      if (typeof window !== 'undefined') {
        localStorage.setItem('token', action.payload.token);
        if (action.payload.user) {
          localStorage.setItem('user', JSON.stringify(action.payload.user));
        }
      }
    },
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      state.isLoading = false;
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      }
    },
    initializeAuth: (state) => {
      // Hydrate state from localStorage on mount (client-side only)
      if (typeof window === 'undefined') {
        state.isLoading = false;
        return;
      }

      const token = localStorage.getItem('token');
      const userStr = localStorage.getItem('user');

      if (token && userStr && userStr !== 'undefined') {
        try {
          state.user = JSON.parse(userStr);
          state.token = token;
          state.isAuthenticated = true;
        } catch (e) {
          console.error('Failed to parse user from local storage', e);
          localStorage.removeItem('user');
        }
      }
      state.isLoading = false;
    }
  },
});

export const { setCredentials, logout, initializeAuth } = authSlice.actions;
export default authSlice.reducer;
