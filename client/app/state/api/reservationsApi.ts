import {
  AvailabilityResponse,
  CheckAvailabilityRequest,
  CreateReservationRequest,
  GetReservationsRequest,
  ReservationListResponse,
  ReservationResponse,
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
        currentCache.reservations.push(...newResponse.reservations);
        currentCache.meta = newResponse.meta;
      },
      forceRefetch: ({ currentArg, previousArg }) => {
        return currentArg?.page !== previousArg?.page;
      },
      providesTags: ["Reservation"],
    }),
    getReservationById: builder.query<ReservationResponse, { id: number }>({
      query: ({ id }) => `/reservations/${id}/user`,
      providesTags: ["Reservation"],
    }),
  }),
});

export const {
  useCheckAvailabilityQuery,
  useCreateReservationMutation,
  useGetAllUserRservationsQuery,
  useGetReservationByIdQuery,
} = reservationApi;
