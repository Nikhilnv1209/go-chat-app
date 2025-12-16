import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Conversation } from '@/types';

interface ConversationState {
  conversations: Conversation[];
  activeConversationId: string | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: ConversationState = {
  conversations: [],
  activeConversationId: null,
  isLoading: false,
  error: null,
};

const conversationSlice = createSlice({
  name: 'conversation',
  initialState,
  reducers: {
    setConversations: (state, action: PayloadAction<Conversation[]>) => {
      state.conversations = action.payload;
      state.isLoading = false;
      state.error = null;
    },
    setActiveConversation: (state, action: PayloadAction<string | null>) => {
      state.activeConversationId = action.payload;
    },
    updateConversation: (state, action: PayloadAction<Conversation>) => {
      const index = state.conversations.findIndex((c) => c.id === action.payload.id);
      if (index !== -1) {
        state.conversations[index] = action.payload;
      } else {
        state.conversations.unshift(action.payload);
      }
      // Sort by last_message_at descending
      state.conversations.sort((a, b) =>
        new Date(b.last_message_at).getTime() - new Date(a.last_message_at).getTime()
      );
    },
    incrementUnread: (state, action: PayloadAction<string>) => {
      const conversation = state.conversations.find((c) => c.id === action.payload);
      if (conversation) {
        conversation.unread_count += 1;
      }
    },
    resetUnread: (state, action: PayloadAction<string>) => {
      const conversation = state.conversations.find((c) => c.id === action.payload);
      if (conversation) {
        conversation.unread_count = 0;
      }
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
      state.isLoading = false;
    },
  },
});

export const {
  setConversations,
  setActiveConversation,
  updateConversation,
  incrementUnread,
  resetUnread,
  setLoading,
  setError,
} = conversationSlice.actions;

export default conversationSlice.reducer;
