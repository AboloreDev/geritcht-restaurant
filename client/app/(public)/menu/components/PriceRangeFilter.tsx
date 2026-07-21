"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Filter } from "@mynaui/icons-react";
import { Button } from "@/components/ui/button";
import { setPriceRange } from "@/app/state/slices/menuSlice";
import {
  useAppDispatch,
  useAppSelector,
  type RootState,
} from "@/app/state/redux";
import { cn } from "@/lib/utils";
import { formatNaira } from "@/app/utils/formatNaira";
import { Slider } from "@/components/ui/slider";

// Adjust these to match your actual menu's real price spread —
// slider bounds are hardcoded since GetMenusRequest has no
// "min/max available price" endpoint to derive them dynamically from.
const PRICE_FLOOR = 0;
const PRICE_CEILING = 300000;

export function PriceRangeFilter() {
  const dispatch = useAppDispatch();
  const [open, setOpen] = useState(false);
  const { minPrice, maxPrice } = useAppSelector((s: RootState) => s.menu);
  const [draft, setDraft] = useState<[number, number]>([
    minPrice ?? PRICE_FLOOR,
    maxPrice ?? PRICE_CEILING,
  ]);

  const isActive = minPrice !== undefined || maxPrice !== undefined;

  function handleOpenChange(next: boolean) {
    if (next) {
      setDraft([minPrice ?? PRICE_FLOOR, maxPrice ?? PRICE_CEILING]);
    }
    setOpen(next);
  }

  function handleApply() {
    dispatch(
      setPriceRange({
        minPrice: draft[0] > PRICE_FLOOR ? draft[0] : undefined,
        maxPrice: draft[1] < PRICE_CEILING ? draft[1] : undefined,
      }),
    );
    setOpen(false);
  }

  function handleClear() {
    setDraft([PRICE_FLOOR, PRICE_CEILING]);
    dispatch(setPriceRange({ minPrice: undefined, maxPrice: undefined }));
    setOpen(false);
  }

  return (
    <div className="relative">
      <Button
        onClick={() => handleOpenChange(!open)}
        className={cn(
          "flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-sm transition-colors",
          isActive
            ? "border-primary bg-primary/5 text-primary"
            : "text-muted-foreground hover:text-foreground",
        )}
      >
        <Filter className="h-4 w-4" />
        {isActive
          ? `${formatNaira(minPrice ?? PRICE_FLOOR)} – ${formatNaira(maxPrice ?? PRICE_CEILING)}`
          : "Price"}
      </Button>

      {open && (
        <div
          className="fixed inset-0 z-10"
          onClick={() => handleOpenChange(false)}
        />
      )}

      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, y: -6, scale: 0.97 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -6, scale: 0.97 }}
            transition={{ duration: 0.15, ease: [0.16, 1, 0.3, 1] }}
            className="absolute right-0 top-11 z-20 w-72 rounded-lg border bg-[#fefae0] p-4 shadow-md"
          >
            <div className="flex items-center justify-between text-sm font-medium">
              <span>{formatNaira(draft[0])}</span>
              <span>{formatNaira(draft[1])}</span>
            </div>

            <Slider
              className="mt-4"
              min={PRICE_FLOOR}
              max={PRICE_CEILING}
              step={500}
              value={draft}
              onValueChange={(value) => setDraft(value as [number, number])}
            />

            <div className="mt-5 flex justify-between gap-2">
              <Button variant="ghost" size="sm" onClick={handleClear}>
                Clear
              </Button>
              <Button size="sm" onClick={handleApply}>
                Apply
              </Button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
