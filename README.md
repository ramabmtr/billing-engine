# Billing Engine

A RESTful API service for managing borrowers, loans, and payments. This application provides a complete solution for loan management, including creating borrowers, issuing loans, and processing payments.

## Features

- **Borrower Management**: Create and list borrowers
- **Loan Management**: Create loan requests, list loans, and view loan details
- **Payment Processing**: Make payments for loans and view payment history

## Tech Stack

- **Language**: Go
- **Web Framework**: Echo
- **Database**: PostgreSQL
- **ORM**: GORM
- **Documentation**: Swagger

## Installation

### Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher
- Git

### Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Create a `.env` file based on the example:
   ```bash
   cp .env.example .env
   ```

3. Update the `.env` file with your configuration.

4. Run database migrations:
   ```bash
   go run cmd/migration/main.go
   ```

5. Start the server:
   ```bash
   go run cmd/main.go
   ```

The server will start on the port specified in your `.env` file (default: 8080).

## Configuration

The application can be configured using environment variables:

### Server Configuration
- `SERVER_PORT`: Port on which the server will run (default: 8080)
- `SERVER_API_KEY`: API key for authentication (default: "")

### Database Configuration
- `DB_HOST`: Database host (default: "localhost")
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: "admin")
- `DB_PASSWORD`: Database password (default: "admin")
- `DB_NAME`: Database name (default: "billing_engine")
- `DB_SSLMODE`: SSL mode for database connection (default: "disable")

## API Documentation

The API documentation is available at `/docs` when the server is running. You can access it by navigating to `http://localhost:8080/docs` in your browser.

### API Endpoints

#### Borrowers
- `POST /api/borrowers`: Create a new borrower
- `GET /api/borrowers`: List all borrowers

#### Loans
- `POST /api/borrowers/:borrowerID/loans`: Create a loan request for a borrower
- `GET /api/borrowers/:borrowerID/loans`: List all loans for a borrower
- `GET /api/borrowers/:borrowerID/loans/:id`: Get detailed information about a loan

#### Payments
- `POST /api/borrowers/:borrowerID/loans/:loanID/payments`: Make a payment for a loan
- `GET /api/borrowers/:borrowerID/loans/:loanID/payments`: List all payments for a loan

### Authentication

All API endpoints (except `/api/ping` and `/docs`) require authentication using an API key. The API key should be provided in the `X-API-KEY` header.

## Contact

For any questions or feedback, please contact:
- **Name**: Rama Bramantara
- **Email**: ramabmtr@gmail.com