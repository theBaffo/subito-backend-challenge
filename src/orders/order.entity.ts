export interface OrderItem {
  productId: string;
  productName: string;
  quantity: number;
  /** Net unit price in euro cents */
  unitPriceInCents: number;
  /** VAT amount per unit in euro cents */
  unitVatInCents: number;
  /** VAT rate as a percentage (e.g. 22 for 22%) */
  vatRate: number;
}

export interface Order {
  id: string;
  items: OrderItem[];
  /** Total net price across all items in euro cents */
  totalPriceInCents: number;
  /** Total VAT across all items in euro cents */
  totalVatInCents: number;
  createdAt: Date;
}
