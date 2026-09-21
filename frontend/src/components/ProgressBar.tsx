interface ProgressBarProps {
  progress: number;
  label?: string;
}

// 圆梦进度条：心愿详情/我的认领复用。
export default function ProgressBar({ progress, label }: ProgressBarProps) {
  const value = Math.max(0, Math.min(100, progress));
  return (
    <div>
      {label && (
        <div className="mb-1 flex justify-between text-xs text-gray-500">
          <span>{label}</span>
          <span>{value}%</span>
        </div>
      )}
      <div className="h-2.5 w-full overflow-hidden rounded-full bg-purple-100">
        <div
          className="h-full rounded-full bg-gradient-to-r from-pink-400 to-purple-500 transition-all"
          style={{ width: `${value}%` }}
        />
      </div>
    </div>
  );
}
