"use client";

import { motion } from "framer-motion";
import { cn } from "@/lib/utils";
import { useAppDispatch, useAppSelector } from "@/app/state/redux";
import { useGetCategoriesQuery } from "@/app/state/api/categoriesApi";
import { setCategory } from "@/app/state/slices/menuSlice";

const container = {
  hidden: {},
  show: {
    transition: { staggerChildren: 0.05 },
  },
};

const pillVariant = {
  hidden: { opacity: 0, y: 12 },
  show: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.35, ease: [0.16, 1, 0.3, 1] },
  },
};

export function CategoryTabs() {
  const dispatch = useAppDispatch();
  const activeId = useAppSelector((state) => state.menu.categoryId);
  const { data, isLoading } = useGetCategoriesQuery({ page: 1, limit: 100 });

  if (isLoading) return <CategoryTabsSkeleton />;

  return (
    <motion.div
      variants={container}
      initial="hidden"
      animate="show"
      className="flex gap-2 overflow-x-auto justify-center"
    >
      {data?.data.map((cat) => {
        const isActive = activeId === cat.id;
        return (
          <motion.button
            key={cat.id}
            // @ts-expect-error "<>"
            variants={pillVariant}
            onClick={() => {
              console.log("clicked", cat.id);
              dispatch(setCategory(cat.id));
            }}
            className={cn(
              "relative rounded-full px-4 py-1.5 cursor-pointer text-sm text-primary whitespace-nowrap transition-colors",
              isActive ? "text-white" : "text-primary-deep hover:text-white",
            )}
          >
            {isActive && (
              <motion.span
                layoutId="active-pill"
                className="absolute inset-0 rounded-full bg-primary"
                transition={{ type: "spring", stiffness: 400, damping: 32 }}
              />
            )}
            <span className="relative z-10">{cat.name}</span>
          </motion.button>
        );
      })}
    </motion.div>
  );
}

function CategoryTabsSkeleton() {
  return (
    <div className="flex gap-2">
      {Array.from({ length: 5 }).map((_, i) => (
        <div
          key={i}
          className="h-8 w-20 animate-pulse rounded-full bg-[#fefee3]"
        />
      ))}
    </div>
  );
}
