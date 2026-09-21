interface EmptyStateProps {
  title?: string;
  description?: string;
  icon?: string;
}

// 空状态占位：所有列表页复用。
export default function EmptyState({ title = "这里还空空的", description = "暂时没有内容", icon = "🫧" }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <div className="text-5xl mb-3 animate-float">{icon}</div>
      <p className="text-gray-600 font-medium">{title}</p>
      <p className="text-gray-400 text-sm mt-1">{description}</p>
    </div>
  );
}
