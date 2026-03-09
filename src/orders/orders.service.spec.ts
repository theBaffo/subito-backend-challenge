import { NotFoundException } from '@nestjs/common';
import { OrdersService } from './orders.service';
import { OrdersRepository } from './orders.repository';
import { ProductsService } from '../products/products.service';
import { Product } from '../products/product.entity';
import { Order } from './order.entity';

const laptop: Product = {
  id: 'prod-1',
  name: 'Laptop',
  priceInCents: 89999,
  vatRate: 0.22,
  category: 'Electronics',
};

const book: Product = {
  id: 'prod-5',
  name: 'Clean Code (Book)',
  priceInCents: 3499,
  vatRate: 0.04,
  category: 'Books',
};

const coffee: Product = {
  id: 'prod-4',
  name: 'Espresso Coffee Beans 1kg',
  priceInCents: 1299,
  vatRate: 0.1,
  category: 'Food & Beverages',
};

const mockOrder: Order = {
  id: 'order-uuid',
  items: [],
  totalPriceInCents: 0,
  totalVatInCents: 0,
  createdAt: new Date(),
};

describe('OrdersService', () => {
  let service: OrdersService;
  let repository: jest.Mocked<OrdersRepository>;
  let productsService: jest.Mocked<ProductsService>;

  beforeEach(() => {
    repository = {
      save: jest.fn(),
      findById: jest.fn(),
      findAll: jest.fn(),
    } as unknown as jest.Mocked<OrdersRepository>;

    productsService = {
      findById: jest.fn(),
    } as unknown as jest.Mocked<ProductsService>;

    // By default, save returns whatever is passed in
    repository.save.mockImplementation((order) => order);

    service = new OrdersService(repository, productsService);
  });

  describe('create', () => {
    it('maps product data to order items correctly', () => {
      productsService.findById.mockReturnValue(laptop);

      const order = service.create({
        items: [{ productId: 'prod-1', quantity: 2 }],
      });

      expect(order.items).toEqual([
        {
          productId: 'prod-1',
          productName: 'Laptop',
          quantity: 2,
          unitPriceInCents: 89999,
          unitVatInCents: 19800, // Math.round(89999 * 0.22)
          vatRate: 22,
        },
      ]);
    });

    it('rounds VAT per unit using Math.round', () => {
      // 3499 * 0.04 = 139.96 → should round to 140
      productsService.findById.mockReturnValue(book);

      const order = service.create({
        items: [{ productId: 'prod-5', quantity: 1 }],
      });

      expect(order.items[0].unitVatInCents).toBe(140);
    });

    it('rounds VAT correctly for 10% rate', () => {
      // 1299 * 0.10 = 129.9 → should round to 130
      productsService.findById.mockReturnValue(coffee);

      const order = service.create({
        items: [{ productId: 'prod-4', quantity: 1 }],
      });

      expect(order.items[0].unitVatInCents).toBe(130);
    });

    it('computes totalPriceInCents as the sum of (unitPrice × quantity) across all items', () => {
      productsService.findById
        .mockReturnValueOnce(laptop) // 89999 × 1
        .mockReturnValueOnce(book); // 3499  × 2

      const order = service.create({
        items: [
          { productId: 'prod-1', quantity: 1 },
          { productId: 'prod-5', quantity: 2 },
        ],
      });

      expect(order.totalPriceInCents).toBe(89999 + 3499 * 2); // 96997
    });

    it('computes totalVatInCents as the sum of (unitVat × quantity) across all items', () => {
      productsService.findById
        .mockReturnValueOnce(laptop) // unitVat: 19800, × 1
        .mockReturnValueOnce(book); // unitVat: 140,   × 2

      const order = service.create({
        items: [
          { productId: 'prod-1', quantity: 1 },
          { productId: 'prod-5', quantity: 2 },
        ],
      });

      expect(order.totalVatInCents).toBe(19800 + 140 * 2); // 20080
    });

    it('assigns a non-empty string ID to the order', () => {
      productsService.findById.mockReturnValue(laptop);

      const order = service.create({
        items: [{ productId: 'prod-1', quantity: 1 }],
      });

      expect(typeof order.id).toBe('string');
      expect(order.id).not.toBe('');
    });

    it('assigns unique IDs to each created order', () => {
      productsService.findById.mockReturnValue(laptop);

      const dto = { items: [{ productId: 'prod-1', quantity: 1 }] };
      const first = service.create(dto);
      const second = service.create(dto);

      expect(first.id).not.toBe(second.id);
    });

    it('persists the order via the repository', () => {
      productsService.findById.mockReturnValue(laptop);

      service.create({ items: [{ productId: 'prod-1', quantity: 1 }] });

      expect(repository.save).toHaveBeenCalledTimes(1);
    });

    it('propagates NotFoundException when a product does not exist', () => {
      productsService.findById.mockImplementation(() => {
        throw new NotFoundException('Product with id "nonexistent" not found');
      });

      expect(() =>
        service.create({ items: [{ productId: 'nonexistent', quantity: 1 }] }),
      ).toThrow(NotFoundException);
    });
  });

  describe('findById', () => {
    it('returns the order when found', () => {
      repository.findById.mockReturnValue(mockOrder);

      const result = service.findById('order-uuid');

      expect(result).toEqual(mockOrder);
      expect(repository.findById).toHaveBeenCalledWith('order-uuid');
    });

    it('throws NotFoundException when the order does not exist', () => {
      repository.findById.mockReturnValue(undefined);

      expect(() => service.findById('nonexistent')).toThrow(
        new NotFoundException('Order with id "nonexistent" not found'),
      );
    });
  });

  describe('findAll', () => {
    it('returns all orders from the repository', () => {
      repository.findAll.mockReturnValue([mockOrder]);

      const result = service.findAll();

      expect(result).toEqual([mockOrder]);
      expect(repository.findAll).toHaveBeenCalledTimes(1);
    });
  });
});
