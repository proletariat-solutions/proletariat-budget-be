## Pull Request: [Brief, Descriptive Title]

**Closes #[Issue Number(s)]** (e.g., Closes #123, Closes #456)

---

### What does this PR do?

* [Concise summary of the changes introduced by this PR. Be specific about the functionality added, modified, or fixed.]
* [If applicable, describe any new features, bug fixes, or performance improvements.]

### Why is this change necessary?

* [Explain the problem this PR solves or the new functionality it provides.]
* [Reference the linked issue(s) and provide context.]

### How has this been tested?

* **Unit Tests:**
    * [ ] All existing unit tests pass.
    * [ ] New unit tests have been added for the new/modified logic.
    * [Describe the scope of your unit tests and what they cover.]
* **Integration Tests:**
    * [ ] Existing integration tests pass.
    * [ ] New integration tests have been added (if applicable, e.g., for new API endpoints).
    * [Describe how you tested the integration of different components/layers.]
* **Manual Testing (if applicable):**
    * [Describe any manual testing performed, including steps to reproduce/verify the changes.]
    * [Provide screenshots/recordings if helpful.]

### Code Quality Checklist

* **Go Linting:**
    * [ ] `golangci-lint run` passes without errors or warnings.
    * [Explain any specific linters you focused on or addressed.]
* **Code Formatting:**
    * [ ] `go fmt` has been run on all modified files.
    * [ ] Code adheres to Go conventions and style guides.
* **Complexity:**
    * [ ] No significant increase in cyclomatic complexity (if applicable, mention tools used like `gocognit` or `go tool vet -composites`).
    * [ ] Refactored complex functions into smaller, more manageable units.
* **Error Handling:**
    * [ ] Appropriate error handling is implemented for all potential failure points.
    * [ ] Errors are returned and propagated correctly.
    * [ ] Context is used for error tracing where appropriate.
* **Logging:**
    * [ ] Meaningful log messages are present for key events and errors.
    * [ ] Sensitive information is not logged.
* **Database Interactions (if applicable):**
    * [ ] SQL queries are optimized and avoid N+1 problems.
    * [ ] Transactions are used appropriately to ensure data integrity.
    * [ ] Database migrations are included (if schema changes are required).
* **Concurrency (if applicable):**
    * [ ] Proper use of goroutines and channels to avoid race conditions.
    * [ ] Mutexes or other synchronization primitives are used correctly.
* **Dependencies:**
    * [ ] `go mod tidy` has been run.
    * [ ] No unnecessary dependencies have been introduced.
    * [ ] Dependencies are up-to-date (within project guidelines).

### Hexagonal Architecture Considerations

* **Domain Layer:**
    * [ ] Business logic resides solely within the domain layer (core).
    * [ ] No external dependencies are introduced into the domain.
    * [ ] Domain models are clean and independent.
* **Ports & Adapters:**
    * [ ] New interfaces (ports) are defined for external interactions.
    * [ ] Concrete implementations (adapters) interact with external services (databases, APIs, etc.).
    * [ ] Adherence to dependency inversion principle (ports depend on domain, adapters depend on ports).
* **Cross-Cutting Concerns:**
    * [ ] Cross-cutting concerns (e.g., logging, metrics, authentication) are handled by appropriate middleware or decorators.

### Screenshots (if applicable)

[Add any relevant screenshots or GIFs to demonstrate the changes.]

---

### Reviewer Checklist (to be filled by reviewer)

* [ ] Code is readable, well-commented, and follows Go best practices.
* [ ] All tests pass and adequately cover the changes.
* [ ] Linting checks pass.
* [ ] Architecture principles (Hexagonal) are maintained.
* [ ] The PR addresses the stated problem/feature.
* [ ] No regressions introduced.
* [ ] Documentation (if any) is updated.
