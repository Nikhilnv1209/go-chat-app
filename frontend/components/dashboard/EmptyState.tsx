import { MessageSquare } from 'lucide-react';

interface EmptyStateProps {
  title?: string;
  description?: string;
  icon?: React.ElementType;
}

export default function EmptyState({
  title = "Select a conversation",
  description = "Choose a chat from the sidebar to start messaging, or create a new conversation to get started.",
  icon: Icon = MessageSquare
}: EmptyStateProps) {
  return (
    <div className="flex items-center justify-center h-full bg-[#f9fafc]">
      {/* Empty State */}
      <div className="text-center max-w-md px-6">
        {/* Icon */}
        <div className="inline-flex items-center justify-center w-20 h-20 rounded-full bg-gradient-to-br from-[#7678ed]/20 to-[#7678ed]/10 border border-[#7678ed]/20 mb-6">
          <Icon className="w-10 h-10 text-[#7678ed]" />
        </div>

        {/* Title */}
        <h2 className="text-2xl font-bold text-[#202022] mb-3">
          {title}
        </h2>

        {/* Description */}
        <p className="text-[#202022]/60 mb-8">
          {description}
        </p>

        {/* Quick Tips (Static for now) */}
        <div className="space-y-3 text-left">
          <div className="flex items-start gap-3 p-3 rounded-xl bg-white border border-[#7678ed]/10 shadow-sm">
            <div className="w-2 h-2 rounded-full bg-[#7678ed] mt-1.5"></div>
            <div>
              <p className="text-sm font-medium text-[#202022]">Real-time messaging</p>
              <p className="text-xs text-[#202022]/50">Messages appear instantly</p>
            </div>
          </div>
          <div className="flex items-start gap-3 p-3 rounded-xl bg-white border border-[#7678ed]/10 shadow-sm">
            <div className="w-2 h-2 rounded-full bg-[#ff7a55] mt-1.5"></div>
            <div>
              <p className="text-sm font-medium text-[#202022]">Read receipts</p>
              <p className="text-xs text-[#202022]/50">See when messages are delivered and read</p>
            </div>
          </div>
          <div className="flex items-start gap-3 p-3 rounded-xl bg-white border border-[#7678ed]/10 shadow-sm">
            <div className="w-2 h-2 rounded-full bg-[#7678ed] mt-1.5"></div>
            <div>
              <p className="text-sm font-medium text-[#202022]">Typing indicators</p>
              <p className="text-xs text-[#202022]/50">Know when someone is typing</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
