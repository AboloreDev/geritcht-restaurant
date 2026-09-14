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
import { useDeleteIngredientMutation } from "@/app/state/api/ingredientApi";
import { closeDeleteIngredientDialog } from "@/app/state/slices/ingredientSlice";
import { RootState, useAppDispatch, useAppSelector } from "@/app/state/redux";
import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";

export default function DeleteIngredientDialog() {
  const dispatch = useAppDispatch();
  const deleteIngredientId = useAppSelector(
    (s: RootState) => s.ingredient.dialogs.deleteIngredientId,
  );
  const [deleteIngredient, { isLoading }] = useDeleteIngredientMutation();

  async function handleConfirm() {
    if (!deleteIngredientId) return;
    try {
      const response = await deleteIngredient({
        id: deleteIngredientId,
      }).unwrap();
      toast.success(response.message);
      dispatch(closeDeleteIngredientDialog());
    } catch (error) {
      toast.error(getApiError(error));
    }
  }

  return (
    <AlertDialog
      open={deleteIngredientId !== null}
      onOpenChange={(open) => !open && dispatch(closeDeleteIngredientDialog())}
    >
      <AlertDialogContent className="!max-w-xl px-4 py-4 bg-white">
        <AlertDialogHeader>
          <AlertDialogTitle>Delete this ingredient?</AlertDialogTitle>
          <AlertDialogDescription>
            This will remove it from any linked recipes and stop stock tracking.
            This cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel className="bg-red-500 text-white border-none hover:bg-red-600">
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            disabled={isLoading}
            onClick={handleConfirm}
            className="text-black hover:bg-destructive/90"
          >
            {isLoading ? "Deleting…" : "Delete"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
