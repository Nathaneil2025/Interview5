import pytest
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)


class TestHealthCheck:
    def test_health_endpoint(self):
        response = client.get("/health")
        assert response.status_code == 200
        assert response.json()["status"] == "healthy"
        assert response.json()["service"] == "transaction-service"


class TestTransactions:
    def test_get_all_transactions(self):
        response = client.get("/transactions")
        assert response.status_code == 200
        assert isinstance(response.json(), list)
        assert len(response.json()) > 0

    def test_get_transaction_by_id(self):
        response = client.get("/transactions/1")
        assert response.status_code == 200
        assert response.json()["id"] == 1

    def test_get_transaction_not_found(self):
        response = client.get("/transactions/9999")
        assert response.status_code == 404

    def test_get_user_transactions(self):
        response = client.get("/transactions/user/1")
        assert response.status_code == 200
        assert isinstance(response.json(), list)

    def test_create_transaction(self):
        new_txn = {
            "user_id": 1,
            "amount": 50.00,
            "description": "Test transaction"
        }
        response = client.post("/transactions", json=new_txn)
        assert response.status_code == 201
        assert response.json()["amount"] == 50.00

    def test_create_transaction_zero_amount(self):
        new_txn = {
            "user_id": 1,
            "amount": 0,
            "description": "Invalid transaction"
        }
        response = client.post("/transactions", json=new_txn)
        assert response.status_code == 400
