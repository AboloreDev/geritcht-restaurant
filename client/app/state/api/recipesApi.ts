import {
  LinkIngredientRequest,
  RecipeListResponse,
  RecipeResponse,
  UpdateLinkItemRequest,
} from "../types/ingredientTypes";
import { MessageResponse } from "../types/userTypes";
import { baseApi } from "./baseApi";

export const recipeApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getRecipes: builder.query<RecipeListResponse, { menuId: number }>({
      query: ({ menuId }) => ({ url: `/recipes/${menuId}` }),
      providesTags: ["Recipe"],
    }),

    addRecipe: builder.mutation<
      RecipeResponse,
      { menuId: number; body: LinkIngredientRequest }
    >({
      query: ({ menuId, body }) => ({
        url: `/recipes/${menuId}`,
        method: "POST",
        body,
      }),
      invalidatesTags: ["Recipe"],
    }),

    updateRecipe: builder.mutation<
      RecipeResponse,
      { menuId: number; ingredientId: number; body: UpdateLinkItemRequest }
    >({
      query: ({ menuId, ingredientId, body }) => ({
        url: `/recipes/${menuId}/${ingredientId}`,
        method: "PATCH",
        body,
      }),
      invalidatesTags: ["Recipe"],
    }),

    deleteRecipe: builder.mutation<
      MessageResponse,
      { menuId: number; ingredientId: number }
    >({
      query: ({ menuId, ingredientId }) => ({
        url: `/recipes/${menuId}/${ingredientId}`,
        method: "DELETE",
      }),
      invalidatesTags: ["Recipe"],
    }),
  }),
});

export const {
  useGetRecipesQuery,
  useAddRecipeMutation,
  useUpdateRecipeMutation,
  useDeleteRecipeMutation,
} = recipeApi;
