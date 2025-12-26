import { WSIncomingEvent, WSOutgoingEvent } from '@/types';

type EventHandler<T = any> = (payload: T) => void;

class SocketService {
  private socket: WebSocket | null = null;
  private messageQueue: string[] = [];
  private listeners: Map<string, EventHandler[]> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 10;
  private reconnectTimeoutRef: NodeJS.Timeout | null = null;
  private token: string | null = null;

  // Singleton instance
  private static instance: SocketService;

  private constructor() {}

  public static getInstance(): SocketService {
    if (!SocketService.instance) {
      SocketService.instance = new SocketService();
    }
    return SocketService.instance;
  }

  public connect(token: string) {
    if (this.socket?.readyState === WebSocket.OPEN) return;


    this.token = token;

    // Determine the WebSocket URL
    // If we are in dev, and API is localhost:8080, we want ws://localhost:8080/ws
    // In production, it might be wss://api.example.com/ws

    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    let wsUrl = apiUrl.replace('http', 'ws');
    wsUrl = `${wsUrl}/ws?token=${token}`;

    console.log('Connecting to WebSocket:', wsUrl);
    this.socket = new WebSocket(wsUrl);

    this.socket.onopen = () => {
      console.log('WebSocket Connected');
      this.reconnectAttempts = 0;
      this.flushQueue();
    };

    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as WSIncomingEvent;
        // console.log('WS Received:', data);
        this.emit(data.type, data.payload);
      } catch (err) {
        console.error('Failed to parse WS message:', err);
      }
    };

    this.socket.onclose = (event) => {
      console.log('WebSocket Closed', event.code, event.reason);
      this.socket = null;
      this.attemptReconnect();
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket Error:', error);
      this.socket?.close();
    };
  }

  public disconnect() {
    if (this.socket) {
      // Prevent reconnect on manual disconnect
      this.reconnectAttempts = this.maxReconnectAttempts;
      if (this.reconnectTimeoutRef) clearTimeout(this.reconnectTimeoutRef);

      this.socket.close();
      this.socket = null;
    }
  }

  // Event Subscription
  public on<T = any>(event: string, callback: EventHandler<T>) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event)!.push(callback);
  }

  public off<T = any>(event: string, callback: EventHandler<T>) {
    if (!this.listeners.has(event)) return;
    const filtered = this.listeners.get(event)!.filter((cb) => cb !== callback);
    this.listeners.set(event, filtered);
  }

  // Internal Event Dispatch (Incoming from Server -> Listeners)
  private emit(event: string, payload: any) {
    const callbacks = this.listeners.get(event);
    if (callbacks) {
      callbacks.forEach((cb) => cb(payload));
    }
  }

  // Sending Messages (Client -> Server)
  private send(event: WSOutgoingEvent) {
    const payload = JSON.stringify(event);
    if (this.socket?.readyState === WebSocket.OPEN) {
      this.socket.send(payload);
    } else {
      this.messageQueue.push(payload);
    }
  }

  private flushQueue() {
    while (this.messageQueue.length > 0 && this.socket?.readyState === WebSocket.OPEN) {
      const msg = this.messageQueue.shift();
      if (msg) this.socket.send(msg);
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.warn('Max reconnect attempts reached');
      return;
    }

    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
    console.log(`Attempting reconnect in ${delay}ms...`);

    this.reconnectTimeoutRef = setTimeout(() => {
        this.reconnectAttempts++;
        if (this.token) {
            this.connect(this.token);
        }
    }, delay);
  }

  // --- Public API Methods ---

  public sendMessage(content: string, targetId: string, type: 'DM' | 'GROUP') {
    const payload: any = { content };
    if (type === 'DM') {
        payload.to_user_id = targetId;
    } else {
        payload.group_id = targetId;
    }

    this.send({
      type: 'send_message',
      payload: payload,
    });
  }

  public sendTyping(targetId: string, type: 'DM' | 'GROUP', isTyping: boolean) {
    this.send({
      type: isTyping ? 'typing_start' : 'typing_stop',
      payload: {
        conversation_type: type,
        target_id: targetId,
      },
    });
  }

  public markDelivered(messageId: string) {
      this.send({
          type: 'message_delivered',
          payload: { message_id: messageId }
      });
  }
}

export const socketService = SocketService.getInstance();
