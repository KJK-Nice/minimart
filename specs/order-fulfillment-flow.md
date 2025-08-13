# Spec-Then-Code: Order Fulfillment Flow

## Overview
This specification defines the implementation of a complete order fulfillment flow for the Minimart B2C ordering system. The feature enables merchants to receive, manage, and fulfill orders while providing real-time status updates to customers. This transforms the current one-way order placement into a complete, bidirectional order lifecycle management system.

### Output Location
Implementation will be distributed across multiple modules:
- `/internal/order/` - Enhanced order entity and use cases
- `/internal/merchant/` - New merchant order management capabilities
- `/internal/notifications/` - Order status change notifications
- `/internal/shared/eventbus/` - Event-driven status updates

## Problem Statement
Currently, when customers place orders through Minimart:
1. Orders are created but merchants have no way to know about them
2. Customers cannot track their order status or know if it's accepted
3. There's no mechanism for merchants to manage order fulfillment
4. Orders lack critical information (merchant assignment, pricing, delivery details)

This makes the system incomplete for real-world use - orders disappear into a void.

## Background
The current implementation has basic order placement (`PlaceOrder`) but lacks:
- Merchant assignment to orders
- Order acceptance/rejection workflow
- Status progression tracking
- Real-time updates to stakeholders
- Order total calculation
- Delivery/pickup options

### Justification
Without order fulfillment flow:
- Merchants cannot operate their business (they don't receive orders)
- Customers lose trust (no confirmation or updates)
- The platform cannot facilitate actual food delivery
- The system is just a demo, not a usable product

This feature is essential to achieve SLC (Simple, Lovable, Complete) status.

## Acceptance Criteria
1. **Order Placement Enhancement**
   - Orders must be associated with a specific merchant
   - Order totals must be calculated from menu item prices
   - Delivery method (pickup/delivery) must be specified
   - Delivery address must be captured for delivery orders

2. **Merchant Order Management**
   - Merchants can view pending orders for their business
   - Merchants can accept or reject orders with reasons
   - Merchants can update order status through fulfillment stages
   - Merchants receive real-time notifications of new orders

3. **Customer Order Tracking**
   - Customers can view current order status
   - Customers receive real-time updates when status changes
   - Customers can see estimated completion time
   - Customers can view order history

4. **Status Workflow**
   - Orders progress through defined states: PENDING → ACCEPTED → PREPARING → READY → COMPLETED
   - Alternative flows: PENDING → REJECTED, or any state → CANCELLED
   - Status changes trigger events for notifications

### Design Goals
- **Simple**: Focus on core fulfillment flow, defer complex features (payments, ratings)
- **Complete**: End-to-end order lifecycle that actually works
- **Lovable**: Real-time updates reduce customer anxiety
- **Event-Driven**: Use existing event bus for loose coupling
- **Testable**: Clear state transitions with business logic in use cases

## Analysis of the Issue

### Root Cause Analysis
The MVP approach focused on technical capability (can place an order) rather than business completeness (can fulfill an order). The order entity lacks merchant context, and there's no merchant-facing API for order management.

### Research Findings
Examining the codebase:
1. Order entity has minimal fields - missing MerchantID, TotalAmount, DeliveryMethod
2. OrderStatus enum exists but only has 4 states - needs more granular fulfillment states
3. Event bus infrastructure exists and is tested - can be leveraged for real-time updates
4. Merchant module exists but has no order management capabilities
5. Authentication middleware exists - can secure merchant endpoints

## Benefits
- **For Customers**: Trust through transparency, reduced anxiety with status updates
- **For Merchants**: Operational capability, business management tools
- **For Platform**: Complete transaction flow, foundation for monetization
- **For Development**: Clear domain boundaries, event-driven extensibility

## Affected Code

### Current Implementation
- `Order` entity: Basic structure with status field
- `OrderUsecase`: Only has `PlaceOrder` method
- `OrderStatus`: Limited enum with 4 states
- No merchant order management endpoints
- No order-merchant relationship

### List of Files
**To Modify:**
- `/internal/order/entity.go` - Enhance Order struct
- `/internal/order/usecase.go` - Add order management methods
- `/internal/order/repository.go` - Add query methods
- `/internal/order/postgres_repository.go` - Implement new queries
- `/internal/order/handler.go` - Add status update endpoints
- `/internal/merchant/handler.go` - Add order management endpoints
- `/internal/merchant/usecase.go` - Add merchant order methods

**To Create:**
- `/internal/order/events.go` - Order status change events
- `/internal/merchant/order_management.go` - Merchant order use cases
- `/internal/notifications/order_subscriber.go` - Handle order events
- `/migrations/002_order_enhancements.sql` - Database schema updates

### Relevant Code Snippets
Current Order entity lacks merchant context:
```go
type Order struct {
    ID         uuid.UUID
    CustomerID uuid.UUID
    Items      []OrderItem
    Status     OrderStatus
    CreatedAt  time.Time
}
```

### Relevant Object/Component Hierarchies
- Order Aggregate: Order → OrderItems → MenuItems
- Event Flow: Order Usecase → Event Bus → Notification Subscribers
- API Routes: /orders (customer) and /merchant/orders (merchant)

## Proposed Solutions

### Solution 1: Enhance Existing Order Module (Recommended)
Extend the current order module with merchant context and management capabilities:
- Add MerchantID and pricing to Order entity
- Create bidirectional API (customer and merchant facing)
- Use event bus for real-time updates
- Implement status state machine in use case layer

**Pros**: Leverages existing code, maintains module boundaries, uses established patterns
**Cons**: Requires database migration, more complex than starting fresh

### Solution 2: Create Separate Fulfillment Module
Build a new module specifically for order fulfillment:
- Keep Order module for placement only
- New Fulfillment module for merchant operations
- Sync via events

**Pros**: Clean separation of concerns
**Cons**: More complexity, potential data consistency issues

## Analysis of Changes Needed

### Entity Enhancements
1. Order entity needs: MerchantID, TotalAmount, DeliveryMethod, DeliveryAddress, EstimatedTime, UpdatedAt
2. OrderStatus needs states: PENDING, ACCEPTED, REJECTED, PREPARING, READY, OUT_FOR_DELIVERY, COMPLETED, CANCELLED
3. New StatusHistory entity for audit trail

### Use Case Enhancements
1. Customer: GetOrderStatus, GetOrderHistory, CancelOrder
2. Merchant: GetPendingOrders, AcceptOrder, RejectOrder, UpdateOrderStatus
3. System: CalculateOrderTotal, ValidateOrderItems, NotifyStatusChange

### Repository Enhancements
1. Queries: FindByMerchant, FindPendingByMerchant, FindByCustomer
2. Updates: UpdateStatus, UpdateEstimatedTime

### API Enhancements
1. Customer endpoints: GET /orders/:id/status, GET /orders/history
2. Merchant endpoints: GET /merchant/orders, PUT /merchant/orders/:id/status
3. WebSocket endpoint for real-time updates (future enhancement)

## Implementation Plan

### Executable Steps for AI Agent

- [ ] Step 1: Update Order entity with merchant context and pricing
  - File: `/internal/order/entity.go`
  - Add MerchantID, TotalAmount, DeliveryMethod, DeliveryAddress fields
  - Expand OrderStatus enum with fulfillment states
  - Add StatusHistory type for audit trail

- [ ] Step 2: Create database migration for order enhancements
  - File: `/migrations/002_order_enhancements.sql`
  - Add merchant_id, total_amount, delivery_method, delivery_address columns
  - Create order_status_history table
  - Add indexes for merchant queries

- [ ] Step 3: Define order event types
  - File: `/internal/order/events.go`
  - Create OrderPlacedEvent, OrderStatusChangedEvent, OrderCompletedEvent
  - Include all necessary data for subscribers

- [ ] Step 4: Enhance order repository interface and implementation
  - Files: `/internal/order/repository.go`, `/internal/order/postgres_repository.go`
  - Add FindByMerchantID, FindPendingByMerchantID methods
  - Add UpdateStatus method with history tracking
  - Implement efficient merchant order queries

- [ ] Step 5: Extend order use case with management methods
  - File: `/internal/order/usecase.go`
  - Add GetOrderByID, UpdateOrderStatus, GetMerchantOrders methods
  - Implement status validation state machine
  - Publish events on status changes
  - Calculate order totals from menu items

- [ ] Step 6: Create merchant order management use case
  - File: `/internal/merchant/order_management.go`
  - Implement AcceptOrder, RejectOrder, UpdateOrderProgress methods
  - Add business logic for time estimates
  - Validate merchant ownership of orders

- [ ] Step 7: Add customer order tracking endpoints
  - File: `/internal/order/handler.go`
  - GET /orders/:id endpoint for order details and status
  - GET /orders endpoint for order history
  - Add proper error handling and auth checks

- [ ] Step 8: Add merchant order management endpoints
  - File: `/internal/merchant/handler.go`
  - GET /merchant/orders endpoint for pending orders
  - PUT /merchant/orders/:id/accept endpoint
  - PUT /merchant/orders/:id/reject endpoint
  - PUT /merchant/orders/:id/status endpoint
  - Ensure merchant auth middleware is applied

- [ ] Step 9: Create order notification subscriber
  - File: `/internal/notifications/order_subscriber.go`
  - Subscribe to order events
  - Log notifications (email/SMS integration deferred)
  - Handle different event types appropriately

- [ ] Step 10: Update main server to wire new components
  - File: `/cmd/server/main.go`
  - Register new endpoints
  - Initialize order event publishers
  - Start notification subscribers

- [ ] Step 11: Write integration tests for order fulfillment flow
  - File: `/internal/order/fulfillment_test.go`
  - Test complete flow: place → accept → prepare → complete
  - Test rejection flow
  - Test concurrent status updates
  - Verify event publishing

- [ ] Step 12: Update API documentation
  - File: `/docs/api/order-fulfillment.md`
  - Document new endpoints
  - Provide example requests/responses
  - Document status workflow

### Required Code Changes

1. **Order Entity Enhancement**:
```go
type Order struct {
    ID              uuid.UUID
    CustomerID      uuid.UUID
    MerchantID      uuid.UUID       // NEW
    Items           []OrderItem
    Status          OrderStatus
    TotalAmount     int64           // NEW: in cents
    DeliveryMethod  DeliveryMethod  // NEW
    DeliveryAddress *Address        // NEW: optional
    EstimatedTime   *time.Time      // NEW: optional
    CreatedAt       time.Time
    UpdatedAt       time.Time       // NEW
    StatusHistory   []StatusChange  // NEW
}
```

2. **Extended Status Enum**:
```go
const (
    PENDING OrderStatus = iota
    ACCEPTED
    REJECTED
    PREPARING
    READY_FOR_PICKUP
    OUT_FOR_DELIVERY
    COMPLETED
    CANCELLED
)
```

## Test-Driven Development

### Verification Criteria
- Order must have valid merchant assignment
- Status transitions must follow valid state machine
- Events must be published for each status change
- Merchants can only manage their own orders
- Customers can only view their own orders
- Order totals must match sum of item prices

### Test Cases

- [ ] Test 1: Place order with merchant assignment and total calculation
  - Expected: Order created with correct merchant, calculated total

- [ ] Test 2: Merchant accepts pending order
  - Expected: Status changes to ACCEPTED, event published, customer notified

- [ ] Test 3: Merchant rejects pending order with reason
  - Expected: Status changes to REJECTED, reason stored, customer notified

- [ ] Test 4: Invalid status transition rejected
  - Expected: Error returned, status unchanged

- [ ] Test 5: Merchant cannot manage another merchant's orders
  - Expected: Authorization error

- [ ] Test 6: Complete order fulfillment flow
  - Expected: Order progresses through all valid states

- [ ] Test 7: Customer cancels pending order
  - Expected: Status changes to CANCELLED, merchant notified

- [ ] Test 8: Concurrent status updates handled correctly
  - Expected: Last valid update wins, history preserved

### Edge Cases
- Order with invalid menu items (removed/out of stock)
- Merchant goes offline during order processing
- Multiple status updates in rapid succession
- Order total calculation with changed prices
- Delivery address required for delivery, not for pickup

### System Integration Points
- Event bus must be running for notifications
- Menu service must be available for price lookups
- Auth middleware must validate merchant identity
- Database must support concurrent updates

## AI-Human Collaboration Notes

### AI-Executable Tasks
- Generate boilerplate code for new entities and methods
- Implement standard CRUD operations
- Create database migrations
- Write unit tests for pure business logic
- Generate API endpoint handlers following existing patterns

### Human Verification Points
- Review state machine logic for business correctness
- Validate merchant user experience flow
- Test real-time notification delivery
- Verify database migration safety
- Approve API design and breaking changes
- Test edge cases in staging environment

