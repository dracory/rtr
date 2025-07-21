**Prompt:**

As an expert in Go HTTP routers, your task is to provide a comprehensive and objective comparison between the `dracory/rtr` project (the current router) and the `gouniverse/router` project (https://github.com/gouniverse/router).

For the `dracory/rtr` project, refer to the `GEMINI.md` file for detailed information on its features, architecture, and design principles.

Present your comparison in a clear, concise Markdown table format, followed by a summary.

### Comparison Criteria:

1.  **Performance**:
    *   How do they handle high traffic and concurrent requests?
    *   Are there any available benchmark results or performance considerations?

2.  **Ease of Use**:
    *   How straightforward is the setup and configuration process for each router?
    *   What is the typical learning curve for developers?

3.  **Features**:
    *   What are the core features offered by each router (e.g., routing capabilities, handler types, error handling)?
    *   Do they support advanced routing features like middleware, route grouping, nested groups, domain-based routing, and declarative configuration?

4.  **Middleware System**:
    *   Describe the middleware architecture and execution order for each.
    *   What types of middleware are supported (e.g., Before, After, Recovery)?

5.  **Extensibility**:
    *   How easy is it to extend the functionality of each router?
    *   Are there provisions for custom handlers, middleware, or plugins?

6.  **Community and Support**:
    *   Assess the activity level of the community for each project.
    *   What kind of documentation, examples, and support resources are available?

7.  **Security**:
    *   What built-in security features do they offer?
    *   How do they approach common web security concerns?

8.  **Use Cases**:
    *   What types of projects or applications are best suited for each router?
    *   Are there any notable real-world examples or companies using them?

9.  **Integration**:
    *   How well do they integrate with other Go libraries, frameworks, and the standard `net/http` package?
    *   Are there any known compatibility issues or limitations?

### Output Format:

Please provide the comparison in a Markdown table with the following columns: "Criteria", "dracory/rtr", and "gouniverse/router".

### Summary:

Conclude with a brief summary (2-3 paragraphs) highlighting the key strengths and weaknesses of each router, and provide a recommendation on when to choose one over the other.