# Proletariat Budget API

## Description

Proletariat Budget is a comprehensive personal finance management system designed to help individuals and households track their finances, manage expenses, plan savings, and achieve financial goals. The API provides a robust backend for financial tracking, offering features such as expense management, income tracking, account balancing, savings goals, and household member management.

This project aims to democratize financial planning tools, making them accessible to everyone regardless of their financial status. The name "Proletariat Budget" reflects our mission to provide powerful budgeting tools to the working class.

To know more, read our [mission statement](./docs/mission-statement.md)

## Tech Stack

- **Backend**: Go (Golang)
- **API Design**: OpenAPI 3.0.2
- **Database**: MySQL (interchangeable using the adapters)
- **API Documentation**: Swagger UI
- **Containerization**: Docker

## Required Tools

- **Go**: Version 1.22 or higher
- **Node.js**: Version 20+ (for bundling OpenAPI specifications using `@redocly/cli`)
- **Docker**: Latest version for containerization
- **MySQL**: Version 8.0 or higher
- **Make**: For running build commands

## Project Structure
```
proletariat-budget-be/
├── build-context/       # Files needed during Docker build
├── config/              # Application configuration
├── core/                # Core business logic
│   ├── domain/          # Domain models and error handling
│   ├── port/            # Interfaces defining application use cases
│   └── usecase/         # Implementation of use cases
├── internal/            # Internal packages
│   ├── adapter/         # Adapters for external systems (MySQL, HTTP)
│   ├── common/          # Common utilities and helpers
│   └── resthttp/        # HTTP server implementation and middleware
├── localenv/            # Local development environment setup
│   └── mysql_data/      # Local MySQL data storage
├── migrations/          # Database migration scripts
├── openapi/             # OpenAPI specification files
│   ├── components/      # Reusable components (schemas, responses)
│   └── paths/           # API path definitions
├── Dockerfile           # Docker configuration
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── main.go              # Application entry point
├── Makefile             # Build and deployment commands
└── README.md            # Project documentation
```

## Local Development

1. Clone the repo
``` bash
git clone https://github.com/ghorkov32/proletariat-budget-be.git
cd proletariat-budget-be
```

2. Start the local development environment:
``` bash
make dev
```
3. Start the local development environment:
``` bash
make run
```

## Building for Production
``` bash
make build
```

## Documentation

This project includes comprehensive documentation to help you understand and work with the Proletariat Budget API:

### API Documentation

The API is fully documented using OpenAPI 3.0.2 specifications:

- **Swagger UI**: When running the application, access the interactive API documentation at [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- **OpenAPI Specs**: Browse the raw OpenAPI specifications in the [`openapi/`](./openapi/) directory
  - [`openapi.yaml`](./openapi/openapi.yaml): Main specification file
  - [`components/`](./openapi/components/): Reusable schema definitions and responses
  - [`paths/`](./openapi/paths/): Individual API endpoint definitions

### Architecture Documentation

- **Hexagonal Architecture**: This project follows a hexagonal (ports and adapters) architecture pattern. Learn more in the [architecture guide](./docs/architecture.md).
- **Domain Model**: Understand the core domain concepts in the [domain model documentation](./docs/domain-model.md).
- **ADRs**: Understand our Architectural decisions in the [ADRs folder](./docs/ADRs)

### Development Guides

- **Contributing Guide**: Guidelines for contributing to the project in [CONTRIBUTING.md](./CONTRIBUTING.md)
- **Code Style**: Go code style conventions in [docs/code-style.md](./docs/code-style.md)
- **Testing Strategy**: Overview of testing approach in [docs/testing.md](./docs/testing.md)

### Database

- **Schema**: Database schema documentation in [migrations/README.md](./migrations/README.md)
- **Migrations**: How to run and manage database migrations in [docs/migrations.md](./docs/migrations.md)

### Deployment

- **Docker**: Instructions for Docker deployment in [docs/deployments/docker-deployment.md](./docs/deployments/docker-deployment.md)
- **AWS**: Instructions for AWS Deployment in [docs/deployments/aws-deployment.md](./docs/deployments/aws-deployment.md)

### Examples

- **API Usage Examples**: Common API usage patterns in [docs/examples/](./docs/examples/)
- **Integration Examples**: Examples of integrating with other systems in [docs/integration/](./docs/integration/)