# Subito Backend Challenge - Purchase Cart Service

A RESTful service that manages a product catalogue and allows creating orders from a cart. Built with [NestJS](https://nestjs.com/) and TypeScript.

See [CODING_CHALLENGE.md](CODING_CHALLENGE.md) for more details of the coding challenge.

---

## Running the service

### Locally

```bash
npm install
npm run start
```

The service listens on port `3000` by default. The port can be overridden with the `PORT` environment variable.

### With Docker

```bash
bash scripts/run.sh
```

This builds the Docker image and runs the container, exposing the service on port `3000`.

---

## Running the tests

### Locally

```bash
# Unit tests
npm test

# E2E tests
npm run test:e2e

# Coverage report
npm run test:cov
```

### With Docker

```bash
bash scripts/tests.sh
```

This builds the image up to the `builder` stage and executes the unit test suite inside the container.

---

## Testing with Postman

A Postman collection and environment file are included in the [`postman/`](postman/) directory:

To use them, open Postman, go to **File → Import**, and import both files. Select the **Backend Challenge — Local** environment before sending requests.

The "Create order" requests include a test script that automatically saves the returned order `id` to the `orderId` collection variable, so the "Get order by ID" request works without manual copy-paste.

---

## API Details

See [API_DETAILS.md](API_DETAILS.md) for more details of the API structure.

## Product catalogue

The service is seeded with six products that intentionally span all three Italian VAT tiers:

| ID | Name | Net price | VAT rate |
|---|---|---|---|
| `prod-1` | Laptop | €899.99 | 22% |
| `prod-2` | Wireless Mouse | €24.99 | 22% |
| `prod-3` | Standing Desk | €349.99 | 22% |
| `prod-4` | Espresso Coffee Beans 1kg | €12.99 | 10% |
| `prod-5` | Clean Code (Book) | €34.99 | 4% |
| `prod-6` | Mechanical Keyboard | €149.99 | 22% |

## Design decisions

### Pricing model

Prices are stored internally as **integer cents** to avoid floating-point arithmetic errors. All calculations (VAT, line totals, order totals) are performed in cents and converted to euros only when building the API response.

VAT is calculated per unit as `Math.round(priceInCents × vatRate)` and then multiplied by the quantity. This means rounding happens once per product line, not once per unit, which is the standard approach.

Product prices represent the **net amount** (VAT excluded). The VAT is computed on top and reported separately, giving the consumer full transparency into the tax breakdown.

### VAT rates

Three rates are supported, reflecting the Italian tax system:

- **22%** — standard rate
- **10%** — reduced rate
- **4%** — super-reduced

The `VatRate` type in [product.entity.ts](src/products/product.entity.ts) is a union of these three literals, making it impossible to seed or introduce a product with an unsupported rate.

### Storage

Both repositories (`ProductsRepository`, `OrdersRepository`) use an in-memory `Map`. This was a deliberate choice for the scope of this challenge: it keeps the service self-contained, requires no infrastructure, and makes tests fast and deterministic.

The repositories are accessed only through their service layer, so swapping in a database-backed implementation would be limited to replacing the repository class — no other code would need to change.

### Request validation

All incoming request bodies are validated by NestJS's global `ValidationPipe`, configured with:

- `whitelist: true` — strips any properties not declared in the DTO
- `forbidNonWhitelisted: true` — returns `400` if unknown properties are sent
- `transform: true` — coerces primitive types (e.g. ensures `quantity` is a number)

### Potential evolutions

- **Persistence** — replace the in-memory repositories with TypeORM or Prisma-backed implementations; the service layer requires no changes.
- **Product management** — add `POST /products`, `PATCH /products/:id`, and `DELETE /products/:id` endpoints to manage the catalogue at runtime.
- **Order status** — extend the `Order` entity with a `status` field (`pending`, `confirmed`, `cancelled`) and expose a `PATCH /orders/:id/status` endpoint.
- **Authentication** — add a JWT guard to protect write endpoints.
- **Pagination** — add cursor- or page-based pagination to `GET /products` and `GET /orders`.
