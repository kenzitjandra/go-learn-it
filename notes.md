## Phase 2: Gin Request Logging Middleware

### 1. Bypassing Default Loggers
*   **Concept:** Standard frameworks often come with built-in text loggers. 
*   **Current Implementation:** Swapped `router := gin.Default()` for `router := gin.New()` (along with `gin.Recovery()`) to completely disable Gin's unstructured terminal output.

### 2. Centralized Middleware Logging
*   **Concept:** It starts a timer when a request arrives, hands control to the handler using `c.Next()`, and logs the metrics when the handler finishes.

### 3. Handling Internal Errors vs. Client Responses
*   **Concept:** When a backend operation fails, the system warns both the client and internal team
*   **Status Code Log Rule:**
    *   **400 (Bad Request) / 401 (Unauthorized) -> `logger.Warn()`**: The client made a mistake (bad JSON, invalid token). The server handled it correctly and is still healthy.
    *   **500 (Internal Server Error) -> `logger.Error()`**: The server failed (database connection dropped). This requires immediate development attention and should trigger an alert.