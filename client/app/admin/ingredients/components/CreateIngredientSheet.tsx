"use client";

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
import { useCreateIngredientMutation } from "@/app/state/api/ingredientApi";
import { getApiError } from "@/app/utils/apiError";
import { closeCreateIngredientSheet } from "@/app/state/slices/ingredientSlice";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";

import IngredientFormFields from "./IngredientFormFields";

export default function CreateIngredientSheet() {
  const dispatch = useAppDispatch();
  const open = useAppSelector((s: RootState) => s.ingredient.create.open);
  const [createIngredient, { isLoading }] = useCreateIngredientMutation();

  const methods = useForm<IngredientFormValues>({
    resolver: zodResolver(ingredientSchema as never),
    defaultValues: ingredientDefaultValues,
  });

  const { handleSubmit, reset } = methods;

  const onSubmit = async (values: IngredientFormValues) => {
    try {
      const response = await createIngredient(values).unwrap();
      toast.success(response.message);
      reset(ingredientDefaultValues);
      dispatch(closeCreateIngredientSheet());
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  const handleClose = (nextOpen: boolean) => {
    if (!nextOpen) {
      reset(ingredientDefaultValues);
      dispatch(closeCreateIngredientSheet());
    }
  };

  return (
    <Sheet open={open} onOpenChange={handleClose}>
      <SheetContent
        side="right"
        className="flex w-full flex-col bg-[#faedcd] !max-w-xl px-4 py-2"
      >
        <SheetHeader>
          <SheetTitle>Add Ingredient</SheetTitle>
          <SheetDescription>
            Track a new ingredient&apos;s stock.
          </SheetDescription>
        </SheetHeader>

        <FormProvider {...methods}>
          <IngredientFormFields
            onSubmit={handleSubmit(onSubmit)}
            isLoading={isLoading}
            submitText="Create Ingredient"
            onCancel={() => handleClose(false)}
            showStockField
          />
        </FormProvider>
      </SheetContent>
    </Sheet>
  );
}
