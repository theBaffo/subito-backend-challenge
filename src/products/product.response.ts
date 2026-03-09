import { Product } from './product.entity';

export interface ProductResponse {
  id: string;
  name: string;
  /** Net price in euros */
  price: number;
  /** VAT rate as a percentage (e.g. 22 for 22%) */
  vatRate: number;
  category: string;
}

export function toProductResponse(product: Product): ProductResponse {
  return {
    id: product.id,
    name: product.name,
    price: product.priceInCents / 100,
    vatRate: product.vatRate * 100,
    category: product.category,
  };
}
