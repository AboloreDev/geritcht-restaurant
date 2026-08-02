import { Menu } from "./menuTypes";

export interface Ingredient {
  id: number;
  name: string;
  unit: string;
  current_stock: number;
  min_threshold: number;
  is_low: boolean;
  created_at: string;
  updated_at: string;
}

export interface IngredientResponse {
  status: boolean;
  message: string;
  data: Ingredient[];
  meta: {
    total: number;
    page: number;
    page_size: number;
  };
}

export interface GetIngredientsRequest {
  page?: number;
  page_size?: number;
}

export interface InventoryAlert {
  LowStockIngredients: Ingredient[];
  OutOfStockItems: Menu[];
  TotalLowStock: number;
  TotalOutOfStock: number;
}

export interface InventoryAlertResponse {
  status: boolean;
  message: string;
  data: InventoryAlert;
}

export interface SearchIngredientRequest {
  q: string;
}
