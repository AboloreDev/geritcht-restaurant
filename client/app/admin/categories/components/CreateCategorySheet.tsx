"use client";

import { motion } from "framer-motion";
import { z } from "zod";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { useCreateCategoryMutation } from "@/app/state/api/categoriesApi";
import { Button } from "@/components/ui/button";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import {
  CreateCategoryFormValues,
  createCategorySchema,
} from "@/schema/categorySchema";
import { getApiError } from "@/app/utils/apiError";

type CreateCategorySheetProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export default function CreateCategorySheet({
  open,
  onOpenChange,
}: CreateCategorySheetProps) {
  const [createCategory, { isLoading }] = useCreateCategoryMutation();
  const { control, handleSubmit, reset } = useForm<CreateCategoryFormValues>({
    resolver: zodResolver(
      // @ts-expect-error "<>"
      createCategorySchema,
    ),
    defaultValues: {
      name: "",
      description: "",
      image_url: "",
      display_order: 0,
      is_active: true,
    },
  });

  const handleOpenChange = (nextOpen: boolean) => {
    if (!nextOpen) reset();
    onOpenChange(nextOpen);
  };

  const onSubmit = async (values: CreateCategoryFormValues) => {
    try {
      const response = await createCategory({ data: values }).unwrap();
      toast.success(response.message);
      handleOpenChange(false);
    } catch (err) {
      toast.error(getApiError(err));
    }
  };

  return (
    <Sheet open={open} onOpenChange={handleOpenChange}>
      <SheetContent
        side="right"
        className="flex h-full w-full flex-col bg-[#fefae0] overflow-y-auto sm:max-w-6xl"
      >
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.38, delay: 0.12, ease: [0.22, 1, 0.36, 1] }}
          className="flex min-h-full flex-col"
        >
          <SheetHeader>
            <SheetTitle className="text-2xl">Create category</SheetTitle>
            <SheetDescription className="text-slate-600 text-xs">
              Create a new section for guests to browse on your menu.
            </SheetDescription>
          </SheetHeader>
          <form
            onSubmit={handleSubmit(onSubmit)}
            className="flex flex-1 flex-col"
          >
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.3, delay: 0.22 }}
              className="flex-1 space-y-5 overflow-y-auto px-6 pb-6"
            >
              <Controller
                name="name"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="new-category-name">
                      Category name
                    </FieldLabel>
                    <Input
                      {...field}
                      id="new-category-name"
                      autoFocus
                      placeholder="e.g. Main Courses"
                      className="placeholder:text-slate-400"
                      aria-invalid={fieldState.invalid}
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />
              <Controller
                name="description"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="new-category-description">
                      Description
                    </FieldLabel>
                    <Textarea
                      {...field}
                      id="new-category-description"
                      placeholder="A short description guests will see."
                      aria-invalid={fieldState.invalid}
                      className="placeholder:text-slate-400"
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />
              <Controller
                name="image_url"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="new-category-image">
                      Image URL (Optional)
                    </FieldLabel>
                    <Input
                      {...field}
                      id="new-category-image"
                      type="url"
                      placeholder="https://example.com/category.jpg"
                      aria-invalid={fieldState.invalid}
                      className="placeholder:text-slate-400"
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />
              <Controller
                name="display_order"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="new-category-order">
                      Display order
                    </FieldLabel>
                    <Input
                      {...field}
                      id="new-category-order"
                      type="number"
                      min="0"
                      onChange={(event) =>
                        field.onChange(event.currentTarget.valueAsNumber)
                      }
                      aria-invalid={fieldState.invalid}
                      className="placeholder:text-slate-400"
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />
            </motion.div>
            <SheetFooter className="">
              <Button
                type="button"
                variant="outline"
                onClick={() => handleOpenChange(false)}
                disabled={isLoading}
                className="bg-white text-slate-700 hover:bg-slate-100 border-none"
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? "Creating…" : "Create category"}
              </Button>
            </SheetFooter>
          </form>
        </motion.div>
      </SheetContent>
    </Sheet>
  );
}
