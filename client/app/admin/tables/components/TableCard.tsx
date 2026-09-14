"use client";

import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  MoreVerticalIcon,
  Edit02Icon,
  Delete02Icon,
  UserGroupIcon,
  Location01Icon,
} from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import { Table } from "@/app/state/types/tableTypes";
import { useAppDispatch } from "@/app/state/redux";
import {
  openEditTableSheet,
  openDeleteTableDialog,
} from "@/app/state/slices/tableSlice";
import { tableStatusStyle } from "@/app/utils/tableStatusHelpers";

export default function TableRow({ table }: { table: Table }) {
  const dispatch = useAppDispatch();
  const router = useRouter();

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.2 }}
      onClick={() => router.push(`/admin/tables/${table.id}`)}
      className="flex cursor-pointer items-center justify-between gap-4 rounded-2xl bg-[#fefae0] p-4 transition-colors hover:bg-[#fefae0]/50"
    >
      {/* Left — identity */}
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <p className="font-serif font-semibold">{table.name}</p>
        </div>

        <div className="mt-1 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted-foreground">
          <div className="flex items-center gap-1.5 text-slate-500 ">
            <HugeiconsIcon icon={UserGroupIcon} strokeWidth={2} size={14} />
            Seats {table.capacity}
          </div>
          {table.location && (
            <div className="flex items-center gap-1.5">
              <HugeiconsIcon icon={Location01Icon} strokeWidth={2} size={14} />
              {table.location}
            </div>
          )}
        </div>
      </div>

      <Badge
        className={`rounded-full text-white ${tableStatusStyle(table.status)}`}
      >
        {table.status}
      </Badge>

      {/* Middle — current activity, if any */}
      <div className="hidden shrink-0 sm:block">
        {table.current_reservation && (
          <span className="rounded-lg bg-amber-50 px-3 py-1.5 text-xs text-amber-800">
            Reserved · {table.current_reservation.party_size} guests
          </span>
        )}
        {table.current_order && (
          <span className="rounded-lg bg-red-50 px-3 py-1.5 text-xs text-red-800">
            Order #{table.current_order.id}
          </span>
        )}
      </div>

      {/* Right — actions */}
      <div className="shrink-0" onClick={(e) => e.stopPropagation()}>
        <DropdownMenu>
          <DropdownMenuTrigger className="cursor-pointer p-1">
            <HugeiconsIcon icon={MoreVerticalIcon} strokeWidth={2} />
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="bg-white">
            <DropdownMenuItem
              onClick={() => dispatch(openEditTableSheet(table.id))}
              className="cursor-pointer"
            >
              <HugeiconsIcon icon={Edit02Icon} strokeWidth={2} />
              Edit
            </DropdownMenuItem>
            <DropdownMenuItem
              className="cursor-pointer text-red-600"
              onClick={() => dispatch(openDeleteTableDialog(table.id))}
            >
              <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </motion.div>
  );
}
