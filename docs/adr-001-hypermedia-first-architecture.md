# ADR-001: Skip Traditional JSON APIs for Hypermedia-First Architecture

## Status
**ACCEPTED** - August 2025

## Context

After completing Phase 3 (Use Cases & Infrastructure), we need to decide between two architectural approaches for Phase 4:

### Option A: Traditional Two-Phase Approach
1. **Phase 4**: Create JSON API handlers working with rich entities
2. **Phase 5**: Build hypermedia UI consuming those JSON APIs

### Option B: Hypermedia-First Approach
1. **Phase 4**: Skip JSON APIs and build hypermedia UI directly with HTML-returning handlers

## Decision

**We choose Option B: Hypermedia-First Architecture**

Skip traditional JSON handler integration and implement hypermedia UI directly.

## Rationale

### 1. **Avoid Duplicate Work**
- **Current State**: Existing handlers work with old anemic domain models
- **Reality**: All handlers need complete rewrites to work with rich entities anyway
- **Conclusion**: Building JSON handlers would be throwaway work

### 2. **Hypermedia Was Always the End Goal**
- **Original Vision**: Datastar-powered reactive UI with server-sent events
- **Design Intent**: Real-time order updates through DOM manipulation
- **Architecture**: Backend-driven interactivity without complex frontend state management

### 3. **Simpler Overall Architecture**
- **No Dual APIs**: HTML templates become the API contract
- **Single Rendering Pipeline**: One way to present data to users
- **Reduced Complexity**: No JSON serialization/deserialization layer needed
- **Fewer Integration Points**: Direct from use cases to HTML responses

### 4. **Superior User Experience**
- **Real-Time Updates**: Server-sent events with automatic DOM updates
- **Instant Feedback**: No page refreshes, no loading states
- **Progressive Enhancement**: Works without JavaScript, enhanced with Datastar
- **Mobile-First**: Optimized for merchant mobile workflows

### 5. **Modern Web Architecture**
- **Hypermedia-Driven Applications**: Following contemporary best practices
- **Reduced JavaScript**: Minimal frontend complexity
- **Server-Side Rendering**: Better SEO, faster initial loads
- **Event-Driven UI**: Natural fit with our domain event architecture

## Implementation Strategy

### Phase 4: Hypermedia UI Implementation
1. **HTML Templates with Datastar**
   - Create reactive templates with `data-*` attributes
   - Progressive enhancement strategy
   
2. **View Models**
   - Transform domain entities for presentation
   - Handle Bitcoin display formatting
   - Presentation logic separation
   
3. **Hypermedia Handlers**
   - Return HTML fragments instead of JSON
   - Use Fiber's HTML rendering capabilities
   - Implement proper error handling with HTML responses
   
4. **Server-Sent Events**
   - Real-time order status updates
   - Merchant notification streams
   - Customer order tracking
   
5. **Static Assets**
   - Datastar JavaScript library
   - Tailwind CSS for styling
   - Progressive web app features

### What We Keep
- **User Authentication Handlers**: Already functional JSON endpoints for login/register
- **Health Check Endpoints**: Simple utility endpoints that don't need rich domain logic
- **API Documentation**: Can expose hypermedia endpoints if external integration needed

### Technical Implementation
```go
// Hypermedia Handler Example
func (h *MerchantHandler) AcceptOrder(c *fiber.Ctx) error {
    orderID := c.Params("id")
    merchantID := getMerchantIDFromContext(c)
    
    // Use case orchestration
    err := h.orderUsecase.AcceptOrder(c.Context(), orderID, merchantID, 30)
    if err != nil {
        return c.Status(400).Render("error", fiber.Map{
            "message": err.Error(),
        })
    }
    
    // Return HTML fragment for Datastar to swap
    order, _ := h.orderUsecase.GetOrderByID(c.Context(), orderID)
    return c.Render("order_card", fiber.Map{
        "order": h.orderViewModel.Transform(order),
    })
}
```

## Consequences

### Positive
- **Faster Development**: No intermediate JSON API layer to build and maintain
- **Better UX**: Real-time, reactive interface from day one
- **Simpler Codebase**: Single rendering pipeline, fewer abstractions
- **Modern Architecture**: Follows hypermedia-driven application principles
- **Mobile-Optimized**: Perfect for merchant mobile workflows
- **Event Integration**: Natural fit with domain events → DOM updates

### Negative
- **Less API Flexibility**: Harder to build separate mobile apps or integrations
- **Learning Curve**: Team needs to understand Datastar and hypermedia patterns
- **SEO Considerations**: Need to ensure proper server-side rendering
- **Debugging**: HTML responses harder to debug than JSON APIs

### Mitigations
- **External Integration**: Can expose specific JSON endpoints if needed for third-party integrations
- **Mobile Apps**: Can build native apps consuming the same use cases directly
- **Development Tools**: Use proper HTML debugging tools and Datastar dev tools
- **Documentation**: Comprehensive hypermedia API documentation

## Related Decisions
- **ADR-002**: Use of Datastar for reactive UI (to be created)
- **ADR-003**: Bitcoin-first pricing display in HTML templates (to be created)

## References
- [Hypermedia-Driven Applications](https://hypermedia.systems/)
- [Datastar Documentation](https://data-star.dev/)
- [Entity-First Development Principles](../specs/order-fulfillment-flow.md)

---

**Decision Date**: August 16, 2025  
**Decision Made By**: Development Team  
**Status**: Accepted and implemented in Phase 4
