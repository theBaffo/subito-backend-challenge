import { NotFoundException } from '@nestjs/common';
import { OrdersController } from './orders.controller';
import { OrdersService } from './orders.service';
import { Order } from './order.entity';

const mockOrder: Order = {
  id: 'order-uuid',
  items: [
    {
      productId: 'prod-1',
      productName: 'Laptop',
      quantity: 1,
      unitPriceInCents: 89999,
      unitVatInCents: 19800,
      vatRate: 22,
    },
  ],
  totalPriceInCents: 89999,
  totalVatInCents: 19800,
  createdAt: new Date('2024-01-01T00:00:00.000Z'),
};

const expectedOrderResponse = {
  id: 'order-uuid',
  items: [
    {
      productId: 'prod-1',
      productName: 'Laptop',
      quantity: 1,
      unitPrice: 899.99,
      unitVat: 198,
      linePrice: 899.99,
      lineVat: 198,
      lineGross: 1097.99,
      vatRate: 22,
    },
  ],
  totalPrice: 899.99,
  totalVat: 198,
  totalGross: 1097.99,
  createdAt: '2024-01-01T00:00:00.000Z',
};

describe('OrdersController', () => {
  let controller: OrdersController;
  let service: jest.Mocked<OrdersService>;

  beforeEach(() => {
    service = {
      create: jest.fn(),
      findById: jest.fn(),
      findAll: jest.fn(),
    } as unknown as jest.Mocked<OrdersService>;

    controller = new OrdersController(service);
  });

  describe('create', () => {
    it('returns the created order mapped to the response format', () => {
      service.create.mockReturnValue(mockOrder);
      const dto = { items: [{ productId: 'prod-1', quantity: 1 }] };

      const result = controller.create(dto);

      expect(result).toEqual(expectedOrderResponse);
      // eslint-disable-next-line @typescript-eslint/unbound-method
      expect(service.create).toHaveBeenCalledWith(dto);
    });

    it('propagates NotFoundException when a product does not exist', () => {
      service.create.mockImplementation(() => {
        throw new NotFoundException('Product with id "nonexistent" not found');
      });

      expect(() =>
        controller.create({
          items: [{ productId: 'nonexistent', quantity: 1 }],
        }),
      ).toThrow(NotFoundException);
    });
  });

  describe('findOne', () => {
    it('returns the order mapped to the response format', () => {
      service.findById.mockReturnValue(mockOrder);

      const result = controller.findOne('order-uuid');

      expect(result).toEqual(expectedOrderResponse);
      // eslint-disable-next-line @typescript-eslint/unbound-method
      expect(service.findById).toHaveBeenCalledWith('order-uuid');
    });

    it('propagates NotFoundException from the service', () => {
      service.findById.mockImplementation(() => {
        throw new NotFoundException('Order with id "nonexistent" not found');
      });

      expect(() => controller.findOne('nonexistent')).toThrow(
        NotFoundException,
      );
    });
  });

  describe('findAll', () => {
    it('returns all orders mapped to the response format', () => {
      service.findAll.mockReturnValue([mockOrder]);

      const result = controller.findAll();

      expect(result).toEqual([expectedOrderResponse]);
    });

    it('returns an empty array when there are no orders', () => {
      service.findAll.mockReturnValue([]);

      expect(controller.findAll()).toEqual([]);
    });
  });
});
