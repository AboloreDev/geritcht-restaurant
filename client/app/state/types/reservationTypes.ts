import { UserResponse } from "./authTypes";

export interface CheckAvailabilityRequest {
  date: string; // "YYYY-MM-DD"
  time_slot: string; // e.g. "18:00:00"
  party_size: number;
}

export interface TableAvailability {
  id: number;
  name: string;
  capacity: number;
  location: string;
  status: string;
}

export interface AvailabilityResponse {
  status: boolean;
  message: string;
  data: {
    date: string;
    time_slot: string;
    party_size: number;
    tables: TableAvailability[];
  };
}

export interface CreateReservationRequest {
  table_id: number;
  date: string;
  time_slot: string;
  party_size: number;
  special_requests?: string;
}

export interface Meta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface ReservationListResponse {
  status: boolean;
  error: string;
  message: string;
  reservations: ReservationResponse[];
  meta: ReservationMeta;
}

export interface ReservationMeta {
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ReservationResponse {
  id: number;
  user_id: number;
  user?: UserResponse;
  table_id: number;
  // table: TableResponse;
  date: string;
  time_slot: string;
  party_size: number;
  status: string;
  special_requests: string;
  checked_in_at: string | null;
  created_at: string;
}

export interface GetReservationsRequest {
  date?: string;
  status?: string;
  page?: number;
  pageSize?: number;
}
