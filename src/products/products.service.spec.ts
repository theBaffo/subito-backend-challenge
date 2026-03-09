import { NotFoundException } from '@nestjs/common';
import { ProductsService } from './products.service';
import { ProductsRepository } from './products.repository';
import { Product } from './product.entity';

const mockProduct: Product = {
  id: 'prod-1',
  name: 'Laptop',
  priceInCents: 89999,
  vatRate: 0.22,
  category: 'Electronics',
};

describe('ProductsService', () => {
  let service: ProductsService;
  let repository: jest.Mocked<ProductsRepository>;

  beforeEach(() => {
    repository = {
      findAll: jest.fn(),
      findById: jest.fn(),
    } as unknown as jest.Mocked<ProductsRepository>;

    service = new ProductsService(repository);
  });

  describe('findAll', () => {
    it('returns all products from the repository', () => {
      repository.findAll.mockReturnValue([mockProduct]);

      const result = service.findAll();

      expect(result).toEqual([mockProduct]);
      // eslint-disable-next-line @typescript-eslint/unbound-method
      expect(repository.findAll).toHaveBeenCalledTimes(1);
    });
  });

  describe('findById', () => {
    it('returns the product when found', () => {
      repository.findById.mockReturnValue(mockProduct);

      const result = service.findById('prod-1');

      expect(result).toEqual(mockProduct);
      // eslint-disable-next-line @typescript-eslint/unbound-method
      expect(repository.findById).toHaveBeenCalledWith('prod-1');
    });

    it('throws NotFoundException when the product does not exist', () => {
      repository.findById.mockReturnValue(undefined);

      expect(() => service.findById('nonexistent')).toThrow(
        new NotFoundException('Product with id "nonexistent" not found'),
      );
    });
  });
});
