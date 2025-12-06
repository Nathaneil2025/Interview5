#!/bin/bash
export PATH="$PATH:/var/lib/jenkins/.local/bin:$HOME/.local/bin:$(go env GOPATH)/bin"
set -e

SERVICE=$1

if [ -z "$SERVICE" ]; then
    echo "Usage: $0 <service-name>"
    exit 1
fi

echo "=== Running tests for $SERVICE ==="

cd "$SERVICE"

case "$SERVICE" in
    user-service)
        echo "Running Jest tests for Node.js..."
        npm install
        npm run test:ci
        ;;
    transaction-service)
        echo "Running pytest for Python..."
        pip install -r requirements.txt -q
        pytest tests/ -v --cov=app --cov-report=xml --cov-report=term --junitxml=reports/junit.xml
        mkdir -p reports
        mv coverage.xml reports/ 2>/dev/null || true
        ;;
    notification-service)
        echo "Running go test for Go..."
        go test -v -coverprofile=coverage.out -covermode=atomic ./...
        go tool cover -func=coverage.out
        # Convert to cobertura format if needed
        if command -v gocover-cobertura &> /dev/null; then
            gocover-cobertura < coverage.out > coverage.xml
        fi
        ;;
    *)
        echo "Unknown service: $SERVICE"
        exit 1
        ;;
esac

echo "=== Tests completed for $SERVICE ==="
