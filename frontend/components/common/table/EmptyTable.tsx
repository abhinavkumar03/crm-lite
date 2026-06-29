import {
  Inbox,
  Plus,
} from "lucide-react";

type Props = {
  title: string;
  description: string;
  actionLabel?: string;
  onAction?: () => void;
};

export default function EmptyTable({
  title,
  description,
  actionLabel,
  onAction,
}: Props) {
  return (
    <tr>
      <td
        colSpan={100}
        className="py-20"
      >
        <div className="flex flex-col items-center">
          <div className="rounded-3xl bg-slate-100 p-5">
            <Inbox
              size={42}
              className="text-slate-400"
            />
          </div>

          <h3 className="mt-6 text-xl font-semibold text-slate-900">
            {title}
          </h3>

          <p className="mt-2 max-w-md text-center text-slate-500">
            {description}
          </p>

          {actionLabel && onAction && (
            <button
              onClick={onAction}
              className="
                mt-8
                inline-flex
                items-center
                gap-2
                rounded-2xl
                bg-emerald-500
                px-5
                py-3
                font-medium
                text-white
                transition
                hover:bg-emerald-600
              "
            >
              <Plus size={18} />

              {actionLabel}
            </button>
          )}
        </div>
      </td>
    </tr>
  );
}