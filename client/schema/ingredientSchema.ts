import { z } from "zod";

export const ingredientSchema = z.object({
  name: z.string().trim().min(1, "Ingredient name is required").max(80),
  unit: z.string().trim().min(1, "Unit is required").max(20),
  current_stock: z.coerce.number().min(0, "Cannot be negative").default(0),
  min_threshold: z.coerce.number().min(0, "Cannot be negative"),
});

export type IngredientFormValues = z.infer<typeof ingredientSchema>;

export const ingredientDefaultValues: IngredientFormValues = {
  name: "",
  unit: "",
  current_stock: 0,
  min_threshold: 0,
};
