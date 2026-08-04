import {
  Allergen,
  AllergenResponse,
  CreateAllergenRequest,
  GetAllergensResponse,
  UpdateAllergenRequest,
} from "../types/allergenTypes";
import { baseApi } from "./baseApi";

export const allergenApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAllergens: builder.query<GetAllergensResponse, void>({
      query: () => "/allergens/",
      providesTags: ["Allergen"],
    }),

    createAllergen: builder.mutation<AllergenResponse, CreateAllergenRequest>({
      query: (body) => ({
        url: "/allergens/",
        method: "POST",
        body,
      }),
      invalidatesTags: ["Allergen"],
    }),
    updateAllergen: builder.mutation<
      AllergenResponse,
      { id: number; body: UpdateAllergenRequest }
    >({
      query: ({ id, body }) => ({
        url: `/allergens/${id}`,
        method: "PATCH",
        body,
      }),
      invalidatesTags: ["Allergen"],
    }),
    deleteAllergen: builder.mutation<void, { id: number }>({
      query: ({ id }) => ({
        url: `/allergens/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["Allergen"],
    }),
  }),
});

export const {
  useGetAllergensQuery,
  useCreateAllergenMutation,
  useUpdateAllergenMutation,
  useDeleteAllergenMutation,
} = allergenApi;
