import { createSlice } from "@reduxjs/toolkit";

interface CartState {
  isDrawerOpen: boolean;
}

const initialState: CartState = {
  isDrawerOpen: false,
};

const cartSlice = createSlice({
  name: "cart",
  initialState,
  reducers: {
    openCartDrawer(state) {
      state.isDrawerOpen = true;
    },
    closeCartDrawer(state) {
      state.isDrawerOpen = false;
    },
    toggleCartDrawer(state) {
      state.isDrawerOpen = !state.isDrawerOpen;
    },
  },
});

export const { openCartDrawer, closeCartDrawer, toggleCartDrawer } =
  cartSlice.actions;
export default cartSlice.reducer;
