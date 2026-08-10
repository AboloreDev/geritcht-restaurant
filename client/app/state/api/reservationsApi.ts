import {
  AvailabilityResponse,
  CheckAvailabilityRequest,
  CreateReservationRequest,
  GetReservationsRequest,
  ReservationListResponse,
  ReservationResponse,
  ReservationSearchResponse,
  SearchReservationRequest,
} from "../types/reservationTypes";
import { baseApi } from "./baseApi";

export const reservationApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    checkAvailability: builder.query<
      AvailabilityResponse,
      CheckAvailabilityRequest
    >({
      query: (params) => ({
        url: "/availability",
        params,
      }),
      providesTags: ["Reservation"],
    }),

    createReservation: builder.mutation<
      { status: boolean; message: string },
      CreateReservationRequest
    >({
      query: (body) => ({
        url: "/reservations/",
        method: "POST",
        body,
      }),
      invalidatesTags: ["Reservation"],
    }),
    getAllUserRservations: builder.query<
      ReservationListResponse,
      GetReservationsRequest
    >({
      query: (params) => ({
        url: "/reservations/user",
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
      providesTags: ["Reservation"],
    }),
    getAllRservations: builder.query<
      ReservationListResponse,
      GetReservationsRequest
    >({
      query: (params) => ({
        url: "/reservations/",
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
      providesTags: ["Reservation"],
    }),
    getReservationById: builder.query<ReservationResponse, { id: number }>({
      query: ({ id }) => `/reservations/${id}`,
      providesTags: ["Reservation"],
    }),
    checkInReservation: builder.mutation<
      { status: boolean; message: string },
      { id: number }
    >({
      query: ({ id }) => ({
        url: `/reservations/${id}`,
        method: "POST",
      }),
      invalidatesTags: ["Reservation"],
    }),
    cancelReservation: builder.mutation<ReservationResponse, { id: number }>({
      query: ({ id }) => ({
        url: `/reservation/${id}/cancel`,
        method: "PATCH",
      }),
    }),
    adminCancelReservation: builder.mutation<
      ReservationResponse,
      { id: number }
    >({
      query: ({ id }) => ({
        url: `/reservation/admin/${id}/cancel`,
        method: "PATCH",
      }),
    }),
    searchReservation: builder.query<
      ReservationSearchResponse,
      SearchReservationRequest
    >({
      query: ({ q }) => ({
        url: "/reservations/search",
        params: { q },
      }),
      providesTags: ["Reservation"],
    }),
  }),
});

export const {
  useCheckAvailabilityQuery,
  useCreateReservationMutation,
  useGetAllUserRservationsQuery,
  useGetReservationByIdQuery,
  useGetAllRservationsQuery,
  useCheckInReservationMutation,
  useSearchReservationQuery,
  useCancelReservationMutation,
  useAdminCancelReservationMutation,
} = reservationApi;
