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
  tableSchema,
  tableDefaultValues,
  TableFormValues,
} from "@/schema/tableSchema";
import { useCreateTableMutation } from "@/app/state/api/tableApi";
import { getApiError } from "@/app/utils/apiError";
import { closeCreateTableSheet } from "@/app/state/slices/tableSlice";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";

import TableFormFields from "./TableFormFields";

export default function CreateTableSheet() {
  const dispatch = useAppDispatch();
  const open = useAppSelector((s: RootState) => s.table.create.open);
  const [createTable, { isLoading }] = useCreateTableMutation();

  const methods = useForm<TableFormValues>({
    resolver: zodResolver(tableSchema as never),
    defaultValues: tableDefaultValues,
  });

  const { handleSubmit, reset } = methods;

  const onSubmit = async (values: TableFormValues) => {
    try {
      const response = await createTable(values).unwrap();
      toast.success(response.message);
      reset(tableDefaultValues);
      dispatch(closeCreateTableSheet());
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  const handleClose = (nextOpen: boolean) => {
    if (!nextOpen) {
      reset(tableDefaultValues);
      dispatch(closeCreateTableSheet());
    }
  };

  return (
    <Sheet open={open} onOpenChange={handleClose}>
      <SheetContent
        side="right"
        className="flex w-full flex-col bg-[#faedcd] sm:max-w-2xl! p-4"
      >
        <SheetHeader>
          <SheetTitle className="text-3xl">Add Table</SheetTitle>
          <SheetDescription>
            Create a new table for your floor plan.
          </SheetDescription>
        </SheetHeader>

        <FormProvider {...methods}>
          <TableFormFields
            onSubmit={handleSubmit(onSubmit)}
            isLoading={isLoading}
            submitText="Create Table"
            onCancel={() => handleClose(false)}
          />
        </FormProvider>
      </SheetContent>
    </Sheet>
  );
}
