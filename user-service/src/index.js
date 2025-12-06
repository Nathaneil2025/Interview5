const express = require('express');

const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.json());

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'healthy', service: 'user-service' });
});

// Get all users
app.get('/users', (req, res) => {
  res.json([
    { id: 1, name: 'John Doe', email: 'john@example.com' },
    { id: 2, name: 'Jane Doe', email: 'jane@example.com' }
  ]);
});

// Get user by ID
app.get('/users/:id', (req, res) => {
  const userId = parseInt(req.params.id);
  if (userId === 1) {
    res.json({ id: 1, name: 'John Doe', email: 'john@example.com' });
  } else {
    res.status(404).json({ error: 'User not found' });
  }
});

// Create user
app.post('/users', (req, res) => {
  const { name, email } = req.body;
  if (!name || !email) {
    return res.status(400).json({ error: 'Name and email are required' });
  }
  res.status(201).json({ id: 3, name, email });
});

// Only start server if not in test mode
if (process.env.NODE_ENV !== 'test') {
  app.listen(PORT, () => {
    console.log(`User service running on port ${PORT}`);
  });
}

module.exports = app;
