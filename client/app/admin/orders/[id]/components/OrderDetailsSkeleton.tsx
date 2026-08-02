export function OrderDetailSkeleton() {
  return (
    <div className="p-6">
      <div className="h-4 w-32 animate-pulse rounded bg-muted" />
      <div className="mt-4 h-8 w-48 animate-pulse rounded bg-muted" />
      <div className="mt-6 grid gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2 h-64 animate-pulse rounded-xl bg-muted" />
        <div className="space-y-4">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="h-24 animate-pulse rounded-xl bg-muted" />
          ))}
        </div>
      </div>
    </div>
  );
}
