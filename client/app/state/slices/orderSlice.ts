import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Order } from "../types/orderTypes";

interface OrderState {
  filterDate?: string;
  filterStatus?: string;
  filterType?: string;
  page?: number;
  pageSize?: number;
  query: string;
  orders: Order[];
}

const initialState: OrderState = {
  filterDate: undefined,
  filterStatus: undefined,
  filterType: undefined,
  page: 1,
  pageSize: 10,
  query: "",
  orders: [],
};

const orderSlice = createSlice({
  name: "order",
  initialState,
  reducers: {
    setPage(state, action: PayloadAction<number>) {
      state.page = action.payload;
    },
    setFilterDate(state, action: PayloadAction<string | undefined>) {
      state.filterDate = action.payload;
      state.page = 1;
    },
    setFilterStatus(state, action: PayloadAction<string | undefined>) {
      state.filterStatus = action.payload;
      state.page = 1;
    },
    setFilterType(state, action: PayloadAction<string | undefined>) {
      state.filterType = action.payload;
      state.page = 1;
    },
    resetOrdersFilters(state) {
      state.filterDate = undefined;
      state.filterStatus = undefined;
      state.filterType = undefined;
      state.page = 1;
      state.query = "";
    },
    setSearch(state, action: PayloadAction<string>) {
      state.query = action.payload;
      state.page = 1;
    },
    appendOrders(state, action) {
      const existing = new Set(state.orders.map((o) => o.id));

      const newOrders = action.payload.filter(
        (order: Order) => !existing.has(order.id),
      );

      state.orders.push(...newOrders);
    },
    setOrders(state, action: PayloadAction<Order[]>) {
      state.orders = action.payload;
    },
  },
});

export const {
  setPage,
  setFilterDate,
  setFilterStatus,
  resetOrdersFilters,
  setFilterType,
  appendOrders,
  setOrders,
} = orderSlice.actions;

export default orderSlice.reducer;
