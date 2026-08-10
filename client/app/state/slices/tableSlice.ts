import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface TableState {
  create: {
    open: boolean;
  };
  edit: {
    open: boolean;
    tableId: number | null;
  };
  dialogs: {
    deleteTableId: number | null;
  };
}

const initialState: TableState = {
  create: { open: false },
  edit: { open: false, tableId: null },
  dialogs: { deleteTableId: null },
};

const tableSlice = createSlice({
  name: "table",
  initialState,
  reducers: {
    openCreateTableSheet(state) {
      state.create.open = true;
    },
    closeCreateTableSheet(state) {
      state.create.open = false;
    },
    openEditTableSheet(state, action: PayloadAction<number>) {
      state.edit.open = true;
      state.edit.tableId = action.payload;
    },
    closeEditTableSheet(state) {
      state.edit.open = false;
      state.edit.tableId = null;
    },
    openDeleteTableDialog(state, action: PayloadAction<number>) {
      state.dialogs.deleteTableId = action.payload;
    },
    closeDeleteTableDialog(state) {
      state.dialogs.deleteTableId = null;
    },
  },
});

export const {
  openCreateTableSheet,
  closeCreateTableSheet,
  openEditTableSheet,
  closeEditTableSheet,
  openDeleteTableDialog,
  closeDeleteTableDialog,
} = tableSlice.actions;

export default tableSlice.reducer;
