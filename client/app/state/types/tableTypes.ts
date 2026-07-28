import { ReservationResponse } from "./reservationTypes";

export interface TableDetailResponse {
  id: number;
  name: string;
  capacity: number;
  location: string;
  status: string;
  qrCode?: string;
  currentReservation?: ReservationResponse | null;
  //   currentOrder?: OrderResponse | null;
}
