import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface CategoryState {
  selectedCategoryId: number | null;
  search: string;
  page: number;
  limit: number;
}
const initialState: CategoryState = {
  selectedCategoryId: null,
  search: "",
  page: 1,
  limit: 100,
};

const categorySlice = createSlice({
  name: "category",
  initialState,
  reducers: {
    setSelectedCategory(state, action: PayloadAction<number | null>) {
      state.selectedCategoryId = action.payload;
    },

    setCategorySearch(state, action: PayloadAction<string>) {
      state.search = action.payload;
      state.page = 1;
    },

    setCategoryPage(state, action: PayloadAction<number>) {
      state.page = action.payload;
    },

    resetCategory(state) {
      state.selectedCategoryId = null;
      state.search = "";
      state.page = 1;
    },
  },
});

export const {
  setSelectedCategory,
  setCategorySearch,
  setCategoryPage,
  resetCategory,
} = categorySlice.actions;

export default categorySlice.reducer;
