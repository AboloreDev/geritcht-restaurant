import type { Menu } from "@/app/state/types/menuTypes";
import { UserResponse } from "./authTypes";
import { PaymentResponse } from "./paymentTypes";

// --- Requests ---

export interface CreateTakeoutOrderRequest {
  notes?: string;
}

export interface CreateOrderItemRequest {
  menu_item_id: number;
  quantity: number;
  special_instructions?: string;
}

export interface CreateDineInOrderRequest {
  table_id: number;
  reservation_id?: number;
  items: CreateOrderItemRequest[];
  notes?: string;
}

export interface UpdateOrderStatusRequest {
  status: string;
}

export interface OrderFilterRequest {
  status?: string;
  type?: string;
  date?: string;
  page?: number;
  page_size?: number;
}

// --- Responses ---

export interface OrderItem {
  id: number;
  menu_item_id: number;
  menu_item: Menu;
  quantity: number;
  price: number;
  subtotal: number;
  special_instructions: string;
}

// placeholder — I don't have your actual UserResponse/PaymentResponse
// DTOs, so these are minimal guesses. Replace with your real shapes.

export interface Order {
  id: number;
  user_id: number | null;
  user?: UserResponse;
  table_id: number | null;
  reservation_id: number | null;
  type: string; // e.g. "takeout" | "dine_in"
  status: string;
  total_amount: number;
  payment_status: string;
  notes: string;
  order_items: OrderItem[];
  payment?: PaymentResponse;
  created_at: string;
  updated_at: string;
}

export interface OrderResponse {
  status: boolean;
  message: string;
  data: Order;
  error: string;
}

export interface OrderListResponse {
  status: boolean;
  message: string;
  data: {
    orders: Order[];
    total: number;
    page: number;
    page_size: number;
  };
}
