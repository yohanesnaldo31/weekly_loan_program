# Weekly Loan Program

A Go-based loan management service with HTTP endpoints for loan requests, payments, delinquency checks, and a cron job that updates delinquent loan statuses.

## Features

- Create a loan and its weekly billings
- Pay an active loan
- Retrieve a user's loans
- Check whether a user has a delinquent loan
- Retrieve a user's current outstanding amount
- Trigger delinquency updates via an internal endpoint
- Run a scheduled delinquency check at midnight with a cron job

## Requirements
- The interests of the loan is a flat 10% of the loaned amount
- When user request loans, the installment will be in weeks
- User has an option to pay the weekly bill or none at all
- When paying bill, users have to pay the full amount, can't be lower and higher
- If user failed to pay for the week billing, the amount user have to pay for the next billing will be accumulation of the previous unpaid billings
- If the user failed to pay for 3 billings (3 weeks worth of bill), user's loan will be flagged as *delinquent*

## Constraints
- Users which have an ongoing loan, unable to request additional loan
- Minimum number of loan is 10000
- Number of installment in weeks can't be > 261 weeks (5 years)

## High Level Design
![High Level Design](https://drive.google.com/uc?export=view&id=1cUcT8mq-vzEJEaq9r2XQmlPJV6ATplsb)

## Architecture

The project follows a layered structure:

- app: application bootstrap and dependency wiring
- cmd/http: HTTP server entrypoint
- cmd/cron: cron scheduler entrypoint
- server/httphandler: HTTP handlers and routes
- server/cronhandler: cron job handlers
- service/loan: service-layer business logic
- usecase/loan: usecase orchestration and business rules
- repo/db: PostgreSQL repository implementation
- repo/cache: Redis cache implementation

## Prerequisites

- Go 1.26+ (only needed if building from source; skip if using a prebuilt binary)
- Docker and Docker Compose
- PostgreSQL and Redis (provided via Docker Compose)

## Local Setup

1. Start the infrastructure services:
   ```bash
   docker compose up -d
   ```

2. Run the HTTP server:
   ```bash
   make run
   ```

4. Run the cron worker:
   ```bash
   make run-cron
   ```

### Running from a prebuilt binary

If you don't have Go installed, you can run the released `weekly_loan` / `weekly_loan_cron` binaries directly (matching your OS/architecture) instead of steps 2 and 4:

```bash
./weekly_loan
./weekly_loan_cron
```

The infrastructure services from step 1 must still be running.

## HTTP API

A Postman collection with ready-to-use requests for all endpoints below is available at [weekly_loan_program.postman_collection.json](weekly_loan_program.postman_collection.json).

### Get user loans
- GET /user/loans?user_id=1

### Check delinquency for a user
- GET /user/delinquent?user_id=1

### Get a user's current outstanding amount
- GET /user/outstanding?user_id=1

### Request a loan
- POST /request/loan
- Body example:
  ```json
  {
    "user_id": 1,
    "loan_amount": 100000,
    "installment_in_weeks": 4
  }
  ```

### Pay a loan
- POST /pay/loan
- Body example:
  ```json
  {
    "user_id": 1,
    "payment_amount": 25000,
    "payment_time": "2026-01-01" // to simulate time for testing, optional
  }
  ```

### Trigger delinquent loan check internally
- POST /internal/trigger/delinquent_check
- Body example:
  ```json
  {
    "time": "2026-01-01"
  }
  ```

## Cron Job

The cron worker runs the delinquency check on a schedule and updates loans that have been inactive for at least three weeks.

Current schedule:
- Midnight every day: `0 0 0 * * *`

## Database

The service uses PostgreSQL and Redis. The initial schema is created from [sql/init/01_create_tables.sql](sql/init/01_create_tables.sql).

