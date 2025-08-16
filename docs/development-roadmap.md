# Minimart Development Roadmap

## Entity-First, Hypermedia-Driven Architecture

### 🎯 **Project Vision**
Build a production-ready Bitcoin-native food ordering platform with real-time order fulfillment, using rich domain models and hypermedia-driven UI.

## Development Phases

### ✅ **Phase 1: Rich Domain Model** - *Complete*
**Goal**: Establish business logic in domain entities
- Rich Order aggregate with state transitions
- Bitcoin-native Money value objects (Satoshis)
- Domain events for integration
- Complete state machine validation
- **Result**: Infrastructure-agnostic business logic

### ✅ **Phase 2: Entity Integration** - *Complete*  
**Goal**: Integrate Order and MenuItem aggregates
- Rich MenuItem entity with stock management
- Order-MenuItem integration with price snapshots
- Stock reservation and release mechanisms
- Comprehensive integration testing
- **Result**: Complete domain model interactions

### ✅ **Phase 3: Use Cases & Infrastructure** - *Complete*
**Goal**: Thin orchestration layer and database foundation
- Merchant entity with operating hours and business rules
- Complete order workflow use cases (Place → Accept → Complete)
- Merchant analytics and order management
- Bitcoin-native database schema with JSONB optimization
- Repository pattern with clean interfaces
- **Result**: Complete backend foundation ready for UI

### 🚧 **Phase 4: Hypermedia UI** - *In Progress*
**Goal**: Real-time, reactive user interface

#### 🎯 **Architectural Decision**
**Skip traditional JSON APIs** and build hypermedia UI directly.

**Why?**
- Existing handlers need complete rewrites anyway
- Hypermedia was always the end goal
- Simpler architecture with single rendering pipeline
- Superior UX with real-time updates

#### Implementation Strategy
1. **HTML Templates with Datastar**
   - Reactive templates with `data-*` attributes
   - Progressive enhancement strategy
   - Server-side rendering with client-side reactivity

2. **View Models & Presentation Logic**
   - Transform domain entities for display
   - Bitcoin display formatting (sats/mBTC/BTC)
   - Status-specific UI components

3. **Hypermedia Handlers**
   - Return HTML fragments instead of JSON  
   - Direct use case → HTML response pipeline
   - Proper error handling with HTML responses

4. **Server-Sent Events (SSE)**
   - Real-time order status updates
   - Merchant notification streams
   - Customer order tracking
   - DOM updates via Datastar

5. **Static Assets & Styling**
   - Datastar JavaScript library
   - Tailwind CSS for responsive design
   - Mobile-first merchant workflows

#### Key Features
- **Merchant Dashboard**: Real-time order management
- **Customer Order Tracking**: Live status updates
- **Menu Management**: Dynamic item availability
- **Bitcoin Pricing**: Native Satoshi display throughout

### 📋 **Phase 5: Production Features** - *Planned*
**Goal**: Production-ready platform
- Advanced analytics dashboards
- Order modification capabilities
- Mobile app optimization
- Performance monitoring
- Security hardening

## Architecture Benefits Achieved

### 🏗️ **Clean Architecture**
- **Domain Layer**: All business logic in entities
- **Use Case Layer**: Pure orchestration, no business rules  
- **Infrastructure Layer**: Pure persistence, no business logic
- **Presentation Layer**: HTML templates with reactive updates

### ⚡ **Performance & Scalability**
- **Fast Tests**: Business logic tests run in milliseconds
- **Efficient Queries**: Optimized database with proper indexes
- **Real-Time Updates**: Server-sent events eliminate polling
- **Bitcoin Precision**: Satoshi-based calculations prevent rounding errors

### 🔧 **Developer Experience**
- **Entity-First Development**: Business rules are clear and testable
- **Type Safety**: Rich domain models prevent invalid states
- **Event-Driven**: Clean integration points via domain events
- **Hypermedia**: Single rendering pipeline reduces complexity

## Technical Stack

### **Backend**
- **Language**: Go 1.23+
- **Framework**: Fiber (HTTP server)
- **Database**: PostgreSQL 15+ with JSONB
- **Cache/Events**: Redis 7+
- **Migrations**: Goose

### **Frontend** 
- **Architecture**: Hypermedia-driven with server-side rendering
- **Reactivity**: Datastar (reactive DOM updates)
- **Styling**: Tailwind CSS
- **Real-Time**: Server-Sent Events (SSE)

### **Financial System**
- **Currency**: Bitcoin-native (Satoshis as base unit)
- **Display**: Smart formatting (sats/mBTC/BTC)
- **Precision**: Integer arithmetic prevents floating-point errors

## Current Status

```
✅ Domain Model     - Rich entities with business logic
✅ Entity Tests     - 100% business rule coverage  
✅ Use Cases        - Complete order workflow
✅ Infrastructure   - Database schema and repositories
🚧 Hypermedia UI    - Real-time merchant dashboard (in progress)
📋 Production       - Advanced features (planned)
```

## Next Milestones

### **Immediate (Phase 4)**
- [ ] Merchant dashboard with real-time order updates
- [ ] Customer order tracking interface
- [ ] Bitcoin pricing display throughout UI
- [ ] Mobile-responsive design

### **Short Term**
- [ ] Menu management interface
- [ ] Advanced merchant analytics
- [ ] Order modification capabilities
- [ ] Performance optimization

### **Long Term**
- [ ] Lightning Network integration
- [ ] Multi-merchant platform features
- [ ] Mobile app development
- [ ] API for third-party integrations

---

**Last Updated**: August 16, 2025  
**Current Phase**: Phase 4 - Hypermedia UI  
**Architecture**: Entity-First, Hypermedia-Driven
