import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface OrderState {
  filterDate?: string;
  filterStatus?: string;
  filterType?: string;
  page?: number;
  pageSize?: number;
}

const initialState: OrderState = {
  filterDate: undefined,
  filterStatus: undefined,
  filterType: undefined,
  page: 1,
  pageSize: 10,
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
    resetReservationFilters(state) {
      state.filterDate = undefined;
      state.filterStatus = undefined;
      state.filterType = undefined;
      state.page = 1;
    },
  },
});

export const {
  setPage,
  setFilterDate,
  setFilterStatus,
  resetReservationFilters,
  setFilterType,
} = orderSlice.actions;

export default orderSlice.reducer;
