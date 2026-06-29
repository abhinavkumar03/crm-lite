import { ReactNode } from "react";

type Props = {
  children: ReactNode;
  footer?: ReactNode;
};

export default function AuthCard({
  children,
  footer,
}: Props) {
  return (
    <div
      className="
      overflow-hidden
      rounded-3xl
      border
      border-slate-200
      bg-white
      shadow-xl
      shadow-slate-200/50
      "
    >
      {/* Body */}

      <div className="p-8 sm:p-10">
        {children}
      </div>

      {/* Footer */}

      {footer && (
        <div
          className="
          border-t
          border-slate-100
          bg-slate-50
          px-8
          py-5
          text-center
          "
        >
          {footer}
        </div>
      )}
    </div>
  );
}