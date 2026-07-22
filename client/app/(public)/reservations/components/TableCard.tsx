"use client";

import { motion } from "framer-motion";
import { Users, MapPin } from "@mynaui/icons-react";
import { cn } from "@/lib/utils";
import { TableAvailability } from "@/app/state/types/reservationTypes";
import { getStatusConfig } from "@/app/utils/tableStatusConfig";

export function TableCard({
  table,
  onSelect,
}: {
  table: TableAvailability;
  onSelect: (table: TableAvailability) => void;
}) {
  const status = getStatusConfig(table.status);

  return (
    <motion.button
      layout
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={status.clickable ? { y: -2 } : undefined}
      transition={{ duration: 0.3, ease: [0.16, 1, 0.3, 1] }}
      disabled={!status.clickable}
      onClick={() => status.clickable && onSelect(table)}
      className={cn(
        "flex flex-col items-start rounded-xl border p-4 text-left transition-all",
        status.card,
      )}
    >
      <div className="flex w-full items-center justify-between">
        <span className="font-serif text-base font-medium">{table.name}</span>
        <span className={cn("h-5 w-5 rounded-full", status.dot)} />
      </div>

      <div className="mt-3 flex items-center gap-1.5 text-sm text-black">
        <Users className="h-3.5 w-3.5" />
        Seats {table.capacity}
      </div>

      <div className="mt-1 flex items-center gap-1.5 text-sm text-muted-foreground">
        <MapPin className="h-3.5 w-3.5" />
        {table.location}
      </div>

      <span className="mt-3 text-xs font-medium text-muted-foreground">
        {status.label}
      </span>
    </motion.button>
  );
}
