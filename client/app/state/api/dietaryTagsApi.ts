import {
  CreateDietaryTagRequest,
  dietaryTagsResponse,
  GetDietaryTagsResponse,
  UpdateDietaryTagRequest,
} from "../types/dietaryTagsTypes";
import { baseApi } from "./baseApi";

export const dietaryTagsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getDietaryTags: builder.query<GetDietaryTagsResponse, void>({
      query: () => "/tags/",
      providesTags: ["DietaryTags"],
    }),

    createDietaryTags: builder.mutation<
      dietaryTagsResponse,
      CreateDietaryTagRequest
    >({
      query: (body) => ({
        url: "/tags/",
        method: "POST",
        body,
      }),
      invalidatesTags: ["DietaryTags"],
    }),
    updateDietaryTags: builder.mutation<
      dietaryTagsResponse,
      { id: number; body: UpdateDietaryTagRequest }
    >({
      query: ({ id, body }) => ({
        url: `/tags/${id}`,
        method: "PATCH",
        body,
      }),
      invalidatesTags: ["DietaryTags"],
    }),
    deleteDietaryTags: builder.mutation<void, { id: number }>({
      query: ({ id }) => ({
        url: `/tags/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["DietaryTags"],
    }),
  }),
});

export const {
  useCreateDietaryTagsMutation,
  useDeleteDietaryTagsMutation,
  useUpdateDietaryTagsMutation,
  useGetDietaryTagsQuery,
} = dietaryTagsApi;
