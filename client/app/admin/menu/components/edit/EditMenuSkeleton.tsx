export default function EditMenuSkeleton() {
  return (
    <div className="space-y-8 py-2">
      {/* form fields */}
      <div className="space-y-5">
        <div className="h-5 w-32 animate-pulse rounded bg-[#fefae0]" />
        <div className="h-10 w-full animate-pulse rounded-lg bg-[#fefae0]" />
        <div className="h-10 w-full animate-pulse rounded-lg bg-[#fefae0]" />
        <div className="h-24 w-full animate-pulse rounded-lg bg-[#fefae0]" />
      </div>

      <div className="grid gap-4 sm:grid-cols-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="h-16 animate-pulse rounded-lg bg-[#fefae0]" />
        ))}
      </div>

      {/* images section */}
      <div className="space-y-4 border-t pt-6">
        <div className="h-5 w-24 animate-pulse rounded bg-[#fefae0]" />
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div
              key={i}
              className="aspect-square animate-pulse rounded-2xl bg-[#fefae0]"
            />
          ))}
        </div>
      </div>
    </div>
  );
}
