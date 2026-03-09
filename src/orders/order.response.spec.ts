import { toOrderResponse } from './order.response';
import { Order } from './order.entity';

const makeOrder = (overrides: Partial<Order> = {}): Order => ({
  id: 'order-uuid',
  items: [
    {
      productId: 'prod-1',
      productName: 'Laptop',
      quantity: 2,
      unitPriceInCents: 89999,
      unitVatInCents: 19800,
      vatRate: 22,
    },
  ],
  totalPriceInCents: 179998,
  totalVatInCents: 39600,
  createdAt: new Date('2024-01-01T00:00:00.000Z'),
  ...overrides,
});

describe('toOrderResponse', () => {
  it('passes the order ID through unchanged', () => {
    expect(toOrderResponse(makeOrder()).id).toBe('order-uuid');
  });

  it('converts totalPriceInCents to euros', () => {
    expect(toOrderResponse(makeOrder()).totalPrice).toBe(1799.98);
  });

  it('converts totalVatInCents to euros', () => {
    expect(toOrderResponse(makeOrder()).totalVat).toBe(396);
  });

  it('serializes createdAt as an ISO string', () => {
    expect(toOrderResponse(makeOrder()).createdAt).toBe(
      '2024-01-01T00:00:00.000Z',
    );
  });

  describe('item mapping', () => {
    it('converts unitPriceInCents to euros', () => {
      const { items } = toOrderResponse(makeOrder());
      expect(items[0].unitPrice).toBe(899.99);
    });

    it('converts unitVatInCents to euros', () => {
      const { items } = toOrderResponse(makeOrder());
      expect(items[0].unitVat).toBe(198);
    });

    it('computes linePrice as unitPrice multiplied by quantity', () => {
      const { items } = toOrderResponse(makeOrder());
      // 89999 cents × 2 = 179998 cents = 1799.98 €
      expect(items[0].linePrice).toBe(1799.98);
    });

    it('computes lineVat as unitVat multiplied by quantity', () => {
      const { items } = toOrderResponse(makeOrder());
      // 19800 cents × 2 = 39600 cents = 396 €
      expect(items[0].lineVat).toBe(396);
    });

    it('passes vatRate through unchanged', () => {
      const { items } = toOrderResponse(makeOrder());
      expect(items[0].vatRate).toBe(22);
    });

    it('passes productId, productName, and quantity through unchanged', () => {
      const { items } = toOrderResponse(makeOrder());
      expect(items[0].productId).toBe('prod-1');
      expect(items[0].productName).toBe('Laptop');
      expect(items[0].quantity).toBe(2);
    });
  });
});
