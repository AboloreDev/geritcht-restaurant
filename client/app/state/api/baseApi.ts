import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type {
  BaseQueryFn,
  FetchArgs,
  FetchBaseQueryError,
} from "@reduxjs/toolkit/query";

const baseQuery = fetchBaseQuery({
  baseUrl: process.env.NEXT_PUBLIC_API_BASE_URL,
  credentials: "include",
  prepareHeaders: (headers) => {
    if (typeof window !== "undefined") {
      const token = localStorage.getItem("accessToken");
      if (token) {
        headers.set("Authorization", `Bearer ${token}`);
      }
    }

    return headers;
  },
});

const baseQueryWithReauth: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  let result = await baseQuery(args, api, extraOptions);

  // Handle 401 Unauthorized errors
  if (result.error && result.error.status === 401) {
    // Get refresh token from localStorage
    const refreshToken =
      typeof window !== "undefined"
        ? localStorage.getItem("refreshToken")
        : null;

    if (refreshToken) {
      const refreshResult = await baseQuery(
        {
          url: "/auth/refresh-token",
          method: "POST",
          body: { refreshToken },
        },
        api,
        extraOptions,
      );

      if (refreshResult.data) {
        const data = refreshResult.data as any;

        const newAccessToken = data?.data?.token?.accessToken;
        const newRefreshToken = data?.data?.token?.refreshToken;

        if (newAccessToken && typeof window !== "undefined") {
          localStorage.setItem("accessToken", newAccessToken);
          if (newRefreshToken) {
            localStorage.setItem("refreshToken", newRefreshToken);
          }

          // Retry the original request with new token
          result = await baseQuery(args, api, extraOptions);
        } else {
          //   api.dispatch(logoutAction());
        }
      } else {
        // api.dispatch(logoutAction());
      }
    } else {
      //   api.dispatch(logoutAction());
    }
  }

  return result;
};

export const baseApi = createApi({
  reducerPath: "baseApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["Category", "Menu", "Reservation"],
  endpoints: (builder) => ({
    // Logout endpoint
    logout: builder.mutation<
      {
        status: boolean;
        code: number;
        message: string;
        data: { message: string };
      },
      void
    >({
      query: () => ({
        url: "/auth/logout",
        method: "DELETE",
      }),
      // Clear all cache when logging out
      invalidatesTags: ["Category", "Menu", "Reservation"],
      async onQueryStarted(_, { dispatch, queryFulfilled }) {
        try {
          await queryFulfilled;

          console.log("logout successful");
        } catch (error) {
          console.error("logout failed, but logging out locally");
        } finally {
          //   dispatch(logoutAction());
        }
      },
    }),
  }),
});

export const { useLogoutMutation } = baseApi;
