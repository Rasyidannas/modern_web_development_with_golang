# Modern Web Development with Golang - Architecture

## Architecture Flow

```
┌─────────────┐    HTTP Request    ┌──────────────────┐
│   Client    │ ──────────────────▶ │  HTTP Server     │
└─────────────┘                    │  (gin)           │
                                   └────────┬─────────┘
                                          │
                              ┌─────────────▼─────────────┐
                              │  Controllers              │
                              │  - RunnersController      │
                              │  - ResultsController      │
                              └─────────────┬─────────────┘
                                          │
                              ┌─────────────▼─────────────┐
                              │  Services                 │
                              │  - RunnersService         │
                              │  - ResultsService         │
                              └─────────────┬─────────────┘
                                          │
                              ┌─────────────▼─────────────┐
                              │  Repositories             │
                              │  - RunnersRepository      │
                              │  - ResultsRepository      │
                              └─────────────┬─────────────┘
                                          │
                              ┌─────────────▼─────────────┐
                              │  Database (PostgreSQL)    │
                              └───────────────────────────┘
```

### Request Flow

1. **Client** sends HTTP request (POST/PUT/DELETE) to `/runner` or `/result` endpoints
2. **HTTP Server** routes request to appropriate **Controller**
3. **Controller** validates input and calls **Service** layer
4. **Service** orchestrates business logic and calls **Repository** layer
5. **Repository** interacts with **Database**
6. **Response** is returned to **Controller** and sent to **Client**

### Key Components

- **Controllers** (`controllers/`): Handle HTTP requests/responses, validate input, call services
- **Services** (`services/`): Business logic, validation, orchestration
- **Repositories** (`repositories/`): Data access, database operations
- **Models** (`models/`): Data structures (Runner, Result, ResponseError)
- **Server** (`server/`): HTTP server setup, route registration
- **Main** (`main.go`): Application entry point

## MVC AntiPatterns Applied

### 1. Controllers (controllers/)
Based on the Rails example:

**✅ What can do:**
- ✅ Call `models/` for data structures used in validation
- ✅ Call `services/` for business logic orchestration
- ✅ Validate request input
- ✅ Return HTTP responses with appropriate status codes

**❌ What cannot do (anti-patterns):**
- ❌ Call `config/` for configuration (it's handled in `server/`)
- ❌ Call `dbscripts/` for migrations (handled by external tools)
- ❌ Call `server/` again (circular dependency)

**Current Implementation:**
```go
// Controllers properly follow MVC:
func CreateRunner(ctx *gin.Context) {
    // Validate input (uses models)
    var runner models.Runner
    
    // Call services for business logic
    response, responseErr := rh.runnersService.CreateRunner(&runner)
    
    // Return HTTP response
    ctx.JSON(http.StatusOK, response)
}
```

### 2. Services (services/)
Based on the Rails example:

**✅ What can do:**
- ✅ Call `repositories/` for data access
- ✅ Call other `services/` for complex orchestration
- ✅ Use `models/` for domain data structures
- ✅ Implement business validation and rules

**❌ What should not do (anti-patterns):**
- ❌ Call `controllers/` (violates separation)
- ❌ Call `server/` (violates separation)
- ❌ Directly call `database/` (should go through repositories)

**Current Implementation:**
```go
// Services properly use repositories:
func CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
    // Validate using business logic
    responseErr := validateRunner(runner)
    
    // Delegate to repository for persistence
    return rs.runnersRepository.CreateRunner(runner)
}
```

### 3. Models (models/)
Based on the Rails example:

**✅ What can be used by:**
- ✅ `controllers/` (for request/response structures)
- ✅ `services/` (for domain logic validation)
- ✅ Used in other `models/` (type references)

**❌ What should not use:**
- ❌ `views/` (not applicable in Go - JSON responses)
- ❌ `helpers/` (not present in Go)

**Current Implementation:**
```go
// Models serve as data structures:
type Runner struct {
    ID string `json:"id"`
    FirstName string `json:"first_name"`
    // ... other fields
}

// Models used throughout:
// - Controllers: for request/response structs
// - Services: for business validation
// - Other models: type references
```

### 4. Repositories (repositories/)
Go's equivalent of Rails 'Models' for database access:

**✅ What can be used by:**
- ✅ `services/` (for data persistence)
- ✅ Used in other `repositories/` (for complex queries)

**❌ What should not use:**
- ❌ `controllers/` (separated for abstraction)
- ❌ `server/` (circular dependency)

**Current Implementation:**
```go
// Repository interface-based design:
type RunnersRepository interface {
    CreateRunner(*Runner) (*Runner, error)
    GetRunner(string) (*Runner, error)
    // ... other methods
}
```

### 5. Server (server/)
Go's equivalent of Rails 'Views' + framework setup:

**✅ What can be used by:**
- ✅ Called by `main.go` (startup)
- ✅ Called by repositories (for database connections)

**❌ What should not use:**
- ❌ Called by controllers (circular dependency)
- ❌ Called by services directly

**Current Implementation:**
```go
// Server handles HTTP setup:
func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
    // Route registration
    router.POST("/runner", runnersController.CreateRunner)
    
    return HttpServer{router: router}
}
```

### 6. Config (config/)
Go's equivalent of Rails 'Config':

**✅ What can be used by:**
- ✅ Called by `main.go` (startup)
- ✅ Called by `server/` (for HTTP settings)

**❌ What should not use:**
- ❌ Called by controllers (should be abstracted)

**✅ What can do:**
- ✅ Load application configuration
- ✅ Provide configuration to other layers

### 7. Main (main.go)
Go's equivalent of Rails 'Entry point':

**✅ What can do:**\n- ✅ Initialize configuration
- ✅ Initialize database
- ✅ Initialize HTTP server
- ✅ Start the application

**❌ What should not do:**
- ❌ Direct business logic
- ❌ Request handling

**✅ What should orchestrate:**
- ✅ Application startup flow
- ✅ Layer dependencies

## Summary

Your Go architecture follows MVC principles but with Go-specific adaptations:

**Good MVC Design:**
- ✅ Clear separation between HTTP, business logic, and data access
- ✅ Repository interfaces for testability
- ✅ Controller validation before service calls
- ✅ Business logic isolation in services

**Areas for Improvement:**
- ❗ Add repository interfaces for proper abstraction
- ❗ Separate HTTP error mapping from domain logic
- ❗ Implement transaction management
- ❗ Add DTOs for request/response boundaries

This maintains clean MVC principles while leveraging Go's idiomatic patterns of explicit dependencies and interface-based design.
