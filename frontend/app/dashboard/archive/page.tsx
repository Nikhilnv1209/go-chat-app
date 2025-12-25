import EmptyState from '@/components/dashboard/EmptyState';
import { Archive } from 'lucide-react';

export default function ArchivePage() {
  return (
    <EmptyState
      title="Archived Chats"
      description="Select an archived conversation to view history."
      icon={Archive}
    />
  );
}
