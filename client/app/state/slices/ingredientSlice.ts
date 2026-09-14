import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface IngredientState {
  page: number;
  query: string;
  create: { open: boolean };
  edit: { open: boolean; ingredientId: number | null };
  dialogs: { deleteIngredientId: number | null };
}

const initialState: IngredientState = {
  page: 1,
  query: "",
  create: { open: false },
  edit: { open: false, ingredientId: null },
  dialogs: { deleteIngredientId: null },
};

const ingredientSlice = createSlice({
  name: "ingredient",
  initialState,
  reducers: {
    setPage(state, action: PayloadAction<number>) {
      state.page = action.payload;
    },
    setSearch(state, action: PayloadAction<string>) {
      state.query = action.payload;
      state.page = 1;
    },
    openCreateIngredientSheet(state) {
      state.create.open = true;
    },
    closeCreateIngredientSheet(state) {
      state.create.open = false;
    },
    openEditIngredientSheet(state, action: PayloadAction<number>) {
      state.edit.open = true;
      state.edit.ingredientId = action.payload;
    },
    closeEditIngredientSheet(state) {
      state.edit.open = false;
      state.edit.ingredientId = null;
    },
    openDeleteIngredientDialog(state, action: PayloadAction<number>) {
      state.dialogs.deleteIngredientId = action.payload;
    },
    closeDeleteIngredientDialog(state) {
      state.dialogs.deleteIngredientId = null;
    },
  },
});

export const {
  setPage,
  setSearch,
  openCreateIngredientSheet,
  closeCreateIngredientSheet,
  openEditIngredientSheet,
  closeEditIngredientSheet,
  openDeleteIngredientDialog,
  closeDeleteIngredientDialog,
} = ingredientSlice.actions;
export default ingredientSlice.reducer;
