import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface userState {
  activeTab: "users" | "staff";
  users: { page: number; query: string };
  staff: { page: number };
  dialogs: { promoteUserId: number | null };
}

const initialState: userState = {
  activeTab: "users",
  users: { page: 1, query: "" },
  staff: { page: 1 },
  dialogs: { promoteUserId: null },
};

const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setUsersPage(state, action: PayloadAction<number>) {
      state.users.page = action.payload;
    },
    setUsersSearch(state, action: PayloadAction<string>) {
      state.users.query = action.payload;
      state.users.page = 1;
    },
    setStaffPage(state, action: PayloadAction<number>) {
      state.staff.page = action.payload;
    },
    openPromoteDialog(state, action: PayloadAction<number>) {
      state.dialogs.promoteUserId = action.payload;
    },
    closePromoteDialog(state) {
      state.dialogs.promoteUserId = null;
    },
    setActiveTab(state, action: PayloadAction<"users" | "staff">) {
      state.activeTab = action.payload;
    },
  },
});

export const {
  setUsersPage,
  setUsersSearch,
  setStaffPage,
  openPromoteDialog,
  closePromoteDialog,
  setActiveTab,
} = userSlice.actions;
export default userSlice.reducer;
