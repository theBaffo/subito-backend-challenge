## API Details

### Products

| Method | Path | Description |
|---|---|---|
| `GET` | `/products` | List all products |
| `GET` | `/products/:id` | Get a single product |

**Product response**
```json
{
  "id": "prod-1",
  "name": "Laptop",
  "price": 899.99,
  "vatRate": 22,
  "category": "Electronics"
}
```

`price` is the net price in euros (VAT excluded). `vatRate` is expressed as a percentage.

---

### Orders

| Method | Path | Description |
|---|---|---|
| `POST` | `/orders` | Create an order from a cart |
| `GET` | `/orders/:id` | Get a single order |
| `GET` | `/orders` | List all orders |

**Create order — request body**
```json
{
  "items": [
    { "productId": "prod-1", "quantity": 1 },
    { "productId": "prod-5", "quantity": 2 }
  ]
}
```

**Order response**
```json
{
  "id": "3f2a1b...",
  "items": [
    {
      "productId": "prod-1",
      "productName": "Laptop",
      "quantity": 1,
      "unitPrice": 899.99,
      "unitVat": 198.00,
      "linePrice": 899.99,
      "lineVat": 198.00,
      "lineGross": 1097.99,
      "vatRate": 22
    },
    {
      "productId": "prod-5",
      "productName": "Clean Code (Book)",
      "quantity": 2,
      "unitPrice": 34.99,
      "unitVat": 1.40,
      "linePrice": 69.98,
      "lineVat": 2.80,
      "lineGross": 72.78,
      "vatRate": 4
    }
  ],
  "totalPrice": 969.97,
  "totalVat": 200.80,
  "totalGross": 1170.77,
  "createdAt": "2024-01-01T00:00:00.000Z"
}
```

All monetary values in responses are in euros. `linePrice` and `lineVat` are the per-line totals (`unit × quantity`); `lineGross` is their sum. `totalPrice`, `totalVat`, and `totalGross` are the respective sums across all lines.