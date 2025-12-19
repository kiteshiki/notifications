## Notifications API

A minimal Gin-based API scaffolded for future extension, with Swagger docs and PostgreSQL database support.

### Requirements

- Go 1.25.5
- PostgreSQL database (default: `fandom_notifications`)

### Structure

```
cmd/server/           # Application entrypoint
internal/config/      # Configuration loading
internal/database/    # Database connection and migrations
internal/middleware/  # HTTP middleware (API key auth)
internal/models/      # Data models
internal/repository/  # Data access layer
internal/service/     # Business logic layer
internal/server/      # Router and HTTP transport
internal/server/transport/  # Handlers and route registration
docs/                 # Swagger docs (generated)
```

### Database Setup

1. Create PostgreSQL database:

```sql
CREATE DATABASE fandom_notifications;
```

2. Configure database connection via environment variables:

```bash
# Option 1: Use DATABASE_URL (recommended)
export DATABASE_URL="postgres://user:password@localhost:5432/fandom_notifications?sslmode=disable"

# Option 2: Use individual variables
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=fandom_notifications
```

3. (Optional) Set master API key for programmatic key generation:

```bash
export MASTER_API_KEY="your-secure-master-key-here"
```

If not set, you can generate API keys using the CLI tool (see API Key Authentication section).

The following tables will be automatically created on server startup:

**bookmarks** table:

- `id` (SERIAL PRIMARY KEY)
- `user_id` (VARCHAR(255))
- `publication_id` (VARCHAR(255))
- `chapter_id` (VARCHAR(255))
- `image` (VARCHAR(255))
- `chapter` (VARCHAR(255))
- `volume` (VARCHAR(255))
- `name` (VARCHAR(255))
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

**api_keys** table:

- `id` (SERIAL PRIMARY KEY)
- `key` (VARCHAR(255) UNIQUE)
- `name` (VARCHAR(255))
- `active` (BOOLEAN)
- `created_at` (TIMESTAMP)
- `last_used_at` (TIMESTAMP)

**request_logs** table:

- `id` (SERIAL PRIMARY KEY)
- `method` (VARCHAR(10))
- `path` (VARCHAR(500))
- `query_params` (TEXT)
- `status_code` (INTEGER)
- `ip_address` (VARCHAR(45))
- `user_agent` (TEXT)
- `api_key` (VARCHAR(255))
- `response_time_ms` (BIGINT)
- `created_at` (TIMESTAMP)

### Development

- Install deps and generate Swagger:

```bash
make deps tidy swag
```

- Run the server:

```bash
make run
```

- Build:

```bash
make build
```

### API Key Authentication

All API endpoints require an API key to be provided as a query parameter `api`. The `/api-keys` endpoint is protected by a master API key.

#### Option 1: Generate Initial API Key via CLI (Recommended)

Generate your first API key using the CLI tool:

```bash
make generate-key
# or with custom name:
go run ./cmd/generate-key -name "My First Key"
```

This will output the generated key. Save it securely!

#### Option 2: Use Master API Key

Set a master API key in your environment:

```bash
export MASTER_API_KEY="your-secure-master-key-here"
```

Then generate API keys via the endpoint:

```bash
curl -X POST "http://localhost:8080/api-keys?api=YOUR_MASTER_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "My API Key"}'
```

Response:

```json
{
  "key": "abc123...",
  "name": "My API Key",
  "created_at": "2025-12-16T09:00:00Z"
}
```

#### Using API Keys

**With Valid API Key:**

```bash
curl "http://localhost:8080/hello?api=YOUR_API_KEY"
```

**Without API Key (403 Forbidden):**

```bash
curl http://localhost:8080/hello
# Returns: {"error": "API key is required. Provide it as 'api' query parameter."}
```

**Invalid API Key (403 Forbidden):**

```bash
curl "http://localhost:8080/hello?api=invalid_key"
# Returns: {"error": "Invalid or inactive API key"}
```

### Authentication

The service supports two authentication methods:

1. **Cookie-based authentication** (recommended for web UI)
   - Visit `/auth` to enter your API key
   - The key is stored in a secure HTTP-only cookie
   - Automatically used for all subsequent requests
   - Valid for 7 days

2. **Query parameter authentication** (for API clients)
   - Include `?api=YOUR_KEY` in the URL
   - Works alongside cookie authentication (cookies take precedence)

### Endpoints

**Public Endpoints (no API key required):**

- `GET /swagger/*` - Swagger UI documentation
- `GET /auth` - Authentication page (HTML form)
- `POST /auth/set` - Set API key cookie (JSON)
- `POST /auth/clear` - Clear API key cookie

**Admin Endpoints (require master API key via `MASTER_API_KEY` env var):**

- `POST /api-keys` - Generate a new API key (cookie or query param)
- `GET /dashboard` - View request logs dashboard (HTML)
- `GET /dashboard/logs` - Get request logs (JSON API)
- `GET /dashboard/stats` - Get log statistics (JSON API)

**Protected Endpoints (require regular API key):**

- `GET /hello` → `{"message":"Hello, World!"}`

All protected endpoints accept authentication via cookie (set via `/auth`) or `api` query parameter.

### Request Logging

All API requests are automatically logged to the database with the following information:
- HTTP method and path
- Query parameters
- Status code
- IP address
- User agent
- API key (partially masked)
- Response time in milliseconds
- Timestamp

Logs are stored asynchronously to avoid impacting request performance.

### Dashboard

Access the dashboard at `/dashboard` (after authenticating at `/auth`) to view:
- Real-time request logs
- Statistics (total requests, average response time)
- Status code distribution
- Top paths and methods
- Filtering by method, status code, path, and date range
- Auto-refresh every 30 seconds

**Dashboard API Endpoints:**

```bash
# First authenticate (or use cookie from /auth page)
# Get logs with filters
curl "http://localhost:8080/dashboard/logs?api=MASTER_KEY&limit=100&offset=0&method=GET&status_code=200"

# Get statistics
curl "http://localhost:8080/dashboard/stats?api=MASTER_KEY&start_date=2025-12-01T00:00:00Z&end_date=2025-12-17T00:00:00Z"
```

**Using the Web Interface:**

1. Visit `http://localhost:8080/auth`
2. Enter your master API key
3. You'll be redirected to the dashboard
4. The API key is stored in a cookie for 7 days
5. Use the "Logout" button to clear the cookie
