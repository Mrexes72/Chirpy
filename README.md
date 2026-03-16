# Chirpy API 🐦

A Twitter-like social media REST API built with Go, featuring JWT authentication, refresh tokens, and webhook integrations.

## 🚀 Features

- **User Management**
  - User registration with secure password hashing (Argon2)
  - Email and password authentication
  - User profile updates
  - Chirpy Red premium membership

- **Authentication & Authorization**
  - JWT-based access tokens (1-hour expiration)
  - Refresh tokens (60-day expiration)
  - Token revocation
  - Bearer token authentication

- **Chirps (Posts)**
  - Create, read, and delete chirps
  - Profanity filtering
  - 140-character limit
  - Author-based filtering
  - Sortable by creation date (ascending/descending)

- **Webhooks**
  - Polka payment integration
  - Idempotent webhook handling
  - User upgrades to Chirpy Red

- **Admin Features**
  - Metrics tracking (file server hits)
  - Health check endpoint
  - Database reset (dev only)

## 📋 Prerequisites

- Go 1.22+
- PostgreSQL 13+
- [SQLC](https://sqlc.dev/) for SQL code generation
- [Goose](https://github.com/pressly/goose) for database migrations
- [Boot.dev CLI](https://github.com/bootdotdev/bootdev) (optional, for testing)

## 🛠️ Installation

1. **Clone the repository**
```bash
   git clone https://github.com/Mrexes72/Chirpy.git
   cd Chirpy
```

2. **Install dependencies**
```bash
   go mod download
```

3. **Install tools**
```bash
   # Install SQLC
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

   # Install Goose
   go install github.com/pressly/goose/v3/cmd/goose@latest
```

4. **Set up environment variables**
```bash
   cp .env.example .env
```

   Edit `.env` with your configuration:
```env
   DB_URL=postgres://postgres:password@localhost:5432/chirpy?sslmode=disable
   PLATFORM=dev
   JWT_SECRET=your-secret-key-here
```

   Generate a secure JWT secret:
```bash
   openssl rand -base64 64
```

5. **Set up the database**
```bash
   # Create database
   createdb chirpy

   # Run migrations
   cd sql/schema
   goose postgres "postgres://postgres:password@localhost:5432/chirpy?sslmode=disable" up
   cd ../..

   # Generate SQLC code
   sqlc generate
```

6. **Build and run**
```bash
   go build -o chirpy && ./chirpy
```

   The server will start on `http://localhost:8080`

## 📁 Project Structure
```
.
├── internal/
│   ├── auth/              # Authentication utilities
│   │   ├── authentication.go
│   │   └── jwt.go
│   └── database/          # SQLC generated code
│       ├── chirps.sql.go
│       ├── users.sql.go
│       ├── models.go
│       └── db.go
├── sql/
│   ├── queries/           # SQL queries for SQLC
│   │   ├── chirps.sql
│   │   ├── users.sql
│   │   └── refresh_tokens.sql
│   └── schema/            # Database migrations
│       ├── 001_users.sql
│       ├── 002_chirps.sql
│       ├── 003_users_hashed_password.sql
│       ├── 004_refresh_tokens.sql
│       └── 005_users_is_red.sql
├── handler_*.go           # HTTP handlers
├── main.go                # Application entry point
├── .env                   # Environment variables (not in Git)
└── sqlc.yaml             # SQLC configuration
```

## 🔌 API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/healthz` | Health check |
| GET | `/admin/metrics` | View metrics (HTML) |

### User Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/users` | Create user | No |
| PUT | `/api/users` | Update user | Yes (JWT) |
| POST | `/api/login` | Login | No |
| POST | `/api/refresh` | Refresh access token | Yes (Refresh Token) |
| POST | `/api/revoke` | Revoke refresh token | Yes (Refresh Token) |

### Chirp Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/chirps` | Create chirp | Yes (JWT) |
| GET | `/api/chirps` | Get all chirps | No |
| GET | `/api/chirps/{chirpID}` | Get chirp by ID | No |
| DELETE | `/api/chirps/{chirpID}` | Delete chirp | Yes (JWT, must be author) |

**Query Parameters for GET /api/chirps:**
- `author_id` - Filter by author UUID
- `sort` - Sort order: `asc` (default) or `desc`

### Webhook Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/polka/webhooks` | Polka payment webhooks |

### Admin Endpoints

| Method | Endpoint | Description | Platform |
|--------|----------|-------------|----------|
| POST | `/admin/reset` | Reset metrics and database | dev only |

## 📝 Example Usage

### Register a new user
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

### Create a chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"body": "Hello, world!"}'
```

### Get all chirps (sorted descending)
```bash
curl "http://localhost:8080/api/chirps?sort=desc"
```

### Get chirps by author
```bash
curl "http://localhost:8080/api/chirps?author_id=USER_UUID&sort=desc"
```

## 🧪 Testing

Run unit tests:
```bash
go test ./... -v
```

Run with coverage:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Run Boot.dev integration tests:
```bash
bootdev run TEST_ID
```

## 🔐 Security Features

- **Password Hashing**: Argon2id with salts
- **JWT Tokens**: HS256 signing with configurable secrets
- **Refresh Tokens**: 256-bit random tokens stored in database
- **Token Revocation**: Explicit revocation with timestamps
- **Authorization**: Resource-level access control
- **Input Validation**: Request body validation and sanitization
- **Profanity Filter**: Automatic content filtering

## 🏗️ Database Schema

### Users Table
```sql
- id (UUID, PK)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
- email (TEXT, UNIQUE)
- hashed_password (TEXT)
- is_red (BOOLEAN)
```

### Chirps Table
```sql
- id (UUID, PK)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
- body (TEXT)
- user_id (UUID, FK → users)
```

### Refresh Tokens Table
```sql
- token (TEXT, PK)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
- user_id (UUID, FK → users)
- expires_at (TIMESTAMP)
- revoked_at (TIMESTAMP, nullable)
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project was created as part of the [Boot.dev](https://boot.dev) backend development course.

## 🙏 Acknowledgments

- Built following [Boot.dev's Learn Web Servers course](https://boot.dev)
- Uses [SQLC](https://sqlc.dev/) for type-safe SQL
- Uses [Goose](https://github.com/pressly/goose) for migrations
- Authentication patterns inspired by industry best practices

## 📧 Contact

Per - [@Mrexes72](https://github.com/Mrexes72)

Project Link: [https://github.com/Mrexes72/Chirpy](https://github.com/Mrexes72/Chirpy)