import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type {
  BaseQueryFn,
  FetchArgs,
  FetchBaseQueryError,
} from "@reduxjs/toolkit/query";
import { Mutex } from "async-mutex";

// ensures only one refresh request fires at a time, even if several
// requests 401 simultaneously — others wait for the in-flight refresh
const mutex = new Mutex();

const baseQuery = fetchBaseQuery({
  baseUrl: process.env.NEXT_PUBLIC_API_BASE_URL,
  credentials: "include", // sends the httpOnly refresh-token cookie automatically
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

function clearAuthAndRedirect() {
  if (typeof window !== "undefined") {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("user");
    window.location.href = "/login";
  }
}

const baseQueryWithReauth: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  // wait if another request is already refreshing
  await mutex.waitForUnlock();

  let result = await baseQuery(args, api, extraOptions);

  if (result.error && result.error.status === 401) {
    if (!mutex.isLocked()) {
      const release = await mutex.acquire();

      try {
        // re-check: another request might have already refreshed
        // while we were acquiring the lock
        const refreshResult = await baseQuery(
          { url: "/auth/refresh", method: "POST" },
          api,
          extraOptions,
        );

        if (refreshResult.data) {
          const data = refreshResult.data as any;
          const newAccessToken = data?.data?.token?.accessToken;

          if (newAccessToken && typeof window !== "undefined") {
            localStorage.setItem("accessToken", newAccessToken);
            result = await baseQuery(args, api, extraOptions);
          } else {
            clearAuthAndRedirect();
          }
        } else {
          clearAuthAndRedirect();
        }
      } finally {
        release();
      }
    } else {
      // someone else is already refreshing — wait, then retry with
      // whatever token they end up setting
      await mutex.waitForUnlock();
      result = await baseQuery(args, api, extraOptions);
    }
  }

  return result;
};

export const baseApi = createApi({
  reducerPath: "baseApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: [
    "Category",
    "Menu",
    "Reservation",
    "Auth",
    "Cart",
    "Orders",
    "User",
  ],
  endpoints: (builder) => ({
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
      invalidatesTags: [
        "Category",
        "Menu",
        "Reservation",
        "Auth",
        "Cart",
        "Orders",
        "User",
      ],
      async onQueryStarted(_, { queryFulfilled }) {
        try {
          await queryFulfilled;
        } finally {
          clearAuthAndRedirect();
        }
      },
    }),
  }),
});

export const { useLogoutMutation } = baseApi;
