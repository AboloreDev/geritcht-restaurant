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
    createTable: builder.mutation<GetTableResponse, CreateTableRequest>({
      query: (body) => ({
        method: "POST",
        url: "/tables/",
        data: body,
      }),
      invalidatesTags: ["Tables"],
    }),
    updateTable: builder.mutation<
      GetTableResponse,
      { id: number; data: UpdateTableRequest }
    >({
      query: ({ id, data }) => ({
        method: "PATCH",
        url: `/tables/${id}`,
        data,
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
    getAllTables: builder.query<GetTablesResponse, void>({
      query: () => ({
        method: "GET",
        url: "/table",
      }),
    }),
    getTable: builder.query<GetTableResponse, { id: number }>({
      query: ({ id }) => ({
        method: "GET",
        url: `/table/${id}`,
      }),
    }),
  }),
});

export const {
  useCreateTableMutation,
  useDeleteTableMutation,
  useGetAllTablesQuery,
  useGetTableQuery,
  useUpdateTableMutation,
} = tableApi;
