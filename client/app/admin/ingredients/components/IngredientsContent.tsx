"use client";

import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  Search01Icon,
  Add01Icon,
  CheeseIcon,
} from "@hugeicons/core-free-icons";
import { cn } from "@/lib/utils";
import { useDebounce } from "@/app/hooks/useDebounce";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import {
  setPage,
  setSearch,
  openCreateIngredientSheet,
} from "@/app/state/slices/ingredientSlice";
import {
  useGetAllIngredientsQuery,
  useSearchIngredientQuery,
} from "@/app/state/api/ingredientApi";
import { Button } from "@/components/ui/button";
import IngredientRow from "./IngredientRow";
import CreateIngredientSheet from "./CreateIngredientSheet";
import EditIngredientSheet from "./EditIngredientSheet";
import DeleteIngredientDialog from "./DeleteIngredientDialog";

export default function IngredientsContent() {
  const dispatch = useAppDispatch();
  const { page, query } = useAppSelector((s: RootState) => s.ingredient);

  const [searchInput, setSearchInput] = useState(query);
  const debouncedSearch = useDebounce(searchInput, 400);

  useEffect(() => {
    dispatch(setSearch(debouncedSearch));
  }, [debouncedSearch, dispatch]);

  const hasQuery = query.trim().length > 0;

  const listQuery = useGetAllIngredientsQuery(
    { page: 1, limit: 15 },
    { skip: hasQuery },
  );
  const searchQuery = useSearchIngredientQuery(
    { q: query, page, limit: 15 },
    { skip: !hasQuery },
  );

  const { data, isLoading, isFetching } = hasQuery ? searchQuery : listQuery;
  const ingredients = data?.data ?? [];

  console.log(ingredients);
  const hasMore =
    !hasQuery && data && "meta" in data ? page < data.meta.total_pages : false;

  return (
    <div className="p-6">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h1 className="font-serif text-2xl font-semibold">Ingredients</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Track stock levels and low-stock thresholds.
          </p>
        </div>
        <Button onClick={() => dispatch(openCreateIngredientSheet())}>
          <HugeiconsIcon icon={Add01Icon} strokeWidth={2} size={16} />
          Add Ingredient
        </Button>
      </div>

      <div className="mt-6 flex items-center gap-3 rounded-2xl px-6 py-3">
        <div className="relative flex-1">
          <HugeiconsIcon
            icon={Search01Icon}
            size={16}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
          />
          <input
            type="text"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            placeholder="Search ingredients…"
            className="w-full rounded-full border border-black py-1.5 pl-9 pr-3 text-sm"
          />
        </div>
        {isFetching && !isLoading && (
          <span className="text-xs text-muted-foreground">Updating…</span>
        )}
      </div>

      {isLoading ? (
        <div className="mt-6 space-y-3">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="h-20 animate-pulse rounded-2xl bg-muted" />
          ))}
        </div>
      ) : ingredients.length === 0 ? (
        <div className="mt-16 flex flex-col items-center gap-3 text-center">
          <HugeiconsIcon
            icon={CheeseIcon}
            size={40}
            className="text-muted-foreground"
          />
          <p className="text-sm text-muted-foreground">No ingredients found.</p>
        </div>
      ) : (
        <>
          <div
            className={cn(
              "mt-6 space-y-3 transition-opacity",
              isFetching && !isLoading && "opacity-50",
            )}
          >
            <AnimatePresence initial={false}>
              {ingredients.map((ingredient) => (
                <IngredientRow key={ingredient.id} ingredient={ingredient} />
              ))}
            </AnimatePresence>
          </div>

          {hasMore && (
            <div className="mt-6 flex justify-center">
              <Button
                variant="outline"
                disabled={isFetching}
                onClick={() => dispatch(setPage(page + 1))}
              >
                {isFetching && page > 1 ? "Loading…" : "Load more"}
              </Button>
            </div>
          )}
        </>
      )}

      <CreateIngredientSheet />
      <EditIngredientSheet />
      <DeleteIngredientDialog />
    </div>
  );
}
