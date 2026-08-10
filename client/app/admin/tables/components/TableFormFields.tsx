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
}

export default function TableFormFields({
  isLoading,
  onSubmit,
  submitText,
  onCancel,
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
              <FieldLabel>Table Name</FieldLabel>
              <Input
                {...field}
                placeholder="Table 4"
                className="placeholder:text-slate-400"
              />
              <FieldError errors={[fieldState.error]} />
            </Field>
          )}
        />

        <Controller
          control={control}
          name="capacity"
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel>Capacity</FieldLabel>
              <Input
                type="number"
                {...field}
                className="placeholder:text-slate-400"
              />
              <FieldError errors={[fieldState.error]} />
            </Field>
          )}
        />

        <Controller
          control={control}
          name="location"
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel>Location</FieldLabel>
              <Input
                {...field}
                placeholder="Rooftop, Indoor, Patio…"
                className="placeholder:text-slate-400"
              />
              <FieldError errors={[fieldState.error]} />
            </Field>
          )}
        />
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
          {isLoading ? "Saving…" : (submitText ?? "Create Table")}
        </Button>
      </div>
    </form>
  );
}
