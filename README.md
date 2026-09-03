## Phase 3: Context Logging & Tracing

### 1. Concurrency Problem
*   **Issue:** Requests and logs can happen at similar times, making it impossible to read chronologically.

### 2. Request IDs (Tracing)
*   **Concept:** Assigning a randomly generated ID (like a UUID or timestamp) to an HTTP request as it hits the server.
*   **Current Implementation:** Created a `requestIDMiddleware` that generates an ID and injects it into Gin's internal `context` (using `c.Set()`), as well as attaching it to the outgoing HTTP headers so the client receives it too.

### 3. Context Extraction
*   **Concept:** Because the Request ID is inside the request context, any handler (like `healthHandler`) or downstream middleware (like `slogMiddleware`) can extract it using `c.Get("request_id")`.
*   **Function:** By appending `slog.Any("request_id", reqID)` to every single log, log aggregators can maybe group logs together. If a request fails, search for that specific ID and see the exact lifecycle of that single request, ignoring the thousands other disconnected requests on the system.
