import { Order, OrderItem } from './order.entity';

export interface OrderItemResponse {
  productId: string;
  productName: string;
  quantity: number;
  /** Net unit price in euros */
  unitPrice: number;
  /** VAT amount per unit in euros */
  unitVat: number;
  /** Total net price for this line (unitPrice × quantity) */
  linePrice: number;
  /** Total VAT for this line (unitVat × quantity) */
  lineVat: number;
  /** Total gross for this line (linePrice + lineVat) */
  lineGross: number;
  /** VAT rate as a percentage (e.g. 22 for 22%) */
  vatRate: number;
}

export interface OrderResponse {
  id: string;
  items: OrderItemResponse[];
  /** Sum of all line prices in euros */
  totalPrice: number;
  /** Sum of all line VATs in euros */
  totalVat: number;
  /** Total gross amount (totalPrice + totalVat) */
  totalGross: number;
  createdAt: string;
}

function toOrderItemResponse(item: OrderItem): OrderItemResponse {
  return {
    productId: item.productId,
    productName: item.productName,
    quantity: item.quantity,
    unitPrice: item.unitPriceInCents / 100,
    unitVat: item.unitVatInCents / 100,
    linePrice: (item.unitPriceInCents * item.quantity) / 100,
    lineVat: (item.unitVatInCents * item.quantity) / 100,
    lineGross: ((item.unitPriceInCents + item.unitVatInCents) * item.quantity) / 100,
    vatRate: item.vatRate,
  };
}

export function toOrderResponse(order: Order): OrderResponse {
  return {
    id: order.id,
    items: order.items.map(toOrderItemResponse),
    totalPrice: order.totalPriceInCents / 100,
    totalVat: order.totalVatInCents / 100,
    totalGross: (order.totalPriceInCents + order.totalVatInCents) / 100,
    createdAt: order.createdAt.toISOString(),
  };
}
