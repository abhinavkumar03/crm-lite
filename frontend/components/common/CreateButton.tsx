import { Plus } from "lucide-react";

type Props = {
  text: string;
  onClick: () => void;
};

export default function CreateButton({
  text,
  onClick,
}: Props) {
  return (
    <button
      onClick={onClick}
      className="
      inline-flex
      items-center
      gap-2
      rounded-2xl
      bg-emerald-500
      px-5
      py-3
      font-medium
      text-white
      shadow-sm
      transition-all
      duration-300
      hover:-translate-y-0.5
      hover:bg-emerald-600
      "
    >
      <Plus size={18} />

      {text}
    </button>
  );
}