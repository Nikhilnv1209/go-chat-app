import { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { socketService } from '@/lib/socketService';
import { Message, Conversation } from '@/types';
import { receiveMessage, setUserOnlineStatus, addConversation } from '@/store/features/conversationSlice';
import { useQueryClient } from '@tanstack/react-query';

export function useSocketConnection() {
  const dispatch = useAppDispatch();
  const queryClient = useQueryClient();
  const { token, isAuthenticated, user } = useAppSelector((state) => state.auth);

  // Combined Effect: Attach listeners THEN connect to prevent race conditions
  useEffect(() => {
    if (!isAuthenticated || !token || !user) {
      socketService.disconnect();
      return;
    }

    // Define event handlers
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

      // Always invalidate conversations list when a new message arrives
      // This ensures the inbox list (sidebar) updates its order, unread counts, and shows new conversations instantly.
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
    };

    const handleMessageSent = (message: Message) => {
        console.log('WS: Message Sent Confirmation:', message);
        handleNewMessage(message);
    };

    const handleUserOnline = (payload: { user_id: string }) => {
      console.log('WS: User Online:', payload.user_id);
      dispatch(setUserOnlineStatus({ userId: payload.user_id, isOnline: true }));
    };

    const handleUserOffline = (payload: { user_id: string }) => {
      console.log('WS: User Offline:', payload.user_id);
      dispatch(setUserOnlineStatus({ userId: payload.user_id, isOnline: false }));
    };

    const handleConversationCreated = (conversation: Conversation) => {
      console.log('WS: Conversation Created:', conversation);
      dispatch(addConversation(conversation));
      // Also invalidate conversations query to ensure consistency
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
    };

    // Step 1: Attach all event listeners FIRST
    console.log('WS: Attaching event listeners...');
    socketService.on('new_message', handleNewMessage);
    socketService.on('message_sent', handleMessageSent);
    socketService.on('user_online', handleUserOnline);
    socketService.on('user_offline', handleUserOffline);
    socketService.on('conversation_created', handleConversationCreated);

    // Step 2: THEN connect (this ensures we don't miss any events)
    console.log('WS: Connecting to socket...');
    socketService.connect(token);

    // Cleanup
    return () => {
      console.log('WS: Cleaning up event listeners...');
      socketService.off('new_message', handleNewMessage);
      socketService.off('message_sent', handleMessageSent);
      socketService.off('user_online', handleUserOnline);
      socketService.off('user_offline', handleUserOffline);
      socketService.off('conversation_created', handleConversationCreated);
      // Note: We don't disconnect here to allow socket to persist across renders
      // Only disconnect when auth state changes (handled by the condition at the top)
    };
  }, [isAuthenticated, token, user, dispatch, queryClient]);
}
