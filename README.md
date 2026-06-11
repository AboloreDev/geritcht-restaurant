# Geritcht — Restaurant Management Platform

A production-grade restaurant platform built in Go, designed to handle real-world concurrency, event-driven workflows, and financial-grade payment processing.

Customers can browse the menu, place takeout orders, and reserve tables for dinner. Staff manage orders and reservations in real time. Admins control everything from menu management to inventory and analytics.

---

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [System Design Decisions](#system-design-decisions)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [API Overview](#api-overview)
- [Background Workers](#background-workers)
- [Email Microservice](#email-microservice)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        Client (REST)                        │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Go + Gin HTTP Server                      │
│  Auth │ Menu │ Cart │ Orders │ Reservations │ Payments       │
└──────┬──────────────────┬──────────────────┬────────────────┘
       │                  │                  │
┌──────▼──────┐  ┌────────▼────────┐  ┌─────▼──────────────┐
│  PostgreSQL  │  │  Redis (Upstash) │  │  Paystack API      │
│  (Neon)     │  │  Cache + Streams │  │  Payments/Refunds  │
└─────────────┘  └────────┬────────┘  └────────────────────┘
                          │
              ┌───────────▼───────────┐
              │  Email Microservice   │
              │  (Separate Go binary) │
              │  Resend API           │
              └───────────────────────┘
```

---

## Features

### Authentication

- JWT access tokens (15 min) + refresh token rotation (7 days)
- Refresh tokens stored as httpOnly cookies — inaccessible to JavaScript
- SHA256 token hashing before database storage
- Email verification, forgot password, reset password
- Reset password invalidates all active sessions across devices
- Role-based access control: **customer**, **staff**, **admin**

### Menu System

- Categories with Cloudinary image upload
- Menu items with multiple images (max 4, primary auto-assigned)
- Many-to-many allergen and dietary tag relationships
- Redis caching on all read endpoints (1hr TTL)
- Server-side filtering: category, price range, spice level, dietary, allergen exclusion, search
- Server-side sorting with column whitelist (SQL injection prevention)
- Cache TTL reduces to 30 minutes for filtered requests
- Pattern-based cache invalidation on every write

### Reservation System

- Time slot based availability (configurable slots)
- **2 DB queries total** regardless of slot count — grouped in memory with `map[tableID]bool` for O(1) lookup per table
- Availability cached per date + time slot + party size (30s TTL)
- **Zero double bookings** via two layers:
  - Redis `SETNX` distributed lock at application level
  - PostgreSQL `FOR UPDATE` inside transaction at database level
- Waitlist system — first-come-first-served via `created_at` ordering
- Staff check-in and no-show marking with automatic table status updates
- DB unique index: `(table_id, date, time_slot)` WHERE status NOT IN (cancelled, no_show)

### Order System

- Takeout orders created from cart
- Dine-in orders linked to reservations (staff only)
- Order state machine: `pending → confirmed → preparing → ready → completed → cancelled`
- Cart system with ownership verification via SQL JOIN

### Payment System (Financial Grade)

- **Idempotency**: unique reference + idempotency key with DB-level unique constraints
- **Paystack webhook** with HMAC-SHA512 signature verification
- **Redis distributed lock** on webhook processing — prevents duplicate processing of retried webhooks
- **Amount verification** — prevents partial payment attacks
- Atomic DB transaction: payment update + order confirmation in one operation
- **Outbox pattern**: events written to DB inside the same transaction, published by background worker — survives Redis outages
- Dedicated **Refunds table** with own idempotency key
- Refund idempotency check before calling Paystack API

### Email Microservice

- Completely separate Go binary
- Main API publishes jobs to **Redis Streams** via Watermill
- Email service subscribes via consumer group — prevents duplicate sends
- `msg.Ack()` on success, `msg.Nack()` on failure (auto-retry)
- **Resend** handles delivery via verified custom domain
- Publisher behind interface — swap Redis Streams for SQS in one line

### Inventory Management

- Ingredients CRUD with minimum threshold tracking
- Recipe management — ingredients linked to menu items with quantities
- Auto stock deduction on order confirmation (inside transaction)
- Menu items auto-disabled when ingredient stock depleted
- Low stock alerts published to Redis Streams → admin email notification

### Background Workers

All workers use `time.Ticker` + `context.Done()` for graceful shutdown:

| Worker                | Interval | Purpose                                                     |
| --------------------- | -------- | ----------------------------------------------------------- |
| Reminder Worker       | 5 min    | Sends reservation reminders 30 min before slot              |
| No-Show Worker        | 5 min    | Auto-marks confirmed reservations as no-show after 15 min   |
| Checkout Worker       | 5 min    | Auto-completes reservations after 90 min slot duration      |
| Outbox Worker         | 30 sec   | Publishes pending outbox events to Redis Streams with retry |
| Reconciliation Worker | 30 min   | Verifies pending payments with Paystack API                 |

---

## Tech Stack

| Layer           | Technology                        |
| --------------- | --------------------------------- |
| Language        | Go 1.23                           |
| Framework       | Gin                               |
| Database        | PostgreSQL (Neon)                 |
| Cache / Streams | Redis (Upstash)                   |
| ORM             | GORM                              |
| Migrations      | golang-migrate                    |
| Event Bus       | Watermill + watermill-redisstream |
| Email           | Resend                            |
| Payments        | Paystack                          |
| Image Storage   | Cloudinary                        |
| Logging         | Zerolog                           |
| Containers      | Docker + Docker Compose           |
| CI/CD           | GitHub Actions                    |
| Deployment      | Railway                           |

---

## System Design Decisions

### Why Go

90% of similar projects use Node.js or Laravel. Go was chosen for its concurrency primitives, performance, and growing adoption in fintech and infrastructure companies. Every engineering decision in this project maps directly to Go's strengths.

### Why Redis Streams over Pub/Sub

Redis Pub/Sub drops messages if the subscriber is offline. Streams persist messages and support consumer groups — only one instance of the email service processes each message even when scaled horizontally.

### Why Outbox Pattern

Without the outbox pattern, a Redis outage between payment confirmation and event publishing would cause silent email failures. Writing the event to the database inside the same transaction as the business operation guarantees the event is never lost.

### Why Distributed Locking on Reservations

A count-based availability check inside a transaction is not sufficient under concurrent load. Two requests can read the same count simultaneously, both see availability, and both proceed to book. Redis `SETNX` + PostgreSQL `FOR UPDATE` provides two independent layers of protection.

### Why Idempotency Keys on Payments

Paystack can deliver the same webhook multiple times. Without idempotency, the same payment could confirm an order twice and deduct stock twice. The unique constraint on `reference` combined with a status check creates a fast-fail path for duplicate webhook deliveries.

---

## Project Structure

```
geritcht-restaurant/
├── cmd/
│   ├── api/              ← main HTTP server
│   └── notifier/         ← email microservice
├── internals/
│   ├── config/
│   ├── database/
│   ├── dto/
│   ├── domain/           ← sentinel errors
│   ├── events/           ← event types + channel constants
│   ├── interfaces/       ← Cacher, Publisher, UploadProvider
│   ├── mapper/           ← pure conversion functions
│   ├── middleware/
│   ├── models/
│   ├── providers/        ← Cloudinary, LocalUpload
│   ├── publisher/        ← Redis Streams publisher
│   ├── redis/            ← RedisStore, NopStore
│   ├── server/           ← handlers + routes
│   ├── services/         ← business logic + workers
│   └── utils/
├── db/
│   └── migrations/       ← 24 up/down SQL migration files
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── Makefile
└── .env.example
```

---

## Getting Started

### Prerequisites

- Go 1.23+
- Docker + Docker Compose
- Make

### Run locally

```bash
# clone the repo
git clone https://github.com/AboloreDev/geritcht-restaurant-web.git
cd geritcht-restaurant-web/server

# copy environment variables
cp .env.example .env
# fill in your values

# start dependencies (postgres, redis, localstack)
make dev-up

# run migrations
make migrate-up-local

# start the API server
make run

# start the email microservice (separate terminal)
make run-notifier
```

---

## Environment Variables

```env
# Server
SERVER_PORT=9090
APP_ENV=development
GIN_MODE=debug

# Database
DATABASE_URL=postgres://postgres:password@localhost:5433/geritcht_db?sslmode=disable

# JWT
JWT_SECRET=
JWT_EXPIRATION_MINUTES=15
JWT_REFRESH_EXPIRATION_DAYS=7

# Redis
REDIS_URL=redis://localhost:6379

# Paystack
PAYSTACK_SECRET_KEY=sk_test_
PAYSTACK_PUBLIC_KEY=pk_test_

# Cloudinary
CLOUDINARY_CLOUD_NAME=
CLOUDINARY_API_KEY=
CLOUDINARY_API_SECRET=
CLOUDINARY_FOLDER=geritcht

# Resend
RESEND_API_KEY=
RESEND_FROM_EMAIL=noreply@yourdomain.com

# AWS (LocalStack for dev)
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
AWS_S3_BUCKET=geritcht-bucket
AWS_S3_ENDPOINT=http://localhost:4567
AWS_EVENT_QUEUE_NAME=geritcht-events
```

---

## API Overview

### Public (no auth)

| Method | Endpoint                   | Description                               |
| ------ | -------------------------- | ----------------------------------------- |
| GET    | `/api/v1/menu`             | Get all menu items (filterable, sortable) |
| GET    | `/api/v1/menu/:id`         | Get single menu item                      |
| GET    | `/api/v1/categories`       | Get all categories                        |
| GET    | `/api/v1/categories/:id`   | Get single category                       |
| GET    | `/api/v1/availability`     | Check table availability                  |
| POST   | `/api/v1/payments/webhook` | Paystack webhook                          |

### Auth

| Method | Endpoint                | Description          |
| ------ | ----------------------- | -------------------- |
| POST   | `/api/v1/auth/register` | Register customer    |
| POST   | `/api/v1/auth/login`    | Login                |
| POST   | `/api/v1/auth/logout`   | Logout               |
| POST   | `/api/v1/auth/refresh`  | Refresh access token |
| POST   | `/api/v1/auth/verify`   | Verify email         |
| POST   | `/api/v1/auth/forgot`   | Forgot password      |
| POST   | `/api/v1/auth/reset`    | Reset password       |

### Customer (auth required)

| Method | Endpoint                          | Description         |
| ------ | --------------------------------- | ------------------- |
| GET    | `/api/v1/cart`                    | Get cart            |
| POST   | `/api/v1/cart`                    | Add item to cart    |
| PATCH  | `/api/v1/cart/:id`                | Update cart item    |
| DELETE | `/api/v1/cart/:id`                | Remove cart item    |
| DELETE | `/api/v1/cart`                    | Clear cart          |
| POST   | `/api/v1/orders/takeout`          | Place takeout order |
| GET    | `/api/v1/orders/my`               | Get my orders       |
| POST   | `/api/v1/reservations`            | Create reservation  |
| GET    | `/api/v1/reservations/my`         | Get my reservations |
| PATCH  | `/api/v1/reservations/:id/cancel` | Cancel reservation  |
| POST   | `/api/v1/payments/initialize`     | Initialize payment  |
| GET    | `/api/v1/payments/verify/:ref`    | Verify payment      |

### Staff (staff + admin)

| Method | Endpoint                                 | Description          |
| ------ | ---------------------------------------- | -------------------- |
| GET    | `/api/v1/staff/orders`                   | Get all orders       |
| PATCH  | `/api/v1/staff/orders/:id/status`        | Update order status  |
| GET    | `/api/v1/staff/reservations/today`       | Today's reservations |
| PATCH  | `/api/v1/staff/reservations/:id/checkin` | Check in reservation |
| PATCH  | `/api/v1/staff/reservations/:id/no-show` | Mark no-show         |

### Admin

| Method | Endpoint                   | Description      |
| ------ | -------------------------- | ---------------- |
| POST   | `/api/v1/admin/menu`       | Create menu item |
| PATCH  | `/api/v1/admin/menu/:id`   | Update menu item |
| DELETE | `/api/v1/admin/menu/:id`   | Delete menu item |
| POST   | `/api/v1/admin/categories` | Create category  |
| GET    | `/api/v1/admin/users`      | Get all users    |
| GET    | `/api/v1/admin/staff`      | Get all staff    |
| GET    | `/api/v1/admin/analytics`  | Get analytics    |

---

## Background Workers

All workers start as goroutines in `main.go` and stop cleanly on `SIGINT`/`SIGTERM`:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go reminderWorker.StartReminderWorker(ctx, log)
go noShowWorker.StartNoShowWorker(ctx, log)
go checkoutWorker.StartCheckoutWorker(ctx, log)
go outboxWorker.StartOutboxWorker(ctx, log)
go reconciliationWorker.StartReconciliationWorker(ctx, log)
```

---

## Email Microservice

The email service runs as a **completely separate binary** (`cmd/notifier`).

```
Main API → publishes to Redis Streams → returns immediately
Email Service → subscribes → processes → calls Resend API → Ack/Nack
```

Events handled:

- `auth.send_verification_email`
- `auth.send_password_reset_email`
- `auth.send_password_changed_email`
- `order.confirmed` (receipt)
- `order.refunded`
- `reservation.confirmation`
- `reservation.reminder`
- `inventory.low_stock_alert`

---

Built by [Alabi Fathiu](https://github.com/AboloreDev) — building in public.
