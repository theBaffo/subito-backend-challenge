import { toProductResponse } from './product.response';
import { Product } from './product.entity';

describe('toProductResponse', () => {
  it('converts priceInCents to euros', () => {
    const product: Product = {
      id: 'prod-1',
      name: 'Laptop',
      priceInCents: 89999,
      vatRate: 0.22,
      category: 'Electronics',
    };

    expect(toProductResponse(product).price).toBe(899.99);
  });

  it('converts the vatRate decimal to a percentage', () => {
    const product: Product = {
      id: 'prod-5',
      name: 'Clean Code (Book)',
      priceInCents: 3499,
      vatRate: 0.04,
      category: 'Books',
    };

    expect(toProductResponse(product).vatRate).toBe(4);
  });

  it('passes all other fields through unchanged', () => {
    const product: Product = {
      id: 'prod-4',
      name: 'Espresso Coffee Beans 1kg',
      priceInCents: 1299,
      vatRate: 0.1,
      category: 'Food & Beverages',
    };

    const response = toProductResponse(product);

    expect(response.id).toBe(product.id);
    expect(response.name).toBe(product.name);
    expect(response.category).toBe(product.category);
  });
});
