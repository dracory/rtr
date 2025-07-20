# GEMINI.md for `dracory/rtr` Project

This file provides context and guidelines for Gemini when interacting with the `dracory/rtr` Go HTTP router project.

## Project Context

This project is an HTTP router package written in Go. Its core functionalities include:
* Route management for all standard HTTP methods.
* Route grouping with shared prefixes and middleware.
* Extensive middleware support (pre-route, post-route, panic recovery).
* Nested routing structures.
* A flexible, chainable API.
* Integration via `http.Handler` interface.
* Declarative configuration options.
* Middleware name detection for readability.

The primary goal of this project is to provide a robust, flexible, and well-organized routing solution for Go web applications.

## Preferred Language

The primary development language for this project is **Go**.

## Code Style Guidelines

When generating or modifying Go code related to this project, please adhere to the following general Go style guidelines:

* **`go fmt` and `goimports`:** Ensure all generated Go code is formatted using `go fmt` and organized with `goimports` conventions (though you don't need to run these tools, just follow their output style).
* **Clear Variable Names:** Use descriptive variable and function names.
* **Error Handling:** Always handle errors explicitly. Do not ignore returned errors.
* **Comments:** Provide clear and concise comments for complex logic, public functions, and exported types.
* **Modularity:** Strive for modular and testable code.
* **Concurrency:** Use Go's concurrency primitives (`goroutines`, `channels`) appropriately and safely.

## Custom Prompts / Interaction Guidelines

When assisting with this project, please consider the following:

* **Focus on Go:** Prioritize Go-specific solutions, best practices, and idiomatic Go code.
* **Router Context:** When discussing routing, assume the context of an HTTP router and its common patterns (middleware, handlers, routes).
* **Performance & Scalability:** Keep performance and scalability in mind for web-related code.
* **Security:** Suggest secure coding practices, especially concerning HTTP request handling.
* **Testing:** Emphasize the importance of testing and suggest appropriate testing strategies for Go.

## Preferred Tools

When answering queries or generating content related to this project, you are encouraged to use:

* `Google Search`: For researching Go language features, standard library usage, or common web development patterns in Go.
* `Browse`: For examining external Go documentation, related GitHub repositories, or articles on web development in Go.
* `Workspace`: If specific code snippets or documentation from external sources are relevant.
