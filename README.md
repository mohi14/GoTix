# GoTix

GoTix is a Go REST API for managing seminars, concerts, and other ticketed
events. Organizers can publish events, authenticated users can reserve tickets,
and administrators can manage user roles.

## Features

- Email/password registration and login
- Google ID-token login
- JWT authentication
- Role-based authorization:
  - `user`
  - `organizer`
  - `admin`
  - `super_admin`
- Event creation, listing, viewing, and updating
- Transaction-safe ticket booking with inventory locking
- User booking history
- PostgreSQL persistence with GORM

## Requirements

- Go 1.27 or newer
- PostgreSQL

## Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

Set the values in `.env`:

```env
DATABASE_URL=postgresql://user:password@localhost:5432/gotix?sslmode=disable
PORT=8080
JWTSECRET=replace-with-a-long-random-secret
GOOGLE_CLIENT_ID=your-google-oauth-client-id
SUPER_ADMIN_EMAIL=admin@example.com
```

`SUPER_ADMIN_EMAIL` is optional. When that email registers with a password, the
account is created with the `super_admin` role. New accounts and Google accounts
otherwise receive the `user` role.

## Running

Install dependencies and start the API:

```bash
go mod download
go run ./cmd
```

The server listens on `http://localhost:8080` by default.

OpenAPI documentation is available at:

- Landing page: `http://localhost:8080/`
- System status page: `http://localhost:8080/health`
- Machine-readable health check: `http://localhost:8080/healthz`
- Swagger UI: `http://localhost:8080/api/v1/docs`
- OpenAPI specification: `http://localhost:8080/api/v1/docs/openapi.yaml`

The Swagger UI is served by the application and loads the official Swagger UI
assets from its CDN. The OpenAPI specification uses a relative server URL, so
it automatically targets the current host in local, staging, and production.
Use the `Authorize` button to test JWT-protected endpoints.

Run checks:

```bash
go test ./...
go vet ./...
```

Database tables are created or updated automatically when the server starts.

## API

The API is versioned under `/api/v1`.

### Public endpoints

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/auth/register` | Register with email and password |
| `POST` | `/api/v1/auth/login` | Login with email and password |
| `POST` | `/api/v1/auth/google` | Login with a Google ID token |
| `GET` | `/api/v1/events` | List events |
| `GET` | `/api/v1/events/:id` | Get an event |

### Authenticated endpoints

Send the JWT in the request header:

```http
Authorization: Bearer <token>
```

| Method | Endpoint | Required role | Description |
| --- | --- | --- | --- |
| `GET` | `/api/v1/auth/me` | Any authenticated user | Get the current user |
| `POST` | `/api/v1/bookings` | Any authenticated user | Book tickets |
| `GET` | `/api/v1/bookings/me` | Any authenticated user | List personal bookings |
| `POST` | `/api/v1/events` | `organizer`, `admin`, `super_admin` | Create an event |
| `PUT` | `/api/v1/events/:id` | `organizer`, `admin`, `super_admin` | Update an event |
| `PATCH` | `/api/v1/auth/users/:id/role` | `admin`, `super_admin` | Change a user's role |

Admins cannot assign the `admin` or `super_admin` roles. Only a
`super_admin` can assign those roles.

## Example requests

Register:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alex","email":"alex@example.com","password":"password123"}'
```

Create an event as an organizer:

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "title":"Go Conference",
    "description":"A backend engineering conference",
    "location":"Berlin",
    "starts_at":"2026-12-01T18:00:00Z",
    "total_tickets":100,
    "price":50
  }'
```

Book tickets:

```bash
curl -X POST http://localhost:8080/api/v1/bookings \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"event_id":1,"quantity":2}'
```

## Project structure

```text
cmd/                 Application entrypoint
internal/auth/       JWT creation and validation
internal/config/     Environment and database configuration
internal/domain/     User, event, and booking modules
internal/middlewares/Authentication and role authorization
internal/server/     Echo server and route registration
```

## Payment status

Bookings currently calculate the total price and reserve inventory
transactionally. A payment provider such as Stripe or PayPal is not integrated
yet; payment processing should be added before treating a booking as paid.
