## Phase 4: Backend Architecture & Dependency Injection

### 1. Architecture

*   **My Implementation (Matching `be-service`):**
    *   `models/`: Defines data shapes (Structs). Has no business logic or HTTP awareness.
    *   `middlewares/`: Houses interceptors like `SlogMiddleware` and `RequestIDMiddleware`.
    *   `utils/`: Houses application-wide configuration logic, like `SetupLogger()`.
    *   `controllers/`: Houses the HTTP handlers. Reads requests, coordinates logic, and returns HTTP responses.
    *   `main.go`: The orchestrator. It imports packages, wires dependencies, and starts the server.

### 2. Dependency Injection (DI)
*   **The Concept:** Global variables (like `var logger *slog.Logger` or `var db *sql.DB`) cause race conditions and make automated testing impossible. Instead, dependencies should be "injected" into the structs that need them.
*   **My Implementation:** 
    *   Created controller structs (e.g., `type TaskController struct { Logger *slog.Logger }`).
    *   Created constructor functions (e.g., `NewTaskController(logger)`).
    *   Attached the HTTP handlers as *methods* on the controller struct (`func (ctrl *TaskController) CreateTask(...)`), allowing the handler to access the injected logger via `ctrl.Logger`.

### 3. Encapsulation & Exporting
*   **The Concept:** Go uses capitalization to determine visibility. 
*   **My Implementation:** Renamed functions from `createTask` to `CreateTask` so they are exported and accessible to the router in `main.go`.