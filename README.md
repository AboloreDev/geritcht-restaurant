# Gericht Restaurant Backend

### A production-oriented restaurant management API built with Go, PostgreSQL, Redis, and event-driven architecture.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go\&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?logo=postgresql\&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-Streams%20%26%20Locks-DC382D?logo=redis\&logoColor=white)](https://redis.io/)
[![Swagger](https://img.shields.io/badge/API-Swagger-85EA2D?logo=swagger\&logoColor=black)](https://swagger.io/)

---

## 📖 Overview

**Gericht Restaurant Backend** is a production-style restaurant management system built with **Go** and **PostgreSQL**.

The project goes beyond traditional CRUD APIs by implementing backend engineering patterns focused on:

* Reliability
* Scalability
* Concurrency
* Performance
* Asynchronous processing
* Data consistency
* Real-time communication

The system manages restaurant operations including **authentication, reservations, orders, payments, inventory, notifications, and real-time order tracking**.

A major goal of the project was to explore how production backend systems handle **concurrent requests, failures, background processing, duplicate operations, and database performance**.

---

## ✨ Features

### 🔐 Authentication & Authorization

* JWT-based authentication
* Access and refresh token flow
* Password hashing
* Role-based authorization
* Protected API endpoints
* Secure token management

### 🍽️ Restaurant Management

* Menu management
* Category management
* Table reservations
* Order management
* Customer reviews
* Inventory management
* Restaurant availability

### 💳 Payments

* Idempotent payment processing
* Refund workflows
* Payment webhook handling
* Duplicate payment protection
* Distributed locking for payment operations

### 📦 Inventory

* Automatic stock deduction
* Low-stock detection
* Low-stock notifications
* Automatic menu availability updates
* Event-driven inventory processing

### 📡 Real-Time Order Tracking

* WebSocket-based communication
* Live order status updates
* Event-driven order notifications
* Real-time client updates

### 🔄 Background Processing

* Redis Streams
* Background workers
* Asynchronous event processing
* Notification processing
* Email delivery
* Inventory updates

### 🔎 Search & Performance

* PostgreSQL full-text search
* Weighted `tsvector`
* PostgreSQL indexes
* Query optimization
* `EXPLAIN ANALYZE`
* Pagination
* Efficient relational queries

---

# 🏗️ System Architecture

The backend follows an **event-driven architecture** where synchronous API operations are separated from asynchronous background processing.

```text
                         ┌────────────────────┐
                         │      Client        │
                         │  Web / Mobile App  │
                         └─────────┬──────────┘
                                   │
                                   ▼
                         ┌────────────────────┐
                         │      Go API        │
                         │       Gin          │
                         └─────────┬──────────┘
                                   │
                    ┌──────────────┴──────────────┐
                    │                             │
                    ▼                             ▼
          ┌──────────────────┐           ┌──────────────────┐
          │   PostgreSQL     │           │      Redis       │
          │                  │           │                  │
          │ Orders           │           │ Streams          │
          │ Users            │           │ Distributed Lock │
          │ Reservations     │           │ Pub/Sub          │
          │ Inventory        │           └────────┬─────────┘
          │ Payments         │                    │
          └────────┬─────────┘                    ▼
                   │                     ┌──────────────────┐
                   │                     │ Background       │
                   │                     │ Workers          │
                   │                     └────────┬─────────┘
                   │                              │
                   │                ┌─────────────┼─────────────┐
                   │                │             │             │
                   ▼                ▼             ▼             ▼
          ┌──────────────┐     ┌─────────┐  ┌──────────┐  ┌───────────┐
          │    Outbox    │     │ Emails  │  │Inventory │  │ Events    │
          │    Events    │     │         │  │ Updates  │  │           │
          └──────────────┘     └─────────┘  └──────────┘  └───────────┘
```

---

# 🔄 Event-Driven Order Flow

When an order is created, operations that require immediate consistency are handled synchronously while secondary operations are processed asynchronously.

```text
Customer
   │
   │ Create Order
   ▼
┌─────────────────┐
│     Go API      │
└────────┬────────┘
         │
         ▼
┌──────────────────────────┐
│ PostgreSQL Transaction   │
│                          │
│ • Create Order           │
│ • Update Inventory       │
│ • Create Outbox Event    │
└────────────┬─────────────┘
             │
             ▼
        Transaction
          COMMIT
             │
             ▼
      ┌───────────────┐
      │ Redis Stream  │
      └───────┬───────┘
              │
              ▼
      ┌─────────────────┐
      │ Background      │
      │ Worker          │
      └────────┬────────┘
               │
       ┌───────┼────────┐
       ▼       ▼        ▼
    Email  Inventory  Notification
       │       │        │
       └───────┼────────┘
               ▼
        WebSocket Update
               │
               ▼
            Client
```

---

# 📦 Transactional Outbox Pattern

The backend uses the **Transactional Outbox Pattern** to make event delivery more reliable.

Instead of directly sending an email or publishing an external event during an HTTP request, the event is stored in the same database transaction as the business operation.

```text
              API Request
                   │
                   ▼
        ┌────────────────────┐
        │ PostgreSQL         │
        │ Transaction        │
        │                    │
        │ Order              │
        │       +            │
        │ Outbox Event       │
        └─────────┬──────────┘
                  │
                  ▼
              COMMIT
                  │
                  ▼
        ┌────────────────────┐
        │ Outbox Worker      │
        └─────────┬──────────┘
                  │
                  ▼
          Redis Stream
                  │
                  ▼
        ┌────────────────────┐
        │ Background Worker  │
        └─────────┬──────────┘
                  │
                  ▼
          Email / Notification
```

### Why use it?

Without an Outbox Pattern, this could happen:

```text
Database Update ✅
       │
       ▼
Email Sending ❌
       │
       ▼
Lost Notification
```

With the Outbox Pattern:

```text
Database + Event
      │
      ▼
   COMMIT
      │
      ▼
Worker retries delivery
```

This makes asynchronous processing more resilient to temporary failures.

---

# 🔒 Distributed Locks

The system uses Redis-based distributed locks to protect operations that cannot safely execute concurrently.

For example, two customers attempting to reserve the same table at almost the same time could create a race condition.

```text
Customer A ───────┐
                  │
                  ▼
            ┌─────────────┐
            │ Redis Lock  │
            └──────┬──────┘
                   │
              Lock Acquired
                   │
                   ▼
            Check Availability
                   │
                   ▼
            Create Reservation
                   │
                   ▼
             Release Lock
                   │
                   ▼
Customer B ────► Check Again
```

Distributed locks are used for concurrency-sensitive operations such as:

* Table reservations
* Payment processing
* Refund processing
* Webhook handling
* Inventory operations

---

# 🔁 Idempotent Payments

Payment operations are designed to be **idempotent**, meaning retrying the same operation does not accidentally create multiple payments or refunds.

```text
Client Request
      │
      ▼
Idempotency Key
      │
      ▼
┌──────────────────┐
│ Check Existing   │
│ Operation        │
└────────┬─────────┘
         │
    ┌────┴─────┐
    │           │
   YES          NO
    │           │
    ▼           ▼
Return       Process
Existing     Payment
Result          │
                ▼
          Store Result
```

This is particularly important for payment systems where network failures can cause clients to retry requests.

---

# 📡 Real-Time Order Tracking

WebSockets allow clients to receive order updates without repeatedly polling the REST API.

```text
Order Status Changes
        │
        ▼
    Go Backend
        │
        ▼
   Event Processing
        │
        ▼
 WebSocket Connection
        │
        ▼
      Client
        │
        ▼
┌────────────────────────┐
│ Preparing              │
│        ↓               │
│ Ready                  │
│        ↓               │
│ Completed              │
└────────────────────────┘
```

This provides a more responsive experience for customers tracking their orders.

---

# 📦 Inventory Automation

Inventory is automatically updated as orders are processed.

```text
             Order Created
                   │
                   ▼
          Deduct Inventory
                   │
                   ▼
            Check Stock
                   │
          ┌────────┴────────┐
          │                 │
          ▼                 ▼
      Stock OK          Low Stock
          │                 │
          │                 ▼
          │          Send Notification
          │
          ▼
   Update Availability
          │
          ▼
   Menu Item Available/
       Unavailable
```

This allows the restaurant's menu availability to automatically reflect current inventory.

---

# 🔎 PostgreSQL Full-Text Search

The backend uses PostgreSQL's native full-text search capabilities instead of relying on simple `LIKE` queries.

Search uses:

* `tsvector`
* Weighted search fields
* PostgreSQL ranking
* Full-text indexes

```text
Search Request
      │
      ▼
PostgreSQL
      │
      ▼
Weighted tsvector
      │
      ▼
Relevance Ranking
      │
      ▼
Relevant Results
```

This allows important fields to have higher search relevance than less important fields.

---

# 🚀 Database Performance

Database performance was analyzed and optimized using PostgreSQL tooling.

The optimization workflow:

```text
       Query
         │
         ▼
  EXPLAIN ANALYZE
         │
         ▼
Identify Bottleneck
         │
         ▼
Optimize Query
         │
         ▼
Add / Improve Index
         │
         ▼
  EXPLAIN ANALYZE
         │
         ▼
Compare Performance
```

Performance techniques include:

* Strategic database indexes
* Query optimization
* Efficient joins
* Full-text indexes
* Pagination
* Reduced unnecessary database queries
* Query analysis using `EXPLAIN ANALYZE`

---

# 🧵 Redis Streams & Background Workers

Redis Streams are used as an asynchronous event-processing mechanism.

```text
Go API
  │
  ▼
Redis Stream
  │
  ├──────────────┐
  │              │
  ▼              ▼
Worker 1       Worker 2
  │              │
  ▼              ▼
Email         Inventory
Processing    Processing
  │              │
  └──────┬───────┘
         ▼
     Completed
```

Background processing keeps expensive or non-critical operations out of the main request-response cycle.

---

# 🛠️ Tech Stack

| Technology     | Purpose                                              |
| -------------- | ---------------------------------------------------- |
| **Go**         | Backend application                                  |
| **Gin**        | HTTP web framework                                   |
| **PostgreSQL** | Primary relational database                          |
| **Redis**      | Streams, distributed locks & asynchronous processing |
| **JWT**        | Authentication                                       |
| **WebSockets** | Real-time communication                              |
| **Swaggo**     | Swagger/OpenAPI documentation                        |
| **Docker**     | Containerization                                     |
| **Render**     | Deployment                                           |

---

# 📂 Project Structure

```text
server/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── handlers/
│   ├── services/
│   ├── repositories/
│   ├── middleware/
│   ├── workers/
│   ├── events/
│   ├── websocket/
│   ├── inventory/
│   └── ...
│
├── migrations/
│
├── docs/
│
├── go.mod
├── go.sum
└── README.md
```

> The exact directory structure may evolve as the project grows.

---

# ⚙️ Getting Started

## Prerequisites

Make sure you have the following installed:

* Go 1.22+
* PostgreSQL
* Redis
* Git

---

## Clone the Repository

```bash
git clone https://github.com/AboloreDev/geritcht-restaurant.git

cd geritcht-restaurant/server
```

---

## Install Dependencies

```bash
go mod download
```

or:

```bash
go mod tidy
```

---

## Environment Variables

Create a `.env` file in the server directory.

Example:

```env
PORT=8080

DATABASE_URL=postgres://postgres:password@localhost:5432/geritcht

REDIS_URL=redis://localhost:6379

JWT_SECRET=your_jwt_secret

```

> Add any additional environment variables required by the payment, email, or deployment configuration.

---

# ▶️ Running the Application

Start the server with:

```bash
go run ./cmd/server
```

The API should then be available at:

```text
http://localhost:8080
```

---

# 📚 API Documentation

The API is documented using **Swaggo/Swagger**.

### Live Documentation

[Open Swagger API Documentation](https://geritcht-restaurant.onrender.com/api-docs)

The documentation provides an interactive interface for exploring and testing available endpoints.

It includes documentation for areas such as:

* Authentication
* Users
* Menu
* Reservations
* Orders
* Payments
* Refunds
* Inventory
* Reviews
* Notifications
* Administrative operations

---

# 🧪 Testing

The project includes automated tests for core backend functionality.

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run tests with race detection:

```bash
go test -race ./...
```

The race detector is particularly useful for identifying concurrency issues in systems involving goroutines, workers, and shared resources.

---

# 🔐 Security Considerations

Security is incorporated throughout the application through:

* JWT authentication
* Password hashing
* Role-based authorization
* Protected endpoints
* Token validation
* Input validation
* Idempotent payment operations
* Distributed locks
* Secure environment configuration

Secrets and credentials should be provided through environment variables and should **never be committed to the repository**.

---

# 📊 Engineering Highlights

The project focuses on several important backend engineering concepts:

### Reliability

* Transactional Outbox Pattern
* Idempotent operations
* Background workers
* Retryable asynchronous processing

### Concurrency

* Redis distributed locks
* Goroutines
* Worker-based processing
* Race-condition prevention

### Scalability

* Event-driven architecture
* Redis Streams
* Asynchronous workloads
* Efficient PostgreSQL queries

### Performance

* Database indexes
* Full-text search
* Query optimization
* `EXPLAIN ANALYZE`
* Pagination

### Real-Time Systems

* WebSockets
* Event-driven order updates
* Background event processing

---

# 🗺️ Architecture at a Glance

```text
                         ┌─────────────────────┐
                         │       CLIENT        │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │      REST API       │
                         │        Go/Gin       │
                         └──────────┬──────────┘
                                    │
                    ┌───────────────┼────────────────┐
                    │               │                │
                    ▼               ▼                ▼
              PostgreSQL          Redis          WebSockets
                    │               │                │
                    │               ▼                │
                    │       Redis Streams            │
                    │               │                │
                    │               ▼                │
                    │       Background Workers       │
                    │               │                │
                    │        ┌──────┼──────┐         │
                    │        ▼      ▼      ▼         │
                    │      Email  Events Inventory   │
                    │                                │
                    └───────────────┬────────────────┘
                                    │
                                    ▼
                              Real-Time Updates
                                    │
                                    ▼
                                  CLIENT
```

---

# 🎯 Project Goals

The main objective of this project was to use **Go to explore real-world backend engineering patterns** rather than building a simple CRUD application.

The project provided hands-on experience with:

* Distributed systems
* Event-driven architecture
* Database transactions
* Concurrency
* Background processing
* Message streams
* Distributed locking
* Idempotency
* Real-time communication
* Database performance optimization

Go has been particularly useful for exploring these concepts because of its lightweight concurrency model, strong standard library, and straightforward approach to building network services.

---

# 🚧 Future Improvements

Potential areas for further development include:

* Comprehensive observability with metrics and tracing
* Structured centralized logging
* Automated CI/CD pipeline
* Horizontal worker scaling
* Improved retry and dead-letter queue strategies
* Load testing and benchmarking
* Container orchestration
* Expanded integration test coverage

---

# 🔗 Links

**GitHub Repository**

https://github.com/AboloreDev/geritcht-restaurant

**Live API Documentation**

https://geritcht-restaurant.onrender.com/api-docs

---

# 👨‍💻 Author

**AboloreDev**

Built with **Go** while exploring production-oriented backend engineering, distributed systems, concurrency, and database performance.

---

## ⭐ If you found this project useful

Feel free to explore the repository, try the API through the Swagger documentation, and share feedback or suggestions.

**Built with Go. Designed for reliability. Engineered for scale.**
