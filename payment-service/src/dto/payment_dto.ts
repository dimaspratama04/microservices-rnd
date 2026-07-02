export interface PaymentRequestDTO {
  invoice_id: string;
  amount: number;
}

export interface PaymentStatusResponseDTO {
  invoice_id: string;
  total: number;
  payment_status: string;
}
