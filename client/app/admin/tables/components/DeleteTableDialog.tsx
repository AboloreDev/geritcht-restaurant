"use client";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useDeleteTableMutation } from "@/app/state/api/tableApi";
import { closeDeleteTableDialog } from "@/app/state/slices/tableSlice";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

export default function DeleteTableDialog() {
  const dispatch = useAppDispatch();
  const deleteTableId = useAppSelector(
    (s: RootState) => s.table.dialogs.deleteTableId,
  );
  const [deleteTable, { isLoading }] = useDeleteTableMutation();

  async function handleConfirm() {
    if (!deleteTableId) return;
    try {
      const response = await deleteTable({ id: deleteTableId }).unwrap();
      toast.success(response.message);
      dispatch(closeDeleteTableDialog());
    } catch (error) {
      toast.error(getApiError(error));
    }
  }

  return (
    <AlertDialog
      open={deleteTableId !== null}
      onOpenChange={(open) => !open && dispatch(closeDeleteTableDialog())}
    >
      <AlertDialogContent className="bg-[#fefae0]">
        <AlertDialogHeader>
          <AlertDialogTitle>Delete this table?</AlertDialogTitle>
          <AlertDialogDescription>
            This cannot be undone. Any reservations linked to this table may be
            affected.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel className="bg-white text-black border-none">
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            disabled={isLoading}
            onClick={handleConfirm}
            className="bg-red-500 text-white hover:bg-red-600"
          >
            {isLoading ? "Deleting…" : "Delete"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
