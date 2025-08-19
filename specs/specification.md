# BitMerchant: Multi-Tenant Food Ordering Platform with Lightning Network

## Executive Summary

BitMerchant is a comprehensive multi-tenant food ordering platform that enables restaurants to accept Bitcoin payments through Lightning Network. The system supports multiple restaurant tenants, each with their independent staff management, menu configuration, and payment processing while sharing common infrastructure.

## System Architecture

### High-Level Architecture (MVP Phase)

- **Frontend**: Progressive Web App (PWA) using Templ + Datastar
- **Backend**: Go-based Modular Monolith with domain-driven modules
- **Database**: PostgreSQL with multi-tenant schema design
- **Payment Integration**: Strike API for Lightning Network payment processing
- **Real-time Communication**: Server-Sent Events (SSE) via Datastar
- **Deployment**: Single binary with embedded static assets

### Modular Monolith Structure

```
/bitmerchant
├── /cmd                    # Application entry points
├── /internal
│   ├── /domain            # Core business logic modules
│   │   ├── /tenant        # Restaurant management
│   │   ├── /auth          # Authentication & authorization
│   │   ├── /menu          # Menu management
│   │   ├── /orders        # Order processing
│   │   ├── /payments      # Payment handling (Strike API)
│   │   └── /analytics     # Reporting & insights
│   ├── /infrastructure    # External integrations
│   │   ├── /database      # PostgreSQL connections
│   │   ├── /strike        # Strike API client
│   │   └── /notifications # SMS/Email services
│   ├── /web               # HTTP handlers & middleware
│   └── /templates         # Templ components
└── /static                # CSS, JS, PWA assets
```

### Technology Stack Rationale

#### Go Modular Monolith Benefits

- **Single Deployment**: One binary, simplified operations
- **Shared Database**: ACID transactions across modules
- **Fast Development**: No network latency between components
- **Easy Debugging**: Single process, unified logging
- **Cost Effective**: Reduced infrastructure complexity

## Multi-Tenant Architecture

### Tenant Isolation Strategy

- **Database Level**: Tenant ID in all tables with row-level security
- **API Level**: URL-based tenant routing (`/api/{tenantId}/...`)
- **Payment Level**: Separate Strike API keys per tenant
- **Configuration Level**: Independent settings per restaurant

### Tenant Onboarding Process

1. **Restaurant Registration**: Business details and Strike API setup
1. **Domain Configuration**: Custom subdomain or path-based routing
1. **Staff Invitation**: Role-based team member setup
1. **Menu Configuration**: Food items, pricing, and categories
1. **Payment Testing**: Lightning Network transaction verification
1. **Go-Live Activation**: System deployment and monitoring setup

## User Types and Roles

### Customer Users

- **Authentication**: SMS OTP-based login
- **Capabilities**: Browse menu, place orders, track status, payment
- **Scope**: Single tenant (restaurant-specific)

### Restaurant Staff Roles

#### Owner

- **Permissions**: Full system access and control
- **Capabilities**:
  - Complete restaurant management
  - Staff invitation and role assignment
  - Financial analytics and reporting
  - Strike API configuration
  - Menu management and pricing
  - Customer data access

#### Manager

- **Permissions**: Operational management without financial settings
- **Capabilities**:
  - Staff management (view only)
  - Comprehensive analytics dashboard
  - Menu management and updates
  - Order oversight and reporting
  - Inventory and operational settings

#### Cashier (Front of House)

- **Permissions**: Customer-facing operations
- **Capabilities**:
  - Active order monitoring
  - Customer assistance and support
  - Payment verification and troubleshooting
  - Order status updates (ready → completed)
  - Customer order history lookup

#### Kitchen Staff

- **Permissions**: Food preparation focused
- **Capabilities**:
  - Incoming order queue visibility
  - Cooking status updates (paid → preparing → ready)
  - Preparation time tracking
  - Menu item availability toggle
  - Kitchen-specific notifications

### Platform Administrator

- **Permissions**: Multi-tenant system oversight
- **Capabilities**:
  - Tenant management and monitoring
  - Platform analytics and health metrics
  - Strike API integration management
  - System configuration and maintenance

## Authentication and Authorization

### Authentication Strategy

- **Token Type**: JWT (JSON Web Tokens)
- **Token Structure**: User ID, Tenant ID, Role, Permissions, Expiration
- **Validation**: Multi-layer verification (signature, tenant match, user status)
- **Refresh Strategy**: Automatic token renewal before expiration

### Authorization Model

- **Role-Based Access Control (RBAC)**: Predefined roles with specific permissions
- **Tenant-Scoped Permissions**: All actions restricted to user’s tenant context
- **API-Level Enforcement**: Middleware validation on every protected endpoint
- **Cross-Tenant Isolation**: Strict prevention of data access across tenants

### Security Measures

- **Rate Limiting**: Per-IP and per-user request throttling
- **CORS Configuration**: Tenant-domain-based origin validation
- **API Key Encryption**: Secure storage of Strike API credentials
- **Audit Logging**: Comprehensive activity tracking and monitoring

## Payment Integration (Strike API)

### Lightning Network Integration

- **Payment Provider**: Strike API for Lightning Network processing
- **Invoice Generation**: Dynamic Bitcoin invoice creation per order
- **Real-time Verification**: Webhook-based payment confirmation
- **Multi-Currency Support**: Automatic fiat conversion capabilities

### Payment Flow

1. **Order Creation**: Calculate total amount in satoshis and fiat
1. **Invoice Generation**: Create Strike Lightning invoice with order correlation
1. **Payment Presentation**: QR code and payment request to customer
1. **Payment Processing**: Customer pays via any Lightning wallet
1. **Confirmation**: Webhook notification triggers order status update
1. **Settlement**: Automatic fund transfer to restaurant (minus commission)

### Financial Management

- **Commission Structure**: Configurable percentage-based fees per tenant
- **Revenue Tracking**: Real-time commission calculation and reporting
- **Settlement Options**: Automatic or manual payout to restaurant accounts
- **Currency Handling**: Multi-currency support with real-time conversion rates

## Order Management System

### Order Lifecycle

1. **Creation**: Customer selects items and creates order
1. **Payment**: Lightning Network invoice generation and payment
1. **Confirmation**: Payment webhook triggers order activation
1. **Preparation**: Kitchen receives order and begins cooking
1. **Completion**: Food ready for pickup/delivery
1. **Fulfillment**: Customer receives order and transaction closes

### Order Status Management

- **Pending**: Created but awaiting payment
- **Paid**: Payment confirmed, sent to kitchen
- **Preparing**: Kitchen actively cooking order
- **Ready**: Food prepared and available for pickup
- **Completed**: Order fulfilled and closed
- **Cancelled**: Order terminated (with refund if applicable)

### Real-time Updates

- **Customer Notifications**: SMS/push notifications for status changes
- **Staff Dashboards**: Live order queues and status boards
- **Kitchen Displays**: Real-time cooking queue management
- **Manager Analytics**: Live operational metrics and insights

## Database Schema Design

### Multi-Tenant Data Model

#### Tenant Management

- **Tenants Table**: Restaurant information, Strike configuration, status
- **Tenant Settings**: Operational preferences, UI customization, business rules
- **Domain Mapping**: URL routing and custom domain configuration

#### User Management

- **Customers Table**: Customer profiles with tenant association
- **Staff Table**: Restaurant employees with role assignments
- **Authentication Tokens**: Session management and security tracking

#### Menu Management

- **Categories Table**: Food category organization per tenant
- **Products Table**: Menu items with pricing, descriptions, availability
- **Modifiers Table**: Add-ons, customizations, and variations

#### Order Management

- **Orders Table**: Order headers with customer, tenant, payment information
- **Order Items Table**: Individual food items within each order
- **Order Status History**: Complete audit trail of status changes

#### Payment Tracking

- **Payment Records**: Strike transaction details and reconciliation
- **Commission Tracking**: Platform revenue calculation and reporting
- **Financial Reconciliation**: Settlement and payout management

## User Interface Design

### Progressive Web App Architecture

- **Server-Side Rendered**: Templ templates with Go backend
- **Enhanced with Datastar**: Real-time updates and interactivity
- **Offline-First**: Service worker for critical functionality
- **Responsive Design**: Single UI adapts to all screen sizes
- **Native-Like Experience**: App installation and push notifications

### Templ Component Structure

#### Shared Components

- **Layout Templates**: Common headers, navigation, footers
- **Form Components**: Reusable input fields and validation
- **Modal Systems**: Overlays for confirmations and forms
- **Notification Banners**: Success/error message displays

#### Role-Specific Templates

- **Customer Views**: Menu browsing, cart, order tracking
- **Kitchen Views**: Order queue, status updates, timers
- **Cashier Views**: Active orders, payment verification, customer help
- **Manager Views**: Analytics dashboard, menu management, staff overview
- **Owner Views**: Complete restaurant management interface

### Real-Time Updates with Datastar

- **Server-Sent Events**: Live order status updates
- **Partial Page Updates**: Update specific components without full reload
- **Form Enhancements**: Validation and submission without page refresh
- **Live Data**: Real-time metrics and queue updates

### PWA Features Implementation

- **Installable**: Web app manifest for home screen installation
- **Offline Capable**: Cache critical pages and functionality
- **Push Notifications**: Order updates and staff alerts
- **Background Sync**: Queue actions when offline, sync when online
- **App-Like Navigation**: Smooth transitions and native-feeling interactions

### Mobile-First Responsive Design

- **Touch-Friendly**: Large tap targets and gesture support
- **Screen Adaptation**: Optimized layouts for phone, tablet, desktop
- **Keyboard Accessibility**: Full functionality without mouse
- **Performance Optimized**: Fast loading on mobile networks

## Technical Requirements

### Performance Specifications

- **Response Time**: Server responses under 100ms for 95th percentile
- **Throughput**: Support 500+ concurrent users on single instance
- **Availability**: 99.5% uptime with automated health checks
- **Scalability**: Vertical scaling initially, horizontal when needed

### Go Application Requirements

- **Binary Size**: Optimized executable under 50MB with embedded assets
- **Memory Usage**: Efficient memory management with garbage collection tuning
- **Concurrency**: Go routines for concurrent request handling
- **Database Connections**: Connection pooling and prepared statements

### PWA Technical Specifications

- **Loading Performance**: First Contentful Paint under 2 seconds
- **Offline Functionality**: Core features available without network
- **Cache Strategy**: Intelligent caching of static and dynamic content
- **Service Worker**: Background sync and push notification handling

### Security Requirements

- **TLS Encryption**: HTTPS everywhere with modern cipher suites
- **Authentication**: Secure JWT handling with proper expiration
- **Input Validation**: Server-side validation for all user inputs
- **SQL Injection Protection**: Parameterized queries and prepared statements
- **XSS Prevention**: Template escaping and Content Security Policy

### Infrastructure Requirements (MVP)

- **Single Server Deployment**: VPS or cloud instance with Go binary
- **PostgreSQL Database**: Single database with backup strategy
- **Reverse Proxy**: Nginx for SSL termination and static file serving
- **Monitoring**: Basic health checks and error logging
- **CI/CD Pipeline**: Automated testing and deployment from Git

## Business Model

### Revenue Streams

- **Transaction Fees**: Percentage-based commission on each order
- **Subscription Plans**: Tiered pricing for advanced features
- **Setup Services**: Onboarding assistance and consultation
- **Premium Support**: Priority customer service and technical assistance

### Pricing Structure

- **Basic Plan**: Core ordering functionality with standard commission
- **Professional Plan**: Advanced analytics and staff management features
- **Enterprise Plan**: Custom integrations and dedicated support
- **Transaction Fees**: 2-3% commission on Lightning Network payments

## Risk Assessment and Mitigation

### Technical Risks

- **Lightning Network Stability**: Multiple payment provider integration
- **Scaling Challenges**: Microservices architecture and load balancing
- **Data Security**: Comprehensive encryption and access controls
- **Third-Party Dependencies**: Fallback systems and monitoring

### Business Risks

- **Market Adoption**: Education and incentive programs for Bitcoin payments
- **Regulatory Changes**: Compliance monitoring and legal consultation
- **Competition**: Unique value proposition and continuous innovation
- **Customer Support**: Scalable support systems and documentation

## Success Metrics

### Technical KPIs

- **Payment Success Rate**: >99% Lightning Network transaction completion
- **System Uptime**: >99.9% availability across all services
- **Response Times**: Sub-200ms API response for standard operations
- **User Adoption**: Monthly active users and transaction volume growth

### Business KPIs

- **Tenant Growth**: Number of active restaurants using the platform
- **Revenue Metrics**: Total payment volume and commission revenue
- **User Engagement**: Order frequency and customer retention rates
- **Market Penetration**: Geographic expansion and market share growth

## Future Enhancements

### Phase 2 Features (6-12 months)

- **Enhanced Analytics**: Advanced reporting and business intelligence
- **Multi-Currency Support**: Additional cryptocurrency payment options
- **Delivery Integration**: Third-party delivery service partnerships
- **Advanced PWA Features**: Better offline support and native integrations

### Phase 3 Expansion (12-18 months)

- **Microservices Migration**: Extract high-load modules to separate services
- **Native Mobile Apps**: iOS and Android apps for enhanced user experience
- **API Marketplace**: Third-party integrations and developer ecosystem
- **International Markets**: Multi-language and regional compliance features

### Scaling Strategy

- **Vertical Scaling**: Optimize single instance performance first
- **Read Replicas**: Database read scaling for analytics and reporting
- **CDN Integration**: Static asset delivery optimization
- **Microservices Extraction**: Graduate to microservices as needed

### Technology Evolution Path

- **Current**: Go Monolith + Templ + Datastar + PWA
- **Phase 2**: Add caching layer (Redis), background jobs
- **Phase 3**: Extract payment service, add event bus
- **Phase 4**: Full microservices with API gateway

-----

**Document Version**: 1.0  
**Last Updated**: August 2025  
**Status**: Draft for Review