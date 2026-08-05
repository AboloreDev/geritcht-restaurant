"use client";

import { useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { FormProvider, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
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
import { CreateStep } from "./CreateSteps";
import { useCreateMenuMutation } from "@/app/state/api/menuApi";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import MenuBasicInfo from "./MenuBasicInfo";

type CreateMenuSheetProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export default function CreateMenuSheet({
  open,
  onOpenChange,
}: CreateMenuSheetProps) {
  const [createMenu, { isLoading }] = useCreateMenuMutation();
  /**
   * Step 1 -> Basic Menu Information
   * Step 2 -> Upload Images
   */
  const [step, setStep] = useState<1 | 2>(1);

  /**
   * We keep the created menu ID here after the
   * createMenu mutation succeeds.
   *
   * Step 2 needs it to upload images.
   */
  const [menuId, setMenuId] = useState<number>();

  const methods = useForm<CreateMenuFormValues>({
    // @ts-expect-error "<>"
    resolver: zodResolver(createMenuSchema),
    defaultValues: createMenuDefaultValues,
    mode: "onChange",
  });

  const { handleSubmit } = methods;

  const onSubmit = async (values: CreateMenuFormValues) => {
    try {
      const response = await createMenu(values).unwrap();
      toast.success(response.message);
      setMenuId(response.data.id);
      setStep(2);
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  const handleClose = (value: boolean) => {
    if (!value) {
      methods.reset(createMenuDefaultValues);
      setStep(1);
      setMenuId(undefined);
    }

    onOpenChange(value);
  };

  return (
    <Sheet open={open} onOpenChange={handleClose}>
      <SheetContent
        side="right"
        className="w-full flex space-y-6 flex-col bg-[#faedcd] sm:max-w-3xl! overflow-y-auto! px-6"
      >
        <SheetHeader>
          <SheetTitle>Create Menu Item</SheetTitle>

          <SheetDescription>
            Add a new menu item and upload beautiful food images for guests.
          </SheetDescription>
        </SheetHeader>

        {/* Step Indicator */}
        <div className=" flex items-center gap-3">
          <CreateStep active={step === 1} completed={step > 1}>
            1
          </CreateStep>

          <div className="h-px flex-1 bg-border" />

          <CreateStep active={step === 2}>2</CreateStep>
        </div>

        <FormProvider {...methods}>
          <div className="overflow-y-auto">
            <AnimatePresence mode="wait">
              {step === 1 ? (
                <motion.div
                  key="step-one"
                  initial={{ opacity: 0, x: 40 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -40 }}
                  transition={{ duration: 0.25 }}
                >
                  <MenuBasicInfo
                    onSubmit={handleSubmit(onSubmit)}
                    isLoading={isLoading}
                    onCancel={() => handleClose(false)}
                  />
                </motion.div>
              ) : (
                <motion.div
                  key="step-two"
                  initial={{ opacity: 0, x: 40 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -40 }}
                  transition={{ duration: 0.25 }}
                >
                  {/* MenuImageUpload goes here */}
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        </FormProvider>
      </SheetContent>
    </Sheet>
  );
}
