export interface Menu {
  id: number;
  name: string;
  description: string;

  price: number;
  prep_time_minutes: number;
  spice_level: number;

  image_url: string;
  is_available: boolean;
  display_order: number;

  category_id: number;
  category: MenuCategory;

  images: MenuImage[];
  allergens: MenuAllergen[];
  dietary_tags: MenuDietaryTag[];

  created_at: string;
  updated_at: string;
}

export interface MenuCategory {
  id: number;
  name: string;
  description: string;
  image_url: string;

  display_order: number;
  is_active: boolean;

  created_at: string;
}

export interface MenuImage {
  id: number;
  url: string;
  alt_text: string;
  is_primary: boolean;
  created_at: string;
}

export interface MenuAllergen {
  id: number;
  name: string;
}

export interface MenuDietaryTag {
  id: number;
  name: string;
}

export interface Meta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface GetMenusResponse {
  status: boolean;
  error: string;
  message: string;
  data: Menu[];
  meta: Meta;
}

export interface GetMenusRequest {
  page?: number;
  limit?: number;

  category_id?: number;
  query?: string;

  min_price?: number;
  max_price?: number;

  prep_time_minutes?: number;
  spice_level?: number;

  sort_by?: "name" | "price" | "created_at";
  sort_order?: "asc" | "desc";
}

export interface SearchMenuRequest {
  q: string;
}
