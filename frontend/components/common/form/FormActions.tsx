type Props = {
  loading?: boolean;
  submitText: string;
  cancelText?: string;
  onCancel?: () => void;
};

export default function FormActions({
  loading = false,
  submitText,
  cancelText = "Cancel",
  onCancel,
}: Props) {
  return (
    <div className="flex flex-col-reverse gap-3 border-t border-slate-200 pt-6 sm:flex-row sm:justify-end">
      {onCancel && (
        <button
          type="button"
          onClick={onCancel}
          className="
            rounded-2xl
            border
            border-slate-300
            bg-white
            px-6
            py-3
            font-medium
            text-slate-700
            transition
            hover:bg-slate-50
          "
        >
          {cancelText}
        </button>
      )}

      <button
        type="submit"
        data-tutorial-action="create-record"
        disabled={loading}
        className="
          rounded-2xl
          bg-emerald-500
          px-6
          py-3
          font-medium
          text-white
          shadow-sm
          transition
          hover:bg-emerald-600
          disabled:cursor-not-allowed
          disabled:opacity-60
        "
      >
        {loading
          ? "Saving..."
          : submitText}
      </button>
    </div>
  );
}