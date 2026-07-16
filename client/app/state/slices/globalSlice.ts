import { createSlice } from "@reduxjs/toolkit";

export interface Global {
  showPassword: boolean;
  sidebarOpen: boolean;
}

export const initialState: Global = {
  showPassword: false,
  sidebarOpen: false,
};

export const globalSlice = createSlice({
  name: "global",
  initialState,
  reducers: {
    setShowPassword: (state) => {
      state.showPassword = !state.showPassword;
    },
    setSidebarOpen: (state) => {
      state.sidebarOpen = !state.sidebarOpen;
    },
  },
});

export const { setShowPassword, setSidebarOpen } = globalSlice.actions;

export default globalSlice.reducer;
