"use client";

import { toast } from "sonner";
import { MenuCategory } from "@/app/state/types/menuTypes";
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
import { getApiError } from "@/app/utils/apiError";
import { useDeleteCategoryMutation } from "@/app/state/api/categoriesApi";

type DeleteCategoryDialogProps = {
  category: MenuCategory;
  onOpenChange: (open: boolean) => void;
};

export default function DeleteCategoryDialog({
  category,
  onOpenChange,
}: DeleteCategoryDialogProps) {
  const [deleteCategory, { isLoading }] = useDeleteCategoryMutation();

  const confirmDelete = async () => {
    try {
      const response = await deleteCategory({ id: category.id }).unwrap();
      toast.success(response.message);
      onOpenChange(false);
    } catch (err) {
      console.error("Error deleting category:", err);
      toast.error(getApiError(err));
    }
  };

  return (
    <AlertDialog open onOpenChange={onOpenChange}>
      <AlertDialogContent className="bg-[#fefae0]">
        <AlertDialogHeader>
          <AlertDialogTitle>Delete {category.name}?</AlertDialogTitle>
          <AlertDialogDescription>
            This will permanently remove the category. This action cannot be
            undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel
            disabled={isLoading}
            className="bg-white hover:bg-gray-100 border-none"
          >
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            variant="destructive"
            onClick={confirmDelete}
            disabled={isLoading}
            className="bg-red-500 hover:bg-red-600 text-white border-none"
          >
            {isLoading ? "Deleting…" : "Delete category"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
