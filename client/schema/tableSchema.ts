import { z } from "zod";

export const tableSchema = z.object({
  name: z.string().trim().min(1, "Table name is required").max(50),
  capacity: z.coerce.number().min(1, "Capacity must be at least 1").max(50),
  location: z.string().trim().max(100).optional().or(z.literal("")),
});

export type TableFormValues = z.infer<typeof tableSchema>;

export const tableDefaultValues: TableFormValues = {
  name: "",
  capacity: 0,
  location: "",
};
