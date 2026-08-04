import {
  CreateMenuRequest,
  GetMenusRequest,
  GetMenusResponse,
  GetSingleMenuResponse,
  SearchMenuRequest,
  ToggleAvailabilityRequest,
  UpdateMenuRequest,
} from "../types/menuTypes";
import { MessageResponse } from "../types/userTypes";
import { baseApi } from "./baseApi";

export const menuApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getMenus: builder.query<GetMenusResponse, GetMenusRequest>({
      query: (params) => ({
        url: "/menu",
        params,
      }),
      serializeQueryArgs: ({ queryArgs }) => {
        const { page, ...stableArgs } = queryArgs;
        return JSON.stringify(stableArgs);
      },
      merge: (currentCache, newResponse, { arg }) => {
        if (!arg.page || arg.page === 1) {
          return newResponse;
        }
        currentCache.data.push(...newResponse.data);
        currentCache.meta = newResponse.meta;
      },
      forceRefetch: ({ currentArg, previousArg }) => {
        return currentArg?.page !== previousArg?.page;
      },
      providesTags: ["Menu"],
    }),

    getSingleMenu: builder.query<GetSingleMenuResponse, { id: string }>({
      query: ({ id }) => `/menu/${id}`,
      providesTags: ["Menu"],
    }),

    searchMenu: builder.query<GetMenusResponse, SearchMenuRequest>({
      query: ({ q }) => ({
        url: "/menu/search",
        params: { q },
      }),
      providesTags: ["Menu"],
    }),

    createMenu: builder.mutation<GetSingleMenuResponse, CreateMenuRequest>({
      query: (body) => ({
        method: "POST",
        url: "/menu/",
        body,
      }),
      invalidatesTags: ["Menu"],
    }),
    updateMenu: builder.mutation<
      GetSingleMenuResponse,
      { id: number; body: UpdateMenuRequest }
    >({
      query: ({ id, body }) => ({
        method: "PATCH",
        url: `/menu/${id}`,
        body,
      }),
      invalidatesTags: ["Menu"],
    }),
    uploadMenuImage: builder.mutation<
      { message: string; data: { url: string } },
      { id: number; image: File; is_primary: boolean }
    >({
      query: ({ id, image, is_primary }) => {
        const formData = new FormData();
        formData.append("image", image);
        formData.append("is_primary", String(is_primary));

        return {
          url: `/menus/${id}/images`,
          method: "POST",
          body: formData,
        };
      },
      invalidatesTags: ["Menu"],
    }),
    deleteImageUpload: builder.mutation<MessageResponse, { id: number }>({
      query: ({ id }) => ({
        method: "DELETE",
        url: `/images/${id}`,
      }),
      invalidatesTags: ["Menu"],
    }),
    toggleMenuAvailability: builder.mutation<
      MessageResponse,
      { id: number; body: ToggleAvailabilityRequest }
    >({
      query: ({ id, body }) => ({
        method: "PATCH",
        url: `/menu/${id}/toggle`,
        body,
      }),
      invalidatesTags: ["Menu"],
    }),
  }),
});

export const {
  useGetMenusQuery,
  useSearchMenuQuery,
  useGetSingleMenuQuery,
  useCreateMenuMutation,
  useUpdateMenuMutation,
  useDeleteImageUploadMutation,
  useToggleMenuAvailabilityMutation,
  useUploadMenuImageMutation,
} = menuApi;
