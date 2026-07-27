import { CreateTakeoutOrderRequest, OrderResponse } from "../types/orderTypes";
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
  }),
});

export const { useCreateTakoutOrderMutation } = orderApi;
