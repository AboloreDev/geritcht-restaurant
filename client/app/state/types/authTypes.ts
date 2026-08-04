import { UserResponse } from "./userTypes";

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone_number?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface MessageResponse {
  status: boolean;
  message: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface VerifyEmailRequest {
  token: string;
}

export interface VerifyResetTokenRequest {
  token: string;
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export interface UpdateProfileRequest {
  first_name?: string;
  last_name?: string;
  phone_number?: string;
}

export interface CreateStaffRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone_number: string;
}

export interface AuthResponse {
  status: boolean;
  message: string;
  data: {
    access_token: string;
    refresh_token: string;
    user: UserResponse;
  };
  error: string;
}

export interface UserSearchResponse extends UserResponse {
  rank: number;
}

export interface ForgotPasswordResponse {
  status: boolean;
  message: string;
  token: string;
}

export interface GetProfileResponse {
  status: boolean;
  message: string;
  data: UserResponse;
}
