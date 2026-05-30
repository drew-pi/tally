# Tally — API Specification

## Project Context

Tally is a local budget tracking app built in Go with PostgreSQL. This document describes all API endpoints needed to support transaction and income management. The app uses `chi` as the HTTP router. All requests and responses use JSON. All amounts are stored as positive `NUMERIC(10,2)`.

---

## Base URL

```
http://localhost:8080
```

## Common Error Shape

All errors follow this structure:

```json
{ "error": "human readable message" }
```

| Status                      | Meaning                                                                         |
| --------------------------- | ------------------------------------------------------------------------------- |
| `400 Bad Request`           | Missing required field, invalid value, or attempt to edit a read-only field     |
| `404 Not Found`             | Resource with given ID does not exist                                           |
| `409 Conflict`              | Action violates a database constraint e.g. deleting a referenced payment method |
| `500 Internal Server Error` | Unexpected server or database error                                             |

---

## Data Models

### Transaction

```json
{
  "id": 1,
  "date": "2026-04-01T00:00:00Z",
  "vendor": "WHOLEFDS #1234",
  "description": "Weekly shop",
  "category": "Groceries",
  "amount": "84.32",
  "payment_method_id": 1
}
```

### Income

```json
{
  "id": 1,
  "date": "2026-04-01T00:00:00Z",
  "amount": "4200.00",
  "source": "Employer",
  "income_type": "salary",
  "description": "April paycheck"
}
```

### Payment Method

```json
{
  "id": 1,
  "name": "Chase Sapphire",
  "type": "credit",
  "base_points_rate": "2.00"
}
```

### CSV Format

```json
{
  "id": 1,
  "payment_method_id": 1,
  "csv_column": "description",
  "column_type": "vendor"
}
```

---

## Transactions

### `GET /transactions`

Returns a list of transactions. Supports filtering and pagination.

**Query Parameters**
| Param | Type | Description |
|-------|------|-------------|
| `payment_method_id` | integer | Filter by payment method |
| `category` | string | Filter by category e.g. `Groceries` |
| `labeled` | boolean | `false` returns only rows where category IS NULL |
| `month` | string | Filter by month e.g. `2026-04` |
| `amount_from` | float | Filter by minimum amount |
| `amount_to` | float | Filter by maximum amount |
| `date_from` | string | Filter from date `YYYY-MM-DD` |
| `date_to` | string | Filter to date `YYYY-MM-DD` |
| `limit` | integer | Number of rows to return, default 50 |
| `offset` | integer | Pagination offset, default 0 |

**Response `200 OK`**

```json
{
  "data": [ ...transactions ],
  "total": 124,
  "limit": 50,
  "offset": 0
}
```

---

### `GET /transactions/{id}`

Returns a single transaction by ID.

**Response `200 OK`**

```json
{
  "data": { ...transaction }
}
```

**Response `404 Not Found`**

```json
{
  "error": "transaction not found"
}
```

---

### `POST /transactions/import`

Uploads a CSV file and imports its contents as transactions.

**Request**

- Content-Type: `multipart/form-data`
- Fields:
  - `file` — the CSV file
  - `payment_method_id` — integer ID of the payment method to use for import

**Behavior**

- Reads the CSV format config from `csv_formats` using `payment_method_id`
- Maps CSV columns to transaction fields using the config
- Skips rows where `amount > 0` (credits/payments)
- Checks each row for duplicates by matching `date + amount + vendor`
- Inserts all non-duplicate rows in a single transaction
- Returns a summary

**Response `200 OK`**

```json
{
  "inserted": 42,
  "skipped_duplicates": 3,
  "skipped_credits": 1,
  "errors": []
}
```

**Response `400 Bad Request`**

```json
{
  "error": "payment_method_id is required"
}
```

---

### `PUT /transactions/{id}`

Updates the editable fields of a transaction. `date`, `vendor`, and `amount` are read-only.

**Request Body**

```json
{
  "category": "Groceries",
  "description": "Weekly shop",
  "payment_method_id": 2
}
```

All fields optional. Only send what you want to update.

**Response `200 OK`**

```json
{
  "data": { ...updated transaction }
}
```

**Response `400 Bad Request`** — returned if request attempts to update a read-only field

```json
{
  "error": "field 'amount' is read-only"
}
```

**Response `404 Not Found`**

```json
{
  "error": "transaction not found"
}
```

---

## Income

### `GET /income`

Returns a list of income entries. Supports filtering and pagination.

**Query Parameters**
| Param | Type | Description |
|-------|------|-------------|
| `month` | string | Filter by month e.g. `2026-04` |
| `income_type` | string | Filter by type: `salary`, `bonus`, `freelance`, `other` |
| `limit` | integer | Number of rows to return, default 50 |
| `offset` | integer | Pagination offset, default 0 |

**Response `200 OK`**

```json
{
  "data": [ ...income entries ],
  "total": 12,
  "limit": 50,
  "offset": 0
}
```

---

### `GET /income/{id}`

Returns a single income entry by ID.

**Response `200 OK`**

```json
{
  "data": { ...income entry }
}
```

**Response `404 Not Found`**

```json
{
  "error": "income entry not found"
}
```

---

### `POST /income`

Creates a new income entry.

**Request Body**

```json
{
  "date": "2026-04-01",
  "amount": "4200.00",
  "source": "Employer",
  "income_type": "salary",
  "description": "April paycheck"
}
```

**Required fields:** `date`, `amount`, `source`, `income_type`

**Valid `income_type` values:** `salary`, `bonus`, `freelance`, `other`

**Response `201 Created`**

```json
{
  "data": { ...created income entry }
}
```

**Response `400 Bad Request`**

```json
{
  "error": "income_type is required"
}
```

---

### `PUT /income/{id}`

Updates any field on an income entry. All fields are editable.

**Request Body**

```json
{
  "date": "2026-04-01",
  "amount": "4200.00",
  "source": "Employer",
  "income_type": "salary",
  "description": "April paycheck"
}
```

All fields optional. Only send what you want to update.

**Response `200 OK`**

```json
{
  "data": { ...updated income entry }
}
```

**Response `404 Not Found`**

```json
{
  "error": "income entry not found"
}
```

---

### `DELETE /income/{id}`

Permanently deletes an income entry.

**Response `200 OK`**

```json
{
  "message": "income entry deleted"
}
```

**Response `404 Not Found`**

```json
{
  "error": "income entry not found"
}
```

---

## Payment Methods

### `GET /payment-methods`

Returns all saved payment methods.

**Response `200 OK`**

```json
{
  "data": [ ...payment methods ]
}
```

---

### `POST /payment-methods`

Creates a new payment method.

**Request Body**

```json
{
  "name": "Chase Sapphire",
  "type": "credit",
  "base_points_rate": "2.00"
}
```

**Required fields:** `name`, `type`

**Valid `type` values:** `credit`, `debit`

**Response `201 Created`**

```json
{
  "data": { ...created payment method }
}
```

**Response `400 Bad Request`**

```json
{
  "error": "type is required"
}
```

---

### `PUT /payment-methods/{id}`

Updates a payment method.

**Request Body**

```json
{
  "name": "Chase Sapphire Reserve",
  "type": "credit",
  "base_points_rate": "3.00"
}
```

All fields optional. Only send what you want to update.

**Response `200 OK`**

```json
{
  "data": { ...updated payment method }
}
```

**Response `404 Not Found`**

```json
{
  "error": "payment method not found"
}
```

---

### `DELETE /payment-methods/{id}`

Deletes a payment method. Fails if any transactions still reference it.

**Response `200 OK`**

```json
{
  "message": "payment method deleted"
}
```

**Response `409 Conflict`**

```json
{
  "error": "cannot delete payment method with existing transactions"
}
```

**Response `404 Not Found`**

```json
{
  "error": "payment method not found"
}
```

---

## CSV Formats

### `GET /csv-formats`

Returns all saved CSV format configs.

**Response `200 OK`**

```json
{
  "data": [ ...csv formats ]
}
```

---

### `GET /csv-formats/{id}`

Returns a single CSV format config by ID.

**Response `200 OK`**

```json
{
  "data": { ...csv format }
}
```

**Response `404 Not Found`**

```json
{
  "error": "csv format not found"
}
```

---

### `POST /csv-formats`

Creates a new CSV format mapping.

**Request Body**

```json
{
  "payment_method_id": 1,
  "csv_column": "description",
  "column_type": "vendor"
}
```

**Required fields:** `payment_method_id`, `csv_column`, `column_type`

**Valid `column_type` values:** `date`, `vendor`, `description`, `category`, `amount`

**Response `201 Created`**

```json
{
  "data": { ...created csv format }
}
```

**Response `400 Bad Request`**

```json
{
  "error": "invalid column_type"
}
```

---

### `PUT /csv-formats/{id}`

Updates a CSV format mapping.

**Request Body**

```json
{
  "csv_column": "memo",
  "column_type": "description"
}
```

All fields optional. Only send what you want to update.

**Response `200 OK`**

```json
{
  "data": { ...updated csv format }
}
```

**Response `404 Not Found`**

```json
{
  "error": "csv format not found"
}
```

---

### `DELETE /csv-formats/{id}`

Deletes a CSV format mapping.

**Response `200 OK`**

```json
{
  "message": "csv format deleted"
}
```

**Response `404 Not Found`**

```json
{
  "error": "csv format not found"
}
```

---

## Schema

### `GET /schema`

Returns the current DB column list and all enum values. Useful for dynamically building frontend forms.

**Response `200 OK`**

```json
{
  "transaction_columns": [
    "id",
    "date",
    "vendor",
    "description",
    "category",
    "amount",
    "payment_method_id"
  ],
  "payment_method_types": ["credit", "debit"],
  "income_types": ["salary", "bonus", "freelance", "other"],
  "valid_categories": [
    "Travel",
    "Groceries",
    "Transit",
    "Entertainment",
    "Shopping",
    "Dining",
    "Laundry",
    "Cell Phone",
    "Health",
    "Other"
  ]
}
```

---

## Valid Category Values

Validated server-side on every transaction PUT. Sending anything outside this list returns a `400`.

```
Travel, Groceries, Transit, Entertainment, Shopping,
Dining, Laundry, Cell Phone, Health, Other
```

# Tally — Database Schema

## Tables

### `payment_methods`

The central config object. Each card you own is one row.

```sql
CREATE TYPE payment_method_type AS ENUM ('credit', 'debit');

CREATE TABLE IF NOT EXISTS payment_methods (
    id               SERIAL PRIMARY KEY,
    name             TEXT NOT NULL,
    type             payment_method_type NOT NULL,
    base_points_rate NUMERIC(5,2)
);
```

| Column             | Type                | Notes                             |
| ------------------ | ------------------- | --------------------------------- |
| `id`               | SERIAL              | Auto-generated primary key        |
| `name`             | TEXT                | e.g. `Chase Sapphire`, `WF Debit` |
| `type`             | payment_method_type | `credit` or `debit`               |
| `base_points_rate` | NUMERIC(5,2)        | Optional, points per dollar       |

---

### `csv_formats`

Maps CSV column names to standardized field types per payment method.

```sql
CREATE TABLE IF NOT EXISTS csv_formats (
    id                SERIAL PRIMARY KEY,
    payment_method_id INTEGER REFERENCES payment_methods(id),
    csv_column        TEXT NOT NULL,
    column_type       TEXT NOT NULL,
    UNIQUE (payment_method_id, csv_column)
);
```

| Column              | Type    | Notes                                             |
| ------------------- | ------- | ------------------------------------------------- |
| `id`                | SERIAL  | Auto-generated primary key                        |
| `payment_method_id` | INTEGER | FK → payment_methods.id                           |
| `csv_column`        | TEXT    | Raw header name in the CSV e.g. `name`, `date`    |
| `column_type`       | TEXT    | Standardized type e.g. `vendor`, `date`, `amount` |

---

### `transactions`

All imported transactions.

```sql
CREATE TABLE IF NOT EXISTS transactions (
    id                SERIAL PRIMARY KEY,
    date              DATE          NOT NULL,
    vendor            TEXT          NOT NULL,
    description       TEXT,
    category          TEXT,
    amount            NUMERIC(10,2) NOT NULL,
    payment_method_id INTEGER REFERENCES payment_methods(id)
);
```

| Column              | Type          | Notes                           |
| ------------------- | ------------- | ------------------------------- |
| `id`                | SERIAL        | Auto-generated primary key      |
| `date`              | DATE          | Transaction date                |
| `vendor`            | TEXT          | Merchant name from CSV          |
| `description`       | TEXT          | Optional user notes             |
| `category`          | TEXT          | Optional, validated server-side |
| `amount`            | NUMERIC(10,2) | Always stored as positive       |
| `payment_method_id` | INTEGER       | FK → payment_methods.id         |

---

### `income`

Manually entered income records.

```sql
CREATE TYPE income_type AS ENUM ('salary', 'bonus', 'freelance', 'other');

CREATE TABLE IF NOT EXISTS income (
    id          SERIAL PRIMARY KEY,
    date        DATE          NOT NULL,
    amount      NUMERIC(10,2) NOT NULL,
    source      TEXT          NOT NULL,
    income_type income_type   NOT NULL,
    description TEXT
);
```

| Column        | Type          | Notes                                   |
| ------------- | ------------- | --------------------------------------- |
| `id`          | SERIAL        | Auto-generated primary key              |
| `date`        | DATE          | Date of income                          |
| `amount`      | NUMERIC(10,2) | Always positive                         |
| `source`      | TEXT          | e.g. `Employer`, `Freelance Client`     |
| `income_type` | income_type   | `salary`, `bonus`, `freelance`, `other` |
| `description` | TEXT          | Optional notes                          |

---

## Relationships

```
payment_methods
  └── csv_formats       (payment_method_id → payment_methods.id)
  └── transactions      (payment_method_id → payment_methods.id)
```

---

## Enum Types

### `payment_method_type`

```sql
CREATE TYPE payment_method_type AS ENUM ('credit', 'debit');
```

### `income_type`

```sql
CREATE TYPE income_type AS ENUM ('salary', 'bonus', 'freelance', 'other');
```

---

## Valid Categories

Validated server-side on every transaction PUT. Sending anything outside this list returns a `400`.

```
Travel, Groceries, Transit, Entertainment, Shopping,
Dining, Laundry, Cell Phone, Health, Other
```

---

## Migration Order

| File                             | Creates                                              |
| -------------------------------- | ---------------------------------------------------- |
| `001_create_payment_methods.sql` | `payment_method_type` enum + `payment_methods` table |
| `002_create_csv_formats.sql`     | `csv_formats` table                                  |
| `003_create_transactions.sql`    | `transactions` table                                 |
| `004_create_income.sql`          | `income_type` enum + `income` table                  |
