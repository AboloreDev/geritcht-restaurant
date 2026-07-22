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
