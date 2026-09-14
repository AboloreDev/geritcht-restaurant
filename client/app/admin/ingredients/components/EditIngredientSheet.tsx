// src/app/admin/ingredients/components/EditIngredientSheet.tsx
"use client";

import { useEffect } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

import {
  ingredientSchema,
  ingredientDefaultValues,
  IngredientFormValues,
} from "@/schema/ingredientSchema";
import {
  useGetIngredientQuery,
  useUpdateIngredientMutation,
} from "@/app/state/api/ingredientApi";
import { getApiError } from "@/app/utils/apiError";
import { closeEditIngredientSheet } from "@/app/state/slices/ingredientSlice";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";

import IngredientFormFields from "./IngredientFormFields";
import IngredientFormSkeleton from "./IngredientFormSkeleton";

export default function EditIngredientSheet() {
  const dispatch = useAppDispatch();
  const { open, ingredientId } = useAppSelector(
    (s: RootState) => s.ingredient.edit,
  );

  const { data, isLoading } = useGetIngredientQuery(
    { id: ingredientId! },
    { skip: !ingredientId },
  );
  const [updateIngredient, { isLoading: isUpdating }] =
    useUpdateIngredientMutation();

  const methods = useForm<IngredientFormValues>({
    resolver: zodResolver(ingredientSchema as never),
    defaultValues: ingredientDefaultValues,
  });

  const { handleSubmit, reset } = methods;

  useEffect(() => {
    if (!data?.data) return;
    reset({
      name: data.data.name,
      unit: data.data.unit,
      current_stock: data.data.current_stock,
      min_threshold: data.data.min_threshold,
    });
  }, [data, reset]);

  useEffect(() => {
    if (!ingredientId) reset(ingredientDefaultValues);
  }, [ingredientId, reset]);

  const onSubmit = async (values: IngredientFormValues) => {
    if (!ingredientId) return;
    try {
      const response = await updateIngredient({
        id: ingredientId,
        body: {
          name: values.name,
          unit: values.unit,
          min_threshold: values.min_threshold,
        },
      }).unwrap();
      toast.success(response.message);
      dispatch(closeEditIngredientSheet());
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  const handleClose = (nextOpen: boolean) => {
    if (!nextOpen) dispatch(closeEditIngredientSheet());
  };

  const showSkeleton = isLoading && !data;

  return (
    <Sheet open={open} onOpenChange={handleClose}>
      <SheetContent
        side="right"
        className="flex w-full flex-col bg-[#faedcd] !max-w-xl px-4 py-2"
      >
        <SheetHeader>
          <SheetTitle>Edit Ingredient</SheetTitle>
          <SheetDescription>
            Update this ingredient&apos;s details.
          </SheetDescription>
        </SheetHeader>

        {showSkeleton ? (
          <IngredientFormSkeleton />
        ) : (
          <FormProvider {...methods}>
            <IngredientFormFields
              onSubmit={handleSubmit(onSubmit)}
              isLoading={isUpdating}
              submitText="Save Changes"
              onCancel={() => handleClose(false)}
              showStockField={false}
            />
          </FormProvider>
        )}
      </SheetContent>
    </Sheet>
  );
}
