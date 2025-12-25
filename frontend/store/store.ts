import { configureStore } from '@reduxjs/toolkit';
import authReducer from './features/authSlice';
import conversationReducer from './features/conversationSlice';
import folderReducer from './features/folderSlice';
import uiReducer from './features/uiSlice';

export const makeStore = () => {
  return configureStore({
    reducer: {
      auth: authReducer,
      conversation: conversationReducer,
      folders: folderReducer,
      ui: uiReducer,
    },
  });
};

// Infer the type of makeStore
export type AppStore = ReturnType<typeof makeStore>;
// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<AppStore['getState']>;
export type AppDispatch = AppStore['dispatch'];
