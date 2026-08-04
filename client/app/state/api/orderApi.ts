import { get } from "http";
import {
  CreateTakeoutOrderRequest,
  GetOrdersRequest,
  OrderListResponse,
  OrderResponse,
  OrderSearchResponse,
  SearchOrderRequest,
} from "../types/orderTypes";
import { baseApi } from "./baseApi";

export const orderApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    createTakoutOrder: builder.mutation<
      OrderResponse,
      CreateTakeoutOrderRequest
    >({
      query: (body) => ({
        url: "/orders/takeout",
        method: "POST",
        body,
      }),
    }),
    getAllUserTakeoutOrders: builder.query<OrderListResponse, GetOrdersRequest>(
      {
        query: (params) => ({
          url: "/orders/takeout/all",
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
          currentCache.data.orders.push(...newResponse.data.orders);
        },
        forceRefetch: ({ currentArg, previousArg }) => {
          return currentArg?.page !== previousArg?.page;
        },
        providesTags: ["Orders"],
      },
    ),
    getAllOrders: builder.query<OrderListResponse, GetOrdersRequest>({
      query: (params) => ({
        url: "/orders/all",
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
        currentCache.data.orders.push(...newResponse.data.orders);
      },
      forceRefetch: ({ currentArg, previousArg }) => {
        return currentArg?.page !== previousArg?.page;
      },
      providesTags: ["Orders"],
    }),
    getOrderById: builder.query<OrderResponse, { id: number }>({
      query: ({ id }) => `/orders/takeout/${id}`,
      providesTags: ["Orders"],
    }),
    adminGetOrderById: builder.query<OrderResponse, { id: number }>({
      query: ({ id }) => `/orders/${id}`,
      providesTags: ["Orders"],
    }),
    searchOrder: builder.query<OrderSearchResponse, SearchOrderRequest>({
      query: ({ q }) => ({
        url: "/orders/search",
        params: { q },
      }),
      providesTags: ["Orders"],
    }),
    adminCancelOrder: builder.mutation<
      { status: boolean; message: string },
      { id: number }
    >({
      query: ({ id }) => ({
        url: `/orders/${id}/cancel`,
        method: "PATCH",
      }),
      invalidatesTags: ["Orders"],
    }),
  }),
});

export const {
  useCreateTakoutOrderMutation,
  useGetAllOrdersQuery,
  useGetAllUserTakeoutOrdersQuery,
  useGetOrderByIdQuery,
  useSearchOrderQuery,
  useAdminCancelOrderMutation,
  useAdminGetOrderByIdQuery,
} = orderApi;
