"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  Delete02Icon,
  CheckmarkCircle01Icon,
} from "@hugeicons/core-free-icons";
import { MenuItemIngredient } from "@/app/state/types/ingredientTypes";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import {
  useDeleteRecipeMutation,
  useUpdateRecipeMutation,
} from "@/app/state/api/recipesApi";

export default function RecipeRow({
  menuId,
  recipe,
}: {
  menuId: number;
  recipe: MenuItemIngredient;
}) {
  const [quantity, setQuantity] = useState(String(recipe.quantity));
  const [isEditing, setIsEditing] = useState(false);

  const [updateRecipe, { isLoading: isUpdating }] = useUpdateRecipeMutation();
  const [deleteRecipe, { isLoading: isDeleting }] = useDeleteRecipeMutation();

  async function handleSave() {
    try {
      const response = await updateRecipe({
        menuId,
        ingredientId: recipe.ingredient_id,
        body: { quantity: Number(quantity) },
      }).unwrap();
      toast.success(response.message);
      setIsEditing(false);
    } catch (error) {
      toast.error(getApiError(error));
    }
  }

  async function handleDelete() {
    try {
      const response = await deleteRecipe({
        menuId,
        ingredientId: recipe.ingredient_id,
      }).unwrap();
      toast.success(response.message);
    } catch (error) {
      toast.error(getApiError(error));
    }
  }

  return (
    <div className="flex items-center justify-between gap-3 rounded-lg  px-3 py-2">
      <span className="text-sm font-medium">{recipe.ingredient.name}</span>

      <div className="flex items-center gap-2">
        {isEditing ? (
          <>
            <Input
              type="number"
              step="0.01"
              value={quantity}
              onChange={(e) => setQuantity(e.target.value)}
              className="h-8 w-20"
            />
            <span className="text-xs text-muted-foreground">
              {recipe.ingredient.unit}
            </span>
            <Button
              size="icon-sm"
              variant="ghost"
              disabled={isUpdating}
              onClick={handleSave}
            >
              <HugeiconsIcon
                icon={CheckmarkCircle01Icon}
                strokeWidth={2}
                size={16}
              />
            </Button>
          </>
        ) : (
          <button
            onClick={() => setIsEditing(true)}
            className="text-sm text-muted-foreground hover:text-foreground"
          >
            {recipe.quantity} {recipe.ingredient.unit}
          </button>
        )}

        <Button
          size="icon-sm"
          variant="ghost"
          className="text-destructive"
          disabled={isDeleting}
          onClick={handleDelete}
        >
          <HugeiconsIcon icon={Delete02Icon} strokeWidth={2} size={16} />
        </Button>
      </div>
    </div>
  );
}
