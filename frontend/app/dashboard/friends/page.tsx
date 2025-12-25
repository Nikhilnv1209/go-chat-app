import EmptyState from '@/components/dashboard/EmptyState';
import { Users2 } from 'lucide-react';

export default function FriendsPage() {
  return (
    <EmptyState
      title="Friends"
      description="Select a conversation with a friend from the sidebar."
      icon={Users2}
    />
  );
}
