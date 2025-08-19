# BitMerchant: Universal Multi-Tenant Commerce Platform with Lightning Network

## Executive Summary

BitMerchant is a comprehensive multi-tenant commerce platform that enables businesses across various industries to accept Bitcoin payments through Lightning Network. The system supports multiple business types including restaurants, retail stores, service providers, and digital goods merchants, each with their independent staff management, product configuration, and payment processing while sharing common infrastructure.

## Supported Business Types

### Food & Beverage

- **Restaurants**: Dine-in, takeout, delivery ordering
- **Cafes & Bars**: Quick service and table ordering
- **Food Trucks**: Mobile ordering and location-based services
- **Catering**: Event-based ordering and scheduling

### Retail & E-commerce

- **Physical Stores**: In-store POS and inventory management
- **Online Shops**: Digital storefronts and order fulfillment
- **Pop-up Shops**: Temporary retail and event sales
- **Marketplaces**: Multi-vendor platforms and commission management

### Service Businesses

- **Salons & Spas**: Appointment booking and service packages
- **Fitness Studios**: Class bookings and membership management
- **Consulting**: Time-based billing and project management
- **Repair Services**: Work orders and parts tracking

### Digital Goods & Subscriptions

- **Software Licenses**: One-time and subscription billing
- **Digital Content**: Media sales and streaming access
- **Online Courses**: Educational content and progress tracking
- **SaaS Products**: Subscription management and usage billing

### Events & Entertainment

- **Event Ticketing**: Concert, conference, and entertainment tickets
- **Tours & Experiences**: Booking and capacity management
- **Sports Facilities**: Court/field reservations and memberships
- **Entertainment Venues**: Show bookings and concession sales

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
│   │   ├── /tenant        # Business management (universal)
│   │   ├── /auth          # Authentication & authorization
│   │   ├── /catalog       # Product/service/menu management
│   │   ├── /orders        # Order/booking/transaction processing
│   │   ├── /payments      # Payment handling (Strike API)
│   │   ├── /inventory     # Stock and availability management
│   │   ├── /scheduling    # Appointments and time-based services
│   │   ├── /subscriptions # Recurring billing and memberships
│   │   └── /analytics     # Reporting & insights
│   ├── /infrastructure    # External integrations
│   │   ├── /database      # PostgreSQL connections
│   │   ├── /strike        # Strike API client
│   │   ├── /notifications # SMS/Email/Push services
│   │   └── /integrations  # Third-party service connectors
│   ├── /web               # HTTP handlers & middleware
│   ├── /templates         # Templ components (business-type aware)
│   └── /business-types    # Business-specific logic and workflows
└── /static                # CSS, JS, PWA assets
```

### Business Type Architecture

#### Universal Core Modules

- **Tenant Management**: Supports any business type with flexible configuration
- **Catalog Management**: Products, services, menus, digital goods, appointments
- **Order Processing**: Universal order/booking/transaction lifecycle
- **Payment Integration**: Lightning Network payments for all business types
- **Staff Management**: Role-based access control adaptable to any industry

#### Business-Specific Extensions

- **Restaurant Module**: Kitchen workflows, table management, food preparation
- **Retail Module**: Inventory tracking, POS integration, shipping management
- **Service Module**: Appointment scheduling, service duration, resource allocation
- **Digital Module**: License generation, download management, subscription billing
- **Event Module**: Ticket generation, capacity management, venue logistics

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

### Universal Tenant Onboarding Process

1. **Business Registration**: Business type selection and details setup
1. **Domain Configuration**: Custom subdomain or path-based routing
1. **Business-Type Setup**: Industry-specific configuration wizard
1. **Staff Invitation**: Role-based team member setup (adapted to business type)
1. **Catalog Configuration**: Products/services/menu items with pricing
1. **Payment Testing**: Lightning Network transaction verification
1. **Integration Setup**: Business-specific third-party integrations
1. **Go-Live Activation**: System deployment and monitoring setup

### Business Type Selection Impact

- **UI Templates**: Industry-appropriate interface and terminology
- **Workflow Configuration**: Business-specific order/booking processes
- **Role Definitions**: Staff roles adapted to industry requirements
- **Feature Enablement**: Relevant modules activated per business type
- **Analytics Dashboards**: Industry-specific KPIs and metrics

## User Types and Roles

### Customer Users (Universal)

- **Authentication**: SMS OTP-based login or email-based registration
- **Capabilities**: Browse catalog, place orders/bookings, track status, payment
- **Scope**: Single tenant (business-specific)
- **Adaptations**: Interface adapts to business type (menu/products/services/tickets)

### Business Staff Roles (Adaptable Framework)

#### Owner (Universal)

- **Permissions**: Full system access and control
- **Capabilities**:
  - Complete business management
  - Staff invitation and role assignment
  - Financial analytics and reporting
  - Strike API configuration
  - Catalog management and pricing
  - Customer data access

#### Manager (Adapted per Business Type)

- **Restaurant**: Kitchen oversight, menu management, staff scheduling
- **Retail**: Inventory management, sales analysis, supplier coordination
- **Service**: Appointment management, service provider scheduling, resource allocation
- **Digital**: Subscription management, content curation, user analytics
- **Events**: Venue management, capacity planning, attendee coordination

#### Staff Level 1 (Customer-Facing)

- **Restaurant**: Cashier/Server - order taking, customer service, payment verification
- **Retail**: Sales Associate - customer assistance, POS operation, returns processing
- **Service**: Receptionist - appointment booking, customer check-in, payment processing
- **Digital**: Support Agent - customer onboarding, technical support, account management
- **Events**: Box Office - ticket sales, attendee check-in, customer service

#### Staff Level 2 (Operations)

- **Restaurant**: Kitchen Staff - food preparation, order status updates, inventory monitoring
- **Retail**: Warehouse Staff - inventory management, order fulfillment, shipping coordination
- **Service**: Service Provider - appointment delivery, service completion, customer interaction
- **Digital**: Content Manager - digital asset management, license provisioning, system monitoring
- **Events**: Event Coordinator - venue setup, logistics coordination, vendor management

#### Staff Level 3 (Specialized)

- **Restaurant**: Chef/Kitchen Manager - menu planning, food quality, kitchen efficiency
- **Retail**: Buyer/Merchandiser - product sourcing, pricing strategy, vendor relations
- **Service**: Senior Practitioner - complex services, training, quality assurance
- **Digital**: Developer/Admin - system configuration, integration management, security
- **Events**: Technical Director - AV setup, production coordination, equipment management

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

## Universal Order/Transaction Management System

### Universal Transaction Lifecycle

1. **Initiation**: Customer browses catalog and creates order/booking/purchase
1. **Configuration**: Item customization, service selection, appointment scheduling
1. **Payment**: Lightning Network invoice generation and payment
1. **Confirmation**: Payment webhook triggers transaction activation
1. **Fulfillment**: Business-specific fulfillment process
1. **Completion**: Transaction fulfilled and closed
1. **Follow-up**: Reviews, support, recurring billing (if applicable)

### Business-Type Specific Workflows

#### Restaurant Orders

- **Pending → Paid → Preparing → Ready → Completed**
- Kitchen queue management and preparation tracking
- Table service and pickup coordination

#### Retail Orders

- **Pending → Paid → Processing → Shipped → Delivered**
- Inventory allocation and warehouse fulfillment
- Shipping integration and tracking updates

#### Service Bookings

- **Pending → Paid → Scheduled → In-Progress → Completed**
- Appointment calendar and resource scheduling
- Service delivery and customer interaction tracking

#### Digital Purchases

- **Pending → Paid → Provisioned → Active → (Renewed/Expired)**
- License generation and delivery automation
- Subscription management and renewal processing

#### Event Tickets

- **Pending → Paid → Issued → Validated → Attended**
- Ticket generation and QR code creation
- Entry validation and attendance tracking

### Status Management Framework

- **Universal States**: All business types support core states (pending, paid, completed)
- **Business-Specific States**: Additional states relevant to each industry
- **State Transitions**: Configurable workflows per business type
- **Notification Triggers**: Automated communications at key transition points

## Database Schema Design

### Universal Multi-Tenant Data Model

#### Tenant Management

- **Tenants Table**: Business information, type, Strike configuration, status
- **Business Types Table**: Industry definitions and configuration templates
- **Tenant Settings**: Operational preferences, UI customization, business rules
- **Domain Mapping**: URL routing and custom domain configuration

#### User Management

- **Customers Table**: Customer profiles with tenant association
- **Staff Table**: Business employees with role assignments (business-type aware)
- **Authentication Tokens**: Session management and security tracking

#### Universal Catalog Management

- **Categories Table**: Product/service/menu category organization per tenant
- **Items Table**: Universal items (products/services/menu/digital goods) with pricing
- **Item Attributes**: Flexible attribute system for business-specific properties
- **Modifiers Table**: Add-ons, customizations, and variations
- **Availability Schedules**: Time-based availability and booking slots

#### Universal Transaction Management

- **Transactions Table**: Universal transaction headers (orders/bookings/purchases)
- **Transaction Items Table**: Individual items within each transaction
- **Transaction Status History**: Complete audit trail of status changes
- **Fulfillment Tracking**: Business-specific fulfillment data

#### Extended Business Modules

- **Appointments Table**: Service bookings and time-slot management
- **Subscriptions Table**: Recurring billing and membership management
- **Inventory Table**: Stock tracking for physical goods
- **Digital Assets Table**: License keys, download links, access tokens
- **Events Table**: Event information, capacity, and ticketing

#### Payment & Financial Tracking

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

#### Business-Type Specific Templates

- **Restaurant Views**: Menu browsing, order tracking, kitchen management, table service
- **Retail Views**: Product catalog, shopping cart, inventory management, shipping
- **Service Views**: Service booking, appointment calendar, provider management, client history
- **Digital Views**: License management, download center, subscription billing, user access
- **Event Views**: Ticket purchasing, event information, attendee management, venue logistics

### Real-Time Updates with Datastar (Business-Aware)

- **Server-Sent Events**: Live transaction status updates across all business types
- **Partial Page Updates**: Update specific components without full reload
- **Business-Specific Updates**: Kitchen queues, inventory levels, appointment schedules
- **Live Data**: Real-time metrics adapted to business type (orders, bookings, sales)

### PWA Features Implementation (Universal)

- **Installable**: Web app manifest adapts to business branding
- **Offline Capable**: Cache critical pages and business-specific functionality
- **Push Notifications**: Transaction updates and staff alerts (business-appropriate)
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

- **Transaction Fees**: Percentage-based commission on each transaction (varies by business type)
- **Subscription Plans**: Tiered pricing for advanced features and higher transaction volumes
- **Setup Services**: Onboarding assistance and business-specific consultation
- **Premium Support**: Priority customer service and technical assistance
- **Integration Fees**: Revenue sharing with third-party service integrations

### Pricing Structure by Business Type

#### Food & Beverage

- **Basic Plan**: Core ordering functionality (2.5% transaction fee)
- **Professional Plan**: Kitchen management + analytics (2.0% + $29/month)
- **Enterprise Plan**: Multi-location + advanced features (1.5% + $99/month)

#### Retail & E-commerce

- **Starter Plan**: Basic e-commerce features (2.0% transaction fee)
- **Growth Plan**: Inventory management + shipping (1.8% + $39/month)
- **Scale Plan**: Multi-channel + automation (1.5% + $149/month)

#### Service Businesses

- **Essential Plan**: Appointment booking (2.5% transaction fee)
- **Professional Plan**: Resource management + client history (2.0% + $49/month)
- **Premium Plan**: Advanced scheduling + integrations (1.5% + $199/month)

#### Digital Goods & SaaS

- **Digital Plan**: License management (1.5% transaction fee)
- **Platform Plan**: Subscription billing + analytics (1.2% + $79/month)
- **Enterprise Plan**: API access + white-label options (1.0% + $299/month)

#### Events & Entertainment

- **Event Basic**: Ticket sales + check-in (2.5% transaction fee)
- **Event Pro**: Capacity management + reporting (2.0% + $59/month)
- **Event Enterprise**: Multi-venue + advanced features (1.5% + $249/month)

### Value Proposition by Industry

- **Universal Lightning Payments**: Bitcoin acceptance across all business types
- **Industry-Specific Workflows**: Optimized processes for each business model
- **Unified Platform**: Single system for businesses with multiple revenue streams
- **Low Transaction Costs**: Significantly lower than traditional payment processors
- **Global Reach**: Bitcoin enables international transactions without currency conversion

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