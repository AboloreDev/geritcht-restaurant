import { int, number } from "zod";
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

export interface CreateIngredientRequest {
  name: string;
  unit: string;
  current_stock?: number;
  min_threshold: number;
}

export interface UpdateIngredientRequest {
  name?: string;
  unit?: string;
  min_threshold?: number;
}

export interface UpdateStockRequest {
  quantity: number;
  reason: string;
  type: "in" | "out" | "waste";
}

export interface LinkIngredientRequest {
  ingredient_id: number;
  quantity: number;
}

export interface ThresholdRequest {
  min_threshold: number;
}

export interface StockMovement {
  id: number;
  ingredient_id: number;
  ingredient: Ingredient;
  type: string;
  quantity: number;
  reason: string;
  created_by: number;
  user: {
    id: number;
    email: string;
    first_name: string;
    last_name: string;
    phone_number: string;
    role: string;
    is_active: boolean;
    email_verified: boolean;
    created_at: string;
  };
  created_at: string;
}

export interface MenuItemIngredient {
  ingredient_id: number;
  ingredient: Ingredient;
  quantity: number;
}

export interface InventoryAlertMenuItem {
  id: number;
  name: string;
  is_available: boolean;
  images: { id: number; url: string; alt_text: string }[];
}

export interface InventoryAlertResponse {
  low_stock_ingredients: Ingredient[];
  out_of_stock_items: InventoryAlertMenuItem[];
  total_low_stock: number;
  total_out_of_stock: number;
}

export interface IngredientSearchRequest {
  q: string;
  page?: number;
  limit?: number;
  min_threshold?: number;
  current_stock?: number;
}

export interface IngredientSearchResult extends Ingredient {
  rank: number;
}

export interface IngredientResponse {
  status: boolean;
  message: string;
  data: Ingredient;
}

export interface IngredientListResponse {
  status: boolean;
  message: string;
  data: Ingredient[];
  meta: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}

export interface IngredientSearchResponse {
  status: boolean;
  message: string;
  data: IngredientSearchResult[];
}

export interface GetInventoryAlertsResponse {
  status: boolean;
  message: string;
  data: InventoryAlertResponse;
}

export interface RecipeIngredient extends MenuItemIngredient {}

export interface RecipeListResponse {
  status: boolean;
  message: string;
  data: MenuItemIngredient[];
}

export interface RecipeResponse {
  status: boolean;
  message: string;
  data: MenuItemIngredient;
}

export interface UpdateLinkItemRequest {
  quantity: number;
}
