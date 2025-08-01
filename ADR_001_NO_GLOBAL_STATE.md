# ADR-001: Prohibit Global Mutable State in Interactive Components

## Status
ACCEPTED

## Context
During development of real-time progress updates, global mutable state was introduced leading to:
- Race conditions with data corruption potential
- Memory leaks from uncleaned resources  
- Security vulnerabilities
- Architectural pattern violations
- Non-deterministic behavior

The original implementation was cleaner, safer, and more maintainable.

## Decision
**PROHIBIT** the use of global mutable state in all interactive components.

### Forbidden Patterns
```go
// ❌ FORBIDDEN: Global mutable state
var globalCleanupState *cleanupState
var sharedProgressData map[string]interface{}
var globalTimers []*time.Timer
```

### Required Patterns
```go
// ✅ REQUIRED: State within Bubble Tea model
type model struct {
    cleanupProgress *cleanupProgress
    // All state encapsulated here
}

// ✅ REQUIRED: Message passing for communication
type progressUpdateMsg struct {
    project string
    progress float64
}
```

## Consequences
### Positive
- Thread safety guaranteed
- Deterministic behavior
- Easy testing and reasoning
- No resource leaks
- Proper framework pattern compliance

### Negative
- Slightly more verbose code
- Need to pass state through message system

## Compliance
All future changes must be reviewed for:
1. No global mutable variables
2. All state in model structs
3. Communication via message passing only
4. Proper resource cleanup patterns

## Enforcement
- Code review checklist includes this check
- Automated linting for global var declarations
- Test coverage for concurrent access patterns

---
**Author:** Cache Remover Utility Team  
**Date:** August 1, 2025  
**Reason for ADR:** Critical security and stability issues from global state