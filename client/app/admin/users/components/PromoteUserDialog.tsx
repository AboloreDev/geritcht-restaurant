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
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { useUpdateUserRoleMutation } from "@/app/state/api/userApi";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import { closePromoteDialog } from "@/app/state/slices/userSlice";

export default function PromoteToStaffDialog() {
  const dispatch = useAppDispatch();
  const userId = useAppSelector((s: RootState) => s.user.dialogs.promoteUserId);
  const [updateRole, { isLoading }] = useUpdateUserRoleMutation();

  async function handleConfirm() {
    if (!userId) return;
    try {
      const response = await updateRole({ id: userId, role: "staff" }).unwrap();
      toast.success(response.message);
      dispatch(closePromoteDialog());
    } catch (err) {
      toast.error(getApiError(err));
    }
  }

  return (
    <AlertDialog
      open={userId !== null}
      onOpenChange={(open) => !open && dispatch(closePromoteDialog())}
    >
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Promote to staff?</AlertDialogTitle>
          <AlertDialogDescription>
            This user will gain staff-level access to the admin dashboard.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction disabled={isLoading} onClick={handleConfirm}>
            {isLoading ? "Promoting…" : "Promote"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
