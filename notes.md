# Go Production Backend Notes: Phase 1 (slog basics)

## 1. Structured Logging
*   **Concept:** Structured logs output machine-readable key-value pairs (JSON).
*   **Current Implementation:** Used `slog.NewJSONHandler(os.Stdout, nil)` to ensure all outputs are strictly JSON.

## 2. Global Attributes
*   **Concept:** Injecting metadata into the logger instance once globally to apply to every output.
*   **Current Implementation:** 
    ```go
    logger.With(
        slog.String("service", "task-api"),
        slog.String("environment", "development"),
    )
    ```

## 3. Business Event Logging
*   **The Concept:** Logs should tell a story about business operations, not narrate every line of code executing.
*   **Current Implementation:** Instead of logging "entering createTask function", better to log the milestone.
    ```go
    logger.Info(
        "task created",
        slog.Int("task_id", newTask.ID),
        slog.String("title", newTask.Title),
    )
    ```

## 4. INFO vs WARN vs ERROR logs
*   **INFO (System):** Used for expected events (`"server starting"`) and business actions (`"health check requested"`).
*   **WARN (Client Errors):** Used when something went wrong, but the backend handled it correctly (system is still healthy).
    *   *My Example:* When a user sends bad JSON, `c.ShouldBindJSON` catches it. The server didn't crash; the client just made a mistake. This is a `WARN`.
*   **ERROR (Server Failures):** For infrastructure failures. In production, `ERROR` logs trigger alerts. Never log an `ERROR` for a client-side mistake.

## 5. Context Logging
*   **Concept:** Functions like `slog.ErrorContext()` accept a `context.Context` object, allowing the logger to automatically extract request data.