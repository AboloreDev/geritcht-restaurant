import { GetMenusRequest, GetMenusResponse } from "../types/menuTypes";
import { baseApi } from "./baseApi";

export const menuApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getMenus: builder.query<GetMenusResponse, GetMenusRequest>({
      query: (params) => ({
        url: "/menu",
        params,
      }),
    }),
  }),
});

export const { useGetMenusQuery } = menuApi;
