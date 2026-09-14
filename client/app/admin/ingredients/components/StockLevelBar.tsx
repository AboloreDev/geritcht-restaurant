"use client";

import { cn } from "@/lib/utils";

export default function StockLevelBar({
  current,
  threshold,
  unit,
}: {
  current: number;
  threshold: number;
  unit: string;
}) {
  // scale the bar against a ceiling of ~2x threshold, so "healthy" stock
  // (well above threshold) fills most of the bar without needing to
  // know the ingredient's actual max capacity
  const ceiling = Math.max(threshold * 2, current, 1);
  const fillPercent = Math.min((current / ceiling) * 100, 100);
  const thresholdPercent = Math.min((threshold / ceiling) * 100, 100);
  const isLow = current <= threshold;

  return (
    <div className="w-full">
      <div className="flex items-baseline justify-between text-xs">
        <span
          className={cn(
            "font-medium",
            isLow ? "text-red-600" : "text-foreground",
          )}
        >
          {current} {unit}
        </span>
        <span className="text-muted-foreground">min {threshold}</span>
      </div>

      <div className="relative mt-1.5 h-2 w-full overflow-hidden rounded-full bg-muted">
        <div
          className={cn(
            "absolute inset-y-0 left-0 rounded-full transition-all duration-500",
            isLow ? "bg-red-500" : "bg-emerald-500",
          )}
          style={{ width: `${fillPercent}%` }}
        />
        {/* threshold marker — a small tick showing where "low" begins */}
        <div
          className="absolute inset-y-0 w-[2px] bg-foreground/40"
          style={{ left: `${thresholdPercent}%` }}
        />
      </div>
    </div>
  );
}
