"use client";

import { useEffect } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { AnimatePresence, motion } from "framer-motion";
import { toast } from "sonner";

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

import {
  createMenuDefaultValues,
  CreateMenuFormValues,
  createMenuSchema,
} from "@/schema/menuSchema";

import {
  useGetSingleMenuQuery,
  useUpdateMenuMutation,
} from "@/app/state/api/menuApi";

import { getApiError } from "@/app/utils/apiError";
import {
  closeEditSheet,
  nextEditStep,
  previousEditStep,
} from "@/app/state/slices/menuSlice";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";

import MenuBasicInfo from "../create/MenuBasicInfo";
import { CreateStep } from "../create/CreateSteps";
import EditImageGrid from "./EditImageGrid";
import AddMoreImages from "./AddMoreImage";
import EditMenuSkeleton from "./EditMenuSkeleton";
import RecipeSection from "./RecipeSection";

export default function EditMenuSheet() {
  const dispatch = useAppDispatch();
  const { open, menuId, step } = useAppSelector(
    (state: RootState) => state.menu.edit,
  );

  const methods = useForm<CreateMenuFormValues>({
    resolver: zodResolver(createMenuSchema as never),
    defaultValues: createMenuDefaultValues,
  });

  const { reset, handleSubmit } = methods;

  const { data, isLoading } = useGetSingleMenuQuery(
    { id: String(menuId) },
    { skip: !menuId },
  );

  const [updateMenu, { isLoading: isUpdating }] = useUpdateMenuMutation();

  useEffect(() => {
    if (!data?.data) return;
    const menu = data.data;

    reset({
      category_id: menu.category.id,
      name: menu.name,
      description: menu.description,
      price: menu.price,
      prep_time_minutes: menu.prep_time_minutes,
      spice_level: menu.spice_level,
      allergen_ids: menu.allergens.map((a) => a.id),
      dietary_tag_ids: menu.dietary_tags.map((d) => d.id),
      display_order: menu.display_order,
    });
  }, [data, reset]);

  // step 1 submit: saves basic info, then advances to images —
  // does NOT close the sheet, mirrors create's flow exactly
  const onSubmit = async (values: CreateMenuFormValues) => {
    if (!menuId) return;
    try {
      const response = await updateMenu({ id: menuId, body: values }).unwrap();
      toast.success(response.message);
      dispatch(nextEditStep());
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  const handleClose = (nextOpen: boolean) => {
    if (!nextOpen) {
      reset(createMenuDefaultValues);
      dispatch(closeEditSheet());
    }
  };

  const images = data?.data.images ?? [];
  const showSkeleton = isLoading && !data;

  return (
    <Sheet open={open} onOpenChange={handleClose}>
      <SheetContent
        side="right"
        className="flex w-full flex-col bg-[#faedcd] sm:max-w-3xl! px-6 py-3 overflow-y-scroll"
      >
        <SheetHeader>
          <SheetTitle>Edit Menu</SheetTitle>
          <SheetDescription>
            Update menu information and images.
          </SheetDescription>
        </SheetHeader>

        {!showSkeleton && (
          <div className="flex items-center gap-3">
            <CreateStep active={step === 1} completed={step > 1}>
              1
            </CreateStep>
            <div className="h-px flex-1 bg-border" />
            <CreateStep active={step === 2}>2</CreateStep>
          </div>
        )}

        {showSkeleton ? (
          <EditMenuSkeleton />
        ) : (
          <FormProvider {...methods}>
            <div className="overflow-y-auto">
              <AnimatePresence mode="wait">
                {step === 1 ? (
                  <motion.div
                    key="basic"
                    initial={{ opacity: 0, x: 40 }}
                    animate={{ opacity: 1, x: 0 }}
                    exit={{ opacity: 0, x: -40 }}
                    transition={{ duration: 0.25 }}
                  >
                    <MenuBasicInfo
                      onSubmit={handleSubmit(onSubmit)}
                      onCancel={() => handleClose(false)}
                      isLoading={isUpdating}
                      submitText="Save & Continue"
                    />
                  </motion.div>
                ) : (
                  <motion.div
                    key="images"
                    initial={{ opacity: 0, x: 40 }}
                    animate={{ opacity: 1, x: 0 }}
                    exit={{ opacity: 0, x: -40 }}
                    transition={{ duration: 0.25 }}
                    className="space-y-6"
                  >
                    <div>
                      <h3 className="text-lg font-semibold">Images</h3>
                      <p className="text-sm text-muted-foreground">
                        {images.length}/4 images
                      </p>
                    </div>

                    {menuId && (
                      <>
                        <EditImageGrid menuId={menuId} images={images} />
                        {images.length < 4 && (
                          <AddMoreImages
                            menuId={menuId}
                            currentCount={images.length}
                          />
                        )}
                        <RecipeSection menuId={menuId} />
                      </>
                    )}

                    <div className="mt-6 flex justify-between gap-3 pt-6">
                      <button
                        type="button"
                        onClick={() => dispatch(previousEditStep())}
                        className="text-sm text-muted-foreground hover:text-foreground"
                      >
                        ← Back to details
                      </button>

                      <button
                        type="button"
                        onClick={() => handleClose(false)}
                        className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
                      >
                        Done
                      </button>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          </FormProvider>
        )}
      </SheetContent>
    </Sheet>
  );
}
