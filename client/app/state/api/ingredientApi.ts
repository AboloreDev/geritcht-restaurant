import {
  CreateIngredientRequest,
  GetInventoryAlertsResponse,
  IngredientListResponse,
  IngredientResponse,
  IngredientSearchRequest,
  IngredientSearchResponse,
  ThresholdRequest,
  UpdateIngredientRequest,
} from "../types/ingredientTypes";
import { MessageResponse } from "../types/userTypes";
import { baseApi } from "./baseApi";

export const ingredientApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAllIngredients: builder.query<
      IngredientListResponse,
      { page?: number; limit?: number }
    >({
      query: (params) => ({ url: "/ingredients/", params }),
      providesTags: ["Ingredient"],
    }),

    getIngredient: builder.query<IngredientResponse, { id: number }>({
      query: ({ id }) => ({ url: `/ingredients/${id}` }),
      providesTags: ["Ingredient"],
    }),

    searchIngredient: builder.query<
      IngredientSearchResponse,
      IngredientSearchRequest
    >({
      query: (params) => ({ url: "/ingredients/search", params }),
      providesTags: ["Ingredient"],
    }),

    getLowStockIngredients: builder.query<IngredientListResponse, void>({
      query: () => ({ url: "/ingredients/low-stock" }),
      providesTags: ["Ingredient"],
    }),

    getInventoryAlerts: builder.query<GetInventoryAlertsResponse, void>({
      query: () => ({ url: "/ingredients/alerts" }),
      providesTags: ["Ingredient"],
    }),

    createIngredient: builder.mutation<
      IngredientResponse,
      CreateIngredientRequest
    >({
      query: (body) => ({
        url: "/ingredients",
        method: "POST",
        body,
      }),
      invalidatesTags: ["Ingredient"],
    }),

    updateIngredient: builder.mutation<
      IngredientResponse,
      { id: number; body: UpdateIngredientRequest }
    >({
      query: ({ id, body }) => ({
        url: `/ingredients/${id}`,
        method: "PATCH",
        body,
      }),
      invalidatesTags: ["Ingredient"],
    }),

    deleteIngredient: builder.mutation<MessageResponse, { id: number }>({
      query: ({ id }) => ({
        url: `/ingredients/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["Ingredient"],
    }),

    setThresholdLimit: builder.mutation<
      IngredientResponse,
      { id: number; body: ThresholdRequest }
    >({
      query: ({ id, body }) => ({
        url: `/ingredients/${id}/limit`,
        method: "POST",
        body,
      }),
      invalidatesTags: ["Ingredient"],
    }),
  }),
});

export const {
  useGetAllIngredientsQuery,
  useGetIngredientQuery,
  useSearchIngredientQuery,
  useGetLowStockIngredientsQuery,
  useGetInventoryAlertsQuery,
  useCreateIngredientMutation,
  useUpdateIngredientMutation,
  useDeleteIngredientMutation,
  useSetThresholdLimitMutation,
} = ingredientApi;
