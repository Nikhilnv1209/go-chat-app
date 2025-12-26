// User Types
export interface User {
  id: string;
  username: string;
  email: string;
  is_online: boolean;
  last_seen: string;
  created_at: string;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

// Message Types
export interface Message {
  id: string;
  sender_id: string;
  receiver_id: string | null;
  group_id: string | null;
  content: string;
  msg_type: 'private' | 'group';
  created_at: string;
  sender?: User;
}

// Conversation Types
export interface Conversation {
  id: string;
  type: 'DM' | 'GROUP';
  target_id: string;
  target_name: string;
  target_avatar?: string;
  last_message: string | null;
  last_message_at: string;
  unread_count: number;
  is_online?: boolean;
  member_count?: number; // Only for GROUP conversations
}

// Group Types
export interface Group {
  id: string;
  name: string;
  created_at: string;
}

// Message Receipt Types
export interface MessageReceipt {
  id: string;
  message_id: string;
  user_id: string;
  status: 'SENT' | 'DELIVERED' | 'READ';
  created_at: string;
  updated_at: string;
}

// WebSocket Event Types
export type WSOutgoingEvent =
  | { type: 'send_message'; payload: { to_user_id?: string; group_id?: string; content: string } }
  | { type: 'typing_start'; payload: { conversation_type: 'DM' | 'GROUP'; target_id: string } }
  | { type: 'typing_stop'; payload: { conversation_type: 'DM' | 'GROUP'; target_id: string } }
  | { type: 'message_delivered'; payload: { message_id: string } };

export type WSIncomingEvent =
  | { type: 'new_message'; payload: Message }
  | { type: 'message_sent'; payload: Message }
  | { type: 'user_typing'; payload: { user_id: string; username: string; conversation_type: string; target_id: string } }
  | { type: 'user_stopped_typing'; payload: { user_id: string; conversation_type: string; target_id: string } }
  | { type: 'receipt_update'; payload: { message_id: string; user_id: string; status: string; updated_at: string } };
