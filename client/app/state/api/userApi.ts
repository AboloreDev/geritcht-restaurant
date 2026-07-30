import {
  ChangePasswordRequest,
  GetProfileResponse,
  UpdateProfileRequest,
} from "../types/authTypes";
import { baseApi } from "./baseApi";

export const userApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getUserProfile: builder.query<GetProfileResponse, void>({
      query: () => ({ url: "/users/profile" }),
      providesTags: ["User"],
    }),

    updateUserProfile: builder.mutation<
      GetProfileResponse,
      UpdateProfileRequest
    >({
      query: (body) => ({
        url: "/users/profile",
        method: "PATCH",
        body,
      }),
      invalidatesTags: ["User"],
    }),

    changePassword: builder.mutation<
      { status: boolean; message: string },
      ChangePasswordRequest
    >({
      query: (body) => ({
        url: "/users/password-change",
        method: "PATCH",
        body,
      }),
    }),

    deactivateAccount: builder.mutation<
      { status: boolean; message: string },
      void
    >({
      query: () => ({
        url: "/users/profile/deactivate",
        method: "PATCH",
      }),
      invalidatesTags: ["User"],
    }),
  }),
});

export const {
  useGetUserProfileQuery,
  useUpdateUserProfileMutation,
  useChangePasswordMutation,
  useDeactivateAccountMutation,
} = userApi;
