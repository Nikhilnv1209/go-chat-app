import api from './api';
import { Conversation, Message, User } from '@/types';

export const conversationApi = {
  /**
   * Search for users
   * GET /users?q=query
   */
  searchUsers: async (query: string): Promise<User[]> => {
    const response = await api.get('/users', {
      params: { q: query },
    });
    return response.data;
  },

  /**
   * Fetch a specific user by ID
   */
  getUser: async (userId: string): Promise<User> => {
    const response = await api.get(`/users/${userId}`);
    return response.data;
  },

  /**
   * Fetch all conversations for the authenticated user
   * GET /conversations
   */
  getConversations: async (): Promise<Conversation[]> => {
    const response = await api.get('/conversations');
    // Backend now returns 'DM' or 'GROUP' directly
    return response.data;
  },

  /**
   * Fetch message history for a specific conversation
   * GET /messages?target_id=<uuid>&type=<DM|GROUP>&limit=<50>&before_id=<uuid>
   */
  getMessages: async (
    targetId: string,
    type: 'DM' | 'GROUP',
    limit: number = 50,
    beforeId?: string
  ): Promise<Message[]> => {
    const params: Record<string, string> = {
      target_id: targetId,
      type,
      limit: limit.toString(),
    };

    if (beforeId) {
      params.before_id = beforeId;
    }

    const response = await api.get('/messages', {
      params,
    });
    return response.data.reverse(); // Backend returns newest first, but UI expects oldest first (bottom up)
  },

  /**
   * Mark messages as read
   * POST /messages/:id/read
   */
  markAsRead: async (messageId: string): Promise<void> => {
    await api.post(`/messages/${messageId}/read`);
  },
};
