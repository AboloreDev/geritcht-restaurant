"use client";

import { useFormContext, Controller } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";

interface Props {
  isLoading?: boolean;
  onSubmit: React.FormEventHandler<HTMLFormElement>;
  submitText?: string;
  onCancel: () => void;
  showStockField?: boolean;
}

export default function IngredientFormFields({
  isLoading,
  onSubmit,
  submitText,
  onCancel,
  showStockField = true,
}: Props) {
  const { control } = useFormContext();

  return (
    <form onSubmit={onSubmit} className="flex h-full flex-col">
      <div className="flex-1 space-y-5 overflow-y-auto pr-2">
        <Controller
          control={control}
          name="name"
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel>Name</FieldLabel>
              <Input
                {...field}
                placeholder="Flour"
                className="placeholder:text-slate-500"
              />
              <FieldError errors={[fieldState.error]} />
            </Field>
          )}
        />

        <div className="grid grid-cols-2 gap-4">
          <Controller
            control={control}
            name="unit"
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel>Unit</FieldLabel>
                <Input
                  {...field}
                  placeholder="kg, litres, pieces…"
                  className="placeholder:text-slate-500"
                />
                <FieldError errors={[fieldState.error]} />
              </Field>
            )}
          />

          <Controller
            control={control}
            name="min_threshold"
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel>Low-stock threshold</FieldLabel>
                <Input
                  type="number"
                  step="0.01"
                  {...field}
                  className="placeholder:text-slate-500"
                />
                <FieldError errors={[fieldState.error]} />
              </Field>
            )}
          />
        </div>

        {showStockField && (
          <Controller
            control={control}
            name="current_stock"
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel>Starting stock</FieldLabel>
                <Input
                  type="number"
                  step="0.01"
                  {...field}
                  className="placeholder:text-slate-500"
                />
                <FieldError errors={[fieldState.error]} />
              </Field>
            )}
          />
        )}
      </div>

      <div className="mt-6 flex justify-end gap-3 pt-6">
        <Button
          variant="outline"
          type="button"
          onClick={onCancel}
          className="border-none bg-red-500 text-white hover:bg-red-600 hover:text-white"
        >
          Cancel
        </Button>
        <Button disabled={isLoading} type="submit">
          {isLoading ? "Saving…" : (submitText ?? "Create Ingredient")}
        </Button>
      </div>
    </form>
  );
}
