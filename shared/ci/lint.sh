#!/bin/bash
set -e

SERVICE=$1

if [ -z "$SERVICE" ]; then
    echo "Usage: $0 <service-name>"
    exit 1
fi

# Add local bin directories to PATH
export PATH="$PATH:/var/lib/jenkins/.local/bin:$HOME/.local/bin:$(go env GOPATH)/bin"

echo "=== Running lint for $SERVICE ==="

cd "$SERVICE"

case "$SERVICE" in
    user-service)
        echo "Running ESLint for Node.js..."
        npm install
        npm run lint
        ;;
    transaction-service)
        echo "Running Flake8 for Python..."
        pip install flake8 -q
        flake8 app/ tests/ --count --show-source --statistics
        ;;
    notification-service)
        echo "Running go vet for Go..."
        go vet ./...
        ;;
    *)
        echo "Unknown service: $SERVICE"
        exit 1
        ;;
esac

echo "=== Lint completed for $SERVICE ==="