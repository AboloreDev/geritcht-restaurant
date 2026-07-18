export interface GetCategoriesRequest {
  page: number;
  limit: number;
  query?: string;
}

export interface Category {
  id: number;
  name: string;
  description: string;
  image_url: string;
  display_order: number;
  is_active: boolean;
  created_at: string;
}

export interface Pagination {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface GetCategoriesResponse {
  data: Category[];
  pagination: Pagination;
}

export interface SearchCategoriesRequest {
  query: string;
}
