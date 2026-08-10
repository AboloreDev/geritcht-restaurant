"use client";

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
import MenuImageUpload from "./image/ImageUpload";

import {
  closeCreateSheet,
  nextCreateStep,
  resetCreateState,
  setCreatedMenu,
} from "@/app/state/slices/menuSlice";

import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";

export default function CreateMenuSheet() {
  const dispatch = useAppDispatch();

  const { open, step } = useAppSelector(
    (state: RootState) => state.menu.create,
  );

  const [createMenu, { isLoading }] = useCreateMenuMutation();

  const methods = useForm<CreateMenuFormValues>({
    resolver: zodResolver(createMenuSchema as never),
    defaultValues: createMenuDefaultValues,
    mode: "onChange",
  });

  const { handleSubmit, reset } = methods;

  const onSubmit = async (values: CreateMenuFormValues) => {
    try {
      const response = await createMenu(values).unwrap();
      toast.success(response.message);
      dispatch(setCreatedMenu(response.data.id));
      dispatch(nextCreateStep());
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  const handleClose = (value: boolean) => {
    if (!value) {
      reset(createMenuDefaultValues);
      dispatch(resetCreateState());
      dispatch(closeCreateSheet());
    }
  };

  // called once all images finish uploading (or the admin chooses to
  // skip images) — closes out the whole create flow the same way
  // cancelling does, since the menu item itself already exists by now
  const handleFinish = () => {
    reset(createMenuDefaultValues);
    dispatch(resetCreateState());
    dispatch(closeCreateSheet());
  };

  return (
    <Sheet open={open} onOpenChange={handleClose}>
      <SheetContent
        side="right"
        className="flex w-full flex-col space-y-6 overflow-y-auto bg-[#faedcd] px-6 sm:max-w-3xl!"
      >
        <SheetHeader>
          <SheetTitle>Create Menu Item</SheetTitle>
          <SheetDescription>
            Add a new menu item and upload beautiful food images.
          </SheetDescription>
        </SheetHeader>

        <div className="flex items-center gap-3">
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
                  key="basic"
                  initial={{ opacity: 0, x: 40 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -40 }}
                  transition={{ duration: 0.25 }}
                >
                  <MenuBasicInfo
                    onSubmit={handleSubmit(onSubmit)}
                    isLoading={isLoading}
                    submitText="Create Menu"
                    onCancel={() => handleClose(false)}
                  />
                </motion.div>
              ) : (
                <motion.div
                  key="images"
                  initial={{ opacity: 0, x: 40 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -40 }}
                  transition={{ duration: 0.25 }}
                >
                  <MenuImageUpload onFinish={handleFinish} />
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        </FormProvider>
      </SheetContent>
    </Sheet>
  );
}
