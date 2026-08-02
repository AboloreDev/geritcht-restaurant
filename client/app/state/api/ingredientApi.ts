import {
  GetIngredientsRequest,
  IngredientResponse,
  InventoryAlertResponse,
} from "../types/ingredientTypes";
import { baseApi } from "./baseApi";

export const ingredientApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAllIngredients: builder.query<IngredientResponse, GetIngredientsRequest>(
      {
        query: (params) => ({
          url: "/ingredients/",
          params,
        }),
        providesTags: ["Ingredient"],
      },
    ),
    getLowStockIngredients: builder.query<IngredientResponse, void>({
      query: () => ({
        url: "/ingredients/low-stock",
      }),
      providesTags: ["Ingredient"],
    }),
    getInventoryAlerts: builder.query<InventoryAlertResponse, void>({
      query: () => ({
        url: "/ingredients/alerts",
      }),
      providesTags: ["Ingredient"],
    }),
  }),
});

export const {
  useGetAllIngredientsQuery,
  useGetLowStockIngredientsQuery,
  useGetInventoryAlertsQuery,
} = ingredientApi;
