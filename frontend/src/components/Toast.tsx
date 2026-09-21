import { createContext, useCallback, useContext, useRef, useState, type ReactNode } from "react";

interface ToastContextValue {
  show: (message: string, type?: "success" | "error") => void;
}

const ToastContext = createContext<ToastContextValue>({ show: () => {} });

export function useToast() {
  return useContext(ToastContext);
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [msg, setMsg] = useState<{ id: number; text: string; type: "success" | "error" } | null>(null);
  const timer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const show = useCallback((text: string, type: "success" | "error" = "success") => {
    if (timer.current) clearTimeout(timer.current);
    setMsg({ id: Date.now(), text, type });
    timer.current = setTimeout(() => setMsg(null), 2600);
  }, []);

  return (
    <ToastContext.Provider value={{ show }}>
      {children}
      {msg && (
        <div className="fixed left-1/2 top-16 z-[100] -translate-x-1/2">
          <div
            key={msg.id}
            className={`rounded-full px-5 py-2.5 text-sm text-white shadow-lg ${
              msg.type === "error" ? "bg-red-500" : "bg-gray-800"
            }`}
          >
            {msg.text}
          </div>
        </div>
      )}
    </ToastContext.Provider>
  );
}
