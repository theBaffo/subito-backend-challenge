import { NotFoundException } from '@nestjs/common';
import { ProductsController } from './products.controller';
import { ProductsService } from './products.service';
import { Product } from './product.entity';

const mockProduct: Product = {
  id: 'prod-1',
  name: 'Laptop',
  priceInCents: 89999,
  vatRate: 0.22,
  category: 'Electronics',
};

describe('ProductsController', () => {
  let controller: ProductsController;
  let service: jest.Mocked<ProductsService>;

  beforeEach(() => {
    service = {
      findAll: jest.fn(),
      findById: jest.fn(),
    } as unknown as jest.Mocked<ProductsService>;

    controller = new ProductsController(service);
  });

  describe('findAll', () => {
    it('returns all products mapped to the response format', () => {
      service.findAll.mockReturnValue([mockProduct]);

      const result = controller.findAll();

      expect(result).toEqual([
        {
          id: 'prod-1',
          name: 'Laptop',
          price: 899.99,
          vatRate: 22,
          category: 'Electronics',
        },
      ]);
    });
  });

  describe('findOne', () => {
    it('returns a single product mapped to the response format', () => {
      service.findById.mockReturnValue(mockProduct);

      const result = controller.findOne('prod-1');

      expect(result).toEqual({
        id: 'prod-1',
        name: 'Laptop',
        price: 899.99,
        vatRate: 22,
        category: 'Electronics',
      });
      // eslint-disable-next-line @typescript-eslint/unbound-method
      expect(service.findById).toHaveBeenCalledWith('prod-1');
    });

    it('propagates NotFoundException from the service', () => {
      service.findById.mockImplementation(() => {
        throw new NotFoundException('Product with id "nonexistent" not found');
      });

      expect(() => controller.findOne('nonexistent')).toThrow(
        NotFoundException,
      );
    });
  });
});
