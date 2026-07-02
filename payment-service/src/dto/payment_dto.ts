export interface PaymentRequestDTO {
    order_id: number;
    amount: number;
}

export interface PaymentStatusResponseDTO {
    order_id: number;
    product_name: string;
    total: number;
}
