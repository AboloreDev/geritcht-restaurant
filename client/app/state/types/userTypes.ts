export interface UserResponse {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  phone_number: string;
  role: string;
  is_active: boolean;
  email_verified: boolean;
  created_at: string;
}

export interface User {
  status: boolean;
  message: string;
  data: UserResponse[];
  error: string;
}

export interface MessageResponse {
  status: boolean;
  message: string;
}

export interface UserSearchRequest {
  q: string;
}

export interface GetAllUsersResponse {
  status: boolean;
  message: string;
  data: UserResponse[];
  meta: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}

export interface GetAllStaffsResponse {
  status: boolean;
  message: string;
  data: UserResponse[];
  meta: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}
