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
      const isGroup = msg.msg_type === 'GROUP';

      let conversationIndex = -1;

      if (isGroup) {
          conversationIndex = state.conversations.findIndex(c => c.target_id === msg.group_id && c.type === 'GROUP');
      } else {
          // Identify the other participant's ID (target_id)
          // For received messages, it's msg.sender_id
          // For sent messages (if we receive our own echo or updated message), it would be receiver_id
          // But usually we just need to find the conversation where the other user is involved.

          // Try finding by sender first
          conversationIndex = state.conversations.findIndex(c => c.target_id === msg.sender_id && c.type === 'DM');

          // If not found and there's a receiver (e.g. we sent it), try receiver_id
          if (conversationIndex === -1 && msg.receiver_id) {
            conversationIndex = state.conversations.findIndex(c => c.target_id === msg.receiver_id && c.type === 'DM');
          }
      }

      if (conversationIndex !== -1) {
          const conversation = state.conversations[conversationIndex];
          conversation.last_message = msg.content;
          conversation.last_message_at = msg.created_at;

          // Increment unread if NOT the active conversation
          // Note: activeConversationId in this slice is the conversation UUID,
          // while conversation.id is also that same UUID.
          if (conversation.id !== state.activeConversationId) {
             conversation.unread_count += 1;
          }

          // Move updated conversation to top
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
