# Go Learning Path: Building a Modular Monolith (Minimart Project)

This learning plan is designed to help you learn Go by building a B2C ordering system called "Minimart". It focuses on creating a modular monolith using Domain-Driven Design (DDD) principles, drawing from your TypeScript background.

You can track your progress by checking off the items as you complete them.

## Phase 1: Go Fundamentals

*   [x] **Go Syntax and Basic Types:**
    *   [x] Variables, constants, and basic types (`string`, `int`, `float`, `bool`).
    *   [x] Pointers: Understand what they are and when to use them (`*` and `&`).
    *   [x] `fmt` package for printing and formatting.
*   [x] **Data Structures:**
    *   [x] Structs (similar to objects/classes in TS).
    *   [x] Slices (dynamic arrays).
    *   [x] Maps (key-value pairs).
*   [x] **Control Flow:**
    *   [x] `if/else` statements.
    *   [x] `for` loops (Go's only looping construct).
    *   [x] `switch` statements.
*   [x] **Functions:**
    *   [x] Defining functions.
    *   [x] Multiple return values.
    *   [x] Variadic functions.
*   [x] **Packages and Modules:**
    *   [x] Understanding `package main`.
    *   [x] Creating and importing your own packages (like you've done with `user`, `merchant`).
    *   [x] Go modules (`go.mod`, `go.sum`).
*   [x] **Error Handling:**
    *   [x] The `error` type.
    *   [x] Returning and checking for errors.
*   [x] **Interfaces:**
    *   [x] Understanding implicit interface implementation (a key difference from TS).
    *   [x] Defining and using interfaces to create abstractions.
*   [ ] **Concurrency:**
    *   [ ] Goroutines (`go` keyword).
    *   [ ] Channels for communication between goroutines.

## Phase 2: Building the Web API with Fiber

*   [x] **Fiber Basics:**
    *   [x] Routing (`app.Get`, `app.Post`, etc.).
    *   [ ] Route parameters, query parameters.
    *   [x] Handling JSON request bodies and responses.
    *   [x] Handlers and the `fiber.Ctx` context.
*   [x] **Structuring Your Fiber App:**
    *   [x] Grouping routes (`app.Group`).
    *   [x] Refactoring handlers into separate files (as you've started).
*   [ ] **Middleware:**
    *   [ ] Using built-in middleware (e.g., for logging, recovery).
    *   [ ] Writing your own custom middleware.

## Phase 3: Testing

*   [x] **Go's `testing` package:**
    *   [x] Writing unit tests (`_test.go` files).
    *   [x] Running tests (`go test ./...`).
*   [x] **Testing Strategies:**
    *   [x] **Unit Tests:** Test individual functions and methods in isolation (e.g., test a use case with a mock repository).
    *   [x] **Integration Tests:** Test the interaction between different parts of your application (e.g., test a handler all the way to the database).
    *   [x] Mocking dependencies using interfaces.

## Phase 4: Persistence with a Real Database

*   [x] **Choosing a Database and Driver:**
    *   [x] PostgreSQL is a great choice.
    *   [x] Select a Go driver (e.g., `pgx`) or an ORM (e.g., `GORM`, `sqlc`).
*   [x] **Connecting to the Database:**
    *   [x] Manage database connection strings.
    *   [x] Create a database connection pool.
*   [x] **Implementing Repositories:**
    *   [x] Replace the `InMemory...Repository` with a `Postgres...Repository`.
    *   [x] Implement the repository interfaces for `user` and `merchant`.
    *   [x] Write SQL queries (or use an ORM) to perform CRUD operations.
*   [x] **Database Migrations:**
    *   [x] Learn how to manage database schema changes over time (e.g., using a library like `golang-migrate/migrate`).

## Phase 5: Deepening Domain-Driven Design (DDD)

*   [x] **Core DDD Concepts:**
    *   [x] **Entities:** Structs with a unique identity (e.g., `User`, `Merchant`, `Order`).
    *   [ ] **Value Objects:** Structs without a unique identity, defined by their attributes (e.g., `Address`, `Money`).
    *   [x] **Aggregates:** A cluster of domain objects that can be treated as a single unit (e.g., an `Order` with its `OrderItems`).
    *   [x] **Repositories:** Mediate between the domain and data mapping layers (what you have started).
    *   [x] **Usecases/Services:** Encapsulate application-specific business logic (what you have started).
*   [x] **Refining the Modules:**
    *   [x] Review the `user` and `merchant` modules. Do they represent clear domain boundaries?
    *   [x] Implement the `order` module. What are the entities and aggregates?
    *   [x] Implement the `menu` module.
    *   [ ] Think about how the modules interact with each other.

## Phase 6: Event-Driven Architecture with a Message Broker

*   [x] **Message Broker Concepts:**
    *   [x] Understand the role of a message broker (e.g., decoupling services, asynchronous communication).
    *   [x] Learn about different messaging patterns (e.g., Pub/Sub).
*   [x] **Designing a Broker Interface:**
    *   [x] Define a generic interface for a message broker (`EventPublisher` and `EventSubscriber`). This is crucial for making the implementation interchangeable.
*   [x] **In-Memory Implementation:**
    *   [x] Create an `InMemoryMessageBroker` that implements the interface for development and testing.
*   [x] **Publishing Events:**
    *   [x] Modify a use case (e.g., `UserUsecase.CreateUser`) to publish an event (e.g., `UserCreatedEvent`) after a successful operation.
*   [x] **Subscribing to Events:**
    *   [x] Create a simple subscriber for testing purposes.
    *   [x] Create a subscriber in a different module that listens for an event and performs an action (e.g., sending a welcome email).
*   [ ] **Integrating a Production Broker:**
    *   [ ] Choose a production-ready message broker (e.g., Redis Pub/Sub, Kafka).
    *   [ ] Add the necessary client library to `go.mod`.
    *   [ ] Implement the message broker interface for your chosen broker (e.g., `RedisMessageBroker`).
    *   [ ] Use configuration to switch between the in-memory and production broker implementations.

## Phase 7: Production-Ready Improvements

*   [ ] **Configuration Management:**
    *   [ ] Use a library like `Viper` to handle configuration from files (e.g., `config.yaml`) and environment variables.
*   [ ] **Structured Logging:**
    *   [ ] Use a library like `slog` (standard library in Go 1.21+) or `zerolog` for structured, leveled logging.
*   [ ] **Authentication & Authorization:**
    *   [ ] Implement JWT-based authentication.
    *   [ ] Create middleware to protect routes.
*   [ ] **Containerization:**
    *   [ ] Write a `Dockerfile` for your Go application.
    *   [ ] Improve your `docker-compose.yml` to run your Go application and a PostgreSQL database together.
*   [ ] **CI/CD:**
    *   [ ] Set up a basic CI pipeline using GitHub Actions to build and test your application on every push.

Good luck with your learning journey! This project is a great way to apply these concepts in a practical way.
