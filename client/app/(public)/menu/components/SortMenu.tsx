// src/app/menu/components/SortMenu.tsx
"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useDispatch, useSelector } from "react-redux";
import { AlignHorizontal, Check } from "@mynaui/icons-react";
import { setSort } from "@/app/state/slices/menuSlice";
import {
  useAppDispatch,
  useAppSelector,
  type RootState,
} from "@/app/state/redux";
import { cn } from "@/lib/utils";

type SortOption = {
  label: string;
  sortBy: "price" | "name" | "created_at";
  order: "asc" | "desc";
};

const SORT_OPTIONS: SortOption[] = [
  { label: "Newest first", sortBy: "created_at", order: "desc" },
  { label: "Name (A–Z)", sortBy: "name", order: "asc" },
  { label: "Name (Z–A)", sortBy: "name", order: "desc" },
  { label: "Price (low to high)", sortBy: "price", order: "asc" },
  { label: "Price (high to low)", sortBy: "price", order: "desc" },
];

export function SortMenu() {
  const dispatch = useAppDispatch();
  const [open, setOpen] = useState(false);
  const { sortBy, order } = useAppSelector((s: RootState) => s.menu);

  const isActive = Boolean(sortBy);
  const activeLabel = SORT_OPTIONS.find(
    (o) => o.sortBy === sortBy && o.order === order,
  )?.label;

  function handleSelect(option: SortOption) {
    dispatch(setSort({ sortBy: option.sortBy, order: option.order }));
    setOpen(false);
  }

  return (
    <div className="relative">
      <button
        onClick={() => setOpen((v) => !v)}
        className={cn(
          "flex items-center cursor-pointer bg-primary gap-1.5 rounded-full border px-3 py-1.5 text-sm transition-colors",
          isActive
            ? " text-black"
            : "text-muted-foreground hover:text-foreground",
        )}
      >
        <AlignHorizontal className="h-4 w-4" />
        {activeLabel ?? "Sort"}
      </button>

      {/* click-outside overlay, matches desktop search pattern */}
      {open && (
        <div className="fixed inset-0 z-10" onClick={() => setOpen(false)} />
      )}

      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, y: -6, scale: 0.97 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -6, scale: 0.97 }}
            transition={{ duration: 0.15, ease: [0.16, 1, 0.3, 1] }}
            className="absolute right-0 top-11 z-20 w-56 rounded-lg border bg-[#fefae0] p-1 shadow-md"
          >
            {SORT_OPTIONS.map((option) => {
              const selected =
                option.sortBy === sortBy && option.order === order;
              return (
                <button
                  key={option.label}
                  onClick={() => handleSelect(option)}
                  className={cn(
                    "flex w-full items-center justify-between rounded-md px-3 py-2 text-left text-sm hover:bg-white cursor-pointer",
                    selected && "text-primary",
                  )}
                >
                  {option.label}
                  {selected && <Check className="h-4 w-4" />}
                </button>
              );
            })}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
