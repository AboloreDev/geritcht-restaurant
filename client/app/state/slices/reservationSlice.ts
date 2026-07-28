import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface ReservationState {
  // booking modal draft state
  isModalOpen: boolean;
  date: string;
  timeSlot: string;
  partySize: number;

  // my-reservations list filter state
  filterDate?: string;
  filterStatus?: string;
  page?: number;
  pageSize?: number;
}

const initialState: ReservationState = {
  isModalOpen: false,
  date: "",
  timeSlot: "",
  partySize: 2,
  filterDate: undefined,
  filterStatus: undefined,
  page: 1,
  pageSize: 10,
};

const reservationSlice = createSlice({
  name: "reservation",
  initialState,
  reducers: {
    openBookingModal(state) {
      state.isModalOpen = true;
    },
    closeBookingModal(state) {
      state.isModalOpen = false;
    },
    setDate(state, action: PayloadAction<string>) {
      state.date = action.payload;
    },
    setTimeSlot(state, action: PayloadAction<string>) {
      state.timeSlot = action.payload;
    },
    setPartySize(state, action: PayloadAction<number>) {
      state.partySize = action.payload;
    },
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
    resetReservationFilters(state) {
      state.filterDate = undefined;
      state.filterStatus = undefined;
      state.page = 1;
    },
  },
});

export const {
  openBookingModal,
  closeBookingModal,
  setDate,
  setTimeSlot,
  setPartySize,
  setPage,
  setFilterDate,
  setFilterStatus,
  resetReservationFilters,
} = reservationSlice.actions;

export default reservationSlice.reducer;
