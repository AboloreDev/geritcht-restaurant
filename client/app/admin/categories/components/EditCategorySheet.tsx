"use client";

import { toast } from "sonner";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEditCategoryMutation } from "@/app/state/api/categoriesApi";
import { Category } from "@/app/state/types/categoriesTypes";
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
  EditCategoryFormValues,
  editCategorySchema,
} from "@/schema/categorySchema";
import { getApiError } from "@/app/utils/apiError";
import { motion } from "framer-motion";

type EditCategorySheetProps = {
  category: Category;
  onOpenChange: (open: boolean) => void;
};

const containerVariants = {
  hidden: {},
  visible: {
    transition: {
      staggerChildren: 0.08,
      delayChildren: 0.1,
    },
  },
};

const itemVariants = {
  hidden: {
    opacity: 0,
    y: 12,
  },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.35,
      ease: "easeOut",
    },
  },
};

export default function EditCategorySheet({
  category,
  onOpenChange,
}: EditCategorySheetProps) {
  const [editCategory, { isLoading }] = useEditCategoryMutation();
  const { control, handleSubmit } = useForm<EditCategoryFormValues>({
    resolver: zodResolver(
      // @ts-expect-error "<>"
      editCategorySchema,
    ),
    defaultValues: {
      name: category.name,
      description: category.description,
      image_url: category.image_url ?? "",
      display_order: category.display_order ?? 0,
      is_active: category.is_active,
    },
  });

  const onSubmit = async (values: EditCategoryFormValues) => {
    try {
      const response = await editCategory({
        id: category.id,
        data: values,
      }).unwrap();
      toast.success(response.message);
      onOpenChange(false);
    } catch (error) {
      console.log(error);
      toast.error(getApiError(error));
    }
  };

  return (
    <Sheet open onOpenChange={onOpenChange}>
      <SheetContent
        side="right"
        className="flex h-full w-full flex-col bg-[#fefae0] overflow-y-auto sm:max-w-2xl"
      >
        <SheetHeader>
          <SheetTitle className="text-2xl">Edit category</SheetTitle>
          <SheetDescription className="text-xs font-thin text-slate-600 ">
            Update the details guests see while browsing your menu.
          </SheetDescription>
        </SheetHeader>

        <form
          onSubmit={handleSubmit(onSubmit)}
          className="flex flex-1 flex-col"
        >
          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
            className="flex-1 space-y-5 overflow-y-auto px-6 pb-6"
          >
            <motion.div
              // @ts-expect-error "<>"
              variants={itemVariants}
            >
              <Controller
                name="name"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="category-name">
                      Category name:
                    </FieldLabel>
                    <Input
                      {...field}
                      id="category-name"
                      autoFocus
                      placeholder="e.g. Main Courses"
                      aria-invalid={fieldState.invalid}
                      className="placeholder:text-slate-400"
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />
            </motion.div>

            <motion.div
              // @ts-expect-error "<>"
              variants={itemVariants}
            >
              <Controller
                name="description"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="category-description">
                      Description:
                    </FieldLabel>
                    <Textarea
                      {...field}
                      id="category-description"
                      placeholder="A short description guests will see."
                      aria-invalid={fieldState.invalid}
                      className="placeholder:text-slate-400"
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />
            </motion.div>

            <motion.div
              // @ts-expect-error "<>"
              variants={itemVariants}
            >
              <Controller
                name="image_url"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="category-image">
                      Image URL: (Optional)
                    </FieldLabel>
                    <Input
                      {...field}
                      id="category-image"
                      type="url"
                      placeholder="https://example.com/category.jpg"
                      aria-invalid={fieldState.invalid}
                      className="placeholder:text-slate-400"
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />
            </motion.div>

            <motion.div
              // @ts-expect-error "<>"
              variants={itemVariants}
            >
              <Controller
                name="display_order"
                control={control}
                render={({ field, fieldState }) => (
                  <Field data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="category-order">
                      Display order:
                    </FieldLabel>
                    <Input
                      {...field}
                      id="category-order"
                      type="number"
                      min="0"
                      onChange={(e) =>
                        field.onChange(e.currentTarget.valueAsNumber)
                      }
                      aria-invalid={fieldState.invalid}
                      className="placeholder:text-slate-400"
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />
            </motion.div>

            <motion.div
              // @ts-expect-error "<>"
              variants={itemVariants}
            >
              <Controller
                name="is_active"
                control={control}
                render={({ field }) => (
                  <div className="flex items-center justify-between rounded-2xl  bg-white/70 p-4">
                    <div>
                      <p className="font-medium text-slate-900">Show on menu</p>
                      <p className="mt-1 text-sm text-slate-400">
                        {field.value
                          ? "Guests can currently browse this category."
                          : "This category is hidden from guests."}
                      </p>
                    </div>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                      aria-label="Show category on menu"
                      className="data-[state=checked]:bg-emerald-600 cursor-pointer data-[state=unchecked]:bg-black [&>span]:bg-slate-500"
                    />
                  </div>
                )}
              />
            </motion.div>
          </motion.div>

          <motion.div
            // @ts-expect-error "<>"
            variants={itemVariants}
            initial="hidden"
            animate="visible"
          >
            <SheetFooter className=" border-border mt-auto  px-6 py-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isLoading}
                className="bg-white text-black border-none hover:bg-white/50"
              >
                Cancel
              </Button>

              <Button type="submit" disabled={isLoading}>
                {isLoading ? "Saving…" : "Save changes"}
              </Button>
            </SheetFooter>
          </motion.div>
        </form>
      </SheetContent>
    </Sheet>
  );
}
