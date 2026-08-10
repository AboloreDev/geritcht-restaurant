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
import { useDeleteMenuMutation } from "@/app/state/api/menuApi";
import { closeDeleteMenuDialog } from "@/app/state/slices/menuSlice";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import { useRouter } from "next/navigation";

export default function DeleteMenuDialog() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const deleteMenuId = useAppSelector(
    (s: RootState) => s.menu.dialogs.deleteMenuId,
  );
  const [deleteMenu, { isLoading }] = useDeleteMenuMutation();

  async function handleConfirm() {
    if (!deleteMenuId) return;
    try {
      const response = await deleteMenu({ id: deleteMenuId }).unwrap();
      toast.success(response.message);
      dispatch(closeDeleteMenuDialog());
      router.push("/admin/menu");
    } catch (error) {
      toast.error(getApiError(error));
    }
  }

  return (
    <AlertDialog
      open={deleteMenuId !== null}
      onOpenChange={(open) => !open && dispatch(closeDeleteMenuDialog())}
    >
      <AlertDialogContent className="bg-[#fefae0]">
        <AlertDialogHeader>
          <AlertDialogTitle>Delete this menu item?</AlertDialogTitle>
          <AlertDialogDescription>
            This will permanently remove the item and all its images. This
            action cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <AlertDialogFooter>
          <AlertDialogCancel className="bg-white text-black border-none">
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            disabled={isLoading}
            onClick={handleConfirm}
            className="bg-red-500 text-white border-none hover:bg-destructive/90"
          >
            {isLoading ? "Deleting…" : "Delete"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
