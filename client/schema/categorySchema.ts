import z from "zod";

export const editCategorySchema = z.object({
  name: z
    .string()
    .trim()
    .min(1, "Category name is required.")
    .max(80, "Category name must be 80 characters or fewer."),
  description: z
    .string()
    .trim()
    .max(500, "Description must be 500 characters or fewer."),
  image_url: z.union([
    z.literal(""),
    z.string().trim().url("Enter a valid image URL."),
  ]),
  display_order: z
    .number()
    .int("Display order must be a whole number.")
    .min(0, "Display order cannot be negative."),
  is_active: z.boolean(),
});

export type EditCategoryFormValues = z.infer<typeof editCategorySchema>;

export const createCategorySchema = z.object({
  name: z
    .string()
    .trim()
    .min(1, "Category name is required.")
    .max(80, "Category name must be 80 characters or fewer."),
  description: z
    .string()
    .trim()
    .max(500, "Description must be 500 characters or fewer."),
  image_url: z.union([
    z.literal(""),
    z.string().trim().url("Enter a valid image URL.").optional(),
  ]),
  display_order: z
    .number()
    .int("Display order must be a whole number.")
    .min(0, "Display order cannot be negative."),
  is_active: z.boolean(),
});

export type CreateCategoryFormValues = z.infer<typeof createCategorySchema>;
