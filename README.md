## UserVault – User API with DOB and Calculated Age

This project is a REST API built in Go to manage users with their **name** and **date of birth (dob)**.

- Data stored in PostgreSQL: `id`, `name`, `dob`.
- For read endpoints, the API also returns an `age` field that is **calculated dynamically** from `dob` using Go’s `time` package (it is not stored in the database).

---

## 1. Tech Stack

- Go (Go modules)
- GoFiber – HTTP web framework
- PostgreSQL – relational database
- SQLC – generates type‑safe DB access code from SQL
- pgx / pgxpool – PostgreSQL driver and connection pool
- Uber Zap – structured logging
- go-playground/validator – request validation library

---

## 2. Project Structure

cmd/server/main.go        # Application entry point

config/                   # Configuration (env variables like PORT, DATABASE_URL)
  config.go

db/
  migrations/             # SQL migrations (schema)
    001_create_users.sql
  sqlc/                   # SQLC config, SQL queries, and generated code
    db.go
    models.go
    queries/users.sql
    users.sql.go

internal/
  models/                 # Domain models + age calculation + unit tests
    user.go
    user_test.go
  repository/             # Data access layer wrapping SQLC Queries
    user_repository.go
  service/                # Business logic, validation, age calculation
    user_service.go
  handler/                # HTTP handlers (Fiber)
    user_handler.go
  routes/                 # Route registration
    routes.go
  middleware/             # Logging middleware (Zap)
    logger.go
  logger/                 # Zap logger initialization
    logger.go---

## 3. Database Schema

The `users` table is created by `db/migrations/001_create_users.sql`:

- `id SERIAL PRIMARY KEY`
- `name TEXT NOT NULL`
- `dob DATE NOT NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`

**Important**:

- `dob` is stored in the database.
- `age` is **not stored**; it is computed in Go when returning responses.

---

## 4. API Endpoints

All endpoints use JSON.

### 4.1 Create User – `POST /users`

**Request:**

{
  "name": "Alice",
  "dob": "1990-05-10"
}**Response (201 Created):**

{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}---

### 4.2 Get User by ID – `GET /users/:id`

**Response (200 OK):**

{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}`age` is calculated dynamically from `dob` and the current date.

---

### 4.3 Update User – `PUT /users/:id`

**Request:**

{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}**Response (200 OK):**

{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}---

### 4.4 Delete User – `DELETE /users/:id`

- **Response**: `204 No Content` on successful deletion.

---

### 4.5 List Users – `GET /users`

Supports **pagination** via query parameters:

- `limit` (default: `50`)
- `offset` (default: `0`)

Example:

GET /users?limit=10&offset=0**Response (200 OK):**

[
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10",
    "age": 34
  }
]---

## 5. Validation, Age Calculation, and Logging

### 5.1 Validation

- Implemented with `go-playground/validator` in the service layer.
- For create and update:
  - `name`: required, length constraints.
  - `dob`: required, must be in `YYYY-MM-DD` format.

Invalid input returns **400 Bad Request** with an error message.

---

### 5.2 Age Calculation

- Implemented in `internal/models/user.go` as `CalculateAge(dob, now time.Time) int`.
- Logic:
  - Computes the difference in years between `now` and `dob`.
  - If the birthday has not happened yet this year, subtracts 1.
  - If `dob` is in the future, returns 0.
- Covered by unit tests in `internal/models/user_test.go`.

---

### 5.3 Logging

- Uber Zap is used for structured logging.
- Logging middleware (`internal/middleware/logger.go`) logs:
  - HTTP method
  - Path
  - Status code
  - Request duration

Errors in handlers and services are also logged with Zap.

---

## 6. How to Run the Project

From the project root: `~/WebDev/JavaScript/UserVault`

### 6.1 Start PostgreSQL

brew services start postgresql@14### 6.2 Create the Database and Run Migrations

createdb uservault
psql uservault < db/migrations/001_create_users.sql### 6.3 Set Environment Variables

export DATABASE_URL="postgres://localhost:5432/uservault?sslmode=disable"
export PORT=8080  # optional, default is 8080### 6.4 Install Go Dependencies

go get github.com/gofiber/fiber/v2 \
       github.com/jackc/pgx/v5/pgxpool \
       github.com/jackc/pgx/v5 \
       github.com/jackc/pgconn \
       go.uber.org/zap \
       github.com/go-playground/validator/v10 \
       github.com/jmoiron/sqlx### 6.5 Install and Run SQLC

Install (choose one):

- Homebrew:

 
  brew install sqlc
  - Or via Go:

 
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  export PATH="$PATH:$(go env GOPATH)/bin"
  Generate the SQLC code:

sqlc generate### 6.6 Run Age Calculation Tests (Bonus)

go test ./internal/models### 6.7 Run the Server

go run ./cmd/serverThe API will listen on `http://localhost:8080` (or on `PORT` if set).

---

## 7. Quick Manual Testing (Examples)

### 7.1 Create User

curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'### 7.2 Get User with Age

curl http://localhost:8080/users/1### 7.3 Update User

curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'### 7.4 List Users with Pagination

curl "http://localhost:8080/users?limit=10&offset=0"### 7.5 Delete User

curl -X DELETE http://localhost:8080/users/1 -i


## 8. Quick Postman Testing (Examples)

1. Create User – POST /users
Method: POST
URL: http://localhost:8080/users
Headers:
Content-Type: application/json
Body → raw → JSON:

{
  "name": "Alice",
  "dob": "1990-05-10
}

Expected status: 201 Created
Expected response:

{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}

2. Get User by ID – GET /users/:id
Method: GET
URL: http://localhost:8080/users/1
Body: none
Expected status: 200 OK
Expected response (age depends on current date):

{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}

If the user does not exist, you should see:
{  "error": "user not found"}

3. Update User – PUT /users/:id
Method: PUT
URL: http://localhost:8080/users/1
Headers:
Content-Type: application/json
Body → raw → JSON:

{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}

Expected status: 200 OK
Expected response:

{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}

. List Users (with pagination) – GET /users
Method: GET
URL: http://localhost:8080/users?limit=10&offset=0
Body: none
Expected status: 200 OK
Expected response (example):

[
  {
    "id": 1,
    "name": "Alice Updated",
    "dob": "1991-03-15",
    "age": 34
  }
]

5. Delete User – DELETE /users/:id
Method: DELETE
URL: http://localhost:8080/users/1
Body: none
Expected status: 204 No Content
Expected response body: empty
After deleting, a GET /users/1 should return:

{
  "error": "user not found"
}
