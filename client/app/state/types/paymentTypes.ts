export interface PaymentResponse {
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

export interface InitializePaymentRequest {
  order_id: number;
}

export interface InitializePaymentResponse {
  status: boolean;
  message: string;
  data: {
    authorization_url: string;
    payment: PaymentResponse;
    reference: string;
  };
}

export interface VerifyPaymentResponse {
  status: boolean;
  message: string;
  data: PaymentResponse;
}
