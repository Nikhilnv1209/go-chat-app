import { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { socketService } from '@/lib/socketService';
import { Message } from '@/types';
import { receiveMessage } from '@/store/features/conversationSlice';
import { useQueryClient } from '@tanstack/react-query';

export function useSocketConnection() {
  const dispatch = useAppDispatch();
  const queryClient = useQueryClient();
  const { token, isAuthenticated, user } = useAppSelector((state) => state.auth);

  // 1. Connection Management Effect
  useEffect(() => {
    if (isAuthenticated && token && user) {
        socketService.connect(token);
    }
    // We do NOT disconnect on cleanup of this effect if dependencies change,
    // because we want the socket to persist across navigations/renders.
    // We only want to disconnect if the user explicitly logs out (token becomes null).

    if (!isAuthenticated || !token) {
        socketService.disconnect();
    }
  }, [isAuthenticated, token, user]);

  // 2. Event Listener Management Effect
  useEffect(() => {
    if (!isAuthenticated || !token || !user) return;

    const handleNewMessage = (message: Message) => {
      dispatch(receiveMessage(message));

      let targetId = '';
      let type = '';

      if (message.msg_type === 'group') {
            targetId = message.group_id!;
            type = 'GROUP';
      } else {
          if (message.sender_id === user.id) {
              targetId = message.receiver_id!;
          } else {
              targetId = message.sender_id;
          }
          type = 'DM';
      }

      queryClient.setQueryData(['messages', targetId, type], (oldData: any) => {
          if (!oldData) return [message];
          if (oldData.some((m: Message) => m.id === message.id)) return oldData;
          return [...oldData, message];
      });
    };

    const handleMessageSent = (message: Message) => {
        handleNewMessage(message);
    };

    socketService.on('new_message', handleNewMessage);
    socketService.on('message_sent', handleMessageSent);

    return () => {
      socketService.off('new_message', handleNewMessage);
      socketService.off('message_sent', handleMessageSent);
    };
  }, [isAuthenticated, token, user, dispatch, queryClient]);
}
