export default function TableListSkeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 8 }).map((_, i) => (
        <div key={i} className="h-20 animate-pulse rounded-2xl bg-[#fefae0]" />
      ))}
    </div>
  );
}
