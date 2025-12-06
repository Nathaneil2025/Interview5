from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Optional
from datetime import datetime

app = FastAPI(title="Transaction Service", version="1.0.0")


class Transaction(BaseModel):
    id: Optional[int] = None
    user_id: int
    amount: float
    description: str
    created_at: Optional[datetime] = None


class TransactionCreate(BaseModel):
    user_id: int
    amount: float
    description: str


# In-memory storage for demo
transactions_db = [
    {
        "id": 1,
        "user_id": 1,
        "amount": 100.50,
        "description": "Initial deposit",
        "created_at": "2024-01-01T00:00:00"
    },
    {
        "id": 2,
        "user_id": 1,
        "amount": -25.00,
        "description": "Coffee shop",
        "created_at": "2024-01-02T00:00:00"
    }
]


@app.get("/health")
def health_check():
    return {"status": "healthy", "service": "transaction-service"}


@app.get("/transactions")
def get_transactions():
    return transactions_db


@app.get("/transactions/{transaction_id}")
def get_transaction(transaction_id: int):
    for txn in transactions_db:
        if txn["id"] == transaction_id:
            return txn
    raise HTTPException(status_code=404, detail="Transaction not found")


@app.get("/transactions/user/{user_id}")
def get_user_transactions(user_id: int):
    user_txns = [txn for txn in transactions_db if txn["user_id"] == user_id]
    return user_txns


@app.post("/transactions", status_code=201)
def create_transaction(transaction: TransactionCreate):
    if transaction.amount == 0:
        raise HTTPException(status_code=400, detail="Amount cannot be zero")

    new_txn = {
        "id": len(transactions_db) + 1,
        "user_id": transaction.user_id,
        "amount": transaction.amount,
        "description": transaction.description,
        "created_at": datetime.now().isoformat()
    }
    transactions_db.append(new_txn)
    return new_txn


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
