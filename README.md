# Purchase Cart Service - Subito Backend Challenge

A RESTful purchase cart service built with Go and [Gin](https://github.com/gin-gonic/gin).

Given a catalogue of products, the service allows clients to create orders and retrieve
them by ID. Every response includes a full price breakdown: gross total, VAT total, and
per-item unit price, line total, VAT rate, and VAT amount.

---

## Table of Contents

- [Quick Start](#quick-start)
- [Running Tests](#running-tests)
- [API Reference](#api-reference)
- [Architecture & Design Decisions](#architecture--design-decisions)
- [Project Structure](#project-structure)
- [Potential Evolutions](#potential-evolutions)

---

## Quick Start

**Prerequisites:** Docker

```sh
# Build and run the service (listens on :8080)
sh scripts/run.sh

# Or with docker-compose
docker-compose up
```

Try it immediately:

```sh
# Health check
curl http://localhost:8080/health

# List products
curl http://localhost:8080/v1/products

# Create an order
curl -X POST http://localhost:8080/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      { "product_id": "prod-001", "quantity": 1 },
      { "product_id": "prod-004", "quantity": 3 }
    ]
  }'

# Retrieve an order (replace with the ID returned above)
curl http://localhost:8080/v1/orders/ord-<id>
```

---

## Running Tests

```sh
# Run full test suite inside Docker (with race detector)
sh scripts/tests.sh

# Run locally (requires Go 1.22+)
go test ./... -v -race
```

---

## API Reference

### `GET /health`
Returns service health status.

```json
{ "status": "ok" }
```

---

### `GET /v1/products`
Returns all available products.

```json
{
  "products": [
    {
      "id": "prod-001",
      "name": "Mechanical Keyboard",
      "description": "Full-size mechanical keyboard with Cherry MX switches",
      "gross_price": "129.99",
      "vat_rate": { "name": "standard", "rate": "0.22" },
      "category": "electronics"
    }
  ]
}
```

---

### `GET /v1/products/:id`
Returns a single product by ID.

**Errors:** `404 Not Found` if the product does not exist.

---

### `POST /v1/orders`
Creates a new order from a set of product IDs and quantities.

**Request body:**
```json
{
  "items": [
    { "product_id": "prod-001", "quantity": 2 },
    { "product_id": "prod-004", "quantity": 1 }
  ]
}
```

**Response `201 Created`:**
```json
{
  "id": "ord-7f3a1b2c",
  "status": "confirmed",
  "created_at": "2026-03-05T10:00:00Z",
  "items": [
    {
      "product_id": "prod-001",
      "name": "Mechanical Keyboard",
      "quantity": 2,
      "unit_price": "129.99",
      "total_price": "259.98",
      "vat_rate": "0.22",
      "vat_amount": "47.00"
    },
    {
      "product_id": "prod-004",
      "name": "Artisan Coffee Beans",
      "quantity": 1,
      "unit_price": "14.90",
      "total_price": "14.90",
      "vat_rate": "0.10",
      "vat_amount": "1.35"
    }
  ],
  "total_price": "274.88",
  "total_vat": "48.35"
}
```

**Errors:**
| Status | Condition |
|--------|-----------|
| `400` | Malformed JSON |
| `404` | A product ID does not exist |
| `422` | Empty items list, or quantity ≤ 0 |

---

### `GET /v1/orders/:id`
Retrieves a previously created order by ID.

**Errors:** `404 Not Found` if the order does not exist.

---

## Architecture & Design Decisions

### Layered / Clean Architecture

The codebase is split into four layers, each with a single responsibility:

```
cmd/server        → entry point, wires dependencies
internal/handler  → HTTP layer (Gin), maps requests/responses
internal/service  → business logic, pricing, validation
internal/domain   → core entities and sentinel errors (no external deps)
internal/repository → storage interfaces + in-memory implementations
```

Layers depend inward only. Each handler depends on a *service interface*, and each
service depends on *repository interfaces* — never on concrete types. There is one
service per resource (`ProductService`, `OrderService`), keeping business logic
consistently organised and independently testable. Swapping storage is a matter of
writing a new `repository/postgres/` package that satisfies the same interfaces —
no service or handler code changes required.

### Money: `shopspring/decimal` not `float64`

Financial calculations must never use `float64`. Binary floating-point arithmetic
produces rounding errors that accumulate across line items. `decimal` gives exact
base-10 arithmetic and banker's rounding (`RoundBank`), which is the accounting standard.

### VAT model (EU / Italian convention)

Prices are stored and displayed **gross** (tax-inclusive). The VAT component is
derived at order time:

```
vat_amount = gross - (gross / (1 + vat_rate))
```

Three VAT rates are pre-defined, matching Italian law:
- **22%** — standard (electronics, accessories)
- **10%** — reduced (food, hospitality)
- **4%** — super-reduced (essential goods, bread)

### Storage: in-memory with interface abstraction

The in-memory store is intentionally simple and self-contained — no database
required to run or test the service. The `ProductStore` is seeded with five
realistic products on startup. Both stores are thread-safe (`sync.RWMutex`).

### Error handling

Domain errors (`ErrProductNotFound`, etc.) are defined as sentinel values in the
`domain` package. A single `domainErrorToStatus` function in the handler layer
maps them to HTTP status codes, keeping handlers free of `switch` statements and
making the mapping easy to audit.

### Pricing snapshot

Order items store a snapshot of the product name and price at the time of purchase.
This is intentional: changing a product's price must not retroactively alter
existing orders.

---

## Project Structure

```
.
├── cmd/server/main.go                    # Entry point & dependency wiring
├── internal/
│   ├── domain/
│   │   ├── errors.go                     # Sentinel domain errors
│   │   ├── order.go                      # Order & OrderItem entities
│   │   └── product.go                    # Product entity & VAT rates
│   ├── repository/
│   │   ├── order_repository.go           # OrderRepository interface
│   │   ├── product_repository.go         # ProductRepository interface
│   │   └── memory/
│   │       ├── order_store.go            # In-memory OrderRepository
│   │       └── product_store.go          # In-memory ProductRepository (seeded)
│   ├── service/
│   │   ├── order_service.go              # Order business logic & pricing
│   │   ├── order_service_test.go         # Unit tests (mocked repos)
│   │   ├── product_service.go            # Product business logic
│   │   └── product_service_test.go       # Unit tests (mocked repo)
│   ├── handler/
│   │   ├── helpers.go                    # Error mapping, response shapes
│   │   ├── order_handler.go              # POST /orders, GET /orders/:id
│   │   ├── order_handler_test.go         # Handler tests (mocked service)
│   │   ├── product_handler.go            # GET /products, GET /products/:id
│   │   └── product_handler_test.go       # Handler tests (mocked service)
│   └── middleware/
│       └── logger.go                     # Request logger
├── scripts/
│   ├── run.sh                            # Build & run in Docker
│   └── tests.sh                          # Run tests in Docker
├── Dockerfile                            # Multi-stage build
├── docker-compose.yml
├── go.mod
└── README.md
```

---

## Potential Evolutions

| Area | Description |
|------|-------------|
| **Persistence** | Implement `repository/postgres/` with `pgx`. No service or handler code changes required. |
| **Authentication** | Add JWT middleware; scope orders to authenticated users. |
| **Cart state** | Introduce a `Cart` resource (create → add items → checkout → `Order`). |
| **Stock management** | Track inventory; prevent overselling with optimistic locking in the DB. |
| **Pagination** | Cursor-based pagination on `GET /v1/orders`. |
| **Events** | Publish `order.created` to Kafka or RabbitMQ for downstream consumers (invoicing, fulfilment). |
| **Discounts** | Apply promo codes or per-user discounts at checkout. |
| **Multi-currency** | Accept ISO 4217 currency codes; store exchange rates; display per-currency totals. |
| **Structured logging** | Replace `log` with `slog` (stdlib) or `zap` for JSON-structured logs. |
| **OpenAPI spec** | Generate a `swagger.yaml` with `swaggo/swag` for client SDK generation. |