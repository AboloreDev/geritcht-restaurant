import { Pagination } from "./categoriesTypes";
import { Order } from "./orderTypes";
import { ReservationResponse } from "./reservationTypes";

export interface CreateTableRequest {
  name: string;
  capacity: number;
  location?: string;
}

export interface UpdateTableRequest {
  name?: string;
  capacity?: number;
  location?: string;
  status?: string;
}

export enum TableStatus {
  Available = "available",
  Occupied = "occupied",
  Reserved = "reserved",
}

export interface UpdateTableStatusRequest {
  status: TableStatus;
}

export interface Table {
  id: number;
  name: string;
  capacity: number;
  location: string;
  status: TableStatus;
  qr_code?: string;
  current_reservation?: ReservationResponse;
  current_order?: Order;
}

export interface GetTablesResponse {
  status: boolean;
  message: string;
  data: Table[];
  meta: Pagination;
}

export interface GetTableResponse {
  status: boolean;
  message: string;
  data: Table;
}
