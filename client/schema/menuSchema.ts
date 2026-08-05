import { z } from "zod";

const menuBaseSchema = z.object({
  category_id: z
    .number({
      error: "Please select a category.",
    })
    .min(1, "Please select a category."),

  name: z
    .string()
    .trim()
    .min(2, "Menu name must be at least 2 characters.")
    .max(100, "Menu name cannot exceed 100 characters."),

  description: z
    .string()
    .trim()
    .max(500, "Description cannot exceed 500 characters.")
    .optional()
    .or(z.literal("")),

  price: z
    .number({
      error: "Price is required.",
    })
    .positive("Price must be greater than zero."),

  prep_time_minutes: z
    .number()
    .min(0, "Preparation time cannot be negative.")
    .max(180, "Preparation time seems too high.")
    .default(0),

  spice_level: z.number().min(0).max(5).default(0),

  allergen_ids: z.array(z.number()).default([]),

  dietary_tag_ids: z.array(z.number()).default([]),

  display_order: z.number().min(0).default(0),
});

export const createMenuSchema = menuBaseSchema;

export const updateMenuSchema = menuBaseSchema.extend({
  is_available: z.boolean().optional(),
});

export type CreateMenuFormValues = z.infer<typeof createMenuSchema>;

export type UpdateMenuFormValues = z.infer<typeof updateMenuSchema>;

export const createMenuDefaultValues: CreateMenuFormValues = {
  category_id: 0,
  name: "",
  description: "",
  price: 0,
  prep_time_minutes: 0,
  spice_level: 0,
  allergen_ids: [],
  dietary_tag_ids: [],
  display_order: 0,
};

export const updateMenuDefaultValues: UpdateMenuFormValues = {
  ...createMenuDefaultValues,
  is_available: true,
};
