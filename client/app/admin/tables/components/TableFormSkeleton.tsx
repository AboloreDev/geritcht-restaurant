export default function TableFormSkeleton() {
  return (
    <div className="space-y-5 py-2">
      {Array.from({ length: 5 }).map((_, i) => (
        <div key={i} className="space-y-2">
          <div className="h-4 w-24 animate-pulse rounded bg-[#fefae0]" />
          <div className="h-10 w-full animate-pulse rounded-lg bg-[#fefae0]" />
        </div>
      ))}
    </div>
  );
}
