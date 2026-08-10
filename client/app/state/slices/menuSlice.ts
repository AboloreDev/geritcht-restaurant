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
  sortOrder?: "asc" | "desc";
  create: {
    open: boolean;
    step: 1 | 2;
    menuId: number | null;
    uploadProgress: number;
  };
  edit: {
    open: boolean;
    menuId: number | null;
    step: 1 | 2;
  };
  dialogs: {
    deleteImageId: number | null;
    deleteMenuId: number | null;
    discardChanges: boolean;
  };
  lightbox: {
    open: boolean;
    imageId: number | null;
  };
}

const initialState: MenuState = {
  page: 1,
  limit: 10,
  query: "",

  create: {
    open: false,
    step: 1,
    menuId: null,
    uploadProgress: 0,
  },

  edit: {
    open: false,
    menuId: null,
    step: 1,
  },

  dialogs: {
    deleteImageId: null,
    deleteMenuId: null,
    discardChanges: false,
  },
  lightbox: {
    open: false,
    imageId: null,
  },
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
        order: MenuState["sortOrder"];
      }>,
    ) {
      state.sortBy = action.payload.sortBy;
      state.sortOrder = action.payload.order;
      state.page = 1;
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
      state.sortOrder = undefined;
      state.spiceLevel = undefined;
    },

    searchAcrossAllCategories(state, action: PayloadAction<string>) {
      state.query = action.payload;
      state.categoryId = undefined;
      state.page = 1;
    },

    openCreateSheet(state) {
      state.create.open = true;
    },

    closeCreateSheet(state) {
      state.create.open = false;
    },

    nextCreateStep(state) {
      if (state.create.step === 1) {
        state.create.step = 2;
      }
    },

    previousCreateStep(state) {
      state.create.step = 1;
    },

    setCreatedMenu(state, action: PayloadAction<number>) {
      state.create.menuId = action.payload;
    },

    setUploadProgress(state, action: PayloadAction<number>) {
      state.create.uploadProgress = action.payload;
    },

    resetCreateState(state) {
      state.create = {
        open: false,
        step: 1,
        menuId: null,
        uploadProgress: 0,
      };
    },

    openEditSheet(state, action: PayloadAction<number>) {
      state.edit.open = true;
      state.edit.menuId = action.payload;
      state.edit.step = 1;
    },

    closeEditSheet(state) {
      state.edit.open = false;
      state.edit.menuId = null;
      state.edit.step = 1;
    },

    nextEditStep(state) {
      if (state.edit.step === 1) {
        state.edit.step = 2;
      }
    },

    previousEditStep(state) {
      state.edit.step = 1;
    },

    openDeleteImageDialog(state, action: PayloadAction<number>) {
      state.dialogs.deleteImageId = action.payload;
    },

    closeDeleteImageDialog(state) {
      state.dialogs.deleteImageId = null;
    },

    openDiscardDialog(state) {
      state.dialogs.discardChanges = true;
    },

    closeDiscardDialog(state) {
      state.dialogs.discardChanges = false;
    },

    openLightbox(state, action: PayloadAction<number>) {
      state.lightbox.open = true;
      state.lightbox.imageId = action.payload;
    },

    closeLightbox(state) {
      state.lightbox.open = false;
      state.lightbox.imageId = null;
    },

    openDeleteMenuDialog(state, action: PayloadAction<number>) {
      state.dialogs.deleteMenuId = action.payload;
    },

    closeDeleteMenuDialog(state) {
      state.dialogs.deleteMenuId = null;
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
  searchAcrossAllCategories,
  openCreateSheet,
  closeCreateSheet,
  nextCreateStep,
  previousCreateStep,
  setCreatedMenu,
  resetCreateState,
  openEditSheet,
  closeEditSheet,
  openDeleteImageDialog,
  closeDeleteImageDialog,
  openDiscardDialog,
  closeDiscardDialog,
  openLightbox,
  closeLightbox,
  setUploadProgress,
  openDeleteMenuDialog,
  closeDeleteMenuDialog,
  previousEditStep,
  nextEditStep,
} = menuSlice.actions;

export default menuSlice.reducer;
