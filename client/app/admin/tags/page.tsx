"use client";

import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import { LabelManagementDashboard } from "@/app/admin/components/LabelManagementDashboard";
import {
  useCreateDietaryTagsMutation,
  useDeleteDietaryTagsMutation,
  useGetDietaryTagsQuery,
  useUpdateDietaryTagsMutation,
} from "@/app/state/api/dietaryTagsApi";

export default function Tags() {
  const { data, isLoading } = useGetDietaryTagsQuery();
  const [createDietaryTag] = useCreateDietaryTagsMutation();
  const [updateDietaryTag] = useUpdateDietaryTagsMutation();
  const [deleteDietaryTag] = useDeleteDietaryTagsMutation();
  return (
    <LabelManagementDashboard
      title="Dietary Tags"
      singular="dietary tag"
      description="Manage the dietary labels that help guests find dishes matching their preferences."
      items={data?.data ?? []}
      isLoading={isLoading}
      accentClass="bg-emerald-100 text-emerald-700"
      onCreate={async (name) => {
        try {
          const response = await createDietaryTag({ name }).unwrap();
          toast.success(response.message || "Dietary tag created.");
        } catch (error) {
          toast.error(getApiError(error));
          throw error;
        }
      }}
      onUpdate={async (id, name) => {
        try {
          const response = await updateDietaryTag({
            id,
            body: { name },
          }).unwrap();
          toast.success(response.message || "Dietary tag updated.");
        } catch (error) {
          toast.error(getApiError(error));
          throw error;
        }
      }}
      onDelete={async (id) => {
        try {
          await deleteDietaryTag({ id }).unwrap();
          toast.success("Dietary tag deleted.");
        } catch (error) {
          toast.error(getApiError(error));
          throw error;
        }
      }}
    />
  );
}
