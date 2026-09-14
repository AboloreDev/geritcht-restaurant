"use client";

import RecipeRow from "./RecipeRow";
import AddRecipeIngredient from "./AddRecipeIngredient";
import { useGetRecipesQuery } from "@/app/state/api/recipesApi";

export default function RecipeSection({ menuId }: { menuId: number }) {
  const { data, isLoading } = useGetRecipesQuery({ menuId });
  const recipes = data?.data ?? [];

  return (
    <section className="space-y-4 border-t pt-6">
      <div>
        <h3 className="text-lg font-semibold">Recipe</h3>
        <p className="text-sm text-muted-foreground">
          Ingredients used per serving — drives automatic stock deduction.
        </p>
      </div>

      {isLoading ? (
        <div className="space-y-2">
          {Array.from({ length: 2 }).map((_, i) => (
            <div key={i} className="h-11 animate-pulse rounded-lg bg-muted" />
          ))}
        </div>
      ) : recipes.length === 0 ? (
        <p className="text-sm text-muted-foreground">
          No ingredients linked yet.
        </p>
      ) : (
        <div className="space-y-2">
          {recipes.map((recipe) => (
            <RecipeRow
              key={recipe.ingredient_id}
              menuId={menuId}
              recipe={recipe}
            />
          ))}
        </div>
      )}

      <AddRecipeIngredient
        menuId={menuId}
        linkedIngredientIds={recipes.map((r) => r.ingredient_id)}
      />
    </section>
  );
}
