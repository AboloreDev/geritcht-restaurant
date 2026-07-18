import {
  GetCategoriesRequest,
  GetCategoriesResponse,
  SearchCategoriesRequest,
} from "../types/categoriesTypes";
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

    searchCategories: builder.query<
      GetCategoriesResponse,
      SearchCategoriesRequest
    >({
      query: ({ query }) => ({
        url: "/categories/search",
        params: { query },
      }),
      providesTags: ["Category"],
    }),
  }),
});

export const { useGetCategoriesQuery, useSearchCategoriesQuery } = categoryApi;
