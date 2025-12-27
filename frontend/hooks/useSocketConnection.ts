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
      console.log('WS: Received Message:', message);
      dispatch(receiveMessage(message));

      let targetId = '';
      let type: 'DM' | 'GROUP' = 'DM';

      if (message.msg_type === 'GROUP') {
          targetId = message.group_id!;
          type = 'GROUP';
      } else {
          // It's a DM (backend uses 'private', but we normalize to 'DM' in types/store)
          if (message.sender_id === user.id) {
              targetId = message.receiver_id!;
          } else {
              targetId = message.sender_id;
          }
          type = 'DM';
      }

      const queryKey = ['messages', targetId, type];
      console.log('WS: Updating Query Cache with key:', queryKey);

      queryClient.setQueryData(queryKey, (oldData: Message[] | undefined) => {
          if (!oldData) {
              console.log('WS: No old data in cache, creating new array');
              return [message];
          }

          if (!Array.isArray(oldData)) {
              console.error('WS: Unexpected cache data format (not an array):', oldData);
              return [message];
          }

          const alreadyExists = oldData.some((m: Message) => m.id === message.id);
          if (alreadyExists) {
              console.log('WS: Message already in cache, ignoring');
              return oldData;
          }

          console.log('WS: Appending message to cache. New count:', oldData.length + 1);
          return [...oldData, message];
      });
    };

    const handleMessageSent = (message: Message) => {
        console.log('WS: Message Sent Confirmation:', message);
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
