import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface ReservationState {
  isModalOpen: boolean;
  date: string;
  timeSlot: string;
  partySize: number;
}

const initialState: ReservationState = {
  isModalOpen: false,
  date: "",
  timeSlot: "",
  partySize: 2,
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
  },
});

export const {
  openBookingModal,
  closeBookingModal,
  setDate,
  setTimeSlot,
  setPartySize,
} = reservationSlice.actions;

export default reservationSlice.reducer;
