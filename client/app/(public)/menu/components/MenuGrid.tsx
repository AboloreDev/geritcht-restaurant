"use client";

import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { useGetMenusQuery } from "@/app/state/api/menuApi";
import { setPage } from "@/app/state/slices/menuSlice";
import { MenuCard } from "./MenuCard";

export function MenuGrid() {
  const dispatch = useAppDispatch();
  const { categoryId, page, limit, query, sortBy, order } = useAppSelector(
    (state: RootState) => state.menu,
  );

  const { data, isFetching, isLoading } = useGetMenusQuery({
    category_id: categoryId,
    page,
    limit,
    query: query || undefined,
    sort_by: sortBy,
    order,
  });

  const items = data?.data ?? [];
  const hasMore = data ? page < data.meta.total_pages : false;

  if (isLoading) return <MenuGridSkeleton />;

  if (!isLoading && items.length === 0) {
    return (
      <div className="py-16 text-center text-sm text-muted-foreground">
        No dishes match that. Try another category or search.
      </div>
    );
  }

  return (
    <div>
      <motion.div
        layout
        className="grid grid-cols-1 gap-10 md:grid-cols-2 xl:grid-cols-3"
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
                // only the newest page's cards stagger in; cards already
                // on screen don't replay their entrance on re-render
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
            {isFetching ? "Loading…" : "Load more"}
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
