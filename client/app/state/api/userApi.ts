import {
  ChangePasswordRequest,
  GetProfileResponse,
  UpdateProfileRequest,
} from "../types/authTypes";
import {
  GetAllStaffsResponse,
  GetAllUsersResponse,
  User,
  UserResponse,
  UserSearchRequest,
} from "../types/userTypes";
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
    activateAccount: builder.mutation<
      { status: boolean; message: string },
      void
    >({
      query: () => ({
        url: "/users/profile/activate",
        method: "PATCH",
      }),
      invalidatesTags: ["User"],
    }),
    getAllUsers: builder.query<
      GetAllUsersResponse,
      { page?: number; limit?: number; q?: string }
    >({
      query: () => ({
        url: `/users/`,
      }),
      providesTags: ["User"],
    }),

    // STAFFS
    deactivateStaffAccount: builder.mutation<
      { status: boolean; message: string },
      void
    >({
      query: () => ({
        url: "/staff/profile/deactivate",
        method: "PATCH",
      }),
      invalidatesTags: ["User"],
    }),
    activateStaffAccount: builder.mutation<
      { status: boolean; message: string },
      void
    >({
      query: () => ({
        url: "/staff/profile/activate",
        method: "PATCH",
      }),
      invalidatesTags: ["User"],
    }),
    getStaffProfile: builder.query<GetProfileResponse, void>({
      query: () => ({ url: "/staff/profile" }),
      providesTags: ["User"],
    }),

    updateStaffProfile: builder.mutation<
      GetProfileResponse,
      UpdateProfileRequest
    >({
      query: (body) => ({
        url: "/staff/profile",
        method: "PATCH",
        body,
      }),
      invalidatesTags: ["User"],
    }),
    getAllStaff: builder.query<
      GetAllStaffsResponse,
      { page: number; limit?: number }
    >({
      query: () => ({
        url: `/staff/`,
      }),
      providesTags: ["User"],
    }),
    searchUser: builder.query<UserResponse, UserSearchRequest>({
      query: ({ q }) => ({
        url: "/users/search",
        params: { q },
      }),
      providesTags: ["User"],
    }),

    updateUserRole: builder.mutation<
      { status: boolean; message: string; data: UserResponse },
      { id: number; role: "user" | "staff" }
    >({
      query: ({ id, role }) => ({
        url: `/users/${id}/role/update`,
        method: "PATCH",
        body: { role },
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
  useGetAllUsersQuery,
  useActivateAccountMutation,
  useGetStaffProfileQuery,
  useUpdateStaffProfileMutation,
  useDeactivateStaffAccountMutation,
  useActivateStaffAccountMutation,
  useGetAllStaffQuery,
  useSearchUserQuery,
  useUpdateUserRoleMutation,
} = userApi;
