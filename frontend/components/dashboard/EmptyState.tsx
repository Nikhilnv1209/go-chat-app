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
    <div className="flex items-center justify-center h-full bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
      {/* Empty State */}
      <div className="text-center max-w-md px-6">
        {/* Icon */}
        <div className="inline-flex items-center justify-center w-20 h-20 rounded-full bg-gradient-to-br from-indigo-500/20 to-purple-600/20 border border-indigo-500/30 mb-6">
          <Icon className="w-10 h-10 text-indigo-400" />
        </div>

        {/* Title */}
        <h2 className="text-2xl font-bold text-white mb-3">
          {title}
        </h2>

        {/* Description */}
        <p className="text-slate-400 mb-8">
          {description}
        </p>

        {/* Quick Tips (Static for now) */}
        <div className="space-y-3 text-left">
          <div className="flex items-start gap-3 p-3 rounded-lg bg-white/[0.03] border border-white/[0.05]">
            <div className="w-2 h-2 rounded-full bg-indigo-500 mt-1.5"></div>
            <div>
              <p className="text-sm font-medium text-slate-300">Real-time messaging</p>
              <p className="text-xs text-slate-500">Messages appear instantly</p>
            </div>
          </div>
          <div className="flex items-start gap-3 p-3 rounded-lg bg-white/[0.03] border border-white/[0.05]">
            <div className="w-2 h-2 rounded-full bg-purple-500 mt-1.5"></div>
            <div>
              <p className="text-sm font-medium text-slate-300">Read receipts</p>
              <p className="text-xs text-slate-500">See when messages are delivered and read</p>
            </div>
          </div>
          <div className="flex items-start gap-3 p-3 rounded-lg bg-white/[0.03] border border-white/[0.05]">
            <div className="w-2 h-2 rounded-full bg-pink-500 mt-1.5"></div>
            <div>
              <p className="text-sm font-medium text-slate-300">Typing indicators</p>
              <p className="text-xs text-slate-500">Know when someone is typing</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
