import { create } from "domain";
import {
  CategoryResponse,
  GetCategoriesRequest,
  GetCategoriesResponse,
  SearchCategoriesRequest,
} from "../types/categoriesTypes";
import { MenuCategory } from "../types/menuTypes";
import { baseApi } from "./baseApi";

export const categoryApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getCategories: builder.query<GetCategoriesResponse, GetCategoriesRequest>({
      query: (params) => ({
        url: "/categories",
        params,
      }),
      providesTags: ["Category"],
    }),
    editCategory: builder.mutation<
      CategoryResponse,
      { id: number; data: Partial<MenuCategory> }
    >({
      query: ({ id, data }) => ({
        url: `/categories/${id}`,
        method: "PATCH",
        body: data,
      }),
      invalidatesTags: ["Category"],
    }),
    deleteCategory: builder.mutation<
      { status: boolean; message: string },
      { id: number }
    >({
      query: ({ id }) => ({
        url: `/categories/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["Category"],
    }),

    createCategory: builder.mutation<
      CategoryResponse,
      { data: Partial<MenuCategory> }
    >({
      query: ({ data }) => ({
        url: "/categories/",
        method: "POST",
        body: data,
      }),
      invalidatesTags: ["Category"],
    }),

    searchCategories: builder.query<
      GetCategoriesResponse,
      SearchCategoriesRequest
    >({
      query: ({ q }) => ({
        url: "/categories/search",
        params: { q },
      }),
      providesTags: ["Category"],
    }),
  }),
});

export const {
  useGetCategoriesQuery,
  useSearchCategoriesQuery,
  useEditCategoryMutation,
  useDeleteCategoryMutation,
  useCreateCategoryMutation,
} = categoryApi;
