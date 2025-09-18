# Amartha - Coding Assessment

This project is a backend service for managing loans and investments, designed with a clean architecture approach. It handles operations related to borrowers, investors, employees, loans, and investments, including email notifications and agreement file generation.

## Features

-   **Loan Management:** Create, list, view details, approve, reject, and disburse loans.
-   **Investment Management:** Add new investments to loans.
-   **User Roles:** Differentiated access for Employees (Loan management) and Investors (Investment management).
-   **Email Notifications:** Automated email notifications for investment confirmations and loan investments.
-   **Agreement File Generation:** Generate PDF agreement files for loans and investments.

## Architecture

The project follows a Domain Driven Design (DDD) architecture to ensure a clear alignment between the software design and the business domain, fostering better communication, maintainability, and scalability by focusing on core business concepts and logic.

-   **`cmd/api`**: Contains the main entry point of the application, responsible for bootstrapping the server, initializing dependencies, and starting the HTTP server.
-   **`internal/application`**: Houses the application-specific logic, including repositories (data access), use cases (business logic), and delivery mechanisms (HTTP handlers, message handlers).
    -   `borrower`, `employee`, `investment`, `investor`, `loan`, `mail`: Each module contains its own `repository`, `usecase`, and `delivery` layers.
-   **`internal/config`**: Manages application configuration loading from environment variables or files.
-   **`internal/domain`**: Defines the core business entities (models), interfaces for repositories and use cases, and DTOs (Data Transfer Objects). This layer is independent of any specific technology.
-   **`internal/infrastructure`**: Provides implementations for external services and technologies, such   as database connections, message bus, mail sender, and tracing.
-   **`internal/presentation`**: Handles external interactions, including REST API routing, middleware, and message bus listeners.
    -   `rest`: Contains HTTP routing and middleware for authentication, error handling, and tracing.
    -   `messaging`: Contains listeners for the internal message bus.
-   **`internal/utils`**: Common utility functions, such as error handling and HTML template processing.

## Dependencies

The project is built with Go and utilizes several key libraries:

-   **`github.com/gin-gonic/gin`**: Web framework for building REST APIs.
-   **`gorm.io/gorm`**: ORM (Object-Relational Mapper) for database interactions.
-   **`gorm.io/driver/postgres`**: PostgreSQL driver for GORM.
-   **`github.com/spf13/viper`**: For configuration management.
-   **`github.com/uber/jaeger-client-go`**: For distributed tracing with Jaeger.
-   **`gopkg.in/gomail.v2`**: For sending emails.
-   **`github.com/Masterminds/squirrel`**: SQL query builder.
-   **`github.com/google/uuid`**: For generating UUIDs.
-   **`github.com/go-playground/validator/v10`**: For request validation.

## API Endpoints

All API endpoints are prefixed with `/api/v1`.

### Loan Management

These endpoints require authentication with `RoleEmployee`.

-   **`POST /api/v1/loan`**
    -   **Description:** Creates a new loan.
    -   **Authentication:** Employee
-   **`GET /api/v1/loan`**
    -   **Description:** Lists all loans.
    -   **Authentication:** Employee
-   **`GET /api/v1/loan/:id`**
    -   **Description:** Retrieves details of a specific loan by ID.
    -   **Authentication:** Employee
-   **`PATCH /api/v1/loan/:id/reject`**
    -   **Description:** Rejects a loan by ID.
    -   **Authentication:** Employee
-   **`PATCH /api/v1/loan/:id/approve`**
    -   **Description:** Approves a loan by ID.
    -   **Authentication:** Employee
-   **`PATCH /api/v1/loan/:id/disburse`**
    -   **Description:** Disburses a loan by ID.
    -   **Authentication:** Employee

### Investment Management

These endpoints require authentication with `RoleInvestor`.

-   **`POST /api/v1/investment`**
    -   **Description:** Adds a new investment to a loan.
    -   **Authentication:** Investor

### Public Endpoints

These endpoints do not require authentication.

-   **`GET /api/v1/loan/agreement/file/:loan_id`**
    -   **Description:** Retrieves the loan agreement file for a given loan ID.
-   **`GET /api/v1/investment/agreement/file/:investment_id`**
    -   **Description:** Retrieves the investment agreement file for a given investment ID.

## API Documentation

Detailed API documentation can be found on Postman: [https://documenter.getpostman.com/view/3187497/2sB3HrnHwc](https://documenter.getpostman.com/view/3187497/2sB3HrnHwc)

## Requirements

-   Go (version 1.24.2 or higher)
-   Docker (recommended for easy setup of PostgreSQL and Jaeger)

## How to Run

### Using Docker Compose (Recommended)

1.  **Ensure Docker is running:** Make sure Docker Desktop or Docker Engine is running on your system.

2.  **Build and run the services:**
    ```bash
    docker-compose up --build -d
    ```

3.  **Access the application and tools:**
    -   The API server will be accessible at `http://localhost:8081`.
    -   Mailhog web UI: `http://localhost:8025`
    -   Jaeger web UI: `http://localhost:16686`

4.  **Stop the services:**
    To stop and remove the containers, networks, and volumes created by `docker-compose up`:
    ```bash
    docker-compose down
    ```
    To stop only the containers but keep them for later restart:
    ```bash
    docker-compose stop
    ```

### Without Docker Compose

1.  **Build the application:**
    ```bash
    go build -o main cmd/api/main.go
    ```

2.  **Run the application:**
    ```bash
    ./main
    ```
    Alternatively, you can run directly:
    ```bash
    go run cmd/api/main.go
    ```

The API server will start on the port specified in your `.env` file (default: `8081`).

## Environment Variables

The application uses the following environment variables, typically loaded from a `.env` file:

-   `APP_PORT`: Port for the HTTP server (e.g., `8081`).
-   `POSTGRES_HOST`: PostgreSQL host.
-   `POSTGRES_PORT`: PostgreSQL port.
-   `POSTGRES_USERNAME`: PostgreSQL username.
-   `POSTGRES_PASSWORD`: PostgreSQL password.
-   `POSTGRES_DATABASE`: PostgreSQL database name.
-   `MAIL_HOST`: SMTP server host for sending emails.
-   `MAIL_PORT`: SMTP server port.
-   `MAIL_USERNAME`: SMTP username.
-   `MAIL_PASSWORD`: SMTP password.
-   `JAEGER_HOST`: Jaeger agent host.
-   `JAEGER_PORT`: Jaeger agent port.
-   `JAEGER_SERVICE_NAME`: Jaeger service name.

## Database Migrations

The project uses `golang-migrate` for managing database migrations.
-   `make migration-apply DRIVER=postgres`: Applies all pending migrations.
-   `make migration-rollback DRIVER=postgres`: Rolls back the last applied migration.
-   `make migration-create DRIVER=postgres NAME=<migration_name>`: Creates new up/down migration files.