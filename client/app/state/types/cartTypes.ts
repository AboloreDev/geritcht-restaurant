export interface MenuItem {
  // reuse your existing Menu type here if it's already defined elsewhere
  id: number;
  name: string;
  price: number;
  image_url: string;
  images?: MenuImages;
}

export interface MenuImages {
  alt_text: string;
  created_at: string;
  id: number;
  is_primary: boolean;
  url: string;
}

export interface AddToCartRequest {
  menu_item_id: number;
  quantity: number;
  special_instructions?: string;
}

export interface UpdateCartItemRequest {
  quantity: number;
  special_instructions?: string;
}

export interface CartItem {
  id: number;
  menu_item_id: number;
  menu_item: MenuItem;
  quantity: number;
  special_instructions: string;
  subtotal: number;
  created_at: string;
}

export interface Cart {
  id: number;
  user_id: number;
  cart_items: CartItem[];
  total: number;
  item_count: number;
  created_at: string;
}

export interface CartResponse {
  status: boolean;
  message: string;
  data: Cart;
}
