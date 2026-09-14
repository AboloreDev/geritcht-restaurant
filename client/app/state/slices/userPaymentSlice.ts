import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { PaymentStatus } from "../types/paymentTypes";

interface UserPaymentState {
  page: number;
  pageSize: number;
  filterStatus?: PaymentStatus;
}

const initialState: UserPaymentState = {
  page: 1,
  pageSize: 15,
  filterStatus: undefined,
};

const userPaymentSlice = createSlice({
  name: "userPayment",
  initialState,
  reducers: {
    setPage(state, action: PayloadAction<number>) {
      state.page = action.payload;
    },
    setFilterStatus(state, action: PayloadAction<PaymentStatus | undefined>) {
      state.filterStatus = action.payload;
      state.page = 1;
    },
  },
});

export const { setPage, setFilterStatus } = userPaymentSlice.actions;
export default userPaymentSlice.reducer;
