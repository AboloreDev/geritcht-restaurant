"use client";

import { motion, AnimatePresence } from "framer-motion";
import { Spinner } from "@mynaui/icons-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { useGetMenusQuery } from "@/app/state/api/menuApi";
import { setPage } from "@/app/state/slices/menuSlice";
import { MenuCard } from "./MenuCard";

export function MenuGrid() {
  const dispatch = useAppDispatch();
  const {
    categoryId,
    page,
    limit,
    query,
    sortBy,
    sortOrder,
    maxPrice,
    minPrice,
  } = useAppSelector((state: RootState) => state.menu);

  const { data, isFetching, isLoading } = useGetMenusQuery({
    category_id: categoryId,
    page,
    limit,
    query: query || undefined,
    sort_by: sortBy,
    sort_order: sortOrder,
    max_price: maxPrice,
    min_price: minPrice,
  });

  const items = data?.data ?? [];
  const hasMore = data ? page < data.meta.total_pages : false;

  // true only on this component's very first ever fetch
  if (isLoading) return <MenuGridSkeleton />;

  // refetching due to category/sort/price/search change — always resets
  // to page 1, so this never fires for a load-more (page > 1) fetch
  const isRefetchingFilters = isFetching && page === 1;

  if (!isLoading && !isFetching && items.length === 0) {
    return (
      <div className="py-16 text-center text-sm text-muted-foreground">
        No dishes match that. Try another category or search.
      </div>
    );
  }

  return (
    <div id="menu-grid" className="relative">
      {/* dims existing grid + shows a status pill while filters refetch,
          instead of flashing back to a skeleton or empty state */}
      <AnimatePresence>
        {isRefetchingFilters && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.15 }}
            className="absolute inset-0 z-10 flex items-start justify-center  pt-16"
          >
            <div className="flex items-center gap-2 rounded-full bg-[#fefae0] px-4 py-2 text-sm shadow-sm">
              <Spinner className="h-4 w-4 animate-spin " />
              Updating menu…
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      <motion.div
        layout
        className={cn(
          "grid grid-cols-2 gap-10 md:grid-cols-3 xl:grid-cols-4 transition-opacity",
          isRefetchingFilters && "opacity-50",
        )}
      >
        <AnimatePresence initial={false}>
          {items.map((item, i) => (
            <motion.div
              key={item.id}
              layout
              initial={{ opacity: 0, y: 24 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{
                duration: 0.4,
                ease: [0.16, 1, 0.3, 1],
                delay: (i % limit) * 0.03,
              }}
            >
              <MenuCard menu={item} />
            </motion.div>
          ))}
        </AnimatePresence>
      </motion.div>

      {hasMore && (
        <div className="mt-8 flex justify-center">
          <Button
            variant="outline"
            disabled={isFetching}
            onClick={() => dispatch(setPage(page + 1))}
          >
            {isFetching && page > 1 ? "Loading…" : "Load more"}
          </Button>
        </div>
      )}
    </div>
  );
}

function MenuGridSkeleton() {
  return (
    <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
      {Array.from({ length: 8 }).map((_, i) => (
        <div key={i} className="h-56 animate-pulse rounded-2xl bg-[#fefee3]" />
      ))}
    </div>
  );
}
