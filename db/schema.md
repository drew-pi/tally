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
