import { get } from "http";
import {
  CreateTakeoutOrderRequest,
  GetOrdersRequest,
  OrderListResponse,
  OrderResponse,
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
          currentCache.data.push(...newResponse.data);
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
        currentCache.data.push(...newResponse.data);
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
  }),
});

export const {
  useCreateTakoutOrderMutation,
  useGetAllOrdersQuery,
  useGetAllUserTakeoutOrdersQuery,
  useGetOrderByIdQuery,
} = orderApi;
