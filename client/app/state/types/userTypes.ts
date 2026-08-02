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
  //   page?: number;
  //   limit?: number;
  //   email?: string;
  //   first_name?: string;
  //   last_name?: string;
}
