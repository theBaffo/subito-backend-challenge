import { OrdersRepository } from './orders.repository';
import { Order } from './order.entity';

const makeOrder = (id: string): Order => ({
  id,
  items: [],
  totalPriceInCents: 0,
  totalVatInCents: 0,
  createdAt: new Date(),
});

describe('OrdersRepository', () => {
  let repository: OrdersRepository;

  beforeEach(() => {
    repository = new OrdersRepository();
  });

  describe('save', () => {
    it('persists the order and returns it', () => {
      const order = makeOrder('order-1');
      const saved = repository.save(order);
      expect(saved).toEqual(order);
    });

    it('makes the order retrievable after saving', () => {
      const order = makeOrder('order-1');
      repository.save(order);
      expect(repository.findById('order-1')).toEqual(order);
    });

    it('overwrites an existing order with the same ID', () => {
      const original = makeOrder('order-1');
      const updated = { ...original, totalPriceInCents: 5000 };
      repository.save(original);
      repository.save(updated);
      expect(repository.findById('order-1')?.totalPriceInCents).toBe(5000);
    });
  });

  describe('findById', () => {
    it('returns undefined for an unknown ID', () => {
      expect(repository.findById('nonexistent')).toBeUndefined();
    });
  });

  describe('findAll', () => {
    it('returns an empty array when no orders exist', () => {
      expect(repository.findAll()).toEqual([]);
    });

    it('returns all saved orders', () => {
      repository.save(makeOrder('order-1'));
      repository.save(makeOrder('order-2'));
      const ids = repository.findAll().map((o) => o.id);
      expect(ids).toEqual(expect.arrayContaining(['order-1', 'order-2']));
      expect(ids).toHaveLength(2);
    });
  });
});
