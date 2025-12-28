'use client';

import { useRef, useEffect } from 'react';
import { Provider } from 'react-redux';
import { makeStore, AppStore } from '../store/store';
import { initializeAuth, setToken } from '../store/features/authSlice';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

export default function Providers({ children }: { children: React.ReactNode }) {
  const storeRef = useRef<AppStore>(null);
  if (!storeRef.current) {
    // Create the store instance the first time this renders
    storeRef.current = makeStore();
  }

  // Initialize auth state from local storage on mount (client-side only)
  useEffect(() => {
    if (storeRef.current) {
      storeRef.current.dispatch(initializeAuth());
    }

    const handleTokenRefresh = (event: CustomEvent<string>) => {
      if (storeRef.current) {
        // console.log('Syncing refreshed token to Redux store');
        storeRef.current.dispatch(setToken(event.detail));
      }
    };

    window.addEventListener('auth:token-refreshed', handleTokenRefresh as EventListener);

    return () => {
      window.removeEventListener('auth:token-refreshed', handleTokenRefresh as EventListener);
    };
  }, []);

  const queryClientRef = useRef<QueryClient>(null);
  if (!queryClientRef.current) {
    queryClientRef.current = new QueryClient();
  }

  return (
    <Provider store={storeRef.current}>
      <QueryClientProvider client={queryClientRef.current}>
        {children}
      </QueryClientProvider>
    </Provider>
  );
}
