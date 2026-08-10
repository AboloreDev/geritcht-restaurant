import {
  CreateTableRequest,
  GetTableResponse,
  GetTablesResponse,
  UpdateTableRequest,
} from "../types/tableTypes";
import { MessageResponse } from "../types/userTypes";
import { baseApi } from "./baseApi";

export const tableApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAllTables: builder.query<GetTablesResponse, void>({
      query: () => ({ url: "/tables" }),
      providesTags: ["Tables"],
    }),

    getTable: builder.query<GetTableResponse, { id: number }>({
      query: ({ id }) => ({ url: `/tables/${id}` }),
      providesTags: ["Tables"],
    }),

    createTable: builder.mutation<GetTableResponse, CreateTableRequest>({
      query: (body) => ({
        method: "POST",
        url: "/tables/",
        body,
      }),
      invalidatesTags: ["Tables"],
    }),

    updateTable: builder.mutation<
      GetTableResponse,
      { id: number; body: UpdateTableRequest }
    >({
      query: ({ id, body }) => ({
        method: "PATCH",
        url: `/tables/${id}`,
        body,
      }),
      invalidatesTags: ["Tables"],
    }),

    deleteTable: builder.mutation<MessageResponse, { id: number }>({
      query: ({ id }) => ({
        method: "DELETE",
        url: `/tables/${id}`,
      }),
      invalidatesTags: ["Tables"],
    }),
  }),
});

export const {
  useGetAllTablesQuery,
  useGetTableQuery,
  useCreateTableMutation,
  useUpdateTableMutation,
  useDeleteTableMutation,
} = tableApi;
