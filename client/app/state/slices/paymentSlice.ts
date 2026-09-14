import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { PaymentStatus } from "../types/paymentTypes";

interface PaymentState {
  page: number;
  pageSize: number;
  filterStatus?: string;
  query: string;
}

const initialState: PaymentState = {
  page: 1,
  pageSize: 15,
  filterStatus: undefined,
  query: "",
};

const paymentSlice = createSlice({
  name: "payment",
  initialState,
  reducers: {
    setPage(state, action: PayloadAction<number>) {
      state.page = action.payload;
    },
    setFilterStatus(state, action: PayloadAction<string | undefined>) {
      state.filterStatus = action.payload;
      state.page = 1;
    },
    setSearch(state, action: PayloadAction<string>) {
      state.query = action.payload;
      state.page = 1;
    },
    resetFilters(state) {
      state.filterStatus = undefined;
      state.query = "";
      state.page = 1;
    },
  },
});

export const { setPage, setFilterStatus, setSearch, resetFilters } =
  paymentSlice.actions;

export default paymentSlice.reducer;
