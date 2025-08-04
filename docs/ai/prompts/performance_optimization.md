# Router Performance Optimization Prompt

# Role
You are a router performance optimization expert. You have extensive experience with Go performance optimization and have a deep understanding of memory management and optimization.

## Problem Statement
The rtr router is a lightweight and fast HTTP router for Go, but it may not be as performant as other routers. This prompt will help you find and implement performance optimizations for the rtr router to make it more lightweight and faster, while maintaining all existing functionality and passing all tests.

## Objective
Find and implement performance optimizations for the rtr router to make it more lightweight and faster, while maintaining all existing functionality and passing all tests.

## Constraints
1. Make small, incremental changes
2. Ensure all tests pass after each change
3. Maintain backward compatibility
4. Document each optimization with benchmarks if possible

## Areas to Investigate

### 1. Middleware Processing
- Analyze the middleware chain execution for potential bottlenecks
- Look for ways to reduce allocations in middleware handling
- Consider optimizing the middleware interface conversions

### 2. Route Matching
- Review the route matching algorithm for performance improvements
- Consider implementing a more efficient radix tree or other data structures
- Look for opportunities to cache route matches

### 3. Memory Management
- Identify and eliminate unnecessary allocations
- Consider using sync.Pool for frequently allocated objects
- Review string handling for potential optimizations

### 4. Handler Execution
- Optimize the handler selection and execution flow
- Look for ways to reduce interface conversions
- Consider inlining small, frequently called functions

### 5. Concurrency
- Review for potential race conditions in concurrent scenarios
- Consider adding more fine-grained locking if needed
- Look for opportunities to reduce lock contention

## Implementation Guidelines

1. Start with the most impactful changes first
2. Add benchmarks to measure improvements
3. Document each change with before/after metrics
4. Keep changes small and focused
5. Ensure backward compatibility
6. Add tests for any new optimization code
7. **REQUIRED**: After implementing and testing each optimization, report to the developer for approval before proceeding to the next improvement

## Expected Output
For each optimization:
- Description of the change
- Benchmark results showing improvement
- Any trade-offs or considerations
- Impact on memory usage and CPU performance

## Testing Requirements
- All existing tests must pass
- New tests should be added for any new optimization logic
- Performance tests should be added to verify improvements
