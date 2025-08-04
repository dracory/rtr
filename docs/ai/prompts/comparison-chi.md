# Generate Router Comparison Document

Routers:
Router A: dracory/rtr
Router B: chi

Create a detailed technical comparison document between the two specified Go HTTP routers. The document should follow this structure. In the following sections, use {ROUTER_A} and {ROUTER_B} to refer to the routers being compared.

## 1. Overview
- Brief introduction to both routers
- Project maturity and maintenance status
- Core design philosophies

## 2. Feature Comparison Table
Create a markdown table comparing these aspects:
- Basic routing capabilities
- Path parameters and pattern matching
- Route grouping and middleware support
- Performance characteristics
- Built-in response helpers
- Special features (WebSocket, static file serving, etc.)
- Testing support
- Documentation quality
- Community and ecosystem

## 3. Code Examples
Provide side-by-side code examples for common routing scenarios:

### Basic Route Definition
```go
// {ROUTER_A} example

// {ROUTER_B} example
```

### Middleware Usage
```go
// {ROUTER_A} middleware example

// {ROUTER_B} middleware example
```

### Route Groups
```go
// {ROUTER_A} group example

// {ROUTER_B} group example
```

## 4. Performance Comparison
- Memory usage
- Request handling speed
- Benchmark results (if available)

## 5. Use Cases
- Best use cases for each router
- When to choose one over the other
- Limitations and trade-offs

## 6. Migration Guide
- Key differences to be aware of
- Common patterns that need adjustment
- Helpful tools or utilities for migration

## 7. Conclusion
- Summary of key differences
- Final recommendations based on different project needs

For each code example, include a brief explanation of what the code does and any important differences in implementation between the two routers.

Make sure to:
- Use proper Go code formatting
- Highlight any gotchas or important considerations
- Include relevant version information
- Be objective and factual in the comparison
- Note any significant performance implications

## Usage Instructions:
1. Replace `{ROUTER_A}` and `{ROUTER_B}` with the names of the routers you want to compare
2. Run this prompt through your preferred LLM
3. Save the output to `docs/comparisons/{router-a}-{router-b}-comparison-{timestamp}.md` (use kebab-case for router names)
   Example: For rtr vs gouniverse/router, save as `docs/comparisons/rtr-gouniverse-router-comparison-2025-07-22-18-25-20.md`
