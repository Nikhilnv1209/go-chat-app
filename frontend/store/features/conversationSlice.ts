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
    receiveMessage: (state, action: PayloadAction<import('@/types').Message>) => {
      const msg = action.payload;
      // Identify conversation target
      const isGroup = msg.msg_type === 'group';
      const targetId = isGroup ? msg.group_id : (msg.sender_id === state.activeConversationId /* logic mismatch here, activeConvId is UUID, sender_id is UUID */ ? msg.sender_id : msg.sender_id);

      // Logic fix: In DM, if I receive a message, the conversation is with the Sender.
      // If I sent a message (this event also fires for me?), the conversation is with Receiver.
      // But we receive 'new_message' usually only if we are the recipient or it's a group broadcast.

      let conversationIndex = -1;

      if (isGroup) {
          conversationIndex = state.conversations.findIndex(c => c.target_id === msg.group_id && c.type === 'GROUP');
      } else {
          // DM: Find conversation where target_id is the sender
          conversationIndex = state.conversations.findIndex(c => c.target_id === msg.sender_id && c.type === 'DM');
      }

      if (conversationIndex !== -1) {
          const conversation = state.conversations[conversationIndex];
          conversation.last_message = msg.content;
          conversation.last_message_at = msg.created_at;

          // Increment unread if not active
          if (conversation.id !== state.activeConversationId) {
             conversation.unread_count += 1;
          }

          // Move to top
          state.conversations.splice(conversationIndex, 1);
          state.conversations.unshift(conversation);
      }
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
  receiveMessage,
} = conversationSlice.actions;

export default conversationSlice.reducer;
