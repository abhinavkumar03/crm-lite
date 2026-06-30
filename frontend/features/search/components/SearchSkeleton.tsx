export default function SearchSkeleton() {
  return (
    <div className="space-y-2 p-4">
      {Array.from({
        length: 6,
      }).map((_, i) => (
        <div
          key={i}
          className="
          h-14
          animate-pulse
          rounded-2xl
          bg-slate-100
          "
        />
      ))}
    </div>
  );
}