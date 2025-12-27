import axios from 'axios';
import { Conversation, Message, User } from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const conversationApi = {
  /**
   * Search for users
   * GET /users?q=query
   */
  searchUsers: async (token: string, query: string): Promise<User[]> => {
    const response = await axios.get(`${API_BASE_URL}/users`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
      params: { q: query },
    });
    return response.data;
  },

  /**
   * Fetch a specific user by ID
   */
  getUser: async (token: string, userId: string): Promise<User> => {
    const response = await axios.get(`${API_BASE_URL}/users/${userId}`, {
       headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  },

  /**
   * Fetch all conversations for the authenticated user
   * GET /conversations
   */
  getConversations: async (token: string): Promise<Conversation[]> => {
    const response = await axios.get(`${API_BASE_URL}/conversations`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    // Backend now returns 'DM' or 'GROUP' directly
    return response.data;
  },

  /**
   * Fetch message history for a specific conversation
   * GET /messages?target_id=<uuid>&type=<DM|GROUP>&limit=<50>&before_id=<uuid>
   */
  getMessages: async (
    token: string,
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

    const response = await axios.get(`${API_BASE_URL}/messages`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
      params,
    });
    return response.data.reverse(); // Backend returns newest first, but UI expects oldest first (bottom up)
  },

  /**
   * Mark messages as read
   * POST /messages/:id/read
   */
  markAsRead: async (token: string, messageId: string): Promise<void> => {
    await axios.post(
      `${API_BASE_URL}/messages/${messageId}/read`,
      {},
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
  },
};
