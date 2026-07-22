import {
  AvailabilityResponse,
  CheckAvailabilityRequest,
  CreateReservationRequest,
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
        url: "/reservations",
        method: "POST",
        body,
      }),
      invalidatesTags: ["Reservation"],
    }),
  }),
});

export const { useCheckAvailabilityQuery, useCreateReservationMutation } =
  reservationApi;
