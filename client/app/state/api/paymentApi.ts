import {
  InitializePaymentRequest,
  InitializePaymentResponse,
  PaymentFilterRequest,
  PaymentListResponse,
  VerifyPaymentResponse,
} from "../types/paymentTypes";
import { RefundResponse } from "../types/refundTypes";
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

    getAllUserPayment: builder.query<PaymentListResponse, PaymentFilterRequest>(
      {
        query: () => ({
          url: "/payments/history",
          method: "GET",
        }),
      },
    ),

    adminGetAllPayment: builder.query<
      PaymentListResponse,
      PaymentFilterRequest
    >({
      query: (params) => ({
        url: "/payments/all",
        method: "GET",
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
    }),

    getPaymentDetails: builder.query<PaymentResponse, { id: number }>({
      query: ({ id }) => ({
        url: `/payments/${id}`,
        method: "GET",
      }),
    }),

    // REFUNDS /////
    // paymentApi.ts — additions
    processRefund: builder.mutation<
      RefundResponse,
      { orderId: number; notes?: string }
    >({
      query: ({ orderId, notes }) => ({
        url: `/payments/refund/${orderId}`,
        method: "POST",
        body: { notes },
      }),
      invalidatesTags: ["Refund"],
    }),

    // re-verify, as a mutation-style trigger for a manual admin action —
    // wraps the existing GET /payments/verify/:reference endpoint, since
    // this needs to fire on button click, not passively on page load
    reVerifyPayment: builder.mutation<VerifyPaymentResponse, string>({
      query: (reference) => ({
        url: `/payments/verify/${reference}`,
        method: "GET",
      }),
      invalidatesTags: ["Refund"],
    }),
  }),
});

export const {
  useInitializePaymentMutation,
  useLazyVerifyPaymentQuery,
  useAdminGetAllPaymentQuery,
  useGetAllUserPaymentQuery,
  useGetPaymentDetailsQuery,
  useProcessRefundMutation,
  useVerifyPaymentQuery,
  useReVerifyPaymentMutation,
} = paymentApi;
