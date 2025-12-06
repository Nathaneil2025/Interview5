# Monorepo CI Pipeline

Production-grade CI pipeline for a monorepo with multiple microservices using Jenkins.

## Repository Structure

```
root/
├── user-service/           # Node.js (Express)
│   ├── src/
│   ├── tests/
│   ├── Dockerfile
│   └── package.json
├── transaction-service/    # Python (FastAPI)
│   ├── app/
│   ├── tests/
│   ├── Dockerfile
│   └── requirements.txt
├── notification-service/   # Go (gorilla/mux)
│   ├── main.go
│   ├── main_test.go
│   ├── Dockerfile
│   └── go.mod
├── shared/
│   └── ci/
│       ├── lint.sh
│       ├── test.sh
│       └── scan.sh
├── Jenkinsfile
└── README.md
```

## Pipeline Features

### Intelligent Change Detection
- Automatically detects which services were modified using `git diff`
- Only runs CI stages for modified services
- Changes to `shared/` or `Jenkinsfile` trigger all services

### Pipeline Stages

1. **Detect Changes** - Identifies modified services
2. **Lint & Code Quality** - Language-specific linting
   - ESLint for Node.js
   - Flake8 for Python
   - golangci-lint for Go
3. **Unit Testing** - With coverage reports
   - Jest for Node.js
   - pytest for Python
   - go test for Go
4. **Security Scans**
   - gitleaks for secrets detection
   - npm audit / bandit / go vet for SAST
5. **Docker Build** - Build images tagged with commit SHA
6. **Push to Registry** - Push to configured registry (main/develop only)
7. **Manual Approval** - Gate before deployment (main branch only)

### Pipeline Characteristics
- Parallel execution per microservice
- Retry logic for flaky stages
- Fail-fast on lint/test/scan errors
- Slack/Teams notifications (configurable)
- Modular with reusable Groovy functions

## Local Development

### User Service (Node.js)
```bash
cd user-service
npm install
npm run lint
npm test
npm start
```

### Transaction Service (Python)
```bash
cd transaction-service
pip install -r requirements.txt
flake8 app/
pytest
uvicorn app.main:app --reload
```

### Notification Service (Go)
```bash
cd notification-service
go mod download
go vet ./...
go test -v ./...
go run main.go
```

## Docker Build

```bash
# Build individual services
docker build -t user-service:latest ./user-service
docker build -t transaction-service:latest ./transaction-service
docker build -t notification-service:latest ./notification-service
```

## Jenkins Setup

### Required Plugins
- Pipeline
- Git
- Docker Pipeline
- JUnit
- HTML Publisher
- Slack Notification (optional)
- Office 365 Connector (optional)

### Configuration
1. Create a new Pipeline job
2. Configure SCM to point to this repository
3. Set branch specifier to `*/main` or `*/develop`
4. Enable webhook trigger or poll SCM

## CI Scripts

The `shared/ci/` directory contains reusable scripts:

- `lint.sh <service>` - Run linting for a service
- `test.sh <service>` - Run tests with coverage
- `scan.sh <service>` - Run security scans
