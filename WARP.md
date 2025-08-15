# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Development Setup & Common Commands

### Local Development Environment

**Prerequisites:**
- Go 1.23.0+
- Docker and Docker Compose
- PostgreSQL 15+ (for local development)
- Redis 7+ (for local development)

### Build & Run Commands

```bash
# Build the application
go build -o bin/minimart ./cmd/server

# Run locally (requires PostgreSQL and Redis running)
go run ./cmd/server/main.go

# Run with Docker Compose (recommended for local development)
docker-compose up --build

# Run in background
docker-compose up -d
```

### Testing Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test package
go test ./internal/order/...

# Run integration tests
go test ./internal/integration/...

# Run tests with coverage
go test -cover ./...

# Run a specific test
go test -run TestOrderAccept ./internal/order/
```

### Database Operations

```bash
# Database migrations run automatically on startup
# Migrations are located in: ./migrations/

# To manually run migrations (if needed):
# Install goose: go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir migrations postgres "postgres://minimart:secret@localhost:5432/minimart_dev" up

# Create new migration
goose -dir migrations create migration_name sql
```

### Docker Commands

```bash
# Full environment startup
docker-compose up --build

# View logs
docker-compose logs -f app

# Restart just the app
docker-compose restart app

# Clean rebuild
docker-compose down && docker-compose up --build

# Database shell access
docker-compose exec db psql -U minimart -d minimart_dev
```

## Spec-Then-Code Development Process

This project follows a **spec-then-code** approach where detailed specifications are written before implementation.

### Specification Files
- **Location**: `/specs/` directory
- **Format**: Markdown files with detailed technical specifications
- **Purpose**: Define features completely before coding begins

### Development Workflow
1. Write or review specification in `/specs/`
2. Break specification into phases (entity-first approach)
3. Implement domain entities with business logic first
4. Add comprehensive tests for domain logic
5. Build infrastructure layers (repositories, handlers)
6. Integrate and test end-to-end

### Key Specification Files
- `/specs/order-fulfillment-flow.md` - Core order processing workflow
- `/docs/phase2-summary.md` - Entity-first development summary

### MCP Integration
The project has Model Context Protocol (MCP) integration for accessing external documentation:
- Access to spec-then-code methodology documentation
- GitHub repository searches for code patterns
- Documentation fetching for implementation guidance

## Domain-Driven Design Architecture

### Bounded Contexts
The application is organized around business domains:
- **Order Context** (`/internal/order/`) - Order management and fulfillment
- **Menu Context** (`/internal/menu/`) - Menu items and inventory
- **Merchant Context** (`/internal/merchant/`) - Merchant management
- **User Context** (`/internal/user/`) - Customer authentication and profiles

### Rich Domain Entities
Business logic lives in domain entities, not in services:

**Order Aggregate (`/internal/order/entity.go`)**:
```go
// Business logic methods that return domain events
func (o *Order) Accept(estimatedMinutes int, acceptedBy uuid.UUID) ([]DomainEvent, error)
func (o *Order) Reject(reason string, rejectedBy uuid.UUID) ([]DomainEvent, error)
func (o *Order) StartPreparing(preparedBy uuid.UUID) ([]DomainEvent, error)
```

**Key Principles**:
- Entities encapsulate state and behavior
- All business rules enforced at entity level
- State transitions controlled by domain logic
- Private fields with public getters
- Domain events emitted for integration

### Value Objects
**Money** (`/internal/order/value_objects.go`):
- Bitcoin-based pricing using Satoshis as base unit
- Immutable value object with arithmetic operations
- Smart display formatting (sats, mBTC, BTC)

**Address**:
- Immutable address representation
- Validation on construction
- Formatted string representation

**TimeWindow**:
- Estimated delivery/preparation windows
- Buffer calculations for realistic estimates

### Domain Events
```go
// Examples from /internal/order/domain_events.go
type OrderPlacedEvent struct {
    OrderID     uuid.UUID
    CustomerID  uuid.UUID
    MerchantID  uuid.UUID
    // ...
}

type OrderAcceptedEvent struct {
    OrderID       uuid.UUID
    EstimatedTime time.Time
    // ...
}
```

### Order State Machine
Valid state transitions are enforced at the entity level:
```
PENDING → ACCEPTED → PREPARING → READY → COMPLETED
        ↓         ↓          ↓       ↓
      REJECTED  CANCELLED  CANCELLED CANCELLED
```

## Technical Implementation Patterns

### Bitcoin Pricing System
All monetary values use Bitcoin with Satoshis as the base unit:
```go
// Create money values
price := NewMoney(50000)                    // 50,000 Satoshis
btcPrice := NewMoneyFromBTC(0.001)         // 0.001 BTC
mbtcPrice := NewMoneyFromMilliBTC(1.0)     // 1.0 mBTC

// Display formatting
fmt.Println(NewMoney(5000))                // "5000 sats"
fmt.Println(NewMoney(150000))              // "1.500 mBTC"
fmt.Println(NewMoney(50000000))            // "0.50000000 BTC"
```

### Authentication & Authorization
- JWT-based authentication using middleware in `/internal/shared/middleware/auth.go`
- Bearer token validation
- User context extraction for request handling
- Merchant-specific authorization for order management

### Event-Driven Architecture
Redis-based event bus for real-time updates:
```go
// Event publishing (from use cases)
events, err := order.Accept(30, merchantID)
for _, event := range events {
    eventBus.Publish(context.Background(), event)
}

// Event subscription (from main.go)
go func() {
    pubsub := redisClient.Subscribe(context.Background(), user.UserCreatedTopic)
    // Handle events...
}()
```

### Error Handling Patterns
- Domain errors defined as constants in entity files
- Wrapped errors with context using `fmt.Errorf`
- Validation at entity boundaries
- HTTP error responses with appropriate status codes

### Testing Strategy
**Entity-First Testing**:
- Pure domain logic tests (no infrastructure dependencies)
- Fast unit tests for business rules
- Integration tests for cross-module interactions
- Test files: `*_test.go` alongside implementation

**Test Examples**:
```go
// Pure domain testing
func TestOrderAccept(t *testing.T) {
    // No mocks, no databases - just business logic
    order := createTestOrder()
    events, err := order.Accept(30, merchantID)
    assert.NoError(t, err)
    assert.Equal(t, OrderStatusAccepted, order.Status())
}
```

## Infrastructure & Dependencies

### Database Schema
- PostgreSQL 15+ with migrations in `/migrations/`
- Automatic migration on application startup
- Connection pooling with pgxpool

**Key Tables**:
- `users` - Customer authentication
- `orders` - Order data with status tracking
- `order_items` - Order line items with price snapshots
- `merchants` - Merchant information
- `menu_items` - Menu with pricing and availability

### Redis Configuration
- Event bus for real-time notifications
- Pub/sub pattern for decoupled communication
- Connection configuration via environment variables

### Environment Configuration
**Local Development** (`.env` file):
```env
PORT=3000
DATABASE_URL=postgres://minimart:secret@localhost:5432/minimart_dev
REDIS_URL=localhost:6379
JWT_SECRET=your-secret-key
```

**Docker Configuration** (`docker-compose.yml`):
- Containerized PostgreSQL and Redis
- Volume persistence for data
- Environment variable injection

### Configuration Management
Uses Viper for configuration management:
- Environment variables take precedence
- YAML config file fallback (`config.yaml`)
- Structured configuration binding

## Order Fulfillment Architecture

### Core Order Flow
1. **Order Placement** - Customer creates order with menu items
2. **Merchant Notification** - Real-time event to merchant systems
3. **Order Acceptance/Rejection** - Merchant decision with estimated time
4. **Preparation Tracking** - Status updates through fulfillment
5. **Completion** - Final delivery/pickup completion

### State Machine Implementation
The order state machine is implemented in the Order aggregate:
```go
var validTransitions = map[OrderStatus][]OrderStatus{
    OrderStatusPending:        {OrderStatusAccepted, OrderStatusRejected, OrderStatusCancelled},
    OrderStatusAccepted:       {OrderStatusPreparing, OrderStatusCancelled},
    OrderStatusPreparing:      {OrderStatusReady, OrderStatusCancelled},
    // ...
}
```

### Menu Integration
- MenuItem entities with stock management
- Price snapshots in OrderItems (immutable pricing)
- Availability checks during order creation
- Stock reservation and release mechanisms

### Real-Time Updates
Event-driven notifications using Redis pub/sub:
- Customers receive order status updates
- Merchants receive new order notifications
- System-wide event broadcasting for integration

## Important Files & Directories

### Core Application
- `cmd/server/main.go` - Application entry point with dependency injection
- `config.yaml` - Default configuration
- `docker-compose.yml` - Local development environment

### Domain Layer
- `internal/order/entity.go` - Order aggregate with business logic
- `internal/order/value_objects.go` - Money, Address, TimeWindow value objects
- `internal/order/domain_events.go` - Domain events for integration
- `internal/menu/entity.go` - Menu item management

### Infrastructure
- `internal/shared/eventbus/` - Event bus implementations
- `internal/shared/middleware/` - HTTP middleware (auth, etc.)
- `migrations/` - Database schema migrations

### Testing
- `internal/integration/` - Cross-module integration tests
- `*_test.go` files - Unit tests alongside implementation

### Documentation
- `specs/` - Technical specifications
- `docs/` - Implementation phase summaries
- `LEARNING_PLAN.md` - Project learning objectives

## Development Best Practices

### Entity-First Development
1. Start with rich domain entities containing business logic
2. Write comprehensive tests for entity behavior
3. Build thin use case orchestration layers
4. Implement infrastructure last

### Testing Approach
- Test business logic without infrastructure dependencies
- Use integration tests for cross-module behavior
- Maintain fast feedback loops with pure unit tests

### Bitcoin Integration
- Always use Satoshis for internal calculations
- Provide conversion utilities for different units
- Display amounts in user-friendly formats

### Event-Driven Patterns
- Emit domain events from entity state changes
- Use events for loose coupling between modules
- Implement idempotent event handlers

### Error Handling
- Define domain-specific errors as constants
- Provide context with wrapped errors
- Validate at domain boundaries

This architecture supports rapid development while maintaining clean separation of concerns and testable business logic.
