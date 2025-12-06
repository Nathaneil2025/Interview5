#!/bin/bash
export PATH="$PATH:/var/lib/jenkins/.local/bin:$HOME/.local/bin:$(go env GOPATH)/bin"
set -e

SERVICE=$1

if [ -z "$SERVICE" ]; then
    echo "Usage: $0 <service-name>"
    exit 1
fi

echo "=== Running security scans for $SERVICE ==="

cd "$SERVICE"

# Run secrets detection with gitleaks
echo "Running gitleaks for secrets detection..."
if command -v gitleaks &> /dev/null; then
    gitleaks detect --source=. --no-git -v || true
else
    echo "gitleaks not installed, skipping secrets detection"
fi

# Run language-specific SAST
case "$SERVICE" in
    user-service)
        echo "Running npm audit for Node.js..."
        npm install
        npm audit --audit-level=high || true
        ;;
    transaction-service)
        echo "Running bandit for Python..."
        pip install bandit -q
        bandit -r app/ -f json -o bandit-report.json || true
        bandit -r app/ -ll || true
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

echo "=== Security scans completed for $SERVICE ==="