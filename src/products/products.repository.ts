import { Injectable } from '@nestjs/common';
import { Product } from './product.entity';

const SEED_PRODUCTS: Product[] = [
  {
    id: 'prod-1',
    name: 'Laptop',
    priceInCents: 89999,
    vatRate: 0.22,
    category: 'Electronics',
  },
  {
    id: 'prod-2',
    name: 'Wireless Mouse',
    priceInCents: 2499,
    vatRate: 0.22,
    category: 'Electronics',
  },
  {
    id: 'prod-3',
    name: 'Standing Desk',
    priceInCents: 34999,
    vatRate: 0.22,
    category: 'Furniture',
  },
  {
    id: 'prod-4',
    name: 'Espresso Coffee Beans 1kg',
    priceInCents: 1299,
    vatRate: 0.1,
    category: 'Food & Beverages',
  },
  {
    id: 'prod-5',
    name: 'Clean Code (Book)',
    priceInCents: 3499,
    vatRate: 0.04,
    category: 'Books',
  },
  {
    id: 'prod-6',
    name: 'Mechanical Keyboard',
    priceInCents: 14999,
    vatRate: 0.22,
    category: 'Electronics',
  },
];

@Injectable()
export class ProductsRepository {
  private readonly products = new Map<string, Product>(
    SEED_PRODUCTS.map((p) => [p.id, p]),
  );

  findAll(): Product[] {
    return Array.from(this.products.values());
  }

  findById(id: string): Product | undefined {
    return this.products.get(id);
  }
}
