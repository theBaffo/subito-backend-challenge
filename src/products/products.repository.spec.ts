import { ProductsRepository } from './products.repository';

describe('ProductsRepository', () => {
  let repository: ProductsRepository;

  beforeEach(() => {
    repository = new ProductsRepository();
  });

  describe('findAll', () => {
    it('returns all seeded products', () => {
      const products = repository.findAll();
      expect(products).toHaveLength(6);
    });

    it('returns products with the expected IDs', () => {
      const ids = repository.findAll().map((p) => p.id);
      expect(ids).toEqual(
        expect.arrayContaining([
          'prod-1',
          'prod-2',
          'prod-3',
          'prod-4',
          'prod-5',
          'prod-6',
        ]),
      );
    });
  });

  describe('findById', () => {
    it('returns the correct product for a known ID', () => {
      const product = repository.findById('prod-1');
      expect(product).toMatchObject({
        id: 'prod-1',
        name: 'Laptop',
        priceInCents: 89999,
        vatRate: 0.22,
        category: 'Electronics',
      });
    });

    it('returns undefined for an unknown ID', () => {
      expect(repository.findById('nonexistent')).toBeUndefined();
    });
  });
});
