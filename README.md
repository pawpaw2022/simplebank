<!-- @format -->

# Simple Banking System Backend Service

Welcome to the Simple Banking System backend service repository. This project is built in Go and comprises three main components: the database (DB), REST APIs, and deployment configurations.

> **Annoucement**: As of 2023-09-01, the service of this project will be terminated from AWS as EKS is too expensive to maintain. The project is meant to be a learning experience and not a production-ready application. If you wish to continue using the project, please fork the repository and deploy it on your own.

## Overview

The Simple Banking System is designed to provide basic banking functionalities, including user sign-up, login, account creation, and money transfers between accounts.

### Components

#### Database (DB)

- Primary database: PostgreSQL
- SQL generation: sqlc
- Unit testing: Go testing package
- Transaction handling: Database transactions for atomic actions
- Continuous integration (CI): GitHub Actions for automated testing with Go and PostgreSQL

#### REST APIs

- Framework: Gin (a popular Go framework for building APIs)
- Configuration: Viper for loading configuration
- Integration testing: Mock DB for comprehensive test coverage
- Password hashing: Bcrypt for secure password storage
- Authentication: PASETO tokens for enhanced security
- Middleware: Authentication middleware to protect specific routes

#### Deployment

- Docker: Build Docker images for the application
- Docker Compose: Orchestrating DB and API server startup
- Continuous integration (CI): GitHub Actions for building and pushing images to AWS ECR
- Database migration: Migrate local DB to AWS RDS
- Environment variables: Stored in AWS Secret Manager
- Kubernetes: AWS EKS for container orchestration, managed with kubectl and k9s
- Custom domain: Achieved using AWS Route 53
- Traffic routing: Kubernetes Ingress for directing traffic to different services
- TLS security: Let's Encrypt for obtaining and renewing free SSL/TLS certificates
- Continuous deployment (CD): GitHub Actions for deploying new commits to Kubernetes

## Getting Started

Follow these steps to set up the project locally:

1. Clone this repository to your local machine.
2. Install the necessary dependencies (Go, Docker, kubectl, etc.).
3. Configure your environment variables using AWS Secret Manager.
4. Build the Docker image for the application using the provided Dockerfile.
5. Start the database and API server using Docker Compose.
6. Access the API endpoints and interact with the banking system.

For detailed setup instructions and usage guidelines, refer to the project's documentation.

## Contributing

We welcome contributions to enhance the Simple Banking System. If you want to contribute, please follow the standard GitHub workflow: fork the repository, create a feature branch, make your changes, and submit a pull request. Don't forget to add tests for new features or bug fixes.

## License

This project is licensed under the [MIT License](LICENSE).

---

Thank you for choosing the Simple Banking System backend service for your project. If you have any questions or need assistance, feel free to reach out to us.

[Project Repository](https://github.com/pawpaw2022/simplebank)
