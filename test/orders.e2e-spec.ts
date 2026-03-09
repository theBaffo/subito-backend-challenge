import { INestApplication } from '@nestjs/common';
import request from 'supertest';
import { App } from 'supertest/types';
import { createTestApp } from './test-app.helper';

describe('Orders (e2e)', () => {
  let app: INestApplication<App>;

  beforeAll(async () => {
    app = await createTestApp();
  });

  afterAll(async () => {
    await app.close();
  });

  describe('POST /orders', () => {
    it('returns 201 with a well-formed order response', async () => {
      const { body } = await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'prod-1', quantity: 1 }] })
        .expect(201);

      expect(body).toMatchObject({
        id: expect.any(String),
        totalPrice: expect.any(Number),
        totalVat: expect.any(Number),
        createdAt: expect.any(String),
        items: expect.any(Array),
      });
    });

    it('returns correct prices and VAT for a single item', async () => {
      // Laptop: net 899.99 €, 22% VAT → VAT = Math.round(89999 * 0.22) / 100 = 198.00 €
      const { body } = await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'prod-1', quantity: 1 }] })
        .expect(201);

      expect(body.totalPrice).toBe(899.99);
      expect(body.totalVat).toBe(198);
      expect(body.items[0]).toMatchObject({
        productId: 'prod-1',
        productName: 'Laptop',
        quantity: 1,
        unitPrice: 899.99,
        unitVat: 198,
        linePrice: 899.99,
        lineVat: 198,
        vatRate: 22,
      });
    });

    it('multiplies line totals correctly when quantity > 1', async () => {
      // Book: net 34.99 €, 4% VAT → VAT = Math.round(3499 * 0.04) / 100 = 1.40 €
      // Quantity 3 → linePrice = 104.97 €, lineVat = 4.20 €
      const { body } = await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'prod-5', quantity: 3 }] })
        .expect(201);

      expect(body.items[0].linePrice).toBe(104.97);
      expect(body.items[0].lineVat).toBe(4.2);
      expect(body.totalPrice).toBe(104.97);
      expect(body.totalVat).toBe(4.2);
    });

    it('sums totals correctly across items with different VAT rates', async () => {
      // Laptop  (22%): net 899.99 €, VAT 198.00 € — ×1
      // Book     (4%): net  34.99 €, VAT   1.40 € — ×2 → net 69.98 €, VAT 2.80 €
      // Total net: 969.97 €  |  Total VAT: 200.80 €
      const { body } = await request(app.getHttpServer())
        .post('/orders')
        .send({
          items: [
            { productId: 'prod-1', quantity: 1 },
            { productId: 'prod-5', quantity: 2 },
          ],
        })
        .expect(201);

      expect(body.totalPrice).toBe(969.97);
      expect(body.totalVat).toBe(200.8);
      expect(body.items).toHaveLength(2);
    });

    it('includes the vatRate on each item in the response', async () => {
      const { body } = await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'prod-4', quantity: 1 }] })
        .expect(201);

      expect(body.items[0].vatRate).toBe(10);
    });

    it('returns 400 when items array is empty', async () => {
      const { body } = await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [] })
        .expect(400);

      expect(body.statusCode).toBe(400);
    });

    it('returns 400 when quantity is less than 1', async () => {
      await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'prod-1', quantity: 0 }] })
        .expect(400);
    });

    it('returns 400 when quantity is not an integer', async () => {
      await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'prod-1', quantity: 1.5 }] })
        .expect(400);
    });

    it('returns 400 when productId is missing', async () => {
      await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ quantity: 1 }] })
        .expect(400);
    });

    it('returns 400 when unknown fields are sent', async () => {
      await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'prod-1', quantity: 1 }], unknownField: 'value' })
        .expect(400);
    });

    it('returns 404 when a product does not exist', async () => {
      const { body } = await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'nonexistent', quantity: 1 }] })
        .expect(404);

      expect(body.message).toContain('nonexistent');
    });
  });

  describe('GET /orders/:id', () => {
    it('returns 200 with the order that was just created', async () => {
      const { body: created } = await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'prod-2', quantity: 1 }] })
        .expect(201);

      const { body: fetched } = await request(app.getHttpServer())
        .get(`/orders/${created.id}`)
        .expect(200);

      expect(fetched.id).toBe(created.id);
      expect(fetched.totalPrice).toBe(created.totalPrice);
    });

    it('returns 404 for an unknown order ID', async () => {
      const { body } = await request(app.getHttpServer())
        .get('/orders/nonexistent')
        .expect(404);

      expect(body.message).toContain('nonexistent');
    });
  });

  describe('GET /orders', () => {
    it('returns 200 with an array', async () => {
      const { body } = await request(app.getHttpServer())
        .get('/orders')
        .expect(200);

      expect(Array.isArray(body)).toBe(true);
    });

    it('includes orders that were previously created', async () => {
      const { body: created } = await request(app.getHttpServer())
        .post('/orders')
        .send({ items: [{ productId: 'prod-3', quantity: 1 }] })
        .expect(201);

      const { body: all } = await request(app.getHttpServer())
        .get('/orders')
        .expect(200);

      expect(all.some((o: { id: string }) => o.id === created.id)).toBe(true);
    });
  });
});
