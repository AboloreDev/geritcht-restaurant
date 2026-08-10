export function TableDetailSkeleton() {
  return (
    <div className="p-6">
      <div className="h-4 w-24 animate-pulse rounded bg-[#fefae0]" />
      <div className="mt-4 h-8 w-48 animate-pulse rounded bg-[#fefae0]" />
      <div className="mt-6 grid gap-4 sm:grid-cols-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-[#fefae0]" />
        ))}
      </div>
    </div>
  );
}
