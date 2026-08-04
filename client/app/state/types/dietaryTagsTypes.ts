import { Pagination } from "./categoriesTypes";

export interface CreateDietaryTagRequest {
  name: string;
}

export interface UpdateDietaryTagRequest {
  name: string;
}

export interface dietaryTagsResponse {
  status: boolean;
  message: string;
  data: DietaryTag;
}

export interface DietaryTag {
  id: number;
  name: string;
}

export interface GetDietaryTagsResponse {
  data: DietaryTag[];
  meta: Pagination;
}
