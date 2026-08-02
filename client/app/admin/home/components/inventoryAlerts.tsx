"use client";

import Link from "next/link";
import { DangerTriangle, XCircle } from "@mynaui/icons-react";
import { useGetInventoryAlertsQuery } from "@/app/state/api/ingredientApi";
import { Ingredient } from "@/app/state/types/ingredientTypes";
import { Menu } from "@/app/state/types/menuTypes";

export function InventoryAlerts() {
  const { data, isLoading } = useGetInventoryAlertsQuery();

  const lowStock = data?.data.low_stock_ingredients ?? [];
  const outOfStock = data?.data.out_of_stock_items ?? [];

  if (isLoading) {
    return <div className="h-16 animate-pulse rounded-xl bg-[#faedcd]/70" />;
  }

  if (lowStock.length === 0 && outOfStock.length === 0) return null;

  return (
    <div className="space-y-3">
      {/* out of stock — more severe, since it's actually disabled dishes */}
      {outOfStock.length > 0 && (
        <div className="rounded-xl border border-red-300 bg-red-50 p-4">
          <div className="flex items-start gap-3">
            <XCircle className="mt-0.5 h-5 w-5 shrink-0 text-red-600" />
            <div className="flex-1">
              <p className="text-sm font-medium text-red-900">
                {outOfStock.length} dish{outOfStock.length !== 1 ? "es" : ""}{" "}
                disabled — out of stock
              </p>
              <div className="mt-2 flex flex-wrap gap-2">
                {outOfStock.slice(0, 6).map((item: Menu) => (
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

      {/* low stock — warning tier, still available but running out */}
      {lowStock.length > 0 && (
        <div className="rounded-xl border border-amber-300 bg-amber-50 p-4">
          <div className="flex items-start gap-3">
            <DangerTriangle className="mt-0.5 h-5 w-5 shrink-0 text-amber-600" />
            <div className="flex-1">
              <p className="text-sm font-medium text-amber-900">
                {lowStock.length} ingredient{lowStock.length !== 1 ? "s" : ""}{" "}
                running low
              </p>
              <div className="mt-2 flex flex-wrap gap-2">
                {lowStock.slice(0, 6).map((ing: Ingredient) => (
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
              href="/admin/inventory"
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
