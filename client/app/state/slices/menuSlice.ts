import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface MenuState {
  page: number;
  limit: number;
  categoryId?: number;
  query: string;
  minPrice?: number;
  maxPrice?: number;
  spiceLevel?: number;
  sortBy?: "price" | "name" | "created_at";
  order?: "asc" | "desc";
}

const initialState: MenuState = {
  page: 1,
  limit: 10,
  query: "",
};

const menuSlice = createSlice({
  name: "menu",
  initialState,
  reducers: {
    setPage(state, action: PayloadAction<number>) {
      state.page = action.payload;
    },

    setCategory(state, action: PayloadAction<number | undefined>) {
      state.categoryId = action.payload;
      state.page = 1;
    },

    setSearch(state, action: PayloadAction<string>) {
      state.query = action.payload;
      state.page = 1;
    },

    setSort(
      state,
      action: PayloadAction<{
        sortBy: MenuState["sortBy"];
        order: MenuState["order"];
      }>,
    ) {
      state.sortBy = action.payload.sortBy;
      state.order = action.payload.order;
    },

    setPriceRange(
      state,
      action: PayloadAction<{
        minPrice?: number;
        maxPrice?: number;
      }>,
    ) {
      state.minPrice = action.payload.minPrice;
      state.maxPrice = action.payload.maxPrice;
      state.page = 1;
    },

    resetFilters(state) {
      state.page = 1;
      state.categoryId = undefined;
      state.query = "";
      state.minPrice = undefined;
      state.maxPrice = undefined;
      state.sortBy = undefined;
      state.order = undefined;
      state.spiceLevel = undefined;
    },
  },
});

export const {
  setPage,
  setCategory,
  setSearch,
  setSort,
  setPriceRange,
  resetFilters,
} = menuSlice.actions;

export default menuSlice.reducer;
