export type VatRate = 0.04 | 0.1 | 0.22;

export interface Product {
  id: string;
  name: string;
  /** Net price in euro cents (excludes VAT) */
  priceInCents: number;
  /** VAT rate as a decimal (e.g. 0.22 for 22%) */
  vatRate: VatRate;
  category: string;
}
