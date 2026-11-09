# Ecommerce REST API

A backend e-commerce application built with Go, Gin, GORM, and MySQL. This project provides RESTful APIs for managing products, categories, shopping carts, orders, and simulated payments.

## Overview

This project demonstrates the implementation of an e-commerce backend system using a layered architecture. It focuses on REST API development, authentication and authorization, database relationships, transaction management, and business logic.

**Project Scope:** Backend REST API only. No frontend application or production deployment is included.

## Tech Stack

- **Language:** Go
- **Web Framework:** Gin
- **ORM:** GORM
- **Database:** MySQL
- **Authentication:** JWT
- **Password Hashing:** bcrypt
- **Environment Configuration:** godotenv
- **API Testing:** Postman

## Features

### Authentication & Authorization
- User registration and login
- JWT-based authentication
- Password hashing with bcrypt
- Role-based access control (Customer and Admin)

### Product & Category Management
- Create, read, update, and delete categories
- Create, read, update, and delete products
- Product and category relationships

### Shopping Cart
- Add products to cart
- Update item quantities
- Remove cart items
- Calculate cart totals
- Validate product stock

### Order Management
- Checkout process
- Order history and detail
- Order cancellation
- Stock management
- Order status management

### Payment Simulation
- Create payment transactions
- Support simulated QRIS, E-Wallet, and Virtual Account methods
- Simulate successful and failed payments
- Update order status after successful payment

> Payment functionality is simulated for development and learning purposes. It is not connected to a real payment gateway.

## Project Architecture

The application uses a layered architecture to separate responsibilities.

| Layer | Responsibility |
|---|---|
| Config | Application and environment configuration |
| Database | Database connection and migrations |
| Models | Database entities |
| DTO | Request and response structures |
| Handlers | HTTP request and response handling |
| Services | Business logic |
| Repositories | Database operations |
| Middleware | Authentication and authorization |
| Routes | API endpoint registration |

## Project Structure

```text
Ecommerce/
├── config/
├── database/
├── docs/
│   └── API_TESTING.md
├── dto/
├── handlers/
├── middleware/
├── models/
├── repositories/
├── routes/
├── services/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Requirements

- Go 1.27.3
- MySQL
- Git
- Postman (for API testing)

## Installation

### 1. Clone Repository

```bash
git clone https://github.com/YOUR_USERNAME/ecommerce-go-api.git
cd ecommerce-go-api
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Create Database

Create a MySQL database:

```sql
CREATE DATABASE ecommerce;
```

### 4. Configure Environment

Copy `.env.example` to `.env`, then configure your local database credentials and JWT secret.

### 5. Run Application

```bash
go run .
```

The API will run at:

```text
http://localhost:8080
```

Database tables are created through GORM AutoMigrate when the application starts.

## API Documentation

The API uses the following base URL:

```text
http://localhost:8080/api/v1
```

| Module | Main Endpoints |
|---|---|
| Authentication | `/auth/register`, `/auth/login`, `/auth/profile` |
| Categories | `/categories` |
| Products | `/products` |
| Shopping Cart | `/cart`, `/cart/items` |
| Orders | `/orders`, `/orders/checkout` |
| Payments | `/orders/:id/payments`, `/payments/:id` |
| Admin Orders | `/admin/orders` |

For detailed request payloads, response expectations, and API test cases, see [`docs/API_TESTING.md`](docs/API_TESTING.md).

## API Testing

API endpoints can be tested using Postman.

The testing documentation includes:
- HTTP methods and endpoint URLs
- JSON request payloads
- Expected HTTP status codes
- Expected API behavior
- Authentication and authorization scenarios

## Key Learning Outcomes

- Building RESTful APIs with Go and Gin
- Applying layered backend architecture
- Implementing JWT authentication and role-based authorization
- Managing relational databases with GORM
- Handling database transactions during checkout and payment
- Implementing stock validation and order management
- Structuring API documentation and testing

## License

This project is intended for educational and portfolio purposes.
