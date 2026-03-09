import { INestApplication } from '@nestjs/common';
import request from 'supertest';
import { App } from 'supertest/types';
import { createTestApp } from './test-app.helper';

describe('Products (e2e)', () => {
  let app: INestApplication<App>;

  beforeAll(async () => {
    app = await createTestApp();
  });

  afterAll(async () => {
    await app.close();
  });

  describe('GET /products', () => {
    it('returns 200 with the full product catalogue', async () => {
      const { body } = await request(app.getHttpServer())
        .get('/products')
        .expect(200);

      expect(Array.isArray(body)).toBe(true);
      expect(body).toHaveLength(6);
    });

    it('returns products with the expected shape', async () => {
      const { body } = await request(app.getHttpServer())
        .get('/products')
        .expect(200);

      expect(body[0]).toMatchObject({
        id: expect.any(String),
        name: expect.any(String),
        price: expect.any(Number),
        vatRate: expect.any(Number),
        category: expect.any(String),
      });
    });

    it('returns prices in euros, not cents', async () => {
      const { body } = await request(app.getHttpServer())
        .get('/products')
        .expect(200);

      const laptop = body.find((p: { id: string }) => p.id === 'prod-1');
      expect(laptop.price).toBe(899.99);
    });

    it('returns vatRate as a percentage', async () => {
      const { body } = await request(app.getHttpServer())
        .get('/products')
        .expect(200);

      const laptop = body.find((p: { id: string }) => p.id === 'prod-1');
      expect(laptop.vatRate).toBe(22);
    });
  });

  describe('GET /products/:id', () => {
    it('returns 200 with the correct product', async () => {
      const { body } = await request(app.getHttpServer())
        .get('/products/prod-5')
        .expect(200);

      expect(body).toMatchObject({
        id: 'prod-5',
        name: 'Clean Code (Book)',
        price: 34.99,
        vatRate: 4,
        category: 'Books',
      });
    });

    it('returns 404 for an unknown product ID', async () => {
      const { body } = await request(app.getHttpServer())
        .get('/products/nonexistent')
        .expect(404);

      expect(body.message).toContain('nonexistent');
    });
  });
});
