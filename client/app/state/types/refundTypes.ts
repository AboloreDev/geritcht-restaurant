export interface ProcessRefundRequest {
  notes?: string;
}

export interface Refund {
  id: number;
  order_id: number;
  payment_id: number;
  amount: number;
  reason: string;
  reference: string;
  currency: string;
  status: string;
  processed_at: string | null;
}

export interface RefundResponse {
  status: boolean;
  message: string;
  data: Refund;
}
