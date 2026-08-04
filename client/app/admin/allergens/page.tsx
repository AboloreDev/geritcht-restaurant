"use client";

import { toast } from "sonner";
import { getApiError } from "@/app/utils/apiError";
import { LabelManagementDashboard } from "@/app/admin/components/LabelManagementDashboard";
import {
  useCreateAllergenMutation,
  useDeleteAllergenMutation,
  useGetAllergensQuery,
  useUpdateAllergenMutation,
} from "@/app/state/api/allergenApi";

export default function Allergens() {
  const { data, isLoading } = useGetAllergensQuery();
  const [createAllergen] = useCreateAllergenMutation();
  const [updateAllergen] = useUpdateAllergenMutation();
  const [deleteAllergen] = useDeleteAllergenMutation();

  console.log(data);
  return (
    <LabelManagementDashboard
      title="Allergens"
      singular="allergen"
      description="Maintain clear allergen labels so guests and staff can make safe menu choices."
      items={data?.data ?? []}
      isLoading={isLoading}
      accentClass="bg-red-100 text-red-700"
      onCreate={async (name) => {
        try {
          const response = await createAllergen({ name }).unwrap();
          toast.success(response.message || "Allergen created.");
        } catch (error) {
          toast.error(getApiError(error));
          throw error;
        }
      }}
      onUpdate={async (id, name) => {
        try {
          const response = await updateAllergen({
            id,
            body: { name },
          }).unwrap();
          toast.success(response.message || "Allergen updated.");
        } catch (error) {
          toast.error(getApiError(error));
          throw error;
        }
      }}
      onDelete={async (id) => {
        try {
          await deleteAllergen({ id }).unwrap();
          toast.success("Allergen deleted.");
        } catch (error) {
          toast.error(getApiError(error));
          throw error;
        }
      }}
    />
  );
}
