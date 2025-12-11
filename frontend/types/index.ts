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
