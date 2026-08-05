"use client";

import { useFormContext, Controller } from "react-hook-form";

import { Button } from "@/components/ui/button";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import CategorySelect from "./CategorySelector";
import SpiceLevelSelector from "./SpiceLevelSelector";
import AllergenSelector from "./AllergenSelector";
import DietaryTagSelector from "./DietaryTagSelector";

interface Props {
  isLoading?: boolean;
  onCancel: () => void;
  onSubmit: React.FormEventHandler<HTMLFormElement>;
}

export default function MenuBasicInfo({
  isLoading,
  onCancel,
  onSubmit,
}: Props) {
  const { control } = useFormContext();

  return (
    <form onSubmit={onSubmit} className="flex h-full flex-col">
      <div className="flex-1 space-y-8 overflow-y-auto pr-2">
        {/* Menu Information */}

        <section className="space-y-5">
          <h3 className="text-lg font-semibold">Menu Information</h3>

          <CategorySelect />

          <Controller
            control={control}
            name="name"
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel>Name:</FieldLabel>

                <Input
                  {...field}
                  placeholder="Grilled Chicken Burger"
                  className="placeholder:text-slate-400"
                />

                <FieldError
                  errors={[fieldState.error]}
                  className="text-red-400"
                />
              </Field>
            )}
          />

          <Controller
            control={control}
            name="description"
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel>Description</FieldLabel>

                <Textarea
                  {...field}
                  rows={4}
                  placeholder="Describe your menu item..."
                  className="placeholder:text-slate-400 focus:none"
                />

                <FieldError errors={[fieldState.error]} />
              </Field>
            )}
          />
        </section>

        {/* Pricing */}

        <section className="space-y-5">
          <h3 className="text-lg font-semibold">Pricing</h3>

          <div className="grid gap-4 md:grid-cols-3">
            <Controller
              name="price"
              control={control}
              render={({ field }) => (
                <Field>
                  <FieldLabel>Price</FieldLabel>

                  <Input type="number" {...field} />
                </Field>
              )}
            />

            <Controller
              name="prep_time_minutes"
              control={control}
              render={({ field }) => (
                <Field>
                  <FieldLabel>Prep Time</FieldLabel>

                  <Input type="number" {...field} />
                </Field>
              )}
            />

            <Controller
              name="display_order"
              control={control}
              render={({ field }) => (
                <Field>
                  <FieldLabel>Display Order</FieldLabel>

                  <Input type="number" {...field} />
                </Field>
              )}
            />
          </div>
        </section>

        {/* Spice */}

        <section className="space-y-4">
          <h3 className="text-lg font-semibold">Spice Level</h3>

          <SpiceLevelSelector />
        </section>

        {/* Dietary */}

        <section className="space-y-6">
          <div>
            <h3 className="text-lg font-semibold">Allergens</h3>

            <AllergenSelector />
          </div>

          <div>
            <h3 className="text-lg font-semibold">Dietary Tags</h3>

            <DietaryTagSelector />
          </div>
        </section>
      </div>

      <div className="mt-6 flex justify-end gap-3  pt-6">
        <Button
          variant="outline"
          type="button"
          onClick={onCancel}
          className="bg-red-500 text-white border-none hover:bg-red-600 hover:text-white"
        >
          Cancel
        </Button>

        <Button disabled={isLoading} type="submit">
          {isLoading ? "Creating..." : "Create Menu"}
        </Button>
      </div>
    </form>
  );
}
