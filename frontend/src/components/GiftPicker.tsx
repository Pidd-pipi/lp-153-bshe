import { GIFT_EMOJIS } from "@/constants";

interface GiftPickerProps {
  value: string;
  onChange: (emoji: string) => void;
}

// 虚拟礼物表情选择器：祝福留言板复用。
export default function GiftPicker({ value, onChange }: GiftPickerProps) {
  return (
    <div className="flex flex-wrap gap-2">
      {GIFT_EMOJIS.map((emoji) => (
        <button
          key={emoji}
          type="button"
          onClick={() => onChange(emoji === value ? "" : emoji)}
          className={`h-9 w-9 rounded-full text-lg transition ${value === emoji ? "bg-pink-100 ring-2 ring-pink-400" : "bg-purple-50 hover:bg-purple-100"}`}
        >
          {emoji}
        </button>
      ))}
    </div>
  );
}
