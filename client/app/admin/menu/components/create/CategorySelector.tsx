"use client";

import { Controller, useFormContext } from "react-hook-form";

import { Field, FieldError, FieldLabel } from "@/components/ui/field";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useGetCategoriesQuery } from "@/app/state/api/categoriesApi";
import { Spinner } from "@/components/ui/spinner";

export default function CategorySelect() {
  const { control } = useFormContext();
  const { data, isLoading } = useGetCategoriesQuery({
    page: 1,
    limit: 100,
  });

  const categories = data?.data ?? [];

  if (isLoading) {
    <div className="animate-spin w-10">
      <Spinner />
    </div>;
  }

  return (
    <Controller
      name="category_id"
      control={control}
      render={({ field, fieldState }) => (
        <Field data-invalid={fieldState.invalid}>
          <FieldLabel>Category</FieldLabel>

          <Select
            value={field.value ? String(field.value) : ""}
            onValueChange={(value) => field.onChange(Number(value))}
          >
            <SelectTrigger className="h-11 bg-white">
              <SelectValue placeholder="Select a category" />
            </SelectTrigger>

            <SelectContent className="bg-white shadow-none">
              {categories.map((category) => (
                <SelectItem
                  key={category.id}
                  value={String(category.id)}
                  className="hover:bg-slate-100! cursor-pointer"
                >
                  {category.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <FieldError errors={[fieldState.error]} />
        </Field>
      )}
    />
  );
}
