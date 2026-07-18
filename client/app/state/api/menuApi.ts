import {
  GetMenusRequest,
  GetMenusResponse,
  SearchMenuRequest,
} from "../types/menuTypes";
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

    searchMenu: builder.query<GetMenusResponse, SearchMenuRequest>({
      query: ({ query }) => ({
        url: "/menu/search",
        params: { query },
      }),
      providesTags: ["Menu"],
    }),
  }),
});

export const { useGetMenusQuery } = menuApi;
