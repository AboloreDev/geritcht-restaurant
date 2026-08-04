import { Pagination } from "./categoriesTypes";

export interface CreateAllergenRequest {
  name: string;
}

export interface UpdateAllergenRequest {
  name: string;
}

export interface Allergen {
  id: number;
  name: string;
}

export interface AllergenResponse {
  status: boolean;
  message: string;
  data: Allergen;
}

export interface GetAllergensResponse {
  status: boolean;
  message: string;
  data: Allergen[];
  meta: Pagination;
}
