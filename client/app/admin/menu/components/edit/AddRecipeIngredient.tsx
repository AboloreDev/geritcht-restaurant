"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  useGetAllIngredientsQuery,
  useSearchIngredientQuery,
} from "@/app/state/api/ingredientApi";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import { useDebounce } from "@/app/hooks/useDebounce";
import { useAddRecipeMutation } from "@/app/state/api/recipesApi";

export default function AddRecipeIngredient({
  menuId,
  linkedIngredientIds,
}: {
  menuId: number;
  linkedIngredientIds: number[];
}) {
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 300);
  const hasQuery = debouncedSearch.trim().length > 0;

  const { data: listData } = useGetAllIngredientsQuery(
    { page: 1, limit: 50 },
    { skip: hasQuery },
  );
  const { data: searchData } = useSearchIngredientQuery(
    { q: debouncedSearch },
    { skip: !hasQuery },
  );

  const options = (hasQuery ? searchData?.data : listData?.data) ?? [];
  const availableOptions = options.filter(
    (i) => !linkedIngredientIds.includes(i.id),
  );

  const [selectedId, setSelectedId] = useState<string>("");
  const [quantity, setQuantity] = useState("");

  const [addRecipe, { isLoading }] = useAddRecipeMutation();

  async function handleAdd() {
    if (!selectedId || !quantity || Number(quantity) <= 0) {
      toast.error("Select an ingredient and enter a valid quantity.");
      return;
    }
    try {
      const response = await addRecipe({
        menuId,
        body: { ingredient_id: Number(selectedId), quantity: Number(quantity) },
      }).unwrap();
      toast.success(response.message);
      setSelectedId("");
      setQuantity("");
      setSearch("");
    } catch (error) {
      toast.error(getApiError(error));
    }
  }

  return (
    <div className="flex flex-wrap items-end gap-2 rounded-xl p-3">
      <div className=" flex-1">
        {/*   @ts-expect-error "<>" */}
        <Select value={selectedId} onValueChange={setSelectedId}>
          <SelectTrigger>
            <SelectValue
              placeholder="Search ingredient…"
              className="placeholder:text-slate-500 text-black"
            />
          </SelectTrigger>
          <SelectContent className="bg-white h-[500px]">
            <div className="p-2">
              <Input
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Type to search…"
                onKeyDown={(e) => e.stopPropagation()}
                className="h-8 placeholder:text-slate-500"
              />
            </div>
            {availableOptions.length === 0 ? (
              <p className="p-2 text-xs text-muted-foreground">
                No matching ingredients.
              </p>
            ) : (
              availableOptions.map((ing) => (
                <SelectItem key={ing.id} value={String(ing.id)}>
                  {ing.name} ({ing.unit})
                </SelectItem>
              ))
            )}
          </SelectContent>
        </Select>
      </div>

      <Input
        type="number"
        step="0.01"
        value={quantity}
        onChange={(e) => setQuantity(e.target.value)}
        placeholder="Qty per serving"
        className="w-32 placeholder:text-slate-500"
      />

      <Button size="sm" disabled={isLoading} onClick={handleAdd}>
        {isLoading ? "Adding…" : "Add"}
      </Button>
    </div>
  );
}
