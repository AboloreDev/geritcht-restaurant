import {
  InitializePaymentRequest,
  InitializePaymentResponse,
  VerifyPaymentResponse,
} from "../types/paymentTypes";
import { baseApi } from "./baseApi";

export const paymentApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    initializePayment: builder.mutation<
      InitializePaymentResponse,
      InitializePaymentRequest
    >({
      query: (body) => ({
        url: "/payments/initialize",
        method: "POST",
        body,
      }),
    }),

    verifyPayment: builder.query<VerifyPaymentResponse, string>({
      query: (reference) => ({
        url: `/payments/verify/${reference}`,
      }),
    }),
  }),
});

export const { useInitializePaymentMutation, useLazyVerifyPaymentQuery } =
  paymentApi;
