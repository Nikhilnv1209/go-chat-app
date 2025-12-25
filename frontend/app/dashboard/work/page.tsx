import EmptyState from '@/components/dashboard/EmptyState';
import { Briefcase } from 'lucide-react';

export default function WorkPage() {
  return (
    <EmptyState
      title="Work Chats"
      description="Select a work-related conversation from the sidebar."
      icon={Briefcase}
    />
  );
}
