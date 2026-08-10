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
  tableSchema,
  tableDefaultValues,
  TableFormValues,
} from "@/schema/tableSchema";
import {
  useGetTableQuery,
  useUpdateTableMutation,
} from "@/app/state/api/tableApi";
import { getApiError } from "@/app/utils/apiError";
import { closeEditTableSheet } from "@/app/state/slices/tableSlice";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";

import TableFormFields from "./TableFormFields";
import TableFormSkeleton from "./TableFormSkeleton";
import { useRouter } from "next/navigation";
export default function EditTableSheet() {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const { open, tableId } = useAppSelector((s: RootState) => s.table.edit);

  const { data, isLoading } = useGetTableQuery(
    { id: tableId! },
    { skip: !tableId },
  );
  const [updateTable, { isLoading: isUpdating }] = useUpdateTableMutation();

  const methods = useForm<TableFormValues>({
    resolver: zodResolver(tableSchema as never),
    defaultValues: tableDefaultValues,
  });

  const { handleSubmit, reset } = methods;

  // populate once fresh data arrives
  useEffect(() => {
    if (!data?.data) return;
    reset({
      name: data.data.name,
      capacity: data.data.capacity,
      location: data.data.location ?? "",
    });
  }, [data, reset]);

  // clear back to blanks when the sheet closes, so a different table's
  // edit doesn't briefly flash the previous table's values
  useEffect(() => {
    if (!tableId) reset(tableDefaultValues);
  }, [tableId, reset]);

  const onSubmit = async (values: TableFormValues) => {
    if (!tableId) return;
    try {
      const response = await updateTable({
        id: tableId,
        body: values,
      }).unwrap();
      toast.success(response.message);
      dispatch(closeEditTableSheet());
      router.push("/admin/tables");
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  const handleClose = (nextOpen: boolean) => {
    if (!nextOpen) dispatch(closeEditTableSheet());
  };

  const showSkeleton = isLoading && !data;

  return (
    <Sheet open={open} onOpenChange={handleClose}>
      <SheetContent
        side="right"
        className="flex w-full flex-col bg-[#faedcd] sm:max-w-2xl! p-4"
      >
        <SheetHeader>
          <SheetTitle className="text-2xl">Edit Table</SheetTitle>
          <SheetDescription>Update this table&apos;s details.</SheetDescription>
        </SheetHeader>

        {showSkeleton ? (
          <TableFormSkeleton />
        ) : (
          <FormProvider {...methods}>
            <TableFormFields
              onSubmit={handleSubmit(onSubmit)}
              isLoading={isUpdating}
              submitText="Save Changes"
              onCancel={() => handleClose(false)}
            />
          </FormProvider>
        )}
      </SheetContent>
    </Sheet>
  );
}
