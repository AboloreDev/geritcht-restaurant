// src/app/admin/components/InventoryAlerts.tsx
"use client";

import Link from "next/link";
import { HugeiconsIcon } from "@hugeicons/react";
import { Alert02Icon, CancelCircleIcon } from "@hugeicons/core-free-icons";
import { useGetInventoryAlertsQuery } from "@/app/state/api/ingredientApi";

export default function InventoryAlerts() {
  const { data, isLoading } = useGetInventoryAlertsQuery();

  const lowStock = data?.data.low_stock_ingredients ?? [];
  const outOfStock = data?.data.out_of_stock_items ?? [];

  if (isLoading) {
    return <div className="h-16 animate-pulse rounded-xl bg-muted" />;
  }

  if (lowStock.length === 0 && outOfStock.length === 0) return null;

  return (
    <div className="space-y-3">
      {outOfStock.length > 0 && (
        <div className="rounded-xl border border-red-300 bg-red-50 p-4">
          <div className="flex items-start gap-3">
            <HugeiconsIcon
              icon={CancelCircleIcon}
              size={20}
              className="mt-0.5 shrink-0 text-red-600"
            />
            <div className="flex-1">
              <p className="text-sm font-medium text-red-900">
                {outOfStock.length} dish{outOfStock.length !== 1 ? "es" : ""}{" "}
                disabled — out of stock
              </p>
              <div className="mt-2 flex flex-wrap gap-2">
                {outOfStock.slice(0, 6).map((item) => (
                  <span
                    key={item.id}
                    className="rounded-full bg-red-100 px-2.5 py-1 text-xs font-medium text-red-800"
                  >
                    {item.name}
                  </span>
                ))}
                {outOfStock.length > 6 && (
                  <span className="text-xs text-red-700">
                    +{outOfStock.length - 6} more
                  </span>
                )}
              </div>
            </div>
            <Link
              href="/admin/menu"
              className="shrink-0 text-xs font-medium text-red-700 hover:underline"
            >
              Manage
            </Link>
          </div>
        </div>
      )}

      {lowStock.length > 0 && (
        <div className="rounded-xl p-4">
          <div className="flex items-start gap-3">
            <HugeiconsIcon
              icon={Alert02Icon}
              size={20}
              className="mt-0.5 shrink-0 text-amber-600"
            />
            <div className="flex-1">
              <p className="text-sm font-medium text-amber-900">
                {lowStock.length} ingredient{lowStock.length !== 1 ? "s" : ""}{" "}
                running low
              </p>
              <div className="mt-2 flex flex-wrap gap-2">
                {lowStock.slice(0, 6).map((ing) => (
                  <span
                    key={ing.id}
                    className="rounded-full bg-amber-100 px-2.5 py-1 text-xs font-medium text-amber-800"
                  >
                    {ing.name} ({ing.current_stock} {ing.unit})
                  </span>
                ))}
                {lowStock.length > 6 && (
                  <span className="text-xs text-amber-700">
                    +{lowStock.length - 6} more
                  </span>
                )}
              </div>
            </div>
            <Link
              href="/admin/ingredients"
              className="shrink-0 text-xs font-medium text-amber-700 hover:underline"
            >
              Manage
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}
