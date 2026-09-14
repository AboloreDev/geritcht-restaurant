export interface Payment {
  id: number;
  order_id: number;
  reference: string;
  amount: number;
  currency: string;
  status: string;
  provider: string;
  provider_reference: string;
  failure_reason?: string;
  paid_at: string | null;
  created_at: string;
}

export interface PaymentResponse {
  status: boolean;
  message: string;
  error: string;
  data: Payment;
}

export interface PaymentListResponse {
  status: boolean;
  message: string;
  data: Payment[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
  };
}

export interface InitializePaymentRequest {
  order_id: number;
}

export interface InitializePaymentResponse {
  status: boolean;
  message: string;
  data: {
    authorization_url: string;
    payment: Payment;
    reference: string;
  };
}

export interface VerifyPaymentResponse {
  status: boolean;
  message: string;
  data: Payment;
}

export type PaymentStatus =
  | "unpaid"
  | "paid"
  | "failed"
  | "refunded"
  | "pending";

export type PaymentFilterRequest = {
  status?: string;
  reference?: string;
  amount?: number;
  page?: number;
  page_size?: number;
};
